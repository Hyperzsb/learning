package main

import (
	"concurrency/producerconsumer/queue"
	"concurrency/sleepingbarber/barber"
	"concurrency/sleepingbarber/customer"
	"math/rand"
	"sync"
	"time"
)

const (
	customerNum      = 20
	maxCustomerNum   = 5
	maxSleepDuration = 2000
)

func main() {
	wg := sync.WaitGroup{}

	cq := make(queue.Queue[customer.Customer], 0)
	qm := sync.Mutex{}
	done := false

	wg.Add(1)
	b := barber.NewBarber(0)
	go func() {
		defer wg.Done()
		for !done {
			if b.IsSleeping {
				continue
			}

			qm.Lock()
			if cq.Empty() {
				// Sleep while there is no customer
				b.Sleeping()
				qm.Unlock()
			} else {
				// Get a customer from the queue
				c := cq.Front()
				cq.Pop()
				qm.Unlock()
				b.Cutting(c.Id)
				time.Sleep(time.Duration(rand.Intn(maxSleepDuration)) * time.Millisecond)
			}
		}
	}()

	go func() {
		wgc := sync.WaitGroup{}
		for i := 0; i < customerNum; i++ {
			wg.Add(1)
			wgc.Add(1)
			id := i
			go func() {
				defer wg.Done()
				defer wgc.Done()

				c := customer.NewCustomer(id)
				c.Coming()

				qm.Lock()
				if cq.Length() == maxCustomerNum {
					// If the queue is full, leave
					c.Leaving()
				} else {
					c.Waiting()
					// If the barber is sleeping, wake him up
					if b.IsSleeping {
						b.WakenUp(c.Id)
					}
					cq.Push(c)
				}
				qm.Unlock()
			}()
			time.Sleep(time.Duration(rand.Intn(maxSleepDuration)) * time.Millisecond)
		}

		wgc.Wait()

		for {
			qm.Lock()
			if cq.Empty() {
				done = true
			}
			qm.Unlock()
		}
	}()

	wg.Wait()
}
