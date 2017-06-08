package main

import (
	"fmt"
	"github.com/yuin/gopher-lua"
	bw2 "gopkg.in/immesys/bw2bind.v5"
	"time"
)

func LoadLib(L *lua.LState) {
	L.SetGlobal("getone", L.NewFunction(GetOne))
	L.SetGlobal("subscribe", L.NewFunction(Subscribe))
	L.SetGlobal("publish", L.NewFunction(Publish))
	L.SetGlobal("query", L.NewFunction(Query))
	L.SetGlobal("keepRunning", L.NewFunction(KeepRunning))
	L.SetGlobal("loop", L.NewFunction(KeepRunning))
	L.SetGlobal("sleep", L.NewFunction(Sleep))
	L.SetGlobal("invokePeriodically", L.NewFunction(InvokePeriodically))
	L.SetGlobal("invokeLater", L.NewFunction(InvokeLater))
	L.SetGlobal("dumptable", L.NewFunction(DumpTable))
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
	// new thread?
	go func() {
		L := lua.NewState()
		defer L.Close()
		for msg := range sub {
			L.Push(f)
			L.Push(lua.LString(msg.URI))
			pushMsg(msg, ponum, L)
			if err := L.PCall(2, 0, nil); err != nil {
				L.RaiseError("Error doing func callback (%v)", err)
			}
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

func Publish(L *lua.LState) int {
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
	})
	if err != nil {
		L.RaiseError("Error while publishing (%v)", err)
	}

	return 0
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
		L := lua.NewState()
		defer L.Close()
		for _ = range time.Tick(time.Duration(n) * time.Millisecond) {
			L.Push(f)
			for _, arg := range args {
				L.Push(arg)
			}
			if err := L.PCall(nargs-2, 0, nil); err != nil {
				L.RaiseError("Error doing func callback (%v)", err)
			}
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

func DumpTable(L *lua.LState) int {
	table := L.ToTable(1)
	table.ForEach(func(k, v lua.LValue) {
		fmt.Printf("%s => %s\n", k, v)
	})
	return 0
}
