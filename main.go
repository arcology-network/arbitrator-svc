package main

import (
	"os"

	"github.com/arcology-network/3rd-party/tm/cli"
	"github.com/arcology-network/arbitrator-svc/service"
)

func main() {

	st := service.StartCmd

	cmd := cli.PrepareMainCmd(st, "BC", os.ExpandEnv("$HOME/monacos/arbitrator"))
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
