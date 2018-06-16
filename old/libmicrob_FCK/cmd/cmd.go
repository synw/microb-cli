package cmd

import (
	"github.com/synw/microb/types"
	"github.com/synw/microb-cli/cmd/info"
	cmd_state "github.com/synw/microb-cli/cmd/state"
)

func GetCmds() []*types.Command { 
	ping = info.Ping()
	return [ping]
}
