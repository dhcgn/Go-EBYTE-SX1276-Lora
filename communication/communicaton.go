package communication

import (
	"Go-EBYTE-SX1276-Lora/commands"
	"Go-EBYTE-SX1276-Lora/serial"
)

type com struct {
	serial serial.SerialPortIO
}

func (c *com) RequestConfig() error {
	_, err := c.serial.Write(commands.ReadingOperatingParameters)
	return err
}

func (c *com) RequestVersion() error {
	_, err := c.serial.Write(commands.ReadingVersionNumber)
	return err
}

type Com interface {
	SendBroadcast(data []byte)
	RequestConfig() error
	RequestVersion() error
}

func New(serial serial.SerialPortIO) Com {
	return &com{
		serial: serial,
	}
}

func (c *com) SendBroadcast(data []byte) {

}
