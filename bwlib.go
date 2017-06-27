package main

import (
	"fmt"
	bw2 "github.com/immesys/bw2bind"
	"github.com/yuin/gopher-lua"
	"time"
)

var exports = map[string]lua.LGFunction{
	// bosswave functions
	"getone":    GetOne,
	"subscribe": Subscribe,
	"publish":   Publish,
	"persist":   Persist,
	"query":     Query,
	// timers
	"sleep":              Sleep,
	"invokePeriodically": InvokePeriodically,
	"invokeLater":        InvokeLater,
	"loop":               KeepRunning,
	"every":              ScheduleEvery,
	// utils
	"dumptable":  DumpTable,
	"nargs":      NArgs,
	"arg":        Arg,
	"uriRequire": URIRequire,
}

func LoadLib(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), exports)
	L.Push(mod)
	return 1
}

// Get one message and payload object off of a subscription
// Arguments:
// - URI to subscribe to
// - Payload object
// - (optional) callback with arguments: value, uri, published vk?
func GetOne(L *lua.LState) int {
	uri := L.ToString(1)
	ponum := L.ToString(2)
	f := L.ToFunction(3)

	// subscribe to the URI
	sub, handle, err := client.SubscribeH(&bw2.SubscribeParams{
		URI: uri,
	})
	if err != nil {
		L.RaiseError("Error subscribing to %s (%v)", uri, err)
	}
	// get the first message
	msg := <-sub
	// unsubscribe
	if err := client.Unsubscribe(handle); err != nil {
		L.RaiseError("Error unsubscribing to %s (%v)", uri, err)
	}
	// extract the contents from the indicated payload object
	// and push it onto the stack

	if f != nil {
		L.Push(f)
		L.Push(lua.LString(msg.URI))
		pushMsg(msg, ponum, L)
		if err := L.PCall(2, 0, nil); err != nil {
			L.RaiseError("Error doing func callback (%v)", err)
		}
		return 0
	} else {
		pushMsg(msg, ponum, L)
		return 1
	}
	return 0
}

func Subscribe(L *lua.LState) int {
	uri := L.ToString(1)
	ponum := L.ToString(2)
	f := L.ToFunction(3)

	// subscribe to the URI
	sub, err := client.Subscribe(&bw2.SubscribeParams{
		URI: uri,
	})
	if err != nil {
		L.RaiseError("Error subscribing to %s (%v)", uri, err)
	}

	// get coroutine

	// new thread?
	go func() {
		for msg := range sub {
			//msg.Dump()
			schedMutex.Lock()
			L.Push(f)
			L.Push(lua.LString(msg.URI))
			pushMsg(msg, ponum, L)
			DoCoroutine(L)
			schedMutex.Unlock()
		}
	}()
	return 0
}

func Query(L *lua.LState) int {
	uri := L.ToString(1)
	ponum := L.ToString(2)
	f := L.ToFunction(3)

	// subscribe to the URI
	sub, err := client.Query(&bw2.QueryParams{
		URI: uri,
	})
	if err != nil {
		L.RaiseError("Error querying %s (%v)", uri, err)
	}
	for msg := range sub {
		L.Push(f)
		L.Push(lua.LString(msg.URI))
		pushMsg(msg, ponum, L)
		if err := L.PCall(2, 0, nil); err != nil {
			L.RaiseError("Error doing func callback (%v)", err)
		}
	}
	return 0
}

func publish(L *lua.LState, persist bool) int {
	uri := L.ToString(1)
	ponum := L.ToString(2)
	payload := L.CheckAny(3)
	_, _, payload = uri, ponum, payload
	if uri == "" {
		L.RaiseError("No URI provided")
	}
	if ponum == "" {
		L.RaiseError("No PO num provided")
	}
	//f := L.ToFunction(3)
	po := lvalueToPO(payload, ponum)
	err := client.Publish(&bw2.PublishParams{
		URI:            uri,
		PayloadObjects: []bw2.PayloadObject{po},
		Persist:        persist,
	})
	if err != nil {
		L.RaiseError("Error while publishing (%v)", err)
	}

	return 0
}

func Publish(L *lua.LState) int {
	return publish(L, false)
}
func Persist(L *lua.LState) int {
	return publish(L, true)
}

func KeepRunning(L *lua.LState) int {
	for {
		time.Sleep(300 * time.Second)
	}
}

// milliseconds
func Sleep(L *lua.LState) int {
	n := L.ToNumber(1)
	time.Sleep(time.Duration(n) * time.Millisecond)
	return 0
}

func InvokePeriodically(L *lua.LState) int {
	nargs := L.GetTop()
	n := L.ToNumber(1)
	f := L.ToFunction(2)
	var args []lua.LValue
	for i := 3; i <= nargs; i++ {
		args = append(args, L.CheckAny(i))
	}

	go func() {
		for _ = range time.Tick(time.Duration(n) * time.Millisecond) {
			schedMutex.Lock()
			L.Push(f)
			for _, arg := range args {
				L.Push(arg)
			}
			DoCoroutine(L)
			schedMutex.Unlock()
		}
	}()
	return 0
}

func InvokeLater(L *lua.LState) int {
	nargs := L.GetTop()
	n := L.ToNumber(1)
	f := L.ToFunction(2)
	var args []lua.LValue
	for i := 3; i <= nargs; i++ {
		args = append(args, L.CheckAny(i))
	}

	L = lua.NewState()
	defer L.Close()
	time.AfterFunc(time.Duration(n)*time.Millisecond, func() {
		L.Push(f)
		for _, arg := range args {
			L.Push(arg)
		}
		if err := L.PCall(nargs-2, 0, nil); err != nil {
			L.RaiseError("Error doing func callback (%v)", err)
		}
	})
	return 0
}

func ScheduleEvery(L *lua.LState) int {
	nargs := L.GetTop()
	if nargs < 2 {
		L.RaiseError("Need >=2 arguments: every(sched string, function)")
	}
	schedString := L.ToString(1)
	f := L.ToFunction(2)
	var args []lua.LValue
	for i := 3; i <= nargs; i++ {
		args = append(args, L.CheckAny(i))
	}
	task := func() {
		schedMutex.Lock()
		L.Push(f)
		for _, arg := range args {
			L.Push(arg)
		}
		DoCoroutine(L)
		schedMutex.Unlock()
	}
	ScheduleTask(schedString, task)
	return 0
}

func DumpTable(L *lua.LState) int {
	table := L.ToTable(1)
	table.ForEach(func(k, v lua.LValue) {
		fmt.Printf("%s => %s\n", k, v)
	})
	return 0
}

func NArgs(L *lua.LState) int {
	L.Push(lua.LNumber(_NARGS))
	return 1
}

func Arg(L *lua.LState) int {
	nargs := L.GetTop()
	if nargs < 1 {
		L.Push(lua.LNil)
	}
	n := L.ToNumber(1)
	_n := int(n)
	// 1-indexed
	if _n > 0 && _n-1 <= len(_ARGS) {
		L.Push(lua.LString(_ARGS[_n-1]))
	} else {
		L.Push(lua.LNil)
	}

	return 1
}

func URIRequire(L *lua.LState) int {
	nargs := L.GetTop()
	if nargs < 1 {
		L.RaiseError("URIRequire needs argument")
	}
	uri := L.ToString(1)
	file, err := readURI(uri)
	if err != nil {
		L.RaiseError("Error reading URI %s (%v)", uri, err)
	}
	L.DoString(file)
	return 1
}

func DoCoroutine(L *lua.LState) int {
	nargs := L.GetTop()
	//fmt.Println("NARGS", nargs)
	co, _ := L.NewThread()
	f := L.ToFunction(1)
	var args []lua.LValue
	for i := 2; i <= nargs; i++ {
		args = append(args, L.CheckAny(i))
	}

	frame := coroutine{
		L:    co,
		fxn:  f,
		name: f.String(),
		args: args,
	}

	coroutines <- frame
	return 0
}
