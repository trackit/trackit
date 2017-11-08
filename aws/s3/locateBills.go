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

// BillRepository is a location where the server may look for bill objects.
type BillRepository struct {
	Bucket string
	Prefix string
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
			if isBillKey(b) {
				out <- b
			}
		}
	}()
	return out
}

// billKey is the regexp for the name AWS gives to the bill objects we want to
// ingest.
var billKey = regexp.MustCompile(`\d+-aws-billing-detailed-line-items-with-resources-and-tags-\d{4}-\d{2}.csv.zip$`)

// isBillKey tests whether an S3 key looks like that of a bill object.
func isBillKey(b BillKey) bool {
	return billKey.MatchString(b.Key)
}

// listBillsFromRepository searches a repository for the keys of bill objects.
func listBillsFromRepository(
	ctx context.Context,
	sess *session.Session,
	r BillRepository,
) (<-chan BillKey, error) {
	s3svc, err := serviceForBucketRegion(ctx, sess, r)
	if err != nil {
		return nil, err
	}
	c := make(chan BillKey)
	go func() {
		defer close(c)
		l := jsonlog.LoggerFromContextOrDefault(ctx)
		input := s3.ListObjectsV2Input{
			Bucket: &r.Bucket,
			Prefix: &r.Prefix,
		}
		err := s3svc.ListObjectsV2PagesWithContext(ctx, &input, listBillsFromRepositoryPage(c, r, l))
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
	r BillRepository,
	l jsonlog.Logger,
) func(*s3.ListObjectsV2Output, bool) bool {
	count := 0
	return func(page *s3.ListObjectsV2Output, last bool) bool {
		for _, o := range page.Contents {
			c <- BillKey{
				Key:          *o.Key,
				Bucket:       r.Bucket,
				Region:       "",
				LastModified: *o.LastModified,
			}
		}
		count += len(page.Contents)
		if count < MaxCheckedKeysByRepository {
			return !last
		} else {
			l.Warning("Checked maximum amount of keys for repository.", r)
			return false
		}
	}
}

// serviceForBucketRegion determines the region an S3 bucket resides in an
// returns a S3 service instance for that region.
func serviceForBucketRegion(ctx context.Context, sess *session.Session, r BillRepository) (*s3.S3, error) {
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
		return s3.New(sess.Copy(&aws.Config{Region: &region})), nil
	} else {
		return nil, err
	}
}
