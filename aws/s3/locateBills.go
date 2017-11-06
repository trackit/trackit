package s3

import (
	"context.Context"

	"github.com/aws/aws-sdk-go/service/s3"
)

type BillLocation struct {
	Region string
	Bucket string
	Key    string
}

type bucketLocation struct {
}

func LocateBills(ctx context.Context, ss *s3.S3, cbk func(context.Context, BillLocation) bool) error {
	buckets, error := listBuckets(ctx, ss)
	for b := range listBills(ctx, ss, buckets) {
		if !cbk(ctx, b) {
			break
		}
	}
}

func listBuckets(ctx context.Context, ss *s3.S3) ([]string, error) {
	input := s3.ListBucketsInput{}
	output, err := ss.ListBucketsWithContext(ctx, &input)
	if err == nil {
		buckets := make([]string, len(output.Buckets))
		for i := range buckets {
			buckets[i] = *output.Buckets[i]
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
	go mergecBillLocations(c, cc)
	for _, b := range buckets {
		cc <- listBillsFromBucket(ctx, ss, b)
	}
	return c
}

func listBillsFromBucket(ctx context.Context, ss *s3.S3, b string) {

}
