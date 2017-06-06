package main

import (
	"github.com/yuin/gopher-lua"
)

var TAGS map[string]string

// API to implement

// executes the lua code in the file, starting in a new Lua state
func RunFile(path string) error {
	L := lua.NewState()
	defer L.Close()
	LoadLib(L)
	return L.DoFile(path)
}

// imports the lua code in the given file into the current Lua state
func ImportFile(path string, L *lua.LState) error {
	//return L.LoadFile(path)
	return nil
}

// executes the lua code in the file, starting in a new Lua state
func RunURI(uri string) {
}

// executes the lua code in the file, starting in a new Lua state
func ImportURI(uri string, L *lua.LState) {
}

// saves any subsequent lines under the named tag
func StartSave(tag string) {
}

// stops saving lines to the tag
func EndSave(tag string) {
}

// publish the contents of the tag so far to the given URI
func PublishTag(tag, uri string) {
}

// save the contents of the tag so far to the given file
func SaveTagToFile(tag, file string) {
}
