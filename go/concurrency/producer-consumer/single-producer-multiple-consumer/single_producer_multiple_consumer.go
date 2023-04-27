package single_producer_multiple_consumer

import (
	"fmt"
	"sync"
	"time"
)

const (
	bufferLen   = 5
	productNum  = 20
	consumerNum = 5
)

func producer(buffer chan<- int) {
	for product := 0; product < productNum; product++ {
		fmt.Printf("Product received: %d\n", product)
		buffer <- product
		time.Sleep(1 * time.Second)
	}

	close(buffer)
}

func consumer(id int, buffer <-chan int) {
	for product := range buffer {
		fmt.Printf("Product received by consumer %d: %d\n", id, product)
	}
}

func SingleProducerMultipleConsumer() {
	var wg sync.WaitGroup
	buffer := make(chan int, bufferLen)

	wg.Add(1)
	go func() {
		defer wg.Done()
		producer(buffer)
	}()

	for i := 0; i < consumerNum; i++ {
		wg.Add(1)
		id := i
		go func() {
			defer wg.Done()
			consumer(id, buffer)
		}()
	}

	wg.Wait()
}
