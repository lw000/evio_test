// evio_test project main.go
package main

import (
	"bufio"
	"bytes"
	protocol "demo/evio_test/pb"
	"demo/evio_test/pk"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/tidwall/evio"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type SessionContext struct {
	uuid uint64
	buf  *bytes.Buffer
	bufs []byte
	conn evio.Conn
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
	s := bytes.NewReader(u.bufs)
	br := bufio.NewReader(s)

	tdata, er := br.Peek(n)
	if er != nil {

	}

	d, er := pk.NewPacketWithData(tdata)
	if er != nil {
		log.Println(er)
		return nil, er
	}

	log.Println(d)

	return nil, nil
}

func (u *SessionContext) Parse(data []byte) ([]byte, error) {
	n, er := u.buf.Write(data)
	if er != nil {
		log.Println(er)
	}

	if n > 0 {

	}

	return nil, nil
}

var (
	scheme        = "tcp"
	port          = 9098
	uid    uint64 = 10000
	users  sync.Map
)

func main() {
	var events evio.Events
	events.NumLoops = -1
	events.Serving = func(server evio.Server) (action evio.Action) {
		log.Println("Serving")
		return
	}

	events.Opened = func(c evio.Conn) (out []byte, opts evio.Options, action evio.Action) {
		log.Println("Opened", c.RemoteAddr())

		uc := NewUContext()
		uc.SetUuid(atomic.AddUint64(&uid, 1))
		uc.SetConn(c)

		c.SetContext(uc)

		users.Store(uc.Uuid(), c)
		return
	}

	events.Data = func(c evio.Conn, in []byte) (out []byte, action evio.Action) {
		v := c.Context()
		uc := v.(*SessionContext)

		d, er := pk.NewPacketWithData(in)
		if er != nil {
			log.Println(er)
			return nil, 0
		}

		var req protocol.RequestChat
		er = proto.Unmarshal(d.Data(), &req)
		if er != nil {
			log.Println(er)
		}

		log.Println(uc.Uuid(), req.T, req.Msg)

		rep := protocol.ResponseChat{}
		rep.Msg = req.Msg
		rep.T = time.Now().UnixNano()
		dd := pk.NewPacket(d.Mid(), d.Sid(), d.RequestId())
		er = dd.EncodeProto(&rep)
		if er != nil {
			log.Println(er)
		}
		out = dd.Data()

		return
	}

	events.Closed = func(c evio.Conn, err error) (action evio.Action) {
		log.Println("Closed", c.RemoteAddr())
		return
	}

	//events.Tick = func() (delay time.Duration, action evio.Action) {
	//	delay = time.Second * time.Duration(1)
	//	log.Println("Tick", delay)
	//	return
	//}
	//
	//events.Detached = func(c evio.Conn, rwc io.ReadWriteCloser) (action evio.Action) {
	//	log.Println("Detached")
	//	return
	//}

	if err := evio.Serve(events, fmt.Sprintf("%s://:%d", scheme, port)); err != nil {
		log.Fatal(err)
	}
}
