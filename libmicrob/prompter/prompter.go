package prompter

import (
	"errors"
	"github.com/c-bata/go-prompt"
	"github.com/synw/microb-cli/libmicrob/cmd"
	"github.com/synw/microb-cli/libmicrob/cmd/handler"
	"github.com/synw/terr"
	"strings"
)

var cmds = cmd.GetCmds()

func executor(in string) {
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
	if cmd, ok := cmds[cmdname]; ok {
		if len(cmdargs) > 0 {
			cmd.Args = cmdargs
		}
		// execute locally and exit if the command has an Exec function
		// this is used by the client for its local commands
		if cmd.Exec != nil {
			cmd = cmd.Exec(cmd)
			return
		}
		// otherwise send the command to the handler
		rescmd, timeout, tr := handler.SendCmd(cmd)
		if tr != nil {
			msg := "Can not execute command " + in
			err := errors.New(msg)
			tr := terr.Add("executor", err, tr)
			tr.Printc()
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
		err := errors.New(msg)
		tr := terr.New("executor", err)
		tr.Printc()
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
