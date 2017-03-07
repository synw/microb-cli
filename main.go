package main

import (
	"fmt"
	"errors"
	"flag"
	"github.com/abiosoft/ishell"
	"github.com/synw/terr"
	"github.com/synw/microb-cli/libmicrob/state"
	"github.com/synw/microb-cli/libmicrob/cmd"
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
		fmt.Println(terr.Ok("State initialized"))
	}
	shell.SetHomeHistoryPath(".ishell_history")
	// commands
	shell = cmd.GetCmds(shell)
	// start shell
    shell.Start()
}
