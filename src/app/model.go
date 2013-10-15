package app

import (
	"log"
	"modbus"
	"os"
)

func Write() {
	
}

func Read() {
	fd, err := os.Open("/dev/ttyAM0")
	if err != nil {
		log.Println("unable to open rs485")
		return
	}
	b, err := modbus.ModbusRead(fd, 1, 3, 1, []byte{0, 0})
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(b)
}

func Crc16(){
	b := []byte{2,3,0,0,0,6}
	c := modbus.ModbusCrc(b)
	log.Println(c)
	
}