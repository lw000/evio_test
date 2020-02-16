package main

import (
	"demo/evio_test/packet"
	"demo/evio_test/protos"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type Client struct {
	conn      net.Conn
	connected bool
	done      chan struct{}
	onMessage func(data []byte)
}

func (c *Client) HandleMessage(fn func(data []byte)) {
	c.onMessage = fn
}

func (c Client) Connected() bool {
	return c.connected
}

func (c *Client) Open(address string) error {
	var err error
	c.conn, err = net.Dial("tcp", address)
	if err != nil {
		log.Println(err)
		return err
	}

	c.connected = true

	go c.run()

	return nil
}

func (c *Client) Send(data []byte) error {
	n, err := c.conn.Write(data)
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
	for i := 0; i <= 1; i++ {
		c := &Client{}
		err := c.Open(fmt.Sprintf("%s:%d", "127.0.0.1", 9000+i))
		if err != nil {
			log.Fatal(err)
		}

		c.HandleMessage(func(data []byte) {
			var pk *packet.Packet
			pk, err = packet.NewPacketWithData(data)
			if err != nil {
				log.Println(err)
				return
			}

			var rep msg.AckChat
			err = proto.Unmarshal(pk.Data(), &rep)
			if err != nil {
				log.Println(err)
			} else {
				log.Println(rep.T, rep.Msg)
			}
		})

		go func() {
			for {
				req := msg.ReqChat{}
				req.Uid = "1"
				req.Msg = strings.Repeat(fmt.Sprintf("%d", i), 10)
				req.T = time.Now().UnixNano()
				pk := packet.NewPacket(1, 1, 1)
				err = pk.EncodeProto(&req)
				if err != nil {
					log.Println(err)
					break
				}

				err = c.Send(pk.Data())
				if err != nil {
					log.Println(err)
					break
				}

				time.Sleep(time.Second * time.Duration(1))
			}
		}()

		go c.run()
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	sign := <-c
	log.Println("signal", sign)
	signal.Stop(c)
}
