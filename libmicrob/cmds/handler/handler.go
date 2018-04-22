package handler

import (
	"encoding/json"
	"errors"
	"github.com/SKAhack/go-shortid"
	//"github.com/davecgh/go-spew/spew"
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
	// get an id for the command
	g := shortid.Generator()
	cmd.Id = g.Generate()
	// update the state for pending commands
	if cmd.ExecCli == nil {
		state.CmdIds = append(state.CmdIds, cmd.Id)
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
		// check if this client issued the command
		newCmds, exec := updateCmdIds(cmd, state)
		// update the pending commands list
		state.CmdIds = newCmds
		if exec == false {
			//err := errors.New("Command from another client")
			//tr := terr.New("cmds.handler.SendCmd", err)
			return cmd, false, nil
		}
		// process the command
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
		newCmds, _ := updateCmdIds(cmd, state)
		// update the pending commands list
		state.CmdIds = newCmds
		return cmd, true, nil
	}
	return cmd, false, nil
}

func updateCmdIds(cmd *types.Cmd, state *cliTypes.State) ([]string, bool) {
	var newCmds []string
	isIn := false
	for _, scmdId := range state.CmdIds {
		if scmdId != cmd.Id {
			newCmds = append(newCmds, cmd.Id)
		} else {
			isIn = true
		}
	}
	return newCmds, isIn
}

func sendCommand(cmd *types.Cmd, state *cliTypes.State) *terr.Trace {
	var cmdp types.Cmd
	cmdp.Name = cmd.Name
	cmdp.Date = cmd.Date
	cmdp.Args = cmd.Args
	cmdp.From = cmd.From
	cmdp.Status = cmd.Status
	cmdp.ErrMsg = cmd.ErrMsg
	cmdp.NoLog = cmd.NoLog
	cmdp.Service = cmd.Service
	cmdp.ReturnValues = cmd.ReturnValues
	cmdp.Id = cmd.Id
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
