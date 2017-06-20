package main

import (
	"fmt"
	bw2 "github.com/immesys/bw2bind"
	"github.com/urfave/cli"
	"github.com/yuin/gopher-lua"
	"io/ioutil"
	"log"
	"os"
)

func init() {
	bw2.SilenceLog()
}

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
		log.Fatal(doREPL(L))
		return nil
	}
	path := c.Args().Get(0)

	// check if the file exists locally
	if _, err := os.Stat(path); os.IsNotExist(err) {
		msg := client.QueryOneOrExit(&bw2.QueryParams{
			URI: path,
		})
		po := msg.GetOnePODF("64.0.2.0/24")
		luafile := po.(bw2.TextPayloadObject).Value()
		log.Fatal(L.DoString(luafile))
	}
	log.Fatal(L.DoFile(path))
	return nil
}

func doURI(c *cli.Context) error {
	// bw2 client
	client = bw2.ConnectOrExit(c.String("agent"))
	client.SetEntityFileOrExit(c.String("entity"))
	client.OverrideAutoChainTo(true)

	if c.NArg() != 2 {
		log.Fatal("Need 2 arguments: <file> <uri to publish>")
	}
	path := c.Args().Get(0)
	uri := c.Args().Get(1)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatalf("File %s does not exist", path)
	}

	luafilebytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Could not read file %s: %s", path, err)
	}

	ponum, _ := bw2.PONumFromDotForm("64.0.2.0")
	po, err := bw2.LoadTextPayloadObject(ponum, luafilebytes)
	if err != nil {
		log.Fatalf("Could not create PO: %s", err)
	}

	client.PublishOrExit(&bw2.PublishParams{
		URI:            uri,
		Persist:        true,
		PayloadObjects: []bw2.PayloadObject{po},
	})
	return nil
}

func doCat(c *cli.Context) error {
	// bw2 client
	client = bw2.ConnectOrExit(c.String("agent"))
	client.SetEntityFileOrExit(c.String("entity"))
	client.OverrideAutoChainTo(true)

	if c.NArg() != 1 {
		log.Fatal("Need 1 argument: <uri to cat>")
	}
	msg := client.QueryOneOrExit(&bw2.QueryParams{
		URI: c.Args().Get(0),
	})
	po := msg.GetOnePODF("64.0.2.0/24")
	luafile := po.(bw2.TextPayloadObject).Value()
	fmt.Println(luafile)

	return nil
}

func doLs(c *cli.Context) error {
	// bw2 client
	client = bw2.ConnectOrExit(c.String("agent"))
	client.SetEntityFileOrExit(c.String("entity"))
	client.OverrideAutoChainTo(true)

	if c.NArg() != 1 {
		log.Fatal("Need 1 argument: <uri to ls>")
	}
	msgs, err := client.Query(&bw2.QueryParams{
		URI: c.Args().Get(0),
	})
	if err != nil {
		log.Fatalf("Could not LS %s: %s", c.Args().Get(0), err)
	}
	for msg := range msgs {
		po := msg.GetOnePODF("64.0.2.0/24")
		if po != nil {
			fmt.Println(resolveURInamespace(msg.URI))
		}
	}

	return nil
}
