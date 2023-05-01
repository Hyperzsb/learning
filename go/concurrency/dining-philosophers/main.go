package main

import (
	"concurrency/diningphilosophers/philosopher"
	"sync"
)

const (
	philosopherNum = 5
	diningNum      = 3
)

func main() {
	philosophers := make([]philosopher.Philosopher, philosopherNum)
	for i := range philosophers {
		philosophers[i] = philosopher.NewPhilosopher(i)
	}

	wg := sync.WaitGroup{}
	diningMutex := sync.Mutex{}
	forkMutexes := make([]sync.Mutex, philosopherNum)

	for i := range philosophers {
		wg.Add(1)
		id := i
		go func() {
			defer wg.Done()
			for i := 0; i < diningNum; i++ {
				// Start to starve
				philosophers[id].Starving(i)
				// Acquire the lock to pick up forks
				// Only one philosopher is allowed to pick up forks at the same time
				diningMutex.Lock()
				// Try to pick up the fork on the left
				forkMutexes[id].Lock()
				// Try to pick up the fork on the right
				forkMutexes[(id+1)%philosopherNum].Lock()
				diningMutex.Unlock()
				// Start to eat
				philosophers[id].Eating(i)
				// Finish eating
				forkMutexes[id].Unlock()
				forkMutexes[(id+1)%philosopherNum].Unlock()
				// Start to think
				philosophers[id].Thinking(i)
			}
		}()
	}

	wg.Wait()
}
