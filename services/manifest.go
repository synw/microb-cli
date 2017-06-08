package services

import (
	"github.com/abiosoft/ishell"
	"github.com/synw/microb-cli/services/httpServer"
)

var commands = []*ishell.Cmd{
	httpServer.Cmds(),
}
