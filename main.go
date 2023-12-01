package main

import (
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
		go worker(stopChan, stopMsg.i)
		fmt.Printf("worker id=%d stop with ", stopMsg.i)
		if stopMsg.err != nil {
			fmt.Printf("error is %v\n", stopMsg.err)
			continue
		}
		fmt.Println("no reason")
	}
}

func worker(stopChan chan<- stopMsg, i int) {
	defer func() {
		var err error
		if r := recover(); r != nil {
			err = r.(error)
		}
		stopWithReason(stopChan, i, err)
	}()

	time.Sleep(time.Duration(i) * time.Second)
}

func stopWithReason(stopChan chan<- stopMsg, i int, reason error) {
	stopChan <- stopMsg{i: i, err: reason}
}
