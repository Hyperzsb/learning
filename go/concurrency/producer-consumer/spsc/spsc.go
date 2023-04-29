package spsc

import (
	"concurrency/producerconsumer/consumer"
	"concurrency/producerconsumer/producer"
	"sync"
)

const (
	bufferLen  = 5
	productNum = 20
)

func SPSC() {
	var wg sync.WaitGroup
	buffer := make(chan int, bufferLen)

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(buffer)
		producer.ChannelBased(0, buffer)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		consumer.ChannelBased(0, buffer)
	}()

	wg.Wait()
}
