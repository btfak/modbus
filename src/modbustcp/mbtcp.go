package modbustcp

import (
	"io/ioutil"
	"net"
	"strconv"
)

var _handler func(req []byte) (res []byte)

type Handler interface {
	Server(req []byte)(res []byte)
	
}

type MbTcp struct {
	Addr byte
	Code byte
	Data []byte
}

func (m MbTcp) generate() []byte {
	head := make([]byte, 8, 8)
	l := byte(len(m.Data) + 2)
	head[0] = 0x00
	head[1] = 0x00
	head[2] = 0x00
	head[3] = 0x00
	head[4] = 0x00
	head[5] = l
	head[6] = m.Addr
	head[7] = m.Code
	body := make([]byte, 260)
	body = append(body, head...)
	body = append(body, m.Data...)
	return body
}

func (m *MbTcp) Send(addr string) ([]byte, error) {
	req := m.generate()
	return send(addr, req)
}

func send(a string, d []byte) ([]byte, error) {
	addr, err := net.ResolveTCPAddr("tcp4", a)
	if err == nil {
		c, err := net.DialTCP("tcp", nil, addr)
		if err == nil {
			_, err = c.Write(d)
			if err == nil {
				r, err := ioutil.ReadAll(c)
				if err == nil {
					return r, nil
				}
			}
		}
	}
	return []byte{}, err
}

func SetHandler(h func(req []byte) (res []byte)) {
	_handler = h
}

func ServerCreate(port int) error {
	p := strconv.Itoa(port)
	ln, err := net.Listen("tcp", ":"+p)
	if err != nil {
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			continue
		}
		go handle(conn)
	}
}

func handle(c net.Conn) {
	req, err := ioutil.ReadAll(c)
	if err != nil {
		return
	}
	res := _handler(req)
	c.Write(res)
}
