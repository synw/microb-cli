package infos

import (
	"github.com/synw/microb-http/cmd"
	"github.com/synw/microb-http/state"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/terr"
)

var Service *types.Service = &types.Service{
	"infos",
	[]string{"use", "using", "ping", "services"},
	ini,
	dispatch,
}

func ini(dev bool, verbosity int, start bool) *terr.Trace {
	return state.Init(dev, verbosity, start)
}

func dispatch(c *types.Command) *types.Command {
	return cmd.Dispatch(c)
}
