package prompter

import (
	"github.com/c-bata/go-prompt"
	"github.com/synw/microb-cli/libmicrob/cmds"
	"github.com/synw/microb-cli/libmicrob/cmds/handler"
	"github.com/synw/microb-cli/libmicrob/msgs"
	st "github.com/synw/microb-cli/libmicrob/state"
	cliTypes "github.com/synw/microb-cli/libmicrob/types"
	"github.com/synw/microb-cli/services"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/terr"
	"strings"
	"time"
)

func switchService(sname string) bool {
	/*
		Switch to a service to use in the shell
	*/
	srv, exists := services.Get(sname, st.State)
	if exists == false {
		msgs.Error("Service " + sname + " not found")
		return false
	}
	st.State.CurrentService = srv
	return true
}

func executor(in string) {
	/*
		Execute a command from user input
	*/
	state := st.State
	args := strings.Split(in, " ")
	cmdname := args[0]
	args = args[1:]
	var cmd *types.Cmd
	var isValid bool
	isLocal := false
	// check for local cli cmds
	if cmdname == "set" {
		exists := switchService(args[0])
		if exists == false {
			return
		}
		msgs.Ok("Switched to service " + args[0])
		return
	} else if cmdname == "unset" {
		var srv *types.Service
		state.CurrentService = srv
		msgs.Ok("Service unset")
		return
	} else if cmdname == "use" {
		cmd = cmds.Use()
		isLocal = true
	} else if cmdname == "using" {
		cmd = cmds.Using()
		isLocal = true
	}
	// get cmd args and encode them to an interface
	var cmdargs []interface{}
	if len(args) > 0 {
		var interfaceSlice []interface{} = make([]interface{}, len(args))
		for i, d := range args {
			interfaceSlice[i] = d
		}
		cmdargs = interfaceSlice
	}
	if isLocal == false {
		cmd, isValid = cmds.GetCmd(cmdname, cmdargs, state)
	} else {
		isValid = true
	}
	if isValid == true {
		cmd.Status = "pending"
		cmd.Date = time.Now()
		if len(cmdargs) > 0 {
			cmd.Args = cmdargs
		}
		// execute locally and exit if the command has an ExecCli function
		// this is used by the client for its local cmds
		if cmd.ExecCli != nil {
			run := cmd.ExecCli.(func(*types.Cmd, *cliTypes.State) (*types.Cmd, *terr.Trace))
			_, tr := run(cmd, state)
			if tr != nil {
				msg := "Can not execute local processing function for command " + in
				tr = tr.Add(msg)
				tr.Print()
				return
			}
			return
		}
		// otherwise send the command to the handler
		rescmd, timeout, tr := handler.SendCmd(cmd, state)
		if tr != nil {
			/*msg := "Can not execute command " + in
			tr := terr.Pass("prompter.executor", tr)
			msgs.Error(msg + "\n" + tr.Formatc())*/
			return
		}
		// execute eventual callback
		if rescmd.ExecAfter != nil {
			_, tr = rescmd.ExecAfter.(func(*types.Cmd, *cliTypes.State) (*types.Cmd, *terr.Trace))(rescmd, state)
			if tr != nil {
				msg := "Can not execute callback for command " + in
				tr := tr.Add(msg)
				msgs.Error(msg + "\n" + tr.Msg())
				return
			}
		}
		if timeout == true {
			msg := "Timeout: the server does not respond. Can not execute command"
			msgs.Timeout(msg)
			return
		}
	} else {
		msg := "Command " + in + " not found"
		msgs.Error(msg)
	}
}

func completer(in prompt.Document) []prompt.Suggest {
	/*
		Prompter completer: is disabled for now
	*/
	s := []prompt.Suggest{}
	return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)
}

func Prompt() {
	/*
		Main prompter function
	*/
	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix(">>> "),
		prompt.OptionTitle("Microb cli"),
	)
	p.Run()
}
