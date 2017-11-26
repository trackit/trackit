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
	"context"
	"fmt"
	"regexp"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/trackit/jsonlog"
)

// MaxCheckedKeysByRepository is the amount of keys inspected before we give
// up. If users have a massive bucket where their bills are stored alongside
// other keys, we don't to spend too much time reading the metadata of all
// keys. This means that it is the responsibility of the user to put their
// bills in a place where there isn't much of anything else.
const MaxCheckedKeysByRepository = 1000

// BillKey is a key where a bill object may be found.
type BillKey struct {
	Region       string
	Bucket       string
	Key          string
	LastModified time.Time
}

// LocateBills searches the repositories for bill objects and calls cbk for
// each.
func LocateBills(
	ctx context.Context,
	sess *session.Session,
	repositories []BillRepository,
	cbk func(context.Context, BillKey) bool,
) {
	for b := range listBills(ctx, sess, repositories) {
		if !cbk(ctx, b) {
			break
		}
	}
}

// mergecBillLocation implements the fan-in pattern by merging to the out
// channel the input from the channels read on cs.
func mergecBillLocation(out chan<- BillKey, cs <-chan <-chan BillKey) {
	var wg sync.WaitGroup
	for c := range cs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for u := range c {
				out <- u
			}
		}()
	}
	go func() {
		wg.Wait()
		close(out)
	}()
}

// listBills returns a channel where bill object locations can be received.
// Each repository will be searched in its own goroutine.
func listBills(
	ctx context.Context,
	sess *session.Session,
	repositories []BillRepository,
) <-chan BillKey {
	c := make(chan BillKey)
	cc := make(chan (<-chan BillKey))
	defer close(cc)
	go mergecBillLocation(c, cc)
	for _, r := range repositories {
		k, err := listBillsFromRepository(ctx, sess, r)
		if err == nil {
			cc <- k
		}
	}
	return c
}

// filterBillsByKey starts a goroutine which filters a channel of BillKey,
// keeping only those whose key indicates it may be a bill object.
func filterBillsByKey(in <-chan BillKey) <-chan BillKey {
	out := make(chan BillKey)
	go func() {
		defer close(out)
		for b := range in {
			if isAwsBillKey(b) {
				out <- b
			}
		}
	}()
	return out
}

// billKey is the regexp for the name AWS gives to the bill objects we want to
// ingest.
var billKey = regexp.MustCompile(`\d+-aws-billing-detailed-line-items-with-resources-and-tags-\d{4}-\d{2}.csv.zip$`)

// isAwsBillKey tests whether an S3 key looks like that of a bill object.
func isAwsBillKey(b BillKey) bool {
	return billKey.MatchString(b.Key)
}

// listBillsFromRepository searches a repository for the keys of bill objects.
func listBillsFromRepository(
	ctx context.Context,
	sess *session.Session,
	r BillRepository,
) (<-chan BillKey, error) {
	region, err := getBucketRegion(ctx, sess, r)
	if err != nil {
		return nil, err
	}
	s3svc := serviceForBucketRegion(sess, region)
	c := make(chan BillKey)
	go func() {
		defer close(c)
		l := jsonlog.LoggerFromContextOrDefault(ctx)
		input := s3.ListObjectsV2Input{
			Bucket: &r.Bucket,
			Prefix: &r.Prefix,
		}
		err := s3svc.ListObjectsV2PagesWithContext(ctx, &input, listBillsFromRepositoryPage(c, r, l, region))
		if err != nil {
			l.Error("Failed to list objects from bucket.", err.Error())
		}
	}()
	return filterBillsByKey(c), nil
}

// listBillsFromRepositoryPage handles a page of results for
// listBillsFromRepository. It will only trigger the processing of the next
// page if less than MaxCheckedKeysByRepository keys where encountered.
func listBillsFromRepositoryPage(
	c chan<- BillKey,
	br BillRepository,
	l jsonlog.Logger,
	region string,
) func(*s3.ListObjectsV2Output, bool) bool {
	count := 0
	return func(page *s3.ListObjectsV2Output, last bool) bool {
		for _, o := range page.Contents {
			c <- BillKey{
				Key:          *o.Key,
				Bucket:       br.Bucket,
				Region:       region,
				LastModified: *o.LastModified,
			}
		}
		count += len(page.Contents)
		if count < MaxCheckedKeysByRepository {
			return !last
		} else {
			l.Warning("Checked maximum amount of keys for repository.", br)
			return false
		}
	}
}

// serviceForBucketRegion determines the region an S3 bucket resides in and
// returns that as a string.
func getBucketRegion(ctx context.Context, sess *session.Session, r BillRepository) (string, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	s3svc := s3.New(sess)
	input := s3.GetBucketLocationInput{
		Bucket: &r.Bucket,
	}
	if output, err := s3svc.GetBucketLocationWithContext(ctx, &input); err == nil {
		region := *output.LocationConstraint
		logger.Debug(fmt.Sprintf("Found bucket region."), map[string]string{
			"bucket": r.Bucket,
			"region": region,
		})
		return region, nil
	} else {
		return "", err
	}
}

// serviceForBucketRegion returns an S3 service for the given region.
func serviceForBucketRegion(sess *session.Session, region string) *s3.S3 {
	return s3.New(sess.Copy(&aws.Config{Region: &region}))
}
