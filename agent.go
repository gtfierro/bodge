// +build !ragent

package main

import (
	bw2 "github.com/immesys/bw2bind"
	"github.com/urfave/cli"
)

func GetClient(c *cli.Context) *bw2.BW2Client {
	client = bw2.ConnectOrExit(c.String("agent"))
	client.SetEntityFileOrExit(c.String("entity"))
	return client
}
