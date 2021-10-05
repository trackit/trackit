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
	"sync"

	"github.com/trackit/trackit/es/indexes/lineItems"
)

// mergecLineItem implements the fan-in pattern by merging to the out
// channel the input from the channels read on cs.
func mergecLineItem(out chan<- lineItems.LineItem, cs <-chan <-chan lineItems.LineItem) {
	var wg sync.WaitGroup
	for c := range cs {
		wg.Add(1)
		go func(c <-chan lineItems.LineItem) {
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
func mergecdLineItem() (chan<- <-chan lineItems.LineItem, <-chan lineItems.LineItem) {
	in := make(chan (<-chan lineItems.LineItem))
	out := make(chan lineItems.LineItem)
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
