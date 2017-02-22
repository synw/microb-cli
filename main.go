package main

import (
	"errors"
	"flag"
	"strconv"
	"github.com/abiosoft/ishell"
	"github.com/synw/microb/libmicrob/events/format"
	"github.com/synw/microb/libmicrob/datatypes"
	"github.com/synw/microb-cli/libmicrob/metadata"
	"github.com/synw/microb-cli/libmicrob/cmds"
	"github.com/synw/microb-cli/libmicrob/metrics/http"
)


var CurrentServer *datatypes.Server

var UseServerFlag = flag.String("s", "", "Use server: -s=server_name")
var Servers map[string]*datatypes.Server = metadata.GetServers()

var shell = ishell.New()

func main(){
	flag.Parse()
	
	// listen to feedback channel and print
	/*
	go func() {
        for msg := range(FeedbackChan) {
            ctx.Println(msg)
        }
    }()*/
	
	// INIT
	var msg string
	if *UseServerFlag != "" {
		_, err  := cmds.ServerExists(*UseServerFlag)
		if err != nil {
			shell.Println(format.ErrorFormated(err))
			return
		}
		CurrentServer= Servers[*UseServerFlag]
		msg = "Using server "+CurrentServer.Domain
	} else {
		msg = "Available servers: "
		for _, server := range Servers {
			msg = msg+server.Domain+" "
		}
	}
    shell.Println(msg)
    
    // USE
    shell.AddCmd(&ishell.Cmd{
		Name: "use",
		Help: "Use server: use server_domain",
		Func: func(ctx *ishell.Context) {
			var err error
			if len(ctx.Args) == 0 {
				err = errors.New("missing server domain")
				return
			}
			if len(ctx.Args) > 1 {
				err = errors.New("please use only one server at the time")
				return
			}
			server_name := ctx.Args[0]
			_, err  = cmds.ServerExists(server_name)
			if err != nil {
				ctx.Println(format.ErrorFormated(err))
				return
			}
			CurrentServer= Servers[server_name]
			ctx.Println("Using server", server_name)
		},
	})
	
	// USING
	shell.AddCmd(&ishell.Cmd{
        Name: "using",
        Help: "Check what server is in use",
        Func: func(ctx *ishell.Context) {
        	if CurrentServer != nil {
            	ctx.Println("Using server", CurrentServer.Domain)
            } else {
            	err := errors.New("No server in use. To select a server type: use server_name")
            	ctx.Println(format.ErrorFormated(err))
            }
        },
    })
    
    // PING
    shell.AddCmd(&ishell.Cmd{
        Name: "ping",
        Help: "Ping the current server",
        Func: func(ctx *ishell.Context) {
        	_, err, msg := cmds.SendCmd(ctx, "ping", CurrentServer)
        	if err != nil {
        		ctx.Println(format.ErrorFormated(err))
        	} else {
        		ctx.Println(msg)
        	}
        },
    })
    
    // TIME
    shell.AddCmd(&ishell.Cmd{
        Name: "time",
        Help: "Time for a request to process: time /mypath/ or time http://myurl.com",
        Func: func(ctx *ishell.Context) {
        	url := "/"
        	if len(ctx.Args) > 0 {
        		url = ctx.Args[0]
        	}
        	metric, err := http_metrics.GetRequestMetric(url, CurrentServer)
        	if err != nil {
        		ctx.Println(format.ErrorFormated(err))
        	} else {
        		msg := metric.Format()
        		ctx.Println(msg)
        	}
        },
    })
    
    // GET
    shell.AddCmd(&ishell.Cmd{
        Name: "get",
        Help: "Get the content of an url: get /mypath/ or time http://myurl.com",
        Func: func(ctx *ishell.Context) {
        	url := "/"
        	if len(ctx.Args) > 0 {
        		url = ctx.Args[0]
        	}
        	httpResp, err := http_metrics.Get(url, CurrentServer)
        	if err != nil {
        		ctx.Println(format.ErrorFormated(err))
        	} else {
        		msg := httpResp.Format()
        		ctx.Println(msg)
        	}
        },
    })
    
     // STRESS
    shell.AddCmd(&ishell.Cmd{
        Name: "stress",
        Help: "Stress the server by sending requests: stress 500 5 50 (interval, number of workers, number of requests)",
        Func: func(ctx *ishell.Context) {
        	interval := 500
        	workers := 5
        	limit := 50
        	var err error
        	if len(ctx.Args) > 0 {
        		if len(ctx.Args) != 3 {
        			err = errors.New("Please provide 3 arguments: interval, number of workers and max number of requests")
        		}
        		if err != nil {
	        		ctx.Println(format.ErrorFormated(err))
	        	}
        		interval, err = strconv.Atoi(ctx.Args[0])
        		if err != nil {
	        		ctx.Println(format.ErrorFormated(err))
	        	}
	        	workers, err = strconv.Atoi(ctx.Args[1])
        		if err != nil {
	        		ctx.Println(format.ErrorFormated(err))
	        	}
	        	limit, err = strconv.Atoi(ctx.Args[2])
        		if err != nil {
	        		ctx.Println(format.ErrorFormated(err))
	        	}
        	}
        	report, err := http_metrics.Stress(CurrentServer, interval, workers, limit)
        	if err != nil {
        		ctx.Println(format.ErrorFormated(err))
        	} else {
        		msg := report.Format()
        		ctx.Println(msg)
        	}
        },
    })
    
    // ROUTES
    shell.AddCmd(&ishell.Cmd{
        Name: "routes",
        Help: "Get the registered routes",
        Func: func(ctx *ishell.Context) {
        	var msg string
        	_, err, msg := cmds.SendCmd(ctx, "routes", CurrentServer)
        	if err != nil {
        		ctx.Println(format.ErrorFormated(err))
        	} else {
        		msg = "Registered routes for server "+CurrentServer.Domain+":\n"+msg
        		ctx.Println(msg)
        	}
        },
    })
    
    /*
    // SYS
    shell.AddCmd(&ishell.Cmd{
        Name: "sys",
        Help: "Get some system information from the current server",
        Func: func(ctx *ishell.Context) {
        	command, err, msg := cmds.SendCmd(ctx, "sys", CurrentServer)
        	if err != nil {
        		ctx.Println(format.ErrorFormated(err))
        	} else {
        		ctx.Println(msg)
        	}
        },
    })
    */
    /*
    // PINGALL
    shell.AddCmd(&ishell.Cmd{
        Name: "pingall",
        Help: "Ping all the servers",
        Func: func(ctx *ishell.Context) {
        	for _, server := range(metadata.GetServers()) {
        		go func() {
        			_, err, msg := cmds.SendCmd(ctx, "ping", server)
	        		if err != nil {
		        		ctx.Println(format.ErrorFormated(err))
		        	} else {
		        		ctx.Println(msg)
		        	}
        		}()
	        }
        },
    })
    */
    // REPARSE_TEMPLATES
    shell.AddCmd(&ishell.Cmd{
        Name: "reparse_templates",
        Help: "Reparse templates",
        Func: func(ctx *ishell.Context) {
        	//events.New("command", "cli", "reparse_templates")
        	_, err, msg := cmds.SendCmd(ctx, "reparse_templates", CurrentServer)
        	if err != nil {
        		ctx.Println(format.ErrorFormated(err))
        	} else {
        		ctx.Println(msg)
        	}
        },
    })
    
    // UPDATE ROUTES
    shell.AddCmd(&ishell.Cmd{
        Name: "update_routes",
        Help: "Update client side routes",
        Func: func(ctx *ishell.Context) {
        	_, err, msg := cmds.SendCmd(ctx, "update_routes", CurrentServer)
        	if err != nil {
        		ctx.Println(format.ErrorFormated(err))
        	} else {
        		ctx.Println(msg)
        	}
        },
    })
    
    // STATE
    shell.AddCmd(&ishell.Cmd{
        Name: "state",
        Help: "Print server state information",
        Func: func(ctx *ishell.Context) {
        	_, err, msg := cmds.SendCmd(ctx, "state", CurrentServer)
        	if err != nil {
        		ctx.Println(format.ErrorFormated(err))
        	} else {
        		ctx.Println(msg)
        	}
        },
    })
    
    // DB STATUS
    shell.AddCmd(&ishell.Cmd{
        Name: "db_status",
        Help: "Reports main database status",
        Func: func(ctx *ishell.Context) {
        	_, err, msg := cmds.SendCmd(ctx, "db_status", CurrentServer)
        	if err != nil {
        		ctx.Println(format.ErrorFormated(err))
        	} else {
        		ctx.Println(msg)
        	}
        },
    })
    
    // DATABASES
    shell.AddCmd(&ishell.Cmd{
        Name: "databases",
        Help: "Display info about the available databases",
        Func: func(ctx *ishell.Context) {
        	_, err, msg := cmds.SendCmd(ctx, "databases", CurrentServer)
        	if err != nil {
        		ctx.Println(format.ErrorFormated(err))
        	} else {
        		ctx.Println(msg)
        	}
        },
    })
     
    shell.Start()
}
