package concutil

type Awaiter struct {
	resp interface{}
	err  error
	ch   chan struct{}
}

// Async executes the code in f() asynchronously as returns an Awaiter as promoise.
func Async(f func() (interface{}, error)) *Awaiter {
	awaiter := &Awaiter{
		ch: make(chan struct{}),
	}
	go func() {
		defer close(awaiter.ch)
		awaiter.resp, awaiter.err = f()
	}()
	return awaiter
}

// Await waits for the task to finish and returns the error if fails.
//
// This function can be called multiple times to get the same value.
func (a *Awaiter) Await() error {
	<-a.ch
	return a.err
}

// Get waits for the task to finish and returns the obtained value.
//
// This function can be called multiple times to get the same value.
func (a *Awaiter) Get() interface{} {
	<-a.ch
	return a.resp
}
