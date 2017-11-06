package s3

import (
	"context"
	"regexp"
	"sync"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/trackit/jsonlog"
)

type BillLocation struct {
	Region string
	Bucket string
	Key    string
}

type bucketLocation struct {
}

func LocateBills(ctx context.Context, ss *s3.S3, cbk func(context.Context, BillLocation) bool) error {
	buckets, err := listBuckets(ctx, ss)
	if err == nil {
		for b := range listBills(ctx, ss, buckets) {
			if !cbk(ctx, b) {
				break
			}
		}
	}
	return err
}

func listBuckets(ctx context.Context, ss *s3.S3) ([]string, error) {
	input := s3.ListBucketsInput{}
	output, err := ss.ListBucketsWithContext(ctx, &input)
	if err == nil {
		buckets := make([]string, len(output.Buckets))
		for i := range buckets {
			buckets[i] = *output.Buckets[i].Name
		}
		return buckets, nil
	} else {
		return nil, err
	}
}

// mergecBillLocation implements the fan-in pattern by merging to the out
// channel the input from the channels read on cs.
func mergecBillLocation(out chan<- BillLocation, cs <-chan <-chan BillLocation) {
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

func listBills(ctx context.Context, ss *s3.S3, buckets []string) <-chan BillLocation {
	c := make(chan BillLocation)
	cc := make(chan (<-chan BillLocation))
	defer close(cc)
	go mergecBillLocation(c, cc)
	for _, b := range buckets {
		cc <- listBillsFromBucket(ctx, ss, b)
	}
	return c
}

func filterBillsByKey(in <-chan BillLocation) <-chan BillLocation {
	out := make(chan BillLocation)
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

func isBillKey(b BillLocation) bool {
	return billKey.MatchString(b.Key)
}

func listBillsFromBucket(ctx context.Context, ss *s3.S3, b string) <-chan BillLocation {
	c := make(chan BillLocation)
	go func() {
		defer close(c)
		input := s3.ListObjectsV2Input{
			Bucket: &b,
		}
		err := ss.ListObjectsV2PagesWithContext(ctx, &input, func(page *s3.ListObjectsV2Output, last bool) bool {
			for _, o := range page.Contents {
				c <- BillLocation{
					Key:    *o.Key,
					Bucket: b,
					Region: "",
				}
			}
			return !last
		})
		if err != nil {
			jsonlog.LoggerFromContextOrDefault(ctx).Error("Failed to list objects from bucket.", err.Error())
		}
	}()
	return c
}
