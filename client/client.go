package main

import (
	"demo/evio_test/pb"
	"demo/evio_test/pk"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"log"
	"net"
	"strings"
	"time"
)

type Client struct {
	conn      net.Conn
	connected bool
	done      chan struct{}
	onMessage func(data []byte)
}

func (c *Client) HandleMessage(onMessage func(data []byte)) {
	c.onMessage = onMessage
}

func (c *Client) Connected() bool {
	return c.connected
}

func (c *Client) Open(address string) error {
	var er error
	c.conn, er = net.Dial("tcp", address)
	if er != nil {
		log.Panic(er)
	}
	c.connected = true

	go c.run()

	return nil
}

func (c *Client) Send(data []byte) error {
	var n int
	var err error
	n, err = c.conn.Write(data)
	if err != nil {
		log.Printf("connected closed")
	}

	if n > 0 {

	}

	return nil
}

func (c *Client) run() {
	var n int
	var err error
	buf := make([]byte, 1024)
	for {
		n, err = c.conn.Read(buf)
		if err != nil {
			log.Printf("connected closed")
			break
		}

		if n > 0 {
			if c.onMessage != nil {
				c.onMessage(buf[0:n])
			}
		}
	}
}

func (c *Client) Close() error {
	err := c.conn.Close()
	if err != nil {

	}
	return nil
}

func main() {
	for i := 0; i <= 0; i++ {
		c := &Client{}
		er := c.Open("127.0.0.1:9098")
		if er != nil {
			log.Panic(er)
		}

		c.HandleMessage(func(data []byte) {
			var d *pk.Packet
			d, er = pk.NewPacketWithData(data)
			if er != nil {
				log.Println(er)
				return
			}

			var rep protocol.ResponseChat
			er = proto.Unmarshal(d.Data(), &rep)
			if er != nil {
				log.Println(er)
			} else {
				log.Println(rep.T, rep.Msg)
			}
		})

		go func() {
			for {
				req := protocol.RequestChat{}
				req.Uid = "1"
				req.Msg = strings.Repeat(fmt.Sprintf("%d", i), 10)
				req.T = time.Now().UnixNano()
				d := pk.NewPacket(1, 1, 1)
				er = d.EncodeProto(&req)
				if er != nil {
					log.Println(er)
					continue
				}

				er = c.Send(d.Data())
				if er != nil {
					log.Println(er)
				}
				time.Sleep(time.Second * time.Duration(1))
			}
		}()

		go c.run()
	}

	select {}
}
