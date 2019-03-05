package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	runLoop(4)
}

func runLoop(concurrency int) int {
	// Waitgroup forr all goroutines to finish
	var wg sync.WaitGroup
	// Holds any errors returned, protected by errorLock Mutex
	errors := make([]error, 0)
	var errorLock sync.Mutex
	// Buffered channel
	semaphore := make(chan int, concurrency)

	// Work loop
	for i := 0; i <= 20; i++ {
		// Start the goroutine that will do the work
		wg.Add(1)
		go func(l int) {
			defer wg.Done()
			semaphore <- 1
			fmt.Printf("running loop: %d\n", l)

			// Simulate some work
			time.Sleep(time.Duration(2) * time.Second)

			if err := fmt.Errorf("omg"); err != nil {
				// Add to errors
				errorLock.Lock()
				defer errorLock.Unlock()
				errors = append(errors, err)
			}
			// Read out of the channel to free up another goroutine
			<-semaphore
		}(i)
	}
	wg.Wait()

	if len(errors) > 0 {
		fmt.Printf("\n%d errors occurred:\n", len(errors))
		for _, err := range errors {
			fmt.Printf("--> %s\n", err)
		}
		return 1
	}
	return 0
}
