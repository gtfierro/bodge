package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	bw2 "github.com/immesys/bw2bind"
	"github.com/yuin/gopher-lua"
)

// takes the contents of the message and pops them onto the Lua stack
// in the following order:
// - object itself (extracted using the ponum)
// - uri it was published on
// -
func pushMsg(msg *bw2.SimpleMessage, ponum string, L *lua.LState) {
	var po bw2.PayloadObject
	if ponum == "" {
		po = msg.GetOnePODF("0.0.0.0/0")
	} else {
		po = msg.GetOnePODF(ponum)
	}
	var value interface{}
	var err error

	if po == nil {
		table := L.NewTable()
		L.Push(table)
		return
	}

	// try different serializations
	switch {
	case po.IsTypeDF("2.0.0.0/8"):
		err = po.(bw2.MsgPackPayloadObject).ValueInto(&value)
	case po.IsTypeDF("64.0.0.0/8"):
		value = po.(bw2.TextPayloadObject).Value()
		err = nil
	case po.IsTypeDF("65.0.0.0/8"):
		fmt.Println("json")
	case po.IsTypeDF("67.0.0.0/8"):
		err = po.(bw2.YAMLPayloadObject).ValueInto(&value)
	}

	if err != nil {
		L.RaiseError("Could not deserialize msg (%v): %+v", err, msg)
		return
	}
	L.Push(toLValue(value, L))
}

// turn an interface{} into a lua.LValue
func toLValue(value interface{}, L *lua.LState) lua.LValue {
	switch v := value.(type) {
	case string:
		return lua.LString(v)
	case int:
		return lua.LNumber(float64(v))
	case uint:
		return lua.LNumber(float64(v))
	case uint64:
		return lua.LNumber(float64(v))
	case int64:
		return lua.LNumber(float64(v))
	case float64:
		return lua.LNumber(v)
	case map[interface{}]interface{}:
		table := L.NewTable()
		for k, v := range v {
			L.SetTable(table, toLValue(k, L), toLValue(v, L))
		}
		return table
	case []interface{}:
		table := L.NewTable()
		for _, val := range v {
			table.Append(toLValue(val, L))
		}
		return table
	case []uint8:
		// TODO: binary array
		table := L.NewTable()
		for _, val := range v {
			table.Append(lua.LNumber(val))
		}
		return table
	case bool:
		if v {
			return lua.LNumber(1)
		} else {
			return lua.LNumber(0)
		}
	default:
		L.RaiseError("Could not figure out type %T of %+v", value, value)
	}
	return lua.LNil
}

func luaToGo(val lua.LValue) interface{} {
	var value interface{}

	switch val.Type() {
	case lua.LTNil:
		value = nil
	case lua.LTBool:
		value = bool(val.(lua.LBool))
	case lua.LTNumber:
		value = float64(val.(lua.LNumber))
	case lua.LTString:
		value = string(val.(lua.LString))
	case lua.LTTable:
		table := val.(*lua.LTable)
		if table.MaxN() == 0 { // map
			m := make(map[interface{}]interface{})
			table.ForEach(func(k, v lua.LValue) {
				m[luaToGo(k)] = luaToGo(v)
			})
			value = m
		} else { // array
			a := make([]interface{}, table.MaxN())
			i := 0
			table.ForEach(func(k, v lua.LValue) {
				a[i] = luaToGo(v)
				i += 1
			})
			value = a
		}
	}

	return value
}

func lvalueToPO(val lua.LValue, ponum string) bw2.PayloadObject {

	var err error
	var po bw2.PayloadObject
	value := luaToGo(val)
	intponum := bw2.FromDotForm(ponum)

	// try different serializations
	switch {
	case isTypeDF(ponum, "2.0.0.0/8"):
		po, err = bw2.CreateMsgPackPayloadObject(intponum, value)
	case isTypeDF(ponum, "64.0.0.0/8"):
		po = bw2.CreateTextPayloadObject(intponum, value.(string))
		err = nil
	case isTypeDF(ponum, "65.0.0.0/8"):
		fmt.Println("json")
	case isTypeDF(ponum, "67.0.0.0/8"):
		po, err = bw2.CreateYAMLPayloadObject(intponum, value)
	}
	if err != nil {
		fmt.Println(err)
	}

	return po
}

func isType(testponum, ponum int, mask int) bool {
	return (ponum >> uint(32-mask)) == (testponum >> uint(32-mask))
}

func isTypeDF(testponum, ponum string) bool {
	parts := strings.SplitN(ponum, "/", 2)
	var mask int
	var err error
	if len(parts) != 2 {
		mask = 32
	} else {
		mask, err = strconv.Atoi(parts[1])
		if err != nil {
			panic("malformed masked dot form")
		}
	}
	_ponum := bw2.FromDotForm(parts[0])
	_testponum := bw2.FromDotForm(testponum)
	return isType(_ponum, _testponum, mask)
}

func resolveURInamespace(uri string) string {
	chunks := strings.Split(uri, "/")
	bytes, err := bw2.FromBase64(chunks[0])
	if err != nil {
		log.Fatal(err)
	}
	alias, _ := client.UnresolveAlias(bytes)
	chunks[0] = alias
	return strings.Join(chunks, "/")
}
