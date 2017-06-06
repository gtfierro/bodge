package main

import (
	"fmt"
	"log"

	"github.com/chzyer/readline"
	"github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"
	bw2 "gopkg.in/immesys/bw2bind.v5"
)

var client *bw2.BW2Client

// do read/eval/print/loop
func doREPL(L *lua.LState) {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:                 "> ",
		HistoryFile:            ".bua-history",
		DisableAutoSaveHistory: false,
	})
	// TODO: cleanup when we exit
	if err != nil {
		log.Fatal(err)
	}
	for {
		if str, err := loadline(rl, L); err == nil {
			if err := L.DoString(str); err != nil {
				fmt.Println(err)
			}
		} else { // error on loadline
			fmt.Println(err)
			return
		}
	}
}

func incomplete(err error) bool {
	if lerr, ok := err.(*lua.ApiError); ok {
		if perr, ok := lerr.Cause.(*parse.Error); ok {
			return perr.Pos.Line == parse.EOF
		}
	}
	return false
}

func loadline(rl *readline.Instance, L *lua.LState) (string, error) {
	rl.SetPrompt("> ")
	if line, err := rl.Readline(); err == nil {
		if _, err := L.LoadString("return " + line); err == nil { // try add return <...> then compile
			return line, nil
		} else {
			return multiline(line, rl, L)
		}
	} else {
		return "", err
	}
}

func multiline(ml string, rl *readline.Instance, L *lua.LState) (string, error) {
	for {
		if _, err := L.LoadString(ml); err == nil { // try compile
			return ml, nil
		} else if !incomplete(err) { // syntax error , but not EOF
			return ml, nil
		} else {
			rl.SetPrompt(">> ")
			if line, err := rl.Readline(); err == nil {
				ml = ml + "\n" + line
			} else {
				return "", err
			}
		}
	}
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

func main() {
	client = bw2.ConnectOrExit("")
	client.SetEntityFromEnvironOrExit()
	client.OverrideAutoChainTo(true)

	L := lua.NewState()
	defer L.Close()
	L.SetGlobal("getone", L.NewFunction(GetOne))
	L.SetGlobal("subscribe", L.NewFunction(Subscribe))
	L.SetGlobal("publish", L.NewFunction(Publish))
	doREPL(L)
}
