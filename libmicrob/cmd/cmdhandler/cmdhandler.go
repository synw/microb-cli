package cmdhandler

import (
	"errors"
	"time"
	"encoding/json"
	"github.com/ventu-io/go-shortid"
	"github.com/synw/terr"
	"github.com/synw/microb/libmicrob/datatypes"
	"github.com/synw/microb/libmicrob/cmd"
	"github.com/synw/microb-cli/libmicrob/state"
)


func SendCmd(command *datatypes.Command) (*datatypes.Command, *terr.Trace) {
	// check if server is set
	if state.Server == nil {
		err := errors.New("No server selected. Set it with: use server_name")
		trace := terr.New("cmd.cmdhandler.SendCmd", err)
		return command, trace
	}
	// check if cli is connected
	if state.Cli.IsConnected == false {
		err := errors.New("No connection to server: use server_name")
		trace := terr.New("cmd.cmdhandler.SendCmd", err)
		return command, trace
	}	
	// check the validity of the command
	if cmd.IsValid(command) != true {
		err := errors.New("Command "+command.Name+" unknown")
		trace := terr.New("cmd.cmdhandler.SendCmd", err)		
		return command, trace
	}
	trace := sendCommand(command)
	if trace != nil {
		trace := terr.Pass("cmd.cmdhandler.SendCmd", trace)
		return command, trace
	}
	return command, nil
}


func New(name string, from string, reason string, args ...interface{}) *datatypes.Command {
	id, _ := shortid.Generate()
	date := time.Now()
	status := "pending"
	var err error
	var rvs []interface{}
	command := &datatypes.Command{
		id,
		name,
		from,
		reason,
		date,
		args,
		status,
		err,
		rvs,
	}
	return command
}

func sendCommand(command *datatypes.Command) *terr.Trace {
	payload, err := json.Marshal(command)
	if err != nil {
		msg := "Unable to marshall json: "+err.Error()
		err := errors.New(msg)
		trace := terr.New("cmd.cmdhandler.sendCommand", err)
		return trace
	}
	_, err = state.Cli.Http.Publish(state.Server.CmdChannel, payload)
	if err != nil {
		trace := terr.New("cmd.cmdhandler.sendCommand", err)
		return trace
	}
	return nil
}
