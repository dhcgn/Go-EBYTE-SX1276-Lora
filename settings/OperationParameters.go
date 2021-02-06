package settings

import (
	"encoding/hex"
	"fmt"
)

type Head int
type UartParityBit int
type TransmissionMode int
// IODriveMode This bit is used to the module internal pull-up
// resistor. It also increases the level's adaptability in
// case of open drain. But in some cases, it may need
// external pull-up resistor.
type IODriveMode int
type BaudRate int
// AirDataRate
// TODO  LoRa Spreading Factors?
type AirDataRate float64
type WirelessWakeUpTime int
type FEC int
type TransmissionPowerDb int
type Address [2]byte
type Channel byte

const (
	HeadSaveUndef Head = iota

	// Save the parameters when power down
	HeadSavePowerDown Head = iota
	// Don't save the parameters when power down
	HeadNoSafePowerDown Head = iota
)

func (value Head) String() string {
	return [...]string{"HeadSaveUndef", "HeadSavePowerDown", "HeadNoSafePowerDown"}[value]
}

const(
	UartParityBit_Undef UartParityBit = iota

	// Conventional notation 8-N-1
	UartParityBit_8N1 UartParityBit = iota
	UartParityBit_8OI UartParityBit = iota
	UartParityBit_8E1 UartParityBit = iota
)

func (value UartParityBit) String() string {
	return [...]string{"UartParityBit_Undef","UartParityBit_8NI", "UartParityBit_8OI", "UartParityBit_8E1", "UartParityBit_8N1"}[value]
}

const(
	BaudRateBps_Undef BaudRate = iota
	// TODO
	BaudRateBps_9600 BaudRate = 9600
)

func (value BaudRate) String() string {
	return fmt.Sprintf("BaudRate %v bps", int(value))
}

const(
	AirDataRate_Undef AirDataRate = iota

	// TODO
	AirDataRate_2_4k AirDataRate = 2.4
)

func (value AirDataRate) String() string {
	f := float64(value)
	return fmt.Sprintf("AirDataRate %.2f bps", f)
}

const(
	TransmissionModeUndef TransmissionMode = iota
	TransmissionModeTransparent TransmissionMode = iota
	TransmissionModeFixed TransmissionMode = iota
)

func (value TransmissionMode) String() string {
	return [...]string{"TransmissionMode Undef","Transparent TransmissionMode","Fixed TransmissionMode" }[value]
}

const(
	IODriveModeUndef IODriveMode = iota
	// IODriveModePushPull TXD and AUX push-pull outputs, RXD pull-up inputs (default value)
	IODriveModePushPull IODriveMode = iota
	// IODriveModeOpenCollector TXD, AUX open-collector outputs, RXD open-collector inputs
	IODriveModeOpenCollector IODriveMode = iota
)

func (value IODriveMode) String() string {
	return [...]string{"IODriveModeUndef","IODriveModePushPull","IODriveModeOpenCollector" }[value]
}

const(
	FEC_Undef FEC = iota

	FEC_Off FEC = iota
	// FEC_On (default value)
	FEC_On FEC = iota
)

func (value FEC) String() string {
	return [...]string{"FEC_Undef","FEC_Off","FEC_On" }[value]
}

const(
	WirelessWakeUpTime_Undef WirelessWakeUpTime = iota

	// WirelessWakeUpTime_250ms (default value)
	WirelessWakeUpTime_250ms WirelessWakeUpTime = 250
	WirelessWakeUpTime_500ms WirelessWakeUpTime = 500
	WirelessWakeUpTime_750ms WirelessWakeUpTime = 750
	WirelessWakeUpTime_1000ms WirelessWakeUpTime = 1000
	WirelessWakeUpTime_1250ms WirelessWakeUpTime = 1250
	WirelessWakeUpTime_1500ms WirelessWakeUpTime = 1500
	WirelessWakeUpTime_1750ms WirelessWakeUpTime = 1750
	WirelessWakeUpTime_2000ms WirelessWakeUpTime = 2000
)

func (value WirelessWakeUpTime) String() string {
	return fmt.Sprintf("WirelessWakeUpTime %vms", int(value))
}


const(
	TransmissionPowerDb_Undef = iota

	// TransmissionPowerDb_30db (default value)
	TransmissionPowerDb_30db TransmissionPowerDb = 30
	TransmissionPowerDb_27db TransmissionPowerDb = 27
	TransmissionPowerDb_24db TransmissionPowerDb = 24
	TransmissionPowerDb_21db TransmissionPowerDb = 21
)

func (value TransmissionPowerDb) String() string {
	return fmt.Sprintf("TransmissionPower %vdBm", int(value))
}

func (value Address) String() string {
	data := [2]byte(value)
	return fmt.Sprintf("Address 0x%v",hex.EncodeToString((data)[:]))
}

func (value Channel) String() string {
	data := byte(value)
	return fmt.Sprintf("Channel 0x%v",hex.EncodeToString([]byte{data}))
}

type OperationParameters struct {
	Head                 Head
	Address              Address
	UartParityBit        UartParityBit
	TtlUartBaudRateBps   BaudRate
	AirDataRateKbps      AirDataRate
	Channel              Channel
	TransmissionMode     TransmissionMode
	IODriveMode          IODriveMode
	WirelessWakeUpTimeMs WirelessWakeUpTime
	FEC                  FEC
	TransmissionPowerDb  TransmissionPowerDb
}

func NewOperationParametersFromData(data []byte) *OperationParameters {
	op := &OperationParameters{	}

	// Set Head
	if data[0] == 0xC0 {
		op.Head = HeadSavePowerDown
	}else if data[0] == 0xC2 {
		op.Head = HeadNoSafePowerDown
	}else {
		panic("Head not valid")
	}

	// Set Address
	op.Address = [2]byte{data[1], data[2]}


	// Set UartParityBit
	uart1 := hasBit(data[3], 7)
	uart2 := hasBit(data[3], 6)

	if uart1 && !uart2  {
		op.UartParityBit = UartParityBit_8OI
	}else if !uart1 && uart2  {
		op.UartParityBit = UartParityBit_8E1
	}else {
		op.UartParityBit = UartParityBit_8N1
	}

	// Set Baud Rate
	if !hasBit(data[3], 5) && hasBit(data[3], 4) && hasBit(data[3], 3) {
		op.TtlUartBaudRateBps = BaudRateBps_9600
	}else {
		panic("Baud Rate not (yet) supported")
	}

	// Set Air data rate
	if !hasBit(data[3], 2) && hasBit(data[3], 1) && !hasBit(data[3], 0) {
		op.AirDataRateKbps = AirDataRate_2_4k
	}else {
		panic("Air data rate(yet) supported ")
	}

	// Set Channel
	op.Channel = Channel(data[4])

	// Set TransmissionMode
	if hasBit(data[5], 7){
		op.TransmissionMode = TransmissionModeFixed
	}else {
		op.TransmissionMode = TransmissionModeTransparent
	}

	// Set IODriveMode
	if hasBit(data[5], 6){
		op.IODriveMode = IODriveModePushPull
	}else {
		op.IODriveMode = IODriveModeOpenCollector
	}

	// Set WirelessWakeUpTimeMs
	if !hasBit(data[5], 5) && !hasBit(data[5], 4) && !hasBit(data[5], 3) {
		op.WirelessWakeUpTimeMs = WirelessWakeUpTime_250ms
	}else {
		panic("WirelessWakeUpTimeMs not (yet) supported")
	}

	// Set IODriveMode
	if hasBit(data[5], 2){
		op.FEC = FEC_On
	}else {
		op.FEC = FEC_Off
	}

	// Set TransmissionPowerDb
	if !hasBit(data[5], 1) && !hasBit(data[5], 0)  {
		op.TransmissionPowerDb= TransmissionPowerDb_30db
	}else {
		panic("TransmissionPowerDb not (yet) supported")
	}

	return op
}

func hasBit(n byte, pos uint) bool {
	val := n & (1 << pos)
	return (val > 0)
}

