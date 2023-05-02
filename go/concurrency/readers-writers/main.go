package main

import (
	"concurrency/readerswriters/reader"
	"concurrency/readerswriters/writer"
	"math/rand"
	"sync"
	"time"
)

const (
	writerNum        = 2
	readerNum        = 5
	writingNum       = 5
	maxSleepDuration = 1000
)

func withMutex(writers []writer.Writer, readers []reader.Reader) {
	wg := sync.WaitGroup{}
	writeMutex, readMutex := sync.Mutex{}, sync.Mutex{}
	readerCnt := 0
	done := false

	go func() {
		wgw := sync.WaitGroup{}
		for i := range writers {
			wg.Add(1)
			wgw.Add(1)
			id := i
			go func() {
				defer wg.Done()
				defer wgw.Done()
				for i := 0; i < writingNum; i++ {
					// The writer wants to write
					writeMutex.Lock()
					// Write
					writers[id].Write(i)
					// Finish writing
					writeMutex.Unlock()
					time.Sleep(time.Duration(rand.Intn(maxSleepDuration)) * time.Millisecond)
				}
			}()
		}

		wgw.Wait()
		done = true
	}()

	for i := range readers {
		wg.Add(1)
		id := i
		go func() {
			defer wg.Done()
			for !done {
				// The reader wants to read
				readMutex.Lock()
				// Increase the reader counter
				readerCnt++
				// If it is the first one
				if readerCnt == 1 {
					// Lock the buffer from being written
					writeMutex.Lock()
				}
				readMutex.Unlock()
				// Start to read
				readers[id].Read()
				// The reader finishes reading
				readMutex.Lock()
				// Decrease the reader counter
				readerCnt--
				// If it is the last one
				if readerCnt == 0 {
					// Allow the buffer to be written
					writeMutex.Unlock()
				}
				readMutex.Unlock()
				time.Sleep(time.Duration(rand.Intn(maxSleepDuration)) * time.Millisecond)
			}
		}()
	}

	wg.Wait()
}

func withRWMutex(writers []writer.Writer, readers []reader.Reader) {
	wg := sync.WaitGroup{}
	rwm := sync.RWMutex{}
	done := false

	go func() {
		wgw := sync.WaitGroup{}

		for i := range writers {
			wg.Add(1)
			wgw.Add(1)
			id := i
			go func() {
				defer wg.Done()
				defer wgw.Done()

				for i := 0; i < writingNum; i++ {
					rwm.Lock()
					writers[id].Write(i)
					rwm.Unlock()
					time.Sleep(time.Duration(rand.Intn(maxSleepDuration)) * time.Millisecond)
				}
			}()
		}

		wgw.Wait()
		done = true
	}()

	for i := range readers {
		wg.Add(1)
		id := i
		go func() {
			defer wg.Done()
			for !done {
				rwm.RLock()
				readers[id].Read()
				rwm.RUnlock()
				time.Sleep(time.Duration(rand.Intn(maxSleepDuration)) * time.Millisecond)
			}
		}()
	}

	wg.Wait()
}

func main() {
	buffer := 0

	writers := make([]writer.Writer, writerNum)
	for i := range writers {
		writers[i] = writer.NewWriter(i, &buffer)
	}

	readers := make([]reader.Reader, readerNum)
	for i := range readers {
		readers[i] = reader.NewReader(i, &buffer)
	}

	//withMutex(writers, readers)
	withRWMutex(writers, readers)
}
