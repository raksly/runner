[![Build Status](https://travis-ci.org/raksly/runner.svg?branch=master)](https://travis-ci.org/raksly/runner)

Runner helps you to not reinvent the wheel on certain type of applications.
Often you begin your application by starting multiple goroutines to do
separate work. Imagine a server accepting http, tcp, and other protocols.
You probably also create WaitGroups and/or channels to sync these goroutines.
Runner does all that for you.
## API
The API is promised to never break existing projects using this library. 
Documentation available on [GoDoc](https://godoc.org/github.com/raksly/runner).
## Examples
### Minimalistic
```golang
var r runner.Runner

r.Run(runHTTP)
r.Run(runTCP)

r.Wait()
```
`runHTTP` and `runTCP` are both of type `func()`. `r.Run` will run both functions in separate goroutines, and `r.Wait()` waits until both functions exit.
### Exit notification
When `runTCP` exits, it might be because the application is supposed to exit alltogether, or there was an irrecoverable error. In that case, you might want HTTP to exit aswell. `Run*` methods return a channel which is closed when its running function returns.
```golang
ctx, cancel := context.WithCancel(context.Background())
r := runner.Runner{Ctx: ctx}

select {
case <-r.RunContext(runHTTP):
case <-r.RunContext(runTCP):
}

cancel()
r.Wait()
```
Both `runHTTP` and `runTCP` are now of type `func(context.Context)` and 
are given the context `Runner.Ctx`. If either `runHTTP` or `runTCP` returns, `select` will break, the context will be cancelled, making the other function exit aswell in due time, and `r.Wait()` waits for that.
### Signals
`Runner` contains a convenience method to work with OS signals
```golang
ctx, cancel := context.WithCancel(context.Background())
r := runner.Runner{Ctx: ctx}

select {
case <-r.RunContext(runHTTP):
case <-r.RunContext(runTCP):
case <-r.RunSigs(syscall.SIGINT, syscall.SIGTERM):
}

cancel()
r.Wait()
```
If either `runHTTP` or `runTCP` returns, or `SIGINT`/`SIGTERM` is received,
the context will be cancelled and we wait for everything to clean up.
### Multiple parameters
You might want to pass more than just a context to your function. To archive
this, you may use closures
```golang
ctx, cancel := context.WithCancel(context.Background())
r := runner.Runner{Ctx: ctx}

select {
case <-r.Run(func() { runSomething(ctx, 1, 2, 3) }):
case <-r.RunSigs(syscall.SIGINT, syscall.SIGTERM):
}

cancel()
r.Wait()
```
## License
MIT