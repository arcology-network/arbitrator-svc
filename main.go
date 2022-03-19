package main

import (
	"os"

	// "net/http"
	// _ "net/http/pprof"

	"github.com/arcology-network/3rd-party/tm/cli"
	"github.com/arcology-network/arbitrator-svc/service"
)

func main() {

	// go func() {
	// 	http.ListenAndServe("0.0.0.0:8090", nil)
	// }()

	st := service.StartCmd

	cmd := cli.PrepareMainCmd(st, "BC", os.ExpandEnv("$HOME/monacos/arbitrator"))
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
