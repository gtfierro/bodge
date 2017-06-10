package main

import (
	"github.com/urfave/cli"
	"github.com/yuin/gopher-lua"
	bw2 "gopkg.in/immesys/bw2bind.v5"
	"log"
)

func doInterpreter(c *cli.Context) error {
	// bw2 client
	client = bw2.ConnectOrExit(c.String("agent"))
	client.SetEntityFileOrExit(c.String("entity"))
	client.OverrideAutoChainTo(true)

	if c.NArg() == 0 {
		// lua state
		L := lua.NewState()
		defer L.Close()
		LoadLib(L)

		return doREPL(L)
	}
	path := c.Args().Get(0)

	log.Fatal(RunFile(path))
	return nil
}
