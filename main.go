package main

import (
	"errors"
	"fmt"
	"time"
)

const (
	workerSize = 5
)

type stopMsg struct {
	i   int
	err error
}

func main() {
	// channel
	stopChan := make(chan stopMsg, workerSize)
	defer close(stopChan)

	// start worker
	for i := 1; i <= workerSize; i++ {
		go worker(stopChan, i)
	}

	// watchdog
	for stopMsg := range stopChan {
		fmt.Printf("stop[%d]: err=%v\n", stopMsg.i, stopMsg.err)

		go worker(stopChan, stopMsg.i)
	}
}

var c int

func worker(stopChan chan<- stopMsg, i int) {
	defer func() {
		var err error
		if r := recover(); r != nil {
			err = r.(error)
			stopWithReason(stopChan, i, err)
			return
		}
		stop(stopChan, i)
	}()

	if c++; c%10 == 0 {
		panic(errors.New("is mod 10"))
	}
	time.Sleep(time.Duration(i) * time.Second)
}

func stopWithReason(stopChan chan<- stopMsg, i int, reason error) {
	stopChan <- stopMsg{i: i, err: reason}
}

func stop(stopChan chan<- stopMsg, i int) {
	stopChan <- stopMsg{i: i}
}
