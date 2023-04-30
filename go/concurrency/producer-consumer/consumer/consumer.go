package consumer

import (
	"concurrency/producerconsumer/queue"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var (
	MaxSleepDuration = 1000
)

func ChannelBased(id int, buffer <-chan int) {
	for product := range buffer {
		fmt.Printf("Product %d consumed by consumer %d\n", product, id)
		time.Sleep(time.Duration(rand.Intn(MaxSleepDuration)) * time.Millisecond)
	}
}

func QueueBased(id int, q *queue.Queue[int], m *sync.Mutex, done *bool) {
	for !*done {
		m.Lock()
		if !q.Empty() {
			product := q.Front()
			q.Pop()
			fmt.Printf("Product %d consumed by consumer %d, current queue %v\n", product, id, *q)
		}
		m.Unlock()
		time.Sleep(time.Duration(rand.Intn(MaxSleepDuration)) * time.Millisecond)
	}
}
