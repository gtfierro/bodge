package main

import (
	"os"

	"github.com/urfave/cli"
	bw2 "gopkg.in/immesys/bw2bind.v5"
)

var client *bw2.BW2Client

func main() {

	app := cli.NewApp()
	app.Name = "bodge"
	app.Version = "0.1.0"
	app.Usage = "Simple BOSSWAVE Lua scripts for interaction, exploration and rule building"

	app.Commands = []cli.Command{
		{
			Name:    "interpreter",
			Aliases: []string{"int", "i"},
			Usage:   "Opens a Lua interpreter",
			Action:  doInterpreter,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "agent,a",
					Value:  "127.0.0.1:28589",
					Usage:  "Local BOSSWAVE Agent",
					EnvVar: "BW2_AGENT",
				},
				cli.StringFlag{
					Name:   "entity,e",
					Value:  "",
					Usage:  "The entity to use",
					EnvVar: "BW2_DEFAULT_ENTITY",
				},
			},
		},
		{
			Name:    "run",
			Aliases: []string{"file", "r"},
			Usage:   "Runs a Lua file",
			Action:  doFile,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "agent,a",
					Value:  "127.0.0.1:28589",
					Usage:  "Local BOSSWAVE Agent",
					EnvVar: "BW2_AGENT",
				},
				cli.StringFlag{
					Name:   "entity,e",
					Value:  "",
					Usage:  "The entity to use",
					EnvVar: "BW2_DEFAULT_ENTITY",
				},
			},
		},
	}

	app.Run(os.Args)
}
