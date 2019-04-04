package pk

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
)

type Packet struct {
	len       uint32
	mid       uint16
	sid       uint16
	requestId uint32
	data      []byte
}

func NewPacket(mid, sid uint16, requestId uint32) *Packet {
	p := &Packet{
		mid:       mid,
		sid:       sid,
		requestId: requestId,
	}
	return p
}

func NewPacketWithData(data []byte) (*Packet, error) {
	if len(data) == 0 {
		return nil, errors.New("data item is zero")
	}
	p := &Packet{}
	buf := bytes.NewBuffer(data)
	err := p.readHead(buf)
	if err != nil {
		return nil, err
	}
	p.data = buf.Bytes()
	return p, nil
}

func (p *Packet) writeHead(buf *bytes.Buffer) (err error) {
	if err = binary.Write(buf, binary.LittleEndian, p.len); err != nil {
		return err
	}
	if err = binary.Write(buf, binary.LittleEndian, p.mid); err != nil {
		return err
	}
	if err = binary.Write(buf, binary.LittleEndian, p.sid); err != nil {
		return err
	}
	if err = binary.Write(buf, binary.LittleEndian, p.requestId); err != nil {
		return err
	}
	return err
}

func (p *Packet) readHead(buf *bytes.Buffer) (err error) {
	if err = binary.Read(buf, binary.LittleEndian, &p.len); err != nil {
		return err
	}
	if err = binary.Read(buf, binary.LittleEndian, &p.mid); err != nil {
		return err
	}
	if err = binary.Read(buf, binary.LittleEndian, &p.sid); err != nil {
		return err
	}
	if err = binary.Read(buf, binary.LittleEndian, &p.requestId); err != nil {
		return err
	}
	return err
}

//Encode 编码数据包
func (p *Packet) encode(data []byte) error {
	buf := &bytes.Buffer{}
	p.len = 12 + uint32(len(data))
	err := p.writeHead(buf)
	if err != nil {
		return err
	}

	if len(data) > 0 {
		var n int
		n, err = buf.Write(data)
		if err != nil {
			return err
		}

		if n < 0 {

		}
	}

	p.data = buf.Bytes()

	return nil
}

//EncodeProto 编码数据包
func (p *Packet) EncodeProto(pb proto.Message) error {
	if pb == nil {
		err := p.encode([]byte{})
		if err != nil {
			return err
		}
	} else {
		data, err := proto.Marshal(pb)
		if err != nil {
			return err
		}
		err = p.encode(data)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p Packet) Mid() uint16 {
	return p.mid
}

func (p Packet) Sid() uint16 {
	return p.sid
}

func (p Packet) RequestId() uint32 {
	return p.requestId
}

func (p Packet) Data() []byte {
	return p.data
}

func (p Packet) String() string {
	return fmt.Sprintf("{mid:%d sid:%d requestId:%d datalen:%d}", p.mid, p.sid, p.requestId, len(p.data))
}
