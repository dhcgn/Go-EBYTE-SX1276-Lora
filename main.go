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

	//receivedSerialData := make(chan []byte)
	msg := make(chan []byte)

	// go createMessagesFromRawData(receivedSerialData, msg, time.Millisecond*15)
	go handleMessages(msg)
	go listeningRawSerialData(stream, msg)

	com := communication.New(stream)

	if err := com.RequestConfig(); err != nil {
		fmt.Println("Error", err)
	}
	time.Sleep(time.Millisecond * 500)

	if err := com.RequestVersion(); err != nil {
		fmt.Println("Error", err)
	}
	time.Sleep(time.Millisecond * 500)

	go sendTime(stream)

	fmt.Println("Waiting for interrupt")
	select {}
}

func handleMessages(msg <-chan []byte) {
	for {
		select {
		case m := <-msg:
			// Operation Parameters, eg 0xc000001a0644
			if len(m) == 6 && (m[0] == 0xC0 || m[0] == 0xC1) {
				fmt.Println("Msg is type Operation Parameters")
				op := settings.NewOperationParametersFromData(m)
				fmt.Println("Operation Parameters Hash:", op.GetShortHash())
				fmt.Println(op)
				break
			}

			// Version, eg 0xc3450d14
			if len(m) == 4 && (m[0] == 0xC3) {
				fmt.Println("Msg is type Version")
				break
			}

			// Get Msg
			toString := hex.EncodeToString(m)
			if len(m) > 3 {
				fmt.Println("Got msg:", "To:", toString[0:4], "Channel:", toString[4:6], "Msg:", string(m[3:]))
				break
			}

			fmt.Println("Got msg:", "0x"+toString)
		}
	}
}

func sendTime(stream serial.SerialPortIO) {
	t := time.NewTicker(time.Second * 5)

	for {
		msg := append(commands.MessagerHeader, []byte(fmt.Sprintf("Time %v from %v", time.Now().Format("2006-01-02 15:04:05.999"), os.Getpid()))...)
		msg = append(msg, 0x0)
		sendData(stream, msg)

		time.Sleep(time.Second * 2)

		msg = append(commands.MessagerBroadcastHeader, []byte(fmt.Sprintf("Hello to broadcast from %v", os.Getpid()))...)
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
	msg := make([]byte, 0, 512)
	for {
		n, err := stream.Read(buf)

		if err != nil {
			log.Fatal(err)
		}

		if n > 0 {
			msg = append(msg, buf[:n]...)
		}

		if n == 0 && len(msg) > 0 {
			data <- msg
			msg = make([]byte, 0, 512)
		}
	}
}
