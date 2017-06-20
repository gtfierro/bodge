package main

import (
	"os"

	bw2 "github.com/immesys/bw2bind"
	"github.com/urfave/cli"
)

var client *bw2.BW2Client

func main() {

	app := cli.NewApp()
	app.Name = "bodge"
	app.Version = "0.1.0"
	app.Usage = "Simple BOSSWAVE Lua scripts for interaction, exploration and rule building"
	app.Action = doInterpreter
	app.Flags = []cli.Flag{
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
	}
	app.Commands = []cli.Command{
		{
			Name:    "publish",
			Aliases: []string{"p", "pub"},
			Usage:   "Publish a file to a given URI",
			Action:  doURI,
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
			Name:   "cat",
			Usage:  "Cat a file on a given URI",
			Action: doCat,
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
			Name:   "ls",
			Usage:  "List bodge files",
			Action: doLs,
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
