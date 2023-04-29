package producer

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	ProductNum       = 10
	MaxSleepDuration = 1000
)

func ChannelBased(id int, buffer chan<- int) {
	for product := 0; product < ProductNum; product++ {
		fmt.Printf("Product %d produced by producer %d\n", product, id)
		buffer <- product
		time.Sleep(time.Duration(rand.Intn(MaxSleepDuration)) * time.Millisecond)
	}
}
