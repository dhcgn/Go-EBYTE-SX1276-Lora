package serial

import (
	"Go-EBYTE-SX1276-Lora/global"
	"github.com/tarm/serial"
)

type SerialPortIO interface {
	Write(buf []byte) (int, error)
	Read(buf []byte) (int, error)
}

func CreateSteam(comPath func() string) (SerialPortIO, error) {
	config := &serial.Config{
		Name:        comPath(),
		Baud:        global.Baud,
		ReadTimeout: global.SerialReadTimeout,
		Size:        8,
	}

	stream, err := serial.OpenPort(config)
	return stream, err
}
