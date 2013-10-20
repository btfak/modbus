package app

import (
	"log"
	"modbustcp"
)

var h *handler

type handler struct {
}

func (h *handler) Server(req []byte) []byte {
	return []byte{}
}

func (h *handler) Fault(detail string) {

}

func TcpServer() {
	modbustcp.SetHandler(h)
	err := modbustcp.ServerCreate(80)
	if err != nil {
		log.Println(err.Error())
	}
}
