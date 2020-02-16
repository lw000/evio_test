// evio_test project main.go
package main

import (
	"demo/evio_test/packet"
	msg "demo/evio_test/protos"
	"demo/evio_test/service/session"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/judwhite/go-svc/svc"
	"github.com/tidwall/evio"
	"io"
	"log"
	"sync/atomic"
	"time"
)

var (
	uid uint64 = 10000
)

type Program struct {
	events evio.Events
}

func (pro *Program) Init(env svc.Environment) error {
	if env.IsWindowsService() {

	} else {

	}

	pro.events.NumLoops = -1

	return nil
}

// Start is called after Init. This method must be non-blocking.
func (pro *Program) Start() error {
	pro.events.Serving = func(server evio.Server) (action evio.Action) {
		for _, s := range server.Addrs {
			log.Printf("server listening %s", s.String())
		}
		return 0
	}

	pro.events.Opened = func(c evio.Conn) (out []byte, opts evio.Options, action evio.Action) {
		log.Println("Opened", c.RemoteAddr())

		u := session.New()
		u.Attach(c, atomic.AddUint64(&uid, 1))
		c.SetContext(u)
		session.Users.Store(u.ClientId(), c)
		return
	}

	pro.events.Data = func(c evio.Conn, in []byte) (out []byte, action evio.Action) {
		if in == nil {
			return
		}
		pk, err := packet.NewPacketWithData(in)
		if err != nil {
			log.Println(err)
			return
		}
		var req msg.ReqChat
		err = proto.Unmarshal(pk.Data(), &req)
		if err != nil {
			log.Println(err)
			return
		}

		ctx := c.Context()
		u, ok := ctx.(*session.Session)
		if !ok {
			return
		}
		log.Println(c.RemoteAddr().String(), u.ClientId(), req.T, req.Msg)

		ack := msg.AckChat{}
		ack.Msg = req.Msg
		ack.T = time.Now().UnixNano()
		dd := packet.NewPacket(pk.Mid(), pk.Sid(), pk.RequestId())
		err = dd.EncodeProto(&ack)
		if err != nil {
			log.Println(err)
			return
		}
		out = dd.Data()

		return
	}

	pro.events.Closed = func(c evio.Conn, err error) (action evio.Action) {
		log.Println("Closed", c.RemoteAddr())
		return
	}

	pro.events.Tick = func() (delay time.Duration, action evio.Action) {
		delay = time.Second * time.Duration(1)
		log.Println("Tick", delay)
		return
	}

	pro.events.Detached = func(c evio.Conn, rwc io.ReadWriteCloser) (action evio.Action) {
		log.Println("Detached")
		return
	}

	go func() {
		var address = []string{
			fmt.Sprintf("tcp://:%d", 9000),
			fmt.Sprintf("tcp://:%d", 9001),
		}
		if err := evio.Serve(pro.events, address...); err != nil {
			log.Println(err)
		}
	}()

	return nil
}

// Stop is called in response to syscall.SIGINT, syscall.SIGTERM, or when a
// Windows Service is stopped.
func (pro *Program) Stop() error {
	log.Println("Program quit")
	return nil
}

func main() {
	pro := &Program{}
	if err := svc.Run(pro); err != nil {
		log.Fatal(err)
	}
}
