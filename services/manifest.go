package services

import (
	"github.com/abiosoft/ishell"
	http "github.com/synw/microb/services/httpServer/cmd/cli"
)

var commands = []*ishell.Cmd{
	http.Cmds(),
}
