package session

import (
	"bufio"
	"bytes"
	"demo/evio_test/packet"
	"github.com/tidwall/evio"
	"log"
)

type Session struct {
	clientId uint64
	buf      *bytes.Buffer
	buffer   []byte
	conn     evio.Conn
}

func New() *Session {
	return &Session{buf: bytes.NewBuffer(nil)}
}

func (u *Session) Attach(c evio.Conn, clientId uint64) {
	u.conn = c
	u.clientId = clientId
}

func (u *Session) ClientId() uint64 {
	return u.clientId
}

func (u *Session) Read(n int) ([]byte, error) {
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

func (u *Session) Parse(data []byte) ([]byte, error) {
	n, err := u.buf.Write(data)
	if err != nil {
		log.Println(err)
	}

	if n > 0 {

	}

	return nil, nil
}
