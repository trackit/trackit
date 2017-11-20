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

package periodic

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/trackit/jsonlog"
)

func init() {
	jsonlog.DefaultLogger = jsonlog.DefaultLogger.WithLogLevel(jsonlog.LogLevelDebug)
}

func messageTask(c chan<- int, i int) Task {
	return func(_ context.Context) error {
		c <- i
		return nil
	}
}

func TestTickingSchedulerIsTicking(t *testing.T) {
	var s Scheduler
	s.Start()
	if s.Ticking() != true {
		t.Errorf("Ticking scheduler should be marked running, isn't.")
	}
	s.Stop()
}

func TestRegisterAtTickingScheduler(t *testing.T) {
	var s Scheduler
	var a int
	c := make(chan int)
	s.Start()
	s.Register(
		messageTask(c, 0),
		100*time.Millisecond,
		"Tenth of a second",
	)
	e := time.After(350 * time.Millisecond)
out:
	for {
		select {
		case <-c:
			a++
		case <-e:
			break out
		}
	}
	if a != 3 {
		t.Errorf("Task should run %d times, ran %d times.", 3, a)
	}
	s.Stop()
}

func TestNotTickingSchedulerIsNotTicking(t *testing.T) {
	var s Scheduler
	if s.Ticking() != false {
		t.Errorf("Never started scheduler should not be marked running. Is.")
	}
	s.Start()
	s.Stop()
	if s.Ticking() != false {
		t.Errorf("Started-then-stopped scheduler should not be marked running. Is.")
	}
}

func TestTaskCount(t *testing.T) {
	var s Scheduler
	var b [10]int
	var be [10]int = [10]int{4, 4, 3, 3, 3, 3, 2, 2, 2, 2}
	c := make(chan int)
	for i := range b {
		d := time.Duration(217 + i*23)
		s.Register(
			messageTask(c, i),
			d*time.Millisecond,
			fmt.Sprintf("Task %2d (%3dms)", b, d),
		)
	}
	s.Start()
	e := time.After(1 * time.Second)
out:
	for {
		select {
		case i := <-c:
			b[i]++
		case <-e:
			break out
		}
	}
	s.Stop()
	if b != be {
		t.Errorf("Expected task count to be %#v, is %#v.", be, b)
	}
}
