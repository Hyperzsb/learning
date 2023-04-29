package spmc

import (
	"concurrency/producerconsumer/consumer"
	"concurrency/producerconsumer/producer"
	"sync"
)

const (
	bufferLen   = 5
	consumerNum = 5
)

func SPMC() {
	var wg sync.WaitGroup
	buffer := make(chan int, bufferLen)

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(buffer)
		producer.ChannelBased(0, buffer)
	}()

	for i := 0; i < consumerNum; i++ {
		wg.Add(1)
		id := i
		go func() {
			defer wg.Done()
			consumer.ChannelBased(id, buffer)
		}()
	}

	wg.Wait()
}
