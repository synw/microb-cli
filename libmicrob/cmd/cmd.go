package cmd

import (
	m "github.com/synw/microb/libmicrob"
	"github.com/synw/microb/libmicrob/types"
)

func GetCmds() map[string]*types.Cmd {
	cmds := make(map[string]*types.Cmd)
	cmds["use"] = Use()
	return cmds
}

func Use() *types.Cmd {
	cmd := &types.Cmd{
		Name: "use",
		Exec: use,
	}
	return cmd
}

func use(cmd *types.Cmd) *types.Cmd {
	server := cmd.Args[0].(string)
	m.Debug(server)
	return cmd
}
