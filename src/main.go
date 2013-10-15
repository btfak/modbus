package main

import (
	"log"
	"modbus"
	"os"
)

func main() {
	fd, err := os.Open("/dev/ttyAM0")
	if err != nil {
		log.Println("unable to open rs485")
		return
	}
	b, err := modbus.ModbusRead(fd, 1, 3, 1, []byte{0, 0})
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println(b)
	err = modbus.ModbusWrite(fd, 1, 3, 1, []byte{0, 1})
	if err != nil {
		log.Println(err.Error())
	}
}
