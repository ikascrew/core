package sync_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ikascrew/core/sync"
)

func TestTenGroup(t *testing.T) {

	start := time.Now()
	wait := sync.NewGroup(10)
	for i := 0; i < 100; i++ {
		wait.Add()
		go func(idx int) {
			defer wait.Done()
			fmt.Println(idx, time.Now())
			time.Sleep(time.Millisecond * 100)
		}(i)
	}

	wait.Wait()
	end := time.Now()

	if errs := wait.Errors(); errs != nil {
		t.Error("wait want not error")
	}

	diff := end.Sub(start)

	if diff.Milliseconds() <= int64(1000) {
		t.Error("error 1 second")
	}

	if diff.Milliseconds() > int64(1500) {
		t.Error("error 1.5 second")
	}

}

func TestZeroGroup(t *testing.T) {
	start := time.Now()
	wait := sync.NewGroup(0)
	for i := 0; i < 100; i++ {
		wait.Add()
		go func(idx int) {
			defer wait.Done()
			fmt.Println(idx, time.Now())
			time.Sleep(time.Millisecond * 100)
		}(i)
	}

	wait.Wait()
	end := time.Now()

	if errs := wait.Errors(); errs != nil {
		t.Error("wait want not error")
	}

	diff := end.Sub(start)
	if diff.Milliseconds() > int64(500) {
		t.Error("error 1 second")
	}

}
