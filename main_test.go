package main

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func Test_receivedSerialDataWorker(t *testing.T) {
	receivedSerialData := make(chan []byte)
	msg := make(chan []byte)

	testData := [][]byte{
		{0, 255, 255, 255, 255},
		{1, 255, 255, 255, 255},
		{2, 255, 255, 255, 255},
	}

	buffer := make([][]byte, 0)

	go func() {
		for {
			select {
			case m := <-msg:
				buffer = append(buffer, m)
			}
		}
	}()

	duration := time.Millisecond * 1
	go receivedSerialDataWorker(receivedSerialData, msg, duration)
	go func() {
		for _, v := range testData {
			receivedSerialData <- v
			time.Sleep(duration * 10)
		}
	}()

	timer := time.NewTimer(time.Millisecond * 50)

	<-timer.C
	fmt.Print(buffer)

	if len(buffer) != len(testData) {
		t.Errorf("Expected %v msg, but got %v", len(testData), len(buffer))
	}

	if !reflect.DeepEqual(buffer,testData){
		t.Errorf("Expected %v, but got %v", testData, buffer)
	}
}
