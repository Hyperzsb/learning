package consumer

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	MaxSleepDuration = 1
)

func ChannelBased(id int, buffer <-chan int) {
	for product := range buffer {
		fmt.Printf("Product %d received by consumer %d\n", product, id)
		time.Sleep(time.Duration(rand.Intn(MaxSleepDuration)) * time.Millisecond)
	}
}
