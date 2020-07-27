//   Copyright 2017 MSolution.IO
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package s3

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/trackit/jsonlog"

	taws "github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/config"
	"github.com/trackit/trackit/es"
	"github.com/trackit/trackit/es/indexes/lineItems"
	"github.com/trackit/trackit/util/csv"
)

const (
	// MaxCheckedKeysByRepository is the amount of keys inspected before we give
	// up. If users have a massive bucket where their bills are stored alongside
	// other keys, we don't to spend too much time reading the metadata of all
	// keys. This means that it is the responsibility of the user to put their
	// bills in a place where there isn't much of anything else.
	MaxCheckedKeysByRepository = 10000

	ReadBillsStsSessionName = "read-bills"
)

var (
	ErrUnsupportedCompression = errors.New("unsupported compression")
	httpClient                = http.Client{}
)

const maxManifestSize = 0x8000

type billTime time.Time

const billTimeFormat = `"20060102T150405Z"`

func (t *billTime) UnmarshalJSON(b []byte) error {
	tt, err := time.Parse(billTimeFormat, string(b))
	if err == nil {
		*t = billTime(tt)
	}
	return err
}

type manifest struct {
	SourceBucket  string   `json:"sourceBucket"`
	Bucket        string   `json:"bucket"`
	ReportKeys    []string `json:"reportKeys"`
	Compression   string   `json:"compression"`
	ReportName    string   `json:"reportName"`
	Account       string   `json:"account"`
	BillingPeriod struct {
		Start billTime `json:"start"`
		End   billTime `json:"end"`
	} `json:"billingPeriod"`
	LastModified time.Time
}

// BillKey is a key where a bill object may be found.
type BillKey struct {
	Region       string
	Bucket       string
	Key          string
	LastModified time.Time
}

type billRepositoryWithRegion struct {
	BillRepository
	Region string `json:"region"`
}

type OnLineItem func(lineItems.LineItem, bool)
type ManifestPredicate func(manifest, bool) bool

// ReadBills reads all LineItems from new bills in a BillRepository, and runs
// `oli` for each one.
func ReadBills(ctx context.Context, aa taws.AwsAccount, br BillRepository, oli OnLineItem, mp ManifestPredicate) (time.Time, error) {
	var lastManifest time.Time
	s3svc, brr, err := getServiceForRepository(ctx, aa, br)
	if err != nil {
		return lastManifest, err
	}
	jsonlog.LoggerFromContextOrDefault(ctx).Debug("Obtained S3 service to read bills.", map[string]interface{}{"account": aa, "billRepository": br})
	mck := getKeys(ctx, s3svc, brr)
	mck = getManifestKeys(ctx, mck)
	mc := getManifests(ctx, s3svc, mck)
	mc, lastManifestPromise := selectManifests(mp, mc)
	es.CleanCurrentMonthBillByBillRepositoryId(ctx, aa.UserId, br.Id)
	importBills(ctx, s3svc, mc, oli, mp)
	return <-lastManifestPromise, nil
}

// selectManifests returns a channel of all AWS manifest files which match
// `mp`.
func selectManifests(mp ManifestPredicate, mc <-chan manifest) (<-chan manifest, <-chan time.Time) {
	out := make(chan manifest)
	lmOut := make(chan time.Time, 1)
	go func() {
		defer close(out)
		defer close(lmOut)
		var lm time.Time
		for m := range mc {
			if mp(m, true) {
				out <- m
				if m.LastModified.After(lm) {
					lm = m.LastModified
				}
			}
		}
		lmOut <- lm
	}()
	return out, lmOut
}

// importBills imports LineItems for bill files described in manifests sent to
// the `manifests` channel.
func importBills(ctx context.Context, s3svc *s3.S3, manifests <-chan manifest, oli OnLineItem, mp ManifestPredicate) {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	outs, out := mergecdLineItem()
	for m := range manifests {
		l.Debug("Will attempt ingesting bills.", m)
		for _, s := range m.ReportKeys {
			l.Debug("Will attempt ingesting bill part.", map[string]interface{}{"key": s, "manifest": m})
			outs <- importBill(ctx, s3svc, s, m, mp)
		}
	}
	close(outs)
	for lineItem := range out {
		oli(lineItem, true)
	}
	oli(lineItems.LineItem{}, false)
}

// importBill imports LineItems for a single bill file.
func importBill(ctx context.Context, s3svc *s3.S3, s string, m manifest, mp ManifestPredicate) <-chan lineItems.LineItem {
	outs, out := mergecdLineItem()
	go func() {
		defer close(outs)
		ctx, cancel := context.WithCancel(ctx)
		l := jsonlog.LoggerFromContextOrDefault(ctx)
		reader, err := getBillReader(ctx, s3svc, s, m)
		if err != nil {
			l.Error("Failed to read bill.", err.Error())
		} else {
			l.Debug("Reading bill.", map[string]interface{}{"key": s, "manifest": m})
			outs <- readBill(ctx, cancel, reader, s, m, mp)
		}
	}()
	return out
}

// readBill returns a channel of all LineItems in a single bill file.
func readBill(ctx context.Context, cancel context.CancelFunc, reader io.ReadCloser, s string, m manifest, mp ManifestPredicate) <-chan lineItems.LineItem {
	out := make(chan lineItems.LineItem)
	go func() {
		defer reader.Close()
		defer close(out)
		csvDecoder := csv.NewDecoder(reader)
		for r := range records(ctx, &csvDecoder) {
			if mp(m, false) || r.InvoiceId == "" {
				out <- r
			}
		}
	}()
	return out
}

func records(ctx context.Context, d *csv.Decoder) <-chan lineItems.LineItem {
	out := make(chan lineItems.LineItem)
	log := jsonlog.LoggerFromContextOrDefault(ctx)
	go func() {
		defer close(out)
		if err := d.ReadHeader(); err != nil {
			log.Error("Failed to read CSV header.", err.Error())
			return
		}
		for {
			record, err := decodeRecord(d)
			if err == io.EOF {
				return // EOF was reached
			} else if err != nil {
				log.Error("Error reading CSV record.", err.Error())
				return
			} else {
				select {
				case out <- record:
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	return out
}

// decodeRecord decodes a LineItem from a csv.Reader.
func decodeRecord(d *csv.Decoder) (lineItems.LineItem, error) {
	var record lineItems.LineItem
	err := d.ReadRecord(&record)
	return record, err
}

// getBillReader returns a ReadCloser for a const and usage report. It will use
// the object described by the key s and the manifest m.
func getBillReader(ctx context.Context, s3svc *s3.S3, s string, m manifest) (io.ReadCloser, error) {
	switch m.Compression {
	case "GZIP":
		return getGzipBillReader(ctx, s3svc, s, m)
	default:
		jsonlog.LoggerFromContextOrDefault(ctx).Error("Unsupported  compression scheme.", map[string]interface{}{"key": s, "manifest": m})
		return nil, ErrUnsupportedCompression
	}
}

// getGzipBillReader returns a ReadCloser for a GZIP-compressed S3 object which
// is downloaded on the fly.
func getGzipBillReader(ctx context.Context, s3svc *s3.S3, s string, m manifest) (io.ReadCloser, error) {
	if reader, err := getRawBillReader(ctx, s3svc, s, m); err == nil {
		return gzip.NewReader(reader)
	} else {
		return nil, err
	}
}

// getRawBillReader gets an io.ReadCloser for the raw data from a billing
// file.
func getRawBillReader(ctx context.Context, s3svc *s3.S3, s string, m manifest) (io.ReadCloser, error) {
	var bucket = m.Bucket
	var key = s
	var region = *s3svc.Client.Config.Region
	var dumbS3mgr dumbS3Manager
	dumbS3mgr.init(s3svc.Client.Config.Credentials)
	return dumbS3mgr.rawS3GetObjectToReader(ctx, &httpClient, region, bucket, key)
}

// getManifests downloads the manifest whose keys are sent to the in channel.
// It immediately returns with a channel where manifest objects will be sent.
func getManifests(ctx context.Context, s3svc *s3.S3, in <-chan BillKey) <-chan manifest {
	outs, out := mergecdManifest()
	go func() {
		defer close(outs)
		var s3mgr dumbS3Manager
		s3mgr.init(s3svc.Client.Config.Credentials)
		for bk := range in {
			outs <- readManifest(ctx, &s3mgr, bk)
		}
	}()
	return out
}

// readManifest downloads and parses a manifest file asynchronously. Returns a
// channel where at most one manifest object will be sent, then the channel
// will be closed.
func readManifest(ctx context.Context, s3mgr *dumbS3Manager, bk BillKey) <-chan manifest {
	out := make(chan manifest)
	go func() {
		defer close(out)
		logger := jsonlog.LoggerFromContextOrDefault(ctx)
		buf, err := s3mgr.rawS3GetObjectToBuffer(ctx, &httpClient, bk.Region, bk.Bucket, bk.Key)
		if err != nil {
			logger.Error("Failed to download usage and cost manifest.", map[string]interface{}{"billKey": bk, "error": err.Error()})
			return
		} else {
			var m manifest
			err := json.Unmarshal(buf, &m)
			if err != nil {
				logger.Error("Failed to parse usage and cost manifest.", map[string]interface{}{"billKey": bk, "error": err.Error()})
				return
			} else {
				m.LastModified = bk.LastModified
				m.SourceBucket = bk.Bucket
				out <- m
			}
		}

	}()
	return out
}

// getServiceForRepository instantiates an *s3.S3 service from an AwsAccount
// and a billRepositoryWithRegion. It returns a nil error if the operation was
// successful.
func getServiceForRepository(ctx context.Context, aa taws.AwsAccount, br BillRepository) (*s3.S3, billRepositoryWithRegion, error) {
	var brr billRepositoryWithRegion
	creds, err := taws.GetTemporaryCredentials(aa, ReadBillsStsSessionName)
	if err != nil {
		return nil, brr, err
	}
	jsonlog.LoggerFromContextOrDefault(ctx).Debug("Obtained credentials to read bills.", map[string]interface{}{"awsAccount": aa, "billRepository": br})
	sess := session.New(&aws.Config{Credentials: creds, Region: &config.AwsRegion})
	region, err := getBucketRegion(ctx, sess, br)
	if err != nil {
		return nil, brr, err
	}
	brr.BillRepository = br
	brr.Region = region
	return serviceForBucketRegion(sess, region), brr, nil
}

// getKeys returns a channel where all keys from the billRepositoryWithRegion
// will be sent.
func getKeys(ctx context.Context, s3svc *s3.S3, brr billRepositoryWithRegion) <-chan BillKey {
	c := make(chan BillKey)
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	l.Debug("Getting manifest files from repository.", brr)
	go func() {
		defer close(c)
		input := s3.ListObjectsV2Input{
			Bucket: &brr.Bucket,
			Prefix: &brr.Prefix,
		}
		err := s3svc.ListObjectsV2PagesWithContext(ctx, &input, listBillsFromRepositoryPage(ctx, c, brr, l))
		if err != nil {
			l.Error("Failed to list objects from bucket.", err.Error())
		}
	}()
	return c
}

// manifestKeyRegex matches keys which look like manifest keys.
var manifestKeyRegex = regexp.MustCompile(`/\d{8}-\d{8}\/[^/]+-Manifest.json$`)

// getManifestKeys filters a channel of BillKey to only keep those which seem to
// be Cost And Usage manifests.
func getManifestKeys(ctx context.Context, in <-chan BillKey) <-chan BillKey {
	out := make(chan BillKey)
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	go func() {
		defer close(out)
		for bk := range in {
			if manifestKeyRegex.MatchString(bk.Key) {
				l.Debug("Found manifest key.", map[string]interface{}{"billKey": bk})
				select {
				case out <- bk:
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	return out
}

// listBillsFromRepositoryPage handles a page of results for
// listBillsFromRepository. It will only trigger the processing of the next
// page if less than MaxCheckedKeysByRepository keys where encountered.
// It will only send the new objects, thanks to LastModified in brr.
func listBillsFromRepositoryPage(
	ctx context.Context,
	c chan<- BillKey,
	brr billRepositoryWithRegion,
	l jsonlog.Logger,
) func(*s3.ListObjectsV2Output, bool) bool {
	count := 0
	return func(page *s3.ListObjectsV2Output, last bool) bool {
		for _, o := range page.Contents {
			if brr.LastImportedManifest.Before((*o.LastModified).AddDate(0, 1, 0)) {
				count += 1
				select {
				case c <- BillKey{
					Key:          *o.Key,
					Bucket:       brr.Bucket,
					Region:       brr.Region,
					LastModified: *o.LastModified,
				}:
				case <-ctx.Done():
					return false
				}
			}
		}
		if count < MaxCheckedKeysByRepository {
			return !last
		} else {
			l.Warning("Checked maximum amount of keys for repository.", brr)
			return false
		}
	}
}

// getBucketRegion determines the region an S3 bucket resides in and
// returns that as a string.
func getBucketRegion(ctx context.Context, sess *session.Session, r BillRepository) (string, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	s3svc := s3.New(sess)
	input := s3.GetBucketLocationInput{
		Bucket: &r.Bucket,
	}
	logger.Debug("Getting bucket region.", map[string]interface{}{
		"input":          input,
		"billRepository": r,
	})
	if output, err := s3svc.GetBucketLocationWithContext(ctx, &input); err == nil {
		region := getBucketRegionFromGetBucketLocationOutput(output)
		logger.Debug(fmt.Sprintf("Found bucket region."), map[string]string{
			"bucket": r.Bucket,
			"region": region,
		})
		return region, nil
	} else {
		return "", err
	}
}

// getBucketRegionFromGetBucketLocationOutput gets the region name for a bucket
// from a non-null *s3.GetBucketLocationOutput. It handles the API's special
// case where a nil LocationConstraint indicates the bucket is situated in the
// us-east-1 region.
func getBucketRegionFromGetBucketLocationOutput(output *s3.GetBucketLocationOutput) string {
	if output.LocationConstraint == nil || *output.LocationConstraint == "" {
		return "us-east-1"
	} else {
		return *output.LocationConstraint
	}
}

// serviceForBucketRegion returns an S3 service for the given region.
func serviceForBucketRegion(sess *session.Session, region string) *s3.S3 {
	return s3.New(sess.Copy(&aws.Config{Region: &region}))
}
