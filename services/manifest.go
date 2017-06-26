package services

import (
	"github.com/abiosoft/ishell"
	"github.com/synw/microb-cli/services/httpServer"
	dashboard "github.com/synw/microb-dashboard/cmd/cli"
	grg "github.com/synw/microb-goregraph/cmd/cli"
)

var commands = []*ishell.Cmd{
	httpServer.Cmds(),
	grg.Cmds(),
	dashboard.Cmds(),
}
