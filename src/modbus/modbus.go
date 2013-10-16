//package main
//author: Lubia Yang
//create: 2013-10-15
//about: www.lubia.me

package modbus

import (
	"errors"
	"os"
	"time"
)

// ModbusRead
//   parameters
//   int        fd:  file descripter for serial device
//   uint8_t  addr:  slave device address
//   uint16_t sr:    starting register number
//   uint16_t nr:    number of registers to read or write
//   uint8_t data[]: memory area for read data -- the size of footprint must be nr//2
func ModbusRead(fd *os.File, addr byte, sr, nr uint16) ([]byte, error) {
	//Preparation for Sending a Packet
	var send_packet = make([]byte, 8)

	//Packet Construction
	send_packet[0] = addr            // Slave Address
	send_packet[1] = 0x03            // Function Code 0x03 = Multiple Read
	send_packet[2] = byte(sr >> 8)   // Start Register (High Byte)
	send_packet[3] = byte(sr & 0xff) // Start Register (Low Byte)
	send_packet[4] = byte(nr >> 8)   // Number of Registers (High Byte)
	send_packet[5] = byte(nr & 0xff) // Number of Registers (Low Byte)

	//Add CRC16
	send_packet_crc := ModbusCrc(send_packet[:6])
	send_packet[6] = byte(send_packet_crc & 0xff)
	send_packet[7] = byte(send_packet_crc >> 8)

	// Preparation for Receiving a Packet
	var recv_packet = make([]byte, 256)
	_, err := fd.Write(send_packet)
	if err != nil {
		return []byte{}, errors.New("MODBUS_ERROR_COMMUNICATION")
	}
	time.Sleep(300 * time.Millisecond)
	_, err = fd.Read(recv_packet)
	if err != nil {
		return []byte{}, errors.New("MODBUS_ERROR_COMMUNICATION")
	}

	// Parse the Response
	if recv_packet[0] != send_packet[0] || recv_packet[1] != send_packet[1] {
		if recv_packet[0] == send_packet[0] && recv_packet[1]&0x7f == send_packet[1] {
			switch recv_packet[2] {
			case 1:
				return []byte{}, errors.New("MODBUS_ERROR_COMMUNICATION_ILLEGAL_FUNCTION")
			case 2:
				return []byte{}, errors.New("MODBUS_ERROR_COMMUNICATION_ILLEGAL_ADDRESS")
			case 3:
				return []byte{}, errors.New("MODBUS_ERROR_COMMUNICATION_ILLEGAL_VALUE")
			case 4:
				return []byte{}, errors.New("MODBUS_ERROR_COMMUNICATION_ILLEGAL_OPERATION")
			}
		}
		return []byte{}, errors.New("MODBUS_ERROR_COMMUNICATION")
	}

	//CRC check
	l := recv_packet[2]
	recv_packet_crc := ModbusCrc(recv_packet[:3+l])
	if recv_packet[3+l] != byte((recv_packet_crc&0xff)) || recv_packet[3+l+1] != byte((recv_packet_crc>>8)) {
		return []byte{}, errors.New("MODBUS_ERROR_COMMUNICATION")
	}
	return recv_packet[3 : l+3], nil
}

// ModbusWrite
//   parameters
//   int        fd:  file descripter for serial device
//   uint8_t  addr:  slave device address
//   uint16_t sr:    starting register number
//   uint16_t nr:    number of registers to write
//   uint8_t data[]: memory area for writing data -- the size of footprint must be nr//2
func ModbusWrite(fd *os.File, addr byte, sr, nr uint16, data []byte) error {
	var send_packet = make([]byte, 256)

	// Packet Construction
	send_packet[0] = addr            // Slave Address
	send_packet[1] = 0x10            // Function Code 0x10 = Multiple Write
	send_packet[2] = byte(sr >> 8)   // Start Register (High Byte)
	send_packet[3] = byte(sr & 0xff) // Start Register (Low Byte)
	send_packet[4] = byte(nr >> 8)   // Number of Registers (High Byte)
	send_packet[5] = byte(nr & 0xff) // Number of Registers (Low Byte)
	send_packet[6] = byte(nr * 2)

	for i := 0; i < int((nr * 2)); i++ {
		send_packet[7+i] = data[i]
	}

	length := 7 + nr*2 + 2
	// Add CRC16
	send_packet_crc := ModbusCrc(send_packet[:length-2])
	send_packet[length-2] = byte(send_packet_crc & 0xff)
	send_packet[length-1] = byte(send_packet_crc >> 8)

	// Preparation for Receiving a Packet
	var recv_packet = make([]byte, 256)
	_, err := fd.Write(send_packet)
	if err != nil {
		return errors.New("MODBUS_ERROR_COMMUNICATION")
	}
	time.Sleep(300 * time.Millisecond)
	_, err = fd.Read(recv_packet)
	if err != nil {
		return errors.New("MODBUS_ERROR_COMMUNICATION")
	}

	// Parse the Response
	if recv_packet[0] != send_packet[0] || recv_packet[1] != send_packet[1] {
		if recv_packet[0] == send_packet[0] && recv_packet[1]&0x7f == send_packet[1] {
			switch recv_packet[2] {
			case 1:
				return errors.New("MODBUS_ERROR_COMMUNICATION_ILLEGAL_FUNCTION")
			case 2:
				return errors.New("MODBUS_ERROR_COMMUNICATION_ILLEGAL_ADDRESS")
			case 3:
				return errors.New("MODBUS_ERROR_COMMUNICATION_ILLEGAL_VALUE")
			case 4:
				return errors.New("MODBUS_ERROR_COMMUNICATION_ILLEGAL_OPERATION")
			}
		}
		return errors.New("MODBUS_ERROR_COMMUNICATION")
	}

	//Target Data Filed Check
	if recv_packet[2] == send_packet[2] && recv_packet[3] == send_packet[3] && recv_packet[4] == send_packet[4] && recv_packet[5] == send_packet[5] {
		//CRC check
		recv_packet_crc := ModbusCrc(recv_packet[:6])
		if recv_packet[6] == byte((recv_packet_crc&0xff)) && recv_packet[7] == byte((recv_packet_crc>>8)) {
			return nil
		}
	}
	return errors.New("MODBUS_ERROR_COMMUNICATION")
}

func ModbusCrc(data []byte) uint16 {
	var crc16 uint16 = 0xffff
	l := len(data)
	for i := 0; i < l; i++ {
		crc16 ^= uint16(data[i])
		for j := 0; j < 8; j++ {
			if crc16&0x0001 > 0 {
				crc16 = (crc16 >> 1) ^ 0xA001
			} else {
				crc16 >>= 1
			}
		}
	}
	return crc16
}
