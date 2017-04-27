package handler

import (
	"errors"
	"time"
	"encoding/json"
	"github.com/abiosoft/ishell"
	"github.com/synw/terr"
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
	err := state.Cli.Subscribe(state.Server.CmdChannel)
	if err != nil {
		err := errors.New(err.Error())
		trace := terr.New("cmd.handler.SendCmd", err)		
		return command, timeout, trace
	}
	defer state.Cli.Unsubscribe(state.Server.CmdChannel)
	select {
	case returnCmd := <- state.Cli.Channels:
		com, _ := cmd.CmdFromPayload(returnCmd.Payload)
		for _, val := range(com.ReturnValues) {
			ctx.Println(val.(string))
		}
	case <-time.After(10*time.Second):
		return command, true, nil
	}
	return command, timeout, nil
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
