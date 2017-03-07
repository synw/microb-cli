package state

import (
	"fmt"
	"github.com/abiosoft/ishell"
	//"github.com/synw/terr"
)


func Use() *ishell.Cmd {
	command := &ishell.Cmd{
        Name: 	"use",
        Help: 	"Use server: use server_domain",
        Func: 	func(ctx *ishell.Context) {
					fmt.Println("Bar")
				},
    }
	return command
}
