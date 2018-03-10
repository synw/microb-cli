package handler

import (
	"encoding/json"
	"errors"
	"github.com/synw/microb-cli/libmicrob/state"
	m "github.com/synw/microb/libmicrob"
	command "github.com/synw/microb/libmicrob/cmd"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/terr"
	"time"
)

func SendCmd(cmd *types.Cmd) (*types.Cmd, bool, *terr.Trace) {
	timeout := false
	// check if server is set
	if state.Server == nil {
		err := errors.New("No server selected. Set it with: use server_name")
		trace := terr.New("cmd.handler.SendCmd", err)
		return cmd, timeout, trace
	}
	// check if cli is connected
	if state.Cli.IsConnected == false {
		err := errors.New("No connection to server: use server_name")
		trace := terr.New("cmd.handler.SendCmd", err)
		return cmd, timeout, trace
	}
	// send the command
	trace := sendCommand(cmd)
	if trace != nil {
		trace := terr.Pass("cmd.handler.SendCmd from sendCommand", trace)
		return cmd, timeout, trace
	}
	// wait for results
	select {
	case returnCmd := <-state.Cli.Channels:
		com, _ := command.ConvertPayload(returnCmd.Payload)
		//ctx.Println(com.Name, com.Status, com.ErrMsg)
		if com.ErrMsg != "" {
			m.Error(com.ErrMsg)
		} else {
			for i, val := range com.ReturnValues {
				if i == 0 {
					m.Ok(val.(string))
				} else {
					m.Msg(val.(string))
				}
			}
		}
	case <-time.After(10 * time.Second):
		return cmd, true, nil
	}
	return cmd, false, nil
}

func sendCommand(cmd *types.Cmd) *terr.Trace {
	var cmdp types.Cmd
	cmdp.Name = cmd.Name
	cmdp.Date = cmd.Date
	cmdp.Args = cmd.Args
	cmdp.From = cmd.From
	cmdp.Status = cmd.Status
	cmdp.ErrMsg = cmd.ErrMsg
	cmdp.ReturnValues = cmd.ReturnValues
	payload, err := json.Marshal(cmdp)
	if err != nil {
		msg := "Unable to marshall json: " + err.Error()
		err := errors.New(msg)
		trace := terr.New("cmd.handler.sendCommand", err)
		return trace
	}
	_, err = state.Cli.Http.Publish(state.Server.CmdChanIn, payload)
	if err != nil {
		trace := terr.New("cmd.handler.sendCommand", err)
		return trace
	}
	return nil
}
