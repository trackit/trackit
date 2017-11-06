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
	"sync"
	"time"

	"github.com/trackit/jsonlog"
)

type contextKey uint
type taskSignal uint

const (
	Hourly      = 1 * time.Hour
	Daily       = 24 * time.Hour
	TwiceDaily  = Daily / 2
	ThriceDaily = Daily / 3

	TaskTime = contextKey(iota)

	taskStop = taskSignal(iota)
)

type Task func(context.Context) error

type taskRegistration struct {
	Name    string `json:"name"`
	task    Task
	Period  time.Duration
	ticker  *time.Ticker
	control chan taskSignal
}

type Scheduler struct {
	running       bool
	registrations []taskRegistration
	mutex         sync.RWMutex
}

func (s *Scheduler) Running() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.running
}

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

func (t *taskRegistration) start() {
	if t.control == nil {
		t.ticker = time.NewTicker(t.Period)
		t.control = make(chan taskSignal)
		go t.tick()
	} else {
		jsonlog.Error("Attempt to start already started task. Ignoring.", t)
	}
}

func (t *taskRegistration) run(d time.Time) error {
	ctx := context.Background()
	ctx = context.WithValue(ctx, TaskTime, d)
	return t.task(ctx)
}

func (t *taskRegistration) getTickerChan() <-chan time.Time {
	c := t.ticker
	if c == nil {
		return nil
	} else {
		return c.C
	}
}

func (t *taskRegistration) tick() {
tickloop:
	for {
		select {
		case d := <-t.getTickerChan():
			go t.run(d)
		case s := <-t.control:
			switch s {
			case taskStop:
				close(t.control)
				t.control = nil
				t.ticker.Stop()
				t.ticker = nil
				break tickloop
			}
		}
	}
}

func (t *taskRegistration) stop() {
	if t.control != nil {
		t.control <- taskStop
	} else {
		jsonlog.Error("Attempt to stop an already stopped task. Ignoring.", t)
	}
}
