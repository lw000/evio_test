// evio_test project main.go
package main

import (
	"bufio"
	"bytes"
	"demo/evio_test/packet"
	protocol "demo/evio_test/protos"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/judwhite/go-svc/svc"
	"github.com/tidwall/evio"
	"io"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

var (
	scheme        = "tcp"
	port          = 9098
	uid    uint64 = 10000
	users  sync.Map
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
		log.Println("Serving")
		return
	}

	pro.events.Opened = func(c evio.Conn) (out []byte, opts evio.Options, action evio.Action) {
		log.Println("Opened", c.RemoteAddr())

		uc := NewUContext()
		uc.SetUuid(atomic.AddUint64(&uid, 1))
		uc.SetConn(c)

		c.SetContext(uc)

		users.Store(uc.Uuid(), c)
		return
	}

	pro.events.Data = func(c evio.Conn, in []byte) (out []byte, action evio.Action) {
		ctx := c.Context()
		uc := ctx.(*SessionContext)

		pk, err := packet.NewPacketWithData(in)
		if err != nil {
			log.Println(err)
			return nil, 0
		}

		var req protocol.ReqChat
		err = proto.Unmarshal(pk.Data(), &req)
		if err != nil {
			log.Println(err)
			return nil, 0
		}

		log.Println(c.RemoteAddr().String(), uc.Uuid(), req.T, req.Msg)

		rep := protocol.AckChat{}
		rep.Msg = req.Msg
		rep.T = time.Now().UnixNano()
		dd := packet.NewPacket(pk.Mid(), pk.Sid(), pk.RequestId())
		err = dd.EncodeProto(&rep)
		if err != nil {
			log.Println(err)
			return nil, 0
		}
		out = dd.Data()

		return nil, 0
	}

	pro.events.Closed = func(c evio.Conn, err error) (action evio.Action) {
		log.Println("Closed", c.RemoteAddr())
		return
	}

	// events.Tick = func() (delay time.Duration, action evio.Action) {
	// 	delay = time.Second * time.Duration(1)
	// 	log.Println("Tick", delay)
	// 	return
	// }
	//
	pro.events.Detached = func(c evio.Conn, rwc io.ReadWriteCloser) (action evio.Action) {
		log.Println("Detached")
		return
	}

	go func() {
		if err := evio.Serve(pro.events, fmt.Sprintf("tcp://:%d", 9000), fmt.Sprintf("tcp://:%d", 9001)); err != nil {
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

type SessionContext struct {
	uuid   uint64
	buf    *bytes.Buffer
	buffer []byte
	conn   evio.Conn
}

func NewUContext() *SessionContext {
	return &SessionContext{buf: bytes.NewBuffer(nil)}
}

func (u *SessionContext) Conn() evio.Conn {
	return u.conn
}

func (u *SessionContext) SetConn(c evio.Conn) {
	u.conn = c
}

func (u *SessionContext) Uuid() uint64 {
	return u.uuid
}

func (u *SessionContext) SetUuid(uuid uint64) {
	u.uuid = uuid
}

func (u *SessionContext) Read(n int) ([]byte, error) {
	s := bytes.NewReader(u.buffer)
	br := bufio.NewReader(s)

	var (
		err  error
		data []byte
	)
	data, err = br.Peek(n)
	if err != nil {
		return nil, err
	}

	var pk *packet.Packet
	pk, err = packet.NewPacketWithData(data)
	if err != nil {
		return nil, err
	}

	log.Println(pk)

	return nil, nil
}

func (u *SessionContext) Parse(data []byte) ([]byte, error) {
	n, err := u.buf.Write(data)
	if err != nil {
		log.Println(err)
	}

	if n > 0 {

	}

	return nil, nil
}

func main() {
	pro := &Program{}
	if err := svc.Run(pro); err != nil {
		log.Fatal(err)
	}
}
