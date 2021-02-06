package main

import (
	"Go-EBYTE-SX1276-Lora/commands"
	"Go-EBYTE-SX1276-Lora/settings"
	"encoding/hex"
	"fmt"
	"github.com/tarm/serial"
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

	config := &serial.Config{
		Name:        comPath(),
		Baud:        9600,
		ReadTimeout: time.Millisecond * 10,
		Size:        8,
	}

	stream, err := serial.OpenPort(config)
	if err != nil {
		log.Fatal(err)
	}

	receivedSerialData := make(chan []byte)
	msg := make(chan []byte)

	go receivedSerialDataWorker(receivedSerialData, msg, time.Millisecond*15)
	go receivedMsgWorker(msg)
	go listeningLoop(stream, receivedSerialData)

	stream.Write(commands.ReadingOperatingParameters)
	time.Sleep(time.Millisecond * 30)

	stream.Write(commands.ReadingOperatingParameters)
	time.Sleep(time.Millisecond * 30)

	go sendTime(stream)
	// stream.Write(commands.ReadingOperatingParameters)
	// stream.Write(commands.MessagerHeader)

	fmt.Println("Waiting for interrupt")
	select {}
}

func receivedMsgWorker(msg <-chan []byte) {
	for {
		select {
		case m := <-msg:
			fmt.Println("Got msg:", "0x"+hex.EncodeToString(m))
			// Operation Parameters
			if len(m) == 6 && (m[0] == 0xC0 || m[0] == 0xC1) {
				fmt.Println("Msg is type Operation Parameters")
				op := settings.NewOperationParametersFromData(m)
				fmt.Println(op)
			}
		}
	}
}



func receivedSerialDataWorker(data <-chan []byte, msg chan<- []byte, waitDuration time.Duration) {
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

func sendTime(stream *serial.Port) {
	t := time.NewTicker(time.Second * 5)

	for {
		msg := append(commands.MessagerHeader, []byte(fmt.Sprintf("Time %v from %v", time.Now().Format("2006-01-02 15:04:05.999"), os.Getpid()))...)
		sendData(stream, msg)

		time.Sleep(time.Second * 1)

		msg = append(commands.MessagerBroadcastHeader, []byte("Hello")...)
		sendData(stream, msg)

		<-t.C
	}
}

func sendData(stream *serial.Port, msg []byte) {
	fmt.Print("Send Time ... ")
	// fmt.Print(hex.EncodeToString(msg))
	n, err := stream.Write(msg)
	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println("send", n, "bytes")
}

func listeningLoop(stream *serial.Port, data chan []byte) {
	fmt.Println("Listening loop started")

	buf := make([]byte, 128)
	for {
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
