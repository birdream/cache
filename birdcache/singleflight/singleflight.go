package singleflight

import (
	"sync"
)

type Call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

type Group struct {
	mu sync.Mutex
	m  map[string]*Call
}

// Do for the same key, no matter Do been call how much time until fn get the response,
// fn only call one time
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*Call)
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait() // if fn is been called then wait for the result
		// using wg here is because the concurrency thread do not need to communicate with each other
		return c.val, c.err
	}
	c := new(Call)
	c.wg.Add(1) // before call fn, lock down
	g.m[key] = c
	g.mu.Unlock()

	c.val, c.err = fn() // call the fn
	c.wg.Done()         // fn called done and get the response

	g.mu.Lock()
	delete(g.m, key) // update g.m, delete the fn record
	// so the next time, when same fn called in and would get call the fn again to get the result
	// that mean this Do func is only keep fn called one time during fn been calling
	// after the fn get the response and next time someone call same fn again
	// fn will be excuted agagin
	g.mu.Unlock()

	return c.val, c.err
}
