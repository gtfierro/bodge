package main

import (
	"fmt"
	"github.com/yuin/gopher-lua"
	"sync"
)

// we need this mutex every time we're going to add a new coroutine.
// The reason is that we need to push values on the stack and need to know
// that no one else is going to be interacting w/ the stack. Normally
// this would be single threaded and not an issue, but we have asynchronous
// messages coming in so we don't have that guarantee here
var schedMutex sync.Mutex

const MAX_COROUTINE = 256

type coroutine struct {
	L    *lua.LState
	fxn  *lua.LFunction
	name string
	args []lua.LValue
}

var coroutines = make(chan coroutine, MAX_COROUTINE)

func startScheduler(L *lua.LState) {
	// start the scheduler
	go func() {
		for frame := range coroutines {
			if frame.fxn == nil {
				continue // bad record!
			}
			// lock the stack so we can run the function
			schedMutex.Lock()
			st, err, values := L.Resume(frame.L, frame.fxn, frame.args...)
			if st == lua.ResumeError {
				schedMutex.Unlock()
				L.RaiseError("Error doing func callback (%v)", err)
			}
			// push return values onto the stack
			for _, val := range values {
				L.Push(val)
			}
			// this is the FINAL status. if here, then we kill the coroutine and reset the stack
			// otherwise, we put the coroutine back on the stack
			if st == lua.ResumeOK {
				L.SetTop(0)
				frame.L.Close()
			} else {
				fmt.Println("put it back", frame.L.GetTop())
				coroutines <- frame
			}
			// free up the scheduler to accept more coroutines
			schedMutex.Unlock()
		}
	}()
}
