package webservice

import (
	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"strconv"
	"time"
)

type SocketClient struct {
	c         *gosocketio.Client
	url       string
	port      string
	isSecured bool
	mthMap    map[string]chan string
	onConnect func()
	done      chan bool
	rooms     []string
}

func NewSocketClient() *SocketClient {
	sc := &SocketClient{
		mthMap: make(map[string]chan string),
		rooms:  []string{},
	}
	return sc
}
func (sc *SocketClient) SetConnection(url string, port string, isSecured bool) *SocketClient {
	sc.url = url
	sc.port = port
	sc.isSecured = isSecured
	return sc
}

func (sc *SocketClient) SetMethod(method string, ch chan string) *SocketClient {
	sc.mthMap[method] = ch
	return sc
}

func (sc *SocketClient) SetRoom(room string) *SocketClient {
	sc.rooms = append(sc.rooms, room)

	return sc
}

func (sc *SocketClient) SetOnConnect(f func()) *SocketClient {
	sc.onConnect = f
	return sc
}

func (sc *SocketClient) SetOnDisconnect(done chan bool) *SocketClient {
	sc.done = done
	return sc
}

func (sc *SocketClient) Connect() error {
	p, err := strconv.Atoi(sc.port)
	if err != nil {
		return err
	}
	c, err := gosocketio.Dial(

		gosocketio.GetUrl(sc.url, p, sc.isSecured),
		transport.GetDefaultWebsocketTransport())

	if err != nil {

		return err
	}
	for key, value := range sc.mthMap {
		err := c.On("/"+key, func(h *gosocketio.Channel, arg string) {
			value <- arg
		})
		if err != nil {
			return err
		}
	}
	err = c.On(gosocketio.OnDisconnection, func(h *gosocketio.Channel) {
		sc.done <- true

	})

	if err != nil {
		return nil
	}

	err = c.On(gosocketio.OnConnection, func(h *gosocketio.Channel) {
		sc.onConnect()
	})

	if err != nil {
		return nil
	}

	for _, room := range sc.rooms {
		_, err := c.Ack("/join", room, time.Second*5)
		if err != nil {
			return err
		}
	}
	sc.c = c
	return nil
}

func (sc *SocketClient) Disconnect() {
	sc.c.Close()

}
