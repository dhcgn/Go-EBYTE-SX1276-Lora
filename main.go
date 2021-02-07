package main

import (
	"Go-EBYTE-SX1276-Lora/commands"
	"Go-EBYTE-SX1276-Lora/communication"
	"Go-EBYTE-SX1276-Lora/serial"
	"Go-EBYTE-SX1276-Lora/settings"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	comPath := func() string { return os.Args[1] }
	if len(os.Args) != 2 || !strings.HasPrefix(comPath(), "COM") {
		log.Fatal("COM Missing")
	}

	fmt.Println("Expect device on:", comPath())

	stream, err := serial.CreateSteam(comPath)
	if err != nil {
		log.Fatal("CreateSteam on "+comPath()+": ", err)
	}

	receivedSerialData := make(chan []byte)
	msg := make(chan []byte)

	go createMessagesFromRawData(receivedSerialData, msg, time.Millisecond*15)
	go handleMessages(msg)
	go listeningRawSerialData(stream, receivedSerialData)

	com := communication.New(stream)

	if err := com.RequestConfig(); err != nil {
		fmt.Println("Error", err)
	}
	time.Sleep(time.Millisecond * 30)

	if err := com.RequestVersion(); err != nil {
		fmt.Println("Error", err)
	}
	time.Sleep(time.Millisecond * 30)

	go sendTime(stream)
	// stream.Write(commands.ReadingOperatingParameters)
	// stream.Write(commands.MessagerHeader)

	fmt.Println("Waiting for interrupt")
	select {}
}

func handleMessages(msg <-chan []byte) {
	for {
		select {
		case m := <-msg:
			fmt.Println("Got msg:", "0x"+hex.EncodeToString(m))
			// Operation Parameters, eg 0xc000001a0644
			if len(m) == 6 && (m[0] == 0xC0 || m[0] == 0xC1) {
				fmt.Println("Msg is type Operation Parameters")
				op := settings.NewOperationParametersFromData(m)
				fmt.Println("Operation Parameters Hash:", op.GetShortHash())
				fmt.Println(op)
			}

			// Version, eg 0xc3450d14
			if len(m) == 4 && (m[0] == 0xC3) {
				fmt.Println("Msg is type Version")
			}
		}
	}
}

func createMessagesFromRawData(data <-chan []byte, msg chan<- []byte, waitDuration time.Duration) {
	// TODO Is timer needed, can i use Null-Byte from serial?
	timer := time.NewTimer(0)
	timer.Stop()

	buf := make([]byte, 0, 512)
	for {
		select {
		case received := <-data:
			// fmt.Println("received := <-data")
			timer = time.NewTimer(waitDuration)
			buf = append(buf, received...)
		case <-timer.C:
			// fmt.Println("<-timer.C")
			msg <- buf
			buf = make([]byte, 0, 512)
		}
	}
}

func sendTime(stream serial.SerialPortIO) {
	t := time.NewTicker(time.Second * 5)

	for {
		msg := append(commands.MessagerHeader, []byte(fmt.Sprintf("Time %v from %v", time.Now().Format("2006-01-02 15:04:05.999"), os.Getpid()))...)
		msg = append(msg, 0x0)
		sendData(stream, msg)

		time.Sleep(time.Second * 1)

		msg = append(commands.MessagerBroadcastHeader, []byte("Hello")...)
		sendData(stream, msg)

		<-t.C
	}
}

func sendData(stream serial.SerialPortIO, msg []byte) {
	fmt.Print("Send Time ... ")
	// fmt.Print(hex.EncodeToString(msg))
	n, err := stream.Write(msg)
	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println("send", n, "bytes")
}

func listeningRawSerialData(stream serial.SerialPortIO, data chan []byte) {
	fmt.Println("Listening loop started")

	buf := make([]byte, 128)
	for {
		// TODO Is this necessary? Can I just write and then read?
		n, err := stream.Read(buf)

		if err != nil {
			log.Fatal(err)
		}
		if n == 0 {
			continue
		}

		// s := string(buf[:n])

		// s := hex.EncodeToString(buf[:n])
		//fmt.Println(s)

		data <- buf[:n]
	}
}
