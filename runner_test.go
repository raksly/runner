package runner_test

import (
	"context"
	"reflect"
	"syscall"
	"testing"
	"time"

	"github.com/raksly/runner"

	"github.com/stretchr/testify/assert"
)

func run() {
	time.Sleep(time.Second)
}

func runContext(ctx context.Context) {
	<-ctx.Done()
}

func assertChan(assert *assert.Assertions, c interface{}, expectClosed bool) {
	if expectClosed {
		_, ok := reflect.ValueOf(c).Recv()
		assert.False(ok)
	} else {
		chosen, _, _ := reflect.Select([]reflect.SelectCase{
			reflect.SelectCase{Dir: reflect.SelectDefault},
			reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(c)},
		})
		assert.Equal(0, chosen)
	}
}

func Test2(t *testing.T) {
	assert := assert.New(t)

	ctx, cancel := context.WithCancel(context.Background())
	r := runner.Runner{Ctx: ctx}

	c1 := r.Run(run)
	c2 := r.RunContext(runContext)
	c3 := r.RunOtherContext(ctx, runContext)
	c4 := r.RunSigs(syscall.SIGINT)

	select {
	case <-c1:
	case <-c2:
		assert.Fail("should not run")
	case <-c3:
		assert.Fail("should not run")
	case <-c4:
		assert.Fail("should not run")
	}

	assertChan(assert, c1, true)
	assertChan(assert, c2, false)
	assertChan(assert, c3, false)
	assertChan(assert, c4, false)

	cancel()
	r.Wait()

	assertChan(assert, c2, true)
	assertChan(assert, c3, true)
	assertChan(assert, c4, false)

	var r2 runner.Runner

	r2.Run(func() {})
	r2.RunContext(func(ctx context.Context) {
		assert.Nil(ctx)
	})

	r2.Wait()
}

// Keeping Test1 to ensure the old api still works
func Test1(t *testing.T) {
	assert := assert.New(t)

	ctx, cancel := context.WithCancel(context.Background())
	runner := runner.New(ctx)

	assert.Equal(ctx, runner.Context())

	c1 := runner.Run(run)
	c2 := runner.RunContext(runContext)
	c3 := runner.RunOtherContext(ctx, runContext)
	c4 := runner.RunSigs(syscall.SIGINT)

	select {
	case <-c1:
	case <-c2:
		assert.Fail("should not run")
	case <-c3:
		assert.Fail("should not run")
	case <-c4:
		assert.Fail("should not run")
	}

	assertChan(assert, c1, true)
	assertChan(assert, c2, false)
	assertChan(assert, c3, false)
	assertChan(assert, c4, false)

	cancel()
	runner.Wait()

	assertChan(assert, c2, true)
	assertChan(assert, c3, true)
	assertChan(assert, c4, false)
}
