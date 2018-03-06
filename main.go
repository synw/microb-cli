package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/synw/microb-cli/libmicrob/prompter"
	"github.com/synw/microb-cli/libmicrob/state"
	"github.com/synw/terr"
)

var verbosity = flag.Int("v", 1, "Verbosity")

func main() {
	flag.Parse()
	// read conf
	trace := state.Init(*verbosity)
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
	prompter.Prompt()
}
