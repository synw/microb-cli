package handler

import (
	"encoding/json"
	"errors"
	"github.com/synw/microb-cli/libmicrob/msgs"
	cliTypes "github.com/synw/microb-cli/libmicrob/types"
	"github.com/synw/microb/libmicrob/cmds"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/terr"
	"time"
)

func SendCmd(cmd *types.Cmd, state *cliTypes.State) (*types.Cmd, bool, *terr.Trace) {
	timeout := false
	// check if server is set
	if state.WsServer == nil {
		err := errors.New("No server selected. Set it with: use server_name")
		tr := terr.New("cmd.handler.SendCmd", err)
		return cmd, timeout, tr
	}
	// check if cli is connected
	if state.Cli.IsConnected == false {
		err := errors.New("No connection to server: use server_name")
		tr := terr.New("cmd.handler.SendCmd", err)
		return cmd, timeout, tr
	}
	// send the cmds
	tr := sendCommand(cmd, state)
	if tr != nil {
		tr := terr.Pass("cmd.handler.SendCmd", tr)
		return cmd, timeout, tr
	}
	// wait for results
	select {
	case returnCmd := <-state.Cli.Channels:
		cmd := cmds.ConvertPayload(returnCmd.Payload)
		if cmd.Status != "success" {
			msgs.Error(cmd.ErrMsg)
			return cmd, false, cmd.Trace
		} else {
			if cmd.ExecAfter != nil {
				run := cmd.ExecAfter.(func(*types.Cmd) (*types.Cmd, *terr.Trace))
				cmd, tr = run(cmd)
				if tr != nil {
					err := errors.New("Can not execute processing function")
					tr := terr.Add("cmd.handler.SendCmd", err, tr)
					return cmd, false, tr
				}
				return cmd, false, nil
			} else {
				for i, val := range cmd.ReturnValues {
					if i == 0 {
						msgs.Ok(val.(string))
					} else {
						msgs.Msg(val.(string))
					}
				}
			}
		}
	case <-time.After(10 * time.Second):
		return cmd, true, nil
	}
	return cmd, false, nil
}

func sendCommand(cmd *types.Cmd, state *cliTypes.State) *terr.Trace {
	var cmdp types.Cmd
	cmdp.Name = cmd.Name
	cmdp.Date = cmd.Date
	cmdp.Args = cmd.Args
	cmdp.From = cmd.From
	cmdp.Status = cmd.Status
	cmdp.ErrMsg = cmd.ErrMsg
	cmdp.Service = cmd.Service
	cmdp.ReturnValues = cmd.ReturnValues
	payload, err := json.Marshal(cmdp)
	if err != nil {
		msg := "Unable to marshall json: " + err.Error()
		err := errors.New(msg)
		tr := terr.New("cmd.handler.sendCommand", err)
		return tr
	}
	_, err = state.Cli.Http.Publish(state.WsServer.CmdChanIn, payload)
	if err != nil {
		tr := terr.New("cmd.handler.sendCommand", err)
		return tr
	}
	return nil
}
