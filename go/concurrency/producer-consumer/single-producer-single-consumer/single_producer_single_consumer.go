package single_producer_single_consumer

import (
	"fmt"
	"sync"
	"time"
)

const (
	bufferLen  = 5
	productNum = 20
)

func producer(buffer chan<- int) {
	for product := 0; product < productNum; product++ {
		fmt.Printf("Product received: %d\n", product)
		buffer <- product
		time.Sleep(1 * time.Second)
	}

	close(buffer)
}

func consumer(buffer <-chan int) {
	for product := range buffer {
		fmt.Printf("Product received: %d\n", product)
	}
}

func SingleProducerSingleConsumer() {
	var wg sync.WaitGroup
	buffer := make(chan int, bufferLen)

	wg.Add(1)
	go func() {
		defer wg.Done()
		producer(buffer)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		consumer(buffer)
	}()

	wg.Wait()
}
