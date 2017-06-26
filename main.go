package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/abiosoft/ishell"
	"github.com/synw/microb-cli/libmicrob/cmd"
	"github.com/synw/microb-cli/libmicrob/state"
	"github.com/synw/microb-cli/services"
	"github.com/synw/terr"
)

var dev_mode = flag.Bool("d", false, "Dev mode")
var verbosity = flag.Int("v", 1, "Verbosity")
var shell = ishell.New()

func main() {
	flag.Parse()
	// read conf
	trace := state.InitState(*dev_mode, *verbosity)
	if trace != nil {
		err := errors.New("Unable to init state")
		trace := terr.Add("main", err, trace)
		terr.Fatal("main", trace)
	}
	if state.Verbosity > 2 {
		terr.Ok("State initialized")
	}
	srvs := "Available servers:"
	for name, _ := range state.Servers {
		srvs = srvs + " " + name
	}
	fmt.Println(srvs)
	shell.SetHomeHistoryPath(".ishell_history")
	// commands
	shell = cmd.GetCmds(shell)
	shell = services.GetCmds(shell)
	// start shell
	shell.Start()
}
