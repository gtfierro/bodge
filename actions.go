package main

import (
	"github.com/urfave/cli"
	"github.com/yuin/gopher-lua"
	bw2 "gopkg.in/immesys/bw2bind.v5"
)

func doInterpreter(c *cli.Context) error {
	// bw2 client
	client = bw2.ConnectOrExit(c.String("agent"))
	client.SetEntityFileOrExit(c.String("entity"))
	client.OverrideAutoChainTo(true)

	L := lua.NewState()
	defer L.Close()
	startScheduler(L)
	LoadLib(L)

	if c.NArg() == 0 {
		return doREPL(L)
	}
	path := c.Args().Get(0)
	return L.DoFile(path)
}
