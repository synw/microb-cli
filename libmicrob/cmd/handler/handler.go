package handler

import (
	"errors"
	"time"
	"encoding/json"
	"github.com/abiosoft/ishell"
	"github.com/ventu-io/go-shortid"
	"github.com/synw/terr"
	"github.com/synw/centcom"
	"github.com/synw/microb/libmicrob/datatypes"
	"github.com/synw/microb/libmicrob/cmd"
	"github.com/synw/microb-cli/libmicrob/state"
)


func SendCmd(command *datatypes.Command, ctx *ishell.Context) (*datatypes.Command, bool, *terr.Trace) {
	timeout := false
	// check if server is set
	if state.Server == nil {
		err := errors.New("No server selected. Set it with: use server_name")
		trace := terr.New("cmd.handler.SendCmd", err)
		return command, timeout, trace
	}
	// check if cli is connected
	if state.Cli.IsConnected == false {
		err := errors.New("No connection to server: use server_name")
		trace := terr.New("cmd.handler.SendCmd", err)
		return command, timeout, trace
	}	
	// check the validity of the command
	if cmd.IsValid(command) != true {
		err := errors.New("Command "+command.Name+" unknown")
		trace := terr.New("cmd.handler.SendCmd", err)		
		return command, timeout, trace
	}
	trace := sendCommand(command)
	if trace != nil {
		trace := terr.Pass("cmd.handler.SendCmd", trace)
		return command, timeout, trace
	}
	// wait for results
	cli := centcom.New(state.Server.WsHost, state.Server.WsPort, state.Server.WsKey)
	err := centcom.Connect(cli)
	if err != nil {
		err := errors.New(err.Error())
		trace := terr.New("cmd.handler.SendCmd", err)		
		return command, timeout, trace
	}
	defer centcom.Disconnect(cli)
	err = cli.Subscribe(state.Server.CmdChannel)
	if err != nil {
		err := errors.New(err.Error())
		trace := terr.New("cmd.handler.SendCmd", err)		
		return command, timeout, trace
	}
	select {
	case returnCmd := <- cli.Channels:
		ctx.Println(returnCmd)
	case <-time.After(10*time.Second):
		return command, true, nil
	}
	return command, timeout, nil
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
		trace := terr.New("cmd.handler.sendCommand", err)
		return trace
	}
	_, err = state.Cli.Http.Publish(state.Server.CmdChannel, payload)
	if err != nil {
		trace := terr.New("cmd.handler.sendCommand", err)
		return trace
	}
	return nil
}
