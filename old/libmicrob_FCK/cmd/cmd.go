package cmd

import (
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/microb-cli/libmicrob/cmd/info"
	cmd_state "github.com/synw/microb-cli/libmicrob/cmd/state"
)

func GetCmds() []*types.Command { 
	ping = info.Ping()
	return [ping]
}
