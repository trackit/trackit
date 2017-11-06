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

func TestRunningSchedulerIsRunning(t *testing.T) {
	var s Scheduler
	s.Start()
	if s.Running() != true {
		t.Errorf("Running scheduler should be marked running, isn't.")
	}
	s.Stop()
}

func TestRegisterAtRunningScheduler(t *testing.T) {
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

func TestNotRunningSchedulerIsNotRunning(t *testing.T) {
	var s Scheduler
	if s.Running() != false {
		t.Errorf("Never started scheduler should not be marked running. Is.")
	}
	s.Start()
	s.Stop()
	if s.Running() != false {
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
