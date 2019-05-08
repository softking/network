package network

import (
	"net"
	"sync"
	"time"
)

type Context struct {
	Cmd    string
	Method uint32
	Msg    interface{}
	From   interface{}
	Data   map[string]interface{}
}

type Link chan byte

func (link *Link) close() {
	defer func() { recover() }()
	close(*link)
}

type TickFunc func(*Client)
type KickFunc func(*Client)
type AcceptFunc func(*Client)
type HandlerFunc func(*Client, Context)

type ClientHandlerFunc func(Context)

type SocketServer struct {
	Net  string
	Addr string

	SenderBoxQueueSize   int
	ReceiverBoxQueueSize int
	MainBoxQueueSize     int

	ReadDeadLine  time.Duration
	WriteDeadLine time.Duration
	TickInterval  time.Duration

	Tick    TickFunc
	Kick    KickFunc
	Accept  AcceptFunc
	Handler HandlerFunc
}

type SocketClient struct {
	Net  string
	Addr string
	conn *net.TCPConn

	HandlerBoxQueueSize int

	ReadDeadLine  time.Duration
	WriteDeadLine time.Duration

	Handler ClientHandlerFunc

	WhenConnDone func()
	WhenConnExit func()
}

type Client struct {
	Uid int32

	Conn        *net.TCPConn
	SenderBox   chan Context
	ReceiverBox chan Context
	MainBox     chan Context

	Data map[string]interface{}

	TickChan <-chan time.Time
}

var (
	online map[string]*Client
	lock   sync.Mutex
)

func init() {
	online = make(map[string]*Client)
}

func RegesiterOnline(uid string, client *Client) {
	lock.Lock()
	defer lock.Unlock()
	online[uid] = client
}

func UnRegesiterOnline(uid string) {
	lock.Lock()
	defer lock.Unlock()
	delete(online, uid)
}

func QueryOnline(uid string) (*Client, bool) {
	lock.Lock()
	defer lock.Unlock()
	data, ok := online[uid]
	return data, ok
}

func ListAll() []string {
	lock.Lock()
	defer lock.Unlock()

	list := make([]string, len(online))
	idx := 0
	for k := range online {
		list[idx] = k
		idx++
	}
	return list
}
