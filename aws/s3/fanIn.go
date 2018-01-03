package s3

import (
	"sync"
)

// mergecLineItem implements the fan-in pattern by merging to the out
// channel the input from the channels read on cs.
func mergecLineItem(out chan<- LineItem, cs <-chan <-chan LineItem) {
	var wg sync.WaitGroup
	for c := range cs {
		wg.Add(1)
		go func(c <-chan LineItem) {
			defer wg.Done()
			for u := range c {
				out <- u
			}
		}(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
}

// mergecdLineItem wraps mergecLineItem by making the channels and starting the
// goroutine responsible for the fan-in.
func mergecdLineItem() (chan<- <-chan LineItem, <-chan LineItem) {
	in := make(chan (<-chan LineItem))
	out := make(chan LineItem)
	go mergecLineItem(out, in)
	return in, out
}

// mergecManifest implements the fan-in pattern by merging to the out
// channel the input from the channels read on cs.
func mergecManifest(out chan<- manifest, cs <-chan <-chan manifest) {
	var wg sync.WaitGroup
	for c := range cs {
		wg.Add(1)
		go func(c <-chan manifest) {
			defer wg.Done()
			for u := range c {
				out <- u
			}
		}(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
}

// mergecdManifest wraps mergecManifest by making the channels and starting the
// goroutine responsible for the fan-in.
func mergecdManifest() (chan<- <-chan manifest, <-chan manifest) {
	in := make(chan (<-chan manifest))
	out := make(chan manifest)
	go mergecManifest(out, in)
	return in, out
}
