package listeners


import (
	//"log"
	"fmt"
	"time"
	"errors"
	"sync"
	"github.com/centrifugal/centrifuge-go"
	"github.com/centrifugal/centrifugo/libcentrifugo/auth"
	"github.com/synw/microb/libmicrob/datatypes"
	appevents "github.com/synw/microb/libmicrob/events"
	"github.com/synw/microb/libmicrob/datatypes/encoding"
	"github.com/synw/microb-cli/libmicrob/conf"
	"github.com/synw/microb-cli/libmicrob/metadata"
)


var Config = conf.GetCliConf()
var Debug = metadata.IsDebug()

func credentials(server *datatypes.Server) *centrifuge.Credentials {
	secret := server.WebsocketsKey
	user := "microbcli"
	timestamp := centrifuge.Timestamp()
	info := ""
	token := auth.GenerateClientToken(secret, user, timestamp, info)
	return &centrifuge.Credentials{
		User:      user,
		Timestamp: timestamp,
		Info:      info,
		Token:     token,
	}
}

func listenForFeedback(channel_name string, feedback chan string, server *datatypes.Server, wg *sync.WaitGroup) (centrifuge.Centrifuge, *centrifuge.SubEventHandler) {
	creds := credentials(server)
	wsURL := "ws://"+server.WebsocketsHost+":"+server.WebsocketsPort+"/connection/websocket"
	if (Debug == true) {
		msg := "Connecting to "+wsURL
		fmt.Println(msg)
	}
	
	onMessage := func(sub centrifuge.Sub, rawmsg centrifuge.Message) error {
		if (Debug == true) {
			msg := fmt.Sprintf("New message received in channel %s: %#v", sub.Channel(), rawmsg)
			fmt.Println(msg)
		}
		payload, err := encoding.DecodeJsonFeedbackRawMessage(rawmsg.Data)
		var msg string
		if err != nil {
			msg = "Error decoding json raw message: "+err.Error()
			feedback <- msg
		}
		event_class := payload.EventClass
		data := payload.Data
		cmd_status := payload.Status
		cmd_error := payload.Error
		cmd_name := data["command"].(string)
		cmd_id := data["id"].(string)
		cmd_from := data["from"].(string)
		cmd_reason := data["reason"].(string)
		var rvs []interface{}
		if data["return_values"] != nil {
			rvs = data["return_values"].([]interface{})
		}
		if (event_class == "command_feedback") {
			now := time.Now()
			var args []interface{}
			err_ := errors.New(cmd_error)
			command := &datatypes.Command{cmd_id, cmd_name, cmd_from, cmd_reason, now, args, cmd_status, err_, rvs}
			msg := appevents.GetCommandReportMsg(command)
			feedback <- msg
		}
		return nil
	}
	
	onJoin := func(sub centrifuge.Sub, msg centrifuge.ClientInfo) error {
		if (Debug == true) {
			fmt.Println(fmt.Sprintf("User %s joined channel %s with client ID %s", msg.User, sub.Channel(), msg.Client))
		}
		return nil
	}

	onLeave := func(sub centrifuge.Sub, msg centrifuge.ClientInfo) error {
		if (Debug == true) {
			fmt.Println(fmt.Sprintf("User %s with clientID left channel %s with client ID %s", msg.User, msg.Client, sub.Channel()))
		}
		return nil
	}
	
	onPrivateSub := func(c centrifuge.Centrifuge, req *centrifuge.PrivateRequest) (*centrifuge.PrivateSign, error) {
		info := ""
		sign := auth.GenerateChannelSign(server.WebsocketsKey, req.ClientID, req.Channel, info)
		privateSign := &centrifuge.PrivateSign{Sign: sign, Info: info}
		return privateSign, nil
	}

	events := &centrifuge.EventHandler{
		OnPrivateSub: onPrivateSub,
	}
	
	subevents := &centrifuge.SubEventHandler{
		OnMessage: onMessage,
		OnJoin:    onJoin,
		OnLeave:   onLeave,
	}
	conn := centrifuge.NewCentrifuge(wsURL, creds, events, centrifuge.DefaultConfig)
	return conn, subevents
}

func ListenToFeedback(server *datatypes.Server, feedback chan string, done chan bool, wg *sync.WaitGroup) {
	channel_name := "$"+server.Domain+"_feedback"
	// connect to channel
	conn, subevents := listenForFeedback(channel_name, feedback, server, wg)
	defer conn.Close()
	err := conn.Connect()
	if err != nil {
		fmt.Println("Error listening to feedback:", err.Error())
	}
	// suscribe to websockets channel
	_, err = conn.Subscribe(channel_name, subevents)
	if err != nil {
		fmt.Println("Error suscribing to feedback channel:", err.Error())
	}
	wg.Done()
	// sit here and wait
	<- done
}
