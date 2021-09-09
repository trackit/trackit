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

// periodic runs periodic tasks. Tasks can be registered to a Scheduler which
// may or may not be "ticking", i.e. running its tasks.
package periodic

import (
	"context"
	"sync"
	"time"

	"github.com/trackit/jsonlog"
)

// contextKey is an unexported type to avoid collisions in context.Context
// values.
type contextKey uint

// taskSignal is the type for signals sent to task tickers.
type taskSignal uint

const (
	Hourly      = 1 * time.Hour
	Daily       = 24 * time.Hour
	TwiceDaily  = Daily / 2
	ThriceDaily = Daily / 3

	// TaskTime is a key for a context.Context where the starting time of
	// the current request is stored.
	TaskTime = contextKey(iota)

	// taskStop instructs a task ticker to stop running its task and return.
	taskStop = taskSignal(iota)
)

// Task is a task that can be scheduled.
type Task func(context.Context) error

// taskRegistration is a task registration that may or may not be ticking.
type taskRegistration struct {
	Name    string `json:"name"`
	task    Task
	Period  time.Duration
	ticker  *time.Ticker
	control chan taskSignal
}

// Scheduler runs registered periodic tasks. Its zero value is a valid
// Scheduler that doesn't tick and has no registered task. It may be used in
// parallel.
type Scheduler struct {
	running       bool
	registrations []taskRegistration
	mutex         sync.RWMutex
}

// Ticking returns whether the Scheduler is currently ticking.
func (s *Scheduler) Ticking() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.running
}

// Register registers a Task to the Scheduler to be run at period p. If the
// Scheduler is ticking, the task starts ticking immediately.
func (s *Scheduler) Register(t Task, p time.Duration, n string) {
	r := taskRegistration{
		task:   t,
		Period: p,
		Name:   n,
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.running {
		r.start()
	}
	s.registrations = append(s.registrations, r)
}

// Start starts a Scheduler. Starting an already started scheduler is
// functionally a noop, though an error will be logged.
func (s *Scheduler) Start() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.running {
		jsonlog.Error("Attempt to start already started scheduler. Ignoring.", nil)
	} else {
		for i := range s.registrations {
			s.registrations[i].start()
		}
		s.running = true
	}
}

// Stop stops a Scheduler. Stopping an already stopped scheduler is
// functionally a noop, though an error will be logged.
func (s *Scheduler) Stop() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.running {
		for i := range s.registrations {
			s.registrations[i].stop()
		}
		s.running = false
	}
}

// start starts a taskRegistration, having it tick and run its task
// periodically.
func (t *taskRegistration) start() {
	if t.control == nil {
		t.ticker = time.NewTicker(t.Period)
		t.control = make(chan taskSignal)
		go t.tick()
	} else {
		jsonlog.Error("Attempt to start already started task. Ignoring.", t)
	}
}

// run runs the taskRegistration's task in the current goroutine.
func (t *taskRegistration) run(d time.Time) error {
	ctx := context.Background()
	ctx = context.WithValue(ctx, TaskTime, d)
	return t.task(ctx)
}

// getTickerChan gets the taskRegistration's ticker channel, creating it if it
// doesn’t exist yet. This is the channel where a time.Ticker outputs to run
// the periodic tasks.
func (t *taskRegistration) getTickerChan() <-chan time.Time {
	c := t.ticker
	if c == nil {
		return nil
	} else {
		return c.C
	}
}

// tick starts periodic tasks in response to ticks from t.ticker. The tasks are
// started in their own goroutine using t.run.
func (t *taskRegistration) tick() {
	for {
		select {
		case d := <-t.getTickerChan():
			go func() {
				if err := t.run(d); err != nil {
					jsonlog.DefaultLogger.Error("Error while running periodic task", map[string]interface{}{
						"error": err.Error(),
					})
				}
			}()
		case s := <-t.control:
			switch s {
			case taskStop:
				close(t.control)
				t.control = nil
				t.ticker.Stop()
				t.ticker = nil
				return
			}
		}
	}
}

// stop stops a taskRegistration from ticking. Currently running tasks are not
// cancelled.
func (t *taskRegistration) stop() {
	if t.control != nil {
		t.control <- taskStop
	} else {
		jsonlog.Error("Attempt to stop an already stopped task. Ignoring.", t)
	}
}
