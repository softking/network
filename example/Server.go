package main

import (
	"fmt"
	"time"

	"github.com/softking/network"
)

func AcceptFun(client *network.Client) {
	fmt.Println("accept", client)
}

func TickFun(client *network.Client) {
	fmt.Println("tick", client)
}

func KickFun(client *network.Client) {
	fmt.Println("kick", client)
}

func HandlerFun(client *network.Client, context network.Context) {
	fmt.Println("handler", client)
	fmt.Printf("%+v", context)
}

func main() {
	server := network.SocketServer{

		Net:  "tcp4",
		Addr: ":9999",

		SenderBoxQueueSize:   2,
		ReceiverBoxQueueSize: 2,
		MainBoxQueueSize:     64,

		ReadDeadLine:  60 * time.Second,
		WriteDeadLine: 60 * time.Second,
		TickInterval:  5 * time.Second,
		Tick:          TickFun,
		Kick:          KickFun,
		Handler:       HandlerFun,
		Accept:        AcceptFun,
	}

	server.Start()
}
