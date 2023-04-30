package producer

import (
	"concurrency/producerconsumer/queue"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var (
	ProductNum       = 10
	MaxSleepDuration = 1000
)

func ChannelBased(id int, buffer chan<- int) {
	for product := 0; product < ProductNum; product++ {
		buffer <- product
		fmt.Printf("Product %d produced by producer %d\n", product, id)
		time.Sleep(time.Duration(rand.Intn(MaxSleepDuration)) * time.Millisecond)
	}
}

func QueueBased(id int, q *queue.Queue[int], bufferLen int, m *sync.Mutex) {
	for product := 0; product < ProductNum; {
		m.Lock()
		if q.Length() < bufferLen {
			q.Push(product)
			fmt.Printf("Product %d produced by producer %d, current queue %v \n", product, id, *q)
			product++
		}
		m.Unlock()
		time.Sleep(time.Duration(rand.Intn(MaxSleepDuration)) * time.Millisecond)
	}
}
