package network

import (
	"encoding/binary"
	"io"
	"log"
	"net"
	"time"
)

func (nw *SocketServer) Start() {

	nw.check()

	tcpAddr, err := net.ResolveTCPAddr(nw.Net, nw.Addr)
	if err != nil {
		panic(err)
	}
	listener, err := net.ListenTCP(nw.Net, tcpAddr)
	if err != nil {
		panic(err)
	}

	log.Println("Server Start Success", nw.Net, nw.Addr)
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Println("accept err: ", err)
			continue
		}

		nw.handleClient(conn)
	}
}

func (nw *SocketServer) handleClient(connect *net.TCPConn) {
	client := &Client{
		Uid:  -1,
		Conn: connect,

		SenderBox:   make(chan Context, nw.SenderBoxQueueSize),
		ReceiverBox: make(chan Context, nw.ReceiverBoxQueueSize),
		MainBox:     make(chan Context, nw.MainBoxQueueSize),
		Data: 		 make(map[string]interface{}),
	}

	if nw.TickInterval > 0 {
		client.TickChan = time.After(nw.TickInterval)
	} else {
		client.TickChan = make(<-chan time.Time) // 没有空间的chan 用于阻塞
	}

	Accept(client, nw)

	link := make(Link)
	go HandlerSender(client, nw.WriteDeadLine, link)
	go HandlerReceiver(client, nw.ReadDeadLine, 10240, link)
	go MainLoop(client, nw, link)
}

func HandlerSender(client *Client, writeDeadLine time.Duration, link Link) {
	defer link.close()

	conn := client.Conn
	var context Context

	for {
		select {
		case _, ok := <-link:
			if !ok {
				return
			}
		case context = <-client.SenderBox:
			switch context.Cmd {
			case "send":
				if writeDeadLine > 0 {
					conn.SetWriteDeadline(time.Now().Add(writeDeadLine))
				}

				data := context.Msg.([]byte)
				byteArray := CreateByteArray()
				byteArray.WriteU32(uint32(len(data)) + 4) // data + method
				byteArray.WriteU32(context.Method)
				byteArray.WriteRawBytes(data)

				n, err := conn.Write(byteArray.Data())
				if err != nil {
					log.Println("send data err: ", err, n)
					return
				}
			}
		}
	}
}

func HandlerReceiver(client *Client, readDeadLine time.Duration, maxLength uint32, link Link) {
	defer link.close()

	conn := client.Conn

	for {
		select {
		case _, ok := <-link:
			if !ok {
				return
			}
		default:
		}

		if readDeadLine != 0 {
			conn.SetReadDeadline(time.Now().Add(readDeadLine))
		}

		header := make([]byte, 4)
		method := make([]byte, 4)
		_, err1 := io.ReadFull(conn, header)
		_, err2 := io.ReadFull(conn, method)
		if err1 != nil || err2 != nil {
			return
		}

		size := binary.BigEndian.Uint32(header) - 4 // 减去method
		if size > maxLength {
			return
		}
		body := make([]byte, size)
		_, err := io.ReadFull(conn, body)
		if err != nil {
			return
		}

		client.MainBox <- Context{
			Cmd:    "message",
			Method: binary.BigEndian.Uint32(method),
			Data:   make(map[string]interface{}),
			Msg:    body,
		}
	}
}

func MainLoop(client *Client, nw *SocketServer, link Link) {
	defer func() { client.Conn.Close() }()
	defer link.close() //异常退出时候关闭所有进程
	defer Kick(client, nw)

	for {
		select {
		case _, ok := <-link:
			if !ok {
				return
			}
		case context := <-client.MainBox:

			switch context.Cmd {
			case "message":
				nw.Handler(client, context)

			case "kick":
				return
			}

		case <-client.TickChan:
			Tick(client, nw)

		}
	}
}

func Accept(client *Client, nw *SocketServer) {
	defer CatchException()
	RegesiterOnline(client.Conn.RemoteAddr().String(), client)
	nw.Accept(client)
}

func Tick(client *Client, nw *SocketServer) {
	defer CatchException()
	defer func() { client.TickChan = time.After(nw.TickInterval) }()
	nw.Tick(client)

}

func Kick(client *Client, nw *SocketServer) {
	defer CatchException()
	UnRegesiterOnline(client.Conn.RemoteAddr().String())
	nw.Kick(client)
}

func (nw *SocketServer) check() {
	if nw.Net == "" {
		panic("Net empty")
	}
	if nw.Addr == "" {
		panic("Addr empty")
	}
	if nw.SenderBoxQueueSize <= 0 {
		panic("SenderBoxQueueSize error")
	}
	if nw.ReceiverBoxQueueSize <= 0 {
		panic("ReceiverBoxQueueSize error")
	}
	if nw.MainBoxQueueSize <= 0 {
		panic("MainBoxQueueSize error")
	}
	if nw.ReadDeadLine <= 0 {
		panic("ReadDeadLine error")
	}
	if nw.WriteDeadLine <= 0 {
		panic("WriteDeadLine error")
	}
	if nw.TickInterval < 0 {
		panic("TickInterval error")
	}
	if nw.Tick == nil {
		nw.TickInterval = 0
		nw.Tick = func(client *Client) {}
	}
	if nw.Kick == nil {
		nw.Kick = func(client *Client) {}
	}
	if nw.Handler == nil {
		panic("Can not find Handler Func")
	}
	if nw.Accept == nil {
		nw.Accept = func(client *Client) {}
	}
}
