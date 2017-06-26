package services

import (
	"github.com/abiosoft/ishell"
	dashboard "github.com/synw/microb-dashboard/cmd/cli"
	http "github.com/synw/microb-http/cmd/cli"
)

var commands = []*ishell.Cmd{
	dashboard.Cmds(),
	http.Cmds(),
}
