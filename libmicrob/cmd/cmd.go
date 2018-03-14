package cmd

import (
	"errors"
	"github.com/synw/microb-cli/libmicrob/cmd/handler"
	"github.com/synw/microb-cli/libmicrob/msgs"
	"github.com/synw/microb-cli/libmicrob/state"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/terr"
	"time"
)

func New(name string, args ...map[string]interface{}) *types.Cmd {
	date := time.Now()
	var service string
	var params []interface{}
	var from string
	var status string
	var errMsg string
	var trace *terr.Trace
	var returnValues []interface{}
	var execCli interface{}
	var execAfter interface{}
	if len(args) == 1 {
		for k, v := range args[0] {
			if k == "service" {
				service = v.(string)
			} else if k == "params" {
				params = v.([]interface{})
			} else if k == "from" {
				from = v.(string)
			} else if k == "status" {
				status = v.(string)
			} else if k == "errMsg" {
				errMsg = v.(string)
			} else if k == "trace" {
				trace = v.(*terr.Trace)
			} else if k == "returnValues" {
				returnValues = v.([]interface{})
			} else if k == "execCli" {
				execCli = v.(interface{})
			} else if k == "execAfter" {
				execAfter = v.(interface{})
			}
		}
	}
	cmd := &types.Cmd{
		Name:         name,
		Date:         date,
		Service:      service,
		Args:         params,
		From:         from,
		Status:       status,
		ErrMsg:       errMsg,
		Trace:        trace,
		ReturnValues: returnValues,
		ExecCli:      execCli,
		ExecAfter:    execAfter,
	}
	msgs.Debug(cmd)
	return cmd
}

func GetCmds() map[string]*types.Cmd {
	cmds := make(map[string]*types.Cmd)
	cmds["use"] = Use()
	cmds["using"] = Using()
	cmds["ping"] = Ping()
	cmds["services"] = Services()
	return cmds
}

func Use() *types.Cmd {
	cmd := &types.Cmd{
		Name:    "use",
		ExecCli: use,
	}
	return cmd
}

func Using() *types.Cmd {
	cmd := &types.Cmd{
		Name:    "using",
		ExecCli: using,
	}
	return cmd
}

func Ping() *types.Cmd {
	args := make(map[string]interface{})
	args["service"] = "infos"
	args["execAfter"] = afterPing
	cmd := New("ping", args)
	return cmd
}

func Services() *types.Cmd {
	args := make(map[string]interface{})
	args["service"] = "infos"
	cmd := New("services", args)
	return cmd
}

func srv(cmd *types.Cmd) *types.Cmd {
	cmd, timeout, tr := handler.SendCmd(cmd)
	if tr != nil {
		tr = terr.Pass("cmd.srv", tr)
		tr.Print()
		return cmd
	}
	if timeout == true {
		msgs.Timeout("The server is not responding")
	}
	return cmd
}

func ping(cmd *types.Cmd) *types.Cmd {
	cmd, timeout, tr := handler.SendCmd(cmd)
	if tr != nil {
		tr = terr.Pass("cmd.ping", tr)
		tr.Print()
		return cmd
	}
	if timeout == true {
		msgs.Timeout("The server is not responding")
	}
	return cmd
}

func afterPing(cmd *types.Cmd) (*types.Cmd, *terr.Trace) {
	msgs.Debug(cmd.ReturnValues)
	return cmd, nil
}

func using(cmd *types.Cmd) (*types.Cmd, *terr.Trace) {
	if state.Server == nil {
		msgs.Warning("No server selected: try the use command: ex: " + msgs.Bold("use") + " server1")
	} else {
		msgs.Msg("Using server " + msgs.Bold(state.Server.Name))
	}
	return cmd, nil
}

func use(cmd *types.Cmd) (*types.Cmd, *terr.Trace) {
	if len(cmd.Args) != 1 {
		msgs.Warning("Please provide a server name: ex: use localhost")
		err := errors.New("Can not find server")
		tr := terr.New("cmd.use", err)
		return cmd, tr
	}
	server := cmd.Args[0].(string)
	tr := serverExists(server)
	if tr != nil {
		tr = terr.Pass("comd.state.Use", tr)
		return cmd, tr
	}
	state.Server = state.Servers[server]
	tr = state.InitServer()
	if tr != nil {
		tr = terr.Pass("comd.state.Use", tr)
		return cmd, tr
	}
	// init cli and check server
	if tr != nil {
		err := errors.New("can not connect to websockets server: check your config")
		tr := terr.Add("cmd.state.Use", err, tr)
		return cmd, tr
	} else {
		msg := "Connnected to server " + server
		msgs.Ready(msg)
	}
	return cmd, nil
}

func serverExists(server_name string) *terr.Trace {
	for name, _ := range state.Servers {
		if server_name == name {
			return nil
		}
	}
	msg := "Server " + server_name + " not found: please check your config file"
	err := errors.New(msg)
	tr := terr.New("ws.serverExists", err)
	return tr
}
