package prompter

import (
	"errors"
	"github.com/c-bata/go-prompt"
	"github.com/synw/microb-cli/libmicrob/cmds"
	"github.com/synw/microb-cli/libmicrob/cmds/handler"
	"github.com/synw/microb-cli/libmicrob/msgs"
	st "github.com/synw/microb-cli/libmicrob/state"
	cliTypes "github.com/synw/microb-cli/libmicrob/types"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/terr"
	"strings"
)

func executor(in string) {
	state := st.State
	args := strings.Split(in, " ")
	cmdname := args[0]
	args = args[1:]
	// get cmd args and encode them to an interface
	var cmdargs []interface{}
	if len(args) > 0 {
		var interfaceSlice []interface{} = make([]interface{}, len(args))
		for i, d := range args {
			interfaceSlice[i] = d
		}
		cmdargs = interfaceSlice
	}
	cmd := cmds.New(cmdname)
	if cmds.IsValid(cmd, state) == true {
		cmd := state.Cmds[cmdname]
		cmd.Status = "pending"
		if len(cmdargs) > 0 {
			cmd.Args = cmdargs
		}
		// execute locally and exit if the command has an Exec function
		// this is used by the client for its local cmds
		if cmd.ExecCli != nil {
			run := cmd.ExecCli.(func(*types.Cmd, *cliTypes.State) (*types.Cmd, *terr.Trace))
			_, tr := run(cmd, state)
			if tr != nil {
				msg := "Can not execute local processing function for command " + in
				err := errors.New(msg)
				tr = terr.Add("executor", err, tr)
				tr.Printc()
				return
			}
			return
		}
		// otherwise send the command to the handler
		rescmd, timeout, tr := handler.SendCmd(cmd, state)
		if tr != nil {
			msg := "Can not execute command " + in
			err := errors.New(msg)
			tr := terr.Add("executor", err, tr)
			msgs.Error(msg + "\n" + tr.Formatc())
			return
		}
		if timeout == true {
			msg := "Timeout: the server does not respond. Can not execute command"
			err := errors.New(msg)
			tr := terr.New("executor", err)
			tr.Printc()
			return
		}
		terr.Debug(rescmd)
	} else {
		msg := "Command " + in + " not found"
		msgs.Error(msg)
	}
}

func completer(in prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{}
	return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)
}

func Prompt() {
	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix(">>> "),
		prompt.OptionTitle("Microb cli"),
	)
	p.Run()
}
