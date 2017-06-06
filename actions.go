package main

import (
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"github.com/yuin/gopher-lua"
	bw2 "gopkg.in/immesys/bw2bind.v5"
)

func doInterpreter(c *cli.Context) error {
	// bw2 client
	client = bw2.ConnectOrExit(c.String("agent"))
	client.SetEntityFileOrExit(c.String("entity"))
	client.OverrideAutoChainTo(true)

	// lua state
	L := lua.NewState()
	defer L.Close()
	LoadLib(L)

	return doREPL(L)
}

func doFile(c *cli.Context) error {
	// bw2 client
	client = bw2.ConnectOrExit(c.String("agent"))
	client.SetEntityFileOrExit(c.String("entity"))
	client.OverrideAutoChainTo(true)

	if c.NArg() == 0 {
		return errors.New("Need to specify a file to run")
	}
	path := c.Args().Get(0)

	// lua state
	return RunFile(path)
}
