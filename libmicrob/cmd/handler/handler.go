package handler

import (
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/terr"
)

func SendCmd(cmd *types.Cmd) (*types.Cmd, bool, *terr.Trace) {
	return cmd, false, nil
}
