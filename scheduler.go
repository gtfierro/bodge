package main

import (
	"fmt"
	"github.com/yuin/gopher-lua"
	"sync"
)

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
			//fmt.Printf("%+v\n", frame)
			if frame.fxn == nil {
				continue // bad record!
			}
			schedMutex.Lock()
			st, err, values := L.Resume(frame.L, frame.fxn, frame.args...)
			if st == lua.ResumeError {
				schedMutex.Unlock()
				L.RaiseError("Error doing func callback (%v)", err)
			}
			//fmt.Println("len values", len(values))
			for _, val := range values {
				L.Push(val)
			}
			if st == lua.ResumeOK {
				L.SetTop(0)
				frame.L.Close()
				//frame.args = []lua.LValue{}
				// function has ended
				//fmt.Println("yield break?")
				//break
			} else {
				fmt.Println("put it back", frame.L.GetTop())
				coroutines <- frame
			}
			schedMutex.Unlock()
		}
	}()
}
