package mpmc

import (
	"concurrency/producerconsumer/consumer"
	"concurrency/producerconsumer/producer"
	"concurrency/producerconsumer/queue"
	"sync"
)

func MPMCSemaphore() {
	var wg sync.WaitGroup
	q := make(queue.Queue[int], 0)
	var m sync.Mutex
	done := false

	go func() {
		var wgp sync.WaitGroup

		for i := 0; i < producerNum; i++ {
			wg.Add(1)
			wgp.Add(1)
			id := i
			go func() {
				defer wg.Done()
				defer wgp.Done()
				producer.QueueBased(id, &q, bufferLen, &m)
			}()
		}

		wgp.Wait()

		done = true
	}()

	for i := 0; i < consumerNum; i++ {
		wg.Add(1)
		id := i
		go func() {
			defer wg.Done()
			consumer.QueueBased(id, &q, &m, &done)
		}()
	}

	wg.Wait()
}
