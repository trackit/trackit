package s3

import (
	"context"
	"regexp"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/trackit/jsonlog"
)

type BillKey struct {
	Region       string
	Bucket       string
	Key          string
	LastModified time.Time
}

type BillRepository struct {
	Bucket string
	Prefix string
}

func LocateBills(ctx context.Context, ss *s3.S3, repositories []BillRepository, cbk func(context.Context, BillKey) bool) {
	for b := range listBills(ctx, ss, repositories) {
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

func listBills(ctx context.Context, ss *s3.S3, repositories []BillRepository) <-chan BillKey {
	c := make(chan BillKey)
	cc := make(chan (<-chan BillKey))
	defer close(cc)
	go mergecBillLocation(c, cc)
	for _, r := range repositories {
		cc <- listBillsFromRepository(ctx, ss, r)
	}
	return c
}

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

var billKey = regexp.MustCompile(`\d+-aws-billing-detailed-line-items-with-resources-and-tags-\d{4}-\d{2}.csv.zip$`)

func isBillKey(b BillKey) bool {
	return billKey.MatchString(b.Key)
}

func listBillsFromRepository(ctx context.Context, ss *s3.S3, r BillRepository) <-chan BillKey {
	c := make(chan BillKey)
	go func() {
		defer close(c)
		input := s3.ListObjectsV2Input{
			Bucket: &r.Bucket,
			Prefix: &r.Prefix,
		}
		err := ss.ListObjectsV2PagesWithContext(ctx, &input, func(page *s3.ListObjectsV2Output, last bool) bool {
			for _, o := range page.Contents {
				c <- BillKey{
					Key:          *o.Key,
					Bucket:       r.Bucket,
					Region:       "",
					LastModified: *o.LastModified,
				}
			}
			return !last
		})
		if err != nil {
			jsonlog.LoggerFromContextOrDefault(ctx).Error("Failed to list objects from bucket.", err.Error())
		}
	}()
	return filterBillsByKey(c)
}
