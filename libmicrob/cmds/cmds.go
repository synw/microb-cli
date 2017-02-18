package cmds

import (
	"time"
	"errors"
	"fmt"
	"sync"
	"encoding/json"
	"github.com/abiosoft/ishell"
	"github.com/centrifugal/gocent"
	"github.com/synw/microb/libmicrob/datatypes"
	"github.com/synw/microb/libmicrob/datatypes/encoding"
	"github.com/synw/microb/libmicrob/commands/methods"
	"github.com/synw/microb-cli/libmicrob/metadata"
	"github.com/synw/microb-cli/libmicrob/listeners"
	appevents "github.com/synw/microb/libmicrob/events"
)


var Debug = metadata.IsDebug()

func SendCmd(ctx *ishell.Context, name string, server *datatypes.Server) (*datatypes.Command, error, string) {
	// create struct
	var rvs []interface{}
	args := getArgs(ctx)
	id := encoding.GenerateId()
	command := &datatypes.Command{id, name, "cli", "", time.Now(), args, "pending", nil, rvs}
	// check validity
	var msg string
	if commands_methods.IsValid(command) != true {
		msg = "Command "+command.Name+" unknown"
		//c.Err(errors.New(msg))
		return command, errors.New(msg), msg
	}
	// wait for results
	c_feedback := make(chan string)
	c_done := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(1)
	go listeners.ListenToFeedback(server, c_feedback, c_done, &wg)
	wg.Wait()
	// listen to feedback from ws
	var wg_res sync.WaitGroup
	wg_res.Add(1)
	go func() {
		select {
			case fb := <- c_feedback:
				msg = fb
				close(c_feedback)
				close(c_done)
				wg_res.Done()
			case <-time.After(10*time.Second):
				close(c_feedback)
				close(c_done)
				err := errors.New("Timeout: server did not respond")
				msg = appevents.ErrorFormated(err)
				wg_res.Done()
		}
	}()
	// send command
	sendCommand(command, server)
	wg_res.Wait()
	return command, nil, msg
}

func ServerExists(server_name string) (bool, error) {
	for _, s := range(metadata.GetServers()) {
		if server_name == s.Domain {
			return true, nil
		}
	}
	msg := "Server "+server_name+" not found: please check your config file"
	err := errors.New(msg)
	return false, err
}

func sendCommand(command *datatypes.Command, server *datatypes.Server) error {
	secret := server.WebsocketsKey
	host := server.WebsocketsHost
	port := server.WebsocketsPort
	purl := fmt.Sprintf("http://%s:%s", host, port)
	// connect to Centrifugo
	client := gocent.NewClient(purl, secret, 5*time.Second)
	ws_msg := encoding.MakeWsMsg(command)
	enc_msg, err := json.Marshal(ws_msg)
	channel := "$"+server.Domain+"_commands"
	if (Debug == true) {
		fmt.Println("Sending command", command.Name, "into channel", channel)
	}
	_, err = client.Publish(channel, enc_msg)
	if err != nil {
		return err
	}
	return nil
}

func getArgs(ctx *ishell.Context) []interface{} {
	var args []interface{}
	if len(ctx.Args) > 1 {
		for i, _ := range(args) {
			args = append(args, ctx.Args[i])
		}
	}
	return args
}