// +build ragent

package main

import (
	"fmt"
	"log"
	"time"

	bw2 "github.com/immesys/bw2bind"
	"github.com/immesys/ragent/ragentlib"
	"github.com/urfave/cli"
)

var _agent = "127.0.0.1:28588"

const serverVK = "dhv8gHPlQuqs9RDEgH5PzX694YREzCcsocHitVNdZqQ="

func GetClient(c *cli.Context) *bw2.BW2Client {
	time.Sleep(5 * time.Second)
	client = bw2.ConnectOrExit(_agent)
	client.SetEntity([]byte(_entity)[1:])
	return client
}

func init() {
	go func() {
		defer func() {
			r := recover()
			if r != nil {
				log.Fatal(fmt.Sprintf("Failed to connect ragent (%v)", r))
			}
		}()
		fmt.Println("connecting")
		ragentlib.DoClientER([]byte(_entity), "128.32.37.244:28591", serverVK, _agent)
	}()
}
