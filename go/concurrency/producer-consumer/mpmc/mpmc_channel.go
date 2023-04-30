package mpmc

import (
	"concurrency/producerconsumer/consumer"
	"concurrency/producerconsumer/producer"
	"sync"
)

const (
	bufferLen   = 5
	producerNum = 5
	consumerNum = 5
)

func MPMCChannel() {
	var wg sync.WaitGroup

	buffer := make(chan int, bufferLen)

	go func() {
		var wgp sync.WaitGroup

		for i := 0; i < producerNum; i++ {
			wg.Add(1)
			wgp.Add(1)
			id := i
			go func() {
				defer wg.Done()
				defer wgp.Done()
				producer.ChannelBased(id, buffer)
			}()
		}

		wgp.Wait()
		close(buffer)
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
