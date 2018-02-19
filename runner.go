/*

Package runner helps you to not reinvent the wheel on certain type of applications.
Often you begin your application by starting multiple goroutines to do
separate work. Imagine a server accepting http, tcp, and other protocols.
You probably also create WaitGroups and/or channels to sync these goroutines.
Runner does all that for you.

For examples please view https://github.com/raksly/runner#examples

*/
package runner

import (
	"context"
	"os"
	"os/signal"
	"sync"
)

// New creates a new Runner
func New(ctx context.Context) Runner {
	return Runner{
		ctx: ctx,
		wg:  &sync.WaitGroup{},
	}
}

// Runner runs functions in goroutines with `Run`, `RunContext` and/or
// `RunOtherContext`, returning a channel that is closed when the goroutine
// exits.
// You may wait on all goroutines to finish using `Wait`.
type Runner struct {
	ctx context.Context
	wg  *sync.WaitGroup
}

// Run runs a function of type func() in a new goroutine.
// The returned channel is closed when f returns
func (r Runner) Run(f func()) <-chan struct{} {
	done := make(chan struct{})
	r.wg.Add(1)
	go func() {
		defer func() {
			r.wg.Done()
			close(done)
		}()
		f()
	}()
	return done
}

// RunContext is like `Run`, but passes the context given to `New`.
func (r Runner) RunContext(f func(context.Context)) <-chan struct{} {
	return r.RunOtherContext(r.ctx, f)
}

// RunOtherContext is like `RunContext`, except you may specify which
// context to be passed to f.
func (r Runner) RunOtherContext(ctx context.Context, f func(context.Context)) <-chan struct{} {
	return r.Run(func() { f(ctx) })
}

// RunSigs is a convenience method to work with OS signals.
// Unlike the other `Run*` functions, the channel returned
// is not closed, but reads the received signal.
func (r Runner) RunSigs(sigs ...os.Signal) <-chan os.Signal {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, sigs...)
	return sig
}

// Wait waits for all goroutines started by this runner.
func (r Runner) Wait() {
	r.wg.Wait()
}

// Context returns the context given to `New`.
func (r Runner) Context() context.Context {
	return r.ctx
}
