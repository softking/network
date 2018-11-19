package main

import (
	"time"
	"fmt"

	"github.com/softking/network"
)

var Client *network.SocketClient

func HandlerMessage(context network.Context)  {
	fmt.Printf("%+v", context)
}

func ConnDone(){
	fmt.Println("connected")
	context := network.Context{
		Method :12345,
		Msg: []byte("hello"),
	}

	Client.Send(context)

}

func ConnExit(){
	fmt.Println("disconnect")
}

func main() {
	Client = &network.SocketClient{
		Net:  "tcp4",
		Addr: "127.0.0.1:9999",
		HandlerBoxQueueSize: 8,
		ReadDeadLine:  60 * time.Second,
		WriteDeadLine: 60 * time.Second,

		Handler: HandlerMessage,
		WhenConnDone: ConnDone,
		WhenConnExit: ConnExit,
	}

	Client.Start()
}
