package philosopher

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	MaxSleepDuration = 1000
)

type Philosopher struct {
	Id int
}

func NewPhilosopher(id int) Philosopher {
	return Philosopher{Id: id}
}

func (p Philosopher) Thinking(idx int) {
	fmt.Printf("Philosopher %d is thinking [%d]\n", p.Id, idx)
	time.Sleep(time.Duration(rand.Intn(MaxSleepDuration)) * time.Millisecond)
}

func (p Philosopher) Starving(idx int) {
	fmt.Printf("Philosopher %d is starving [%d]\n", p.Id, idx)
}

func (p Philosopher) Eating(idx int) {
	fmt.Printf("Philosopher %d is eating [%d]\n", p.Id, idx)
	time.Sleep(time.Duration(rand.Intn(MaxSleepDuration)) * time.Millisecond)
}
