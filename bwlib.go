package main

import (
	"github.com/yuin/gopher-lua"
	bw2 "gopkg.in/immesys/bw2bind.v5"
	"time"
)

func LoadLib(L *lua.LState) {
	L.SetGlobal("getone", L.NewFunction(GetOne))
	L.SetGlobal("subscribe", L.NewFunction(Subscribe))
	L.SetGlobal("publish", L.NewFunction(Publish))
	L.SetGlobal("keeprunning", L.NewFunction(KeepRunning))
	L.SetGlobal("sleep", L.NewFunction(Sleep))
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
		L.Push(lua.LString(uri))
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
			L.Push(lua.LString(uri))
			pushMsg(msg, ponum, L)
			if err := L.PCall(2, 0, nil); err != nil {
				L.RaiseError("Error doing func callback (%v)", err)
			}
		}
	}()
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
