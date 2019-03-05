package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	runLoop(4, 20)
}

func runLoop(concurrency int, items int) int {
	fmt.Printf("running jobs: %d, concurrency: %d\n", items, concurrency)

	// Waitgroup for all goroutines to finish
	var wg sync.WaitGroup

	// Holds any errors returned, protected by errLock Mutex
	errs := make([]error, 0)
	var errLock sync.Mutex

	// Buffered channel
	semaphore := make(chan int, concurrency)

	// Work loop
	for i := 0; i < items; i++ {

		// Start the goroutines that will do the work
		wg.Add(1)
		go func(loop int) {
			defer wg.Done()
			semaphore <- 1
			fmt.Printf("running loop: %d\n", loop)

			// Simulate some work
			time.Sleep(time.Duration(2) * time.Second)

			if err := fmt.Errorf("error on loop: %d", loop); err != nil {
				// Add to errors
				errLock.Lock()
				defer errLock.Unlock()
				errs = append(errs, err)
			}

			// Read out of the channel to free up another goroutine
			<-semaphore
		}(i)
	}
	wg.Wait()

	if len(errs) > 0 {
		// Report, or deal with errs
		fmt.Printf("\n%d errors occurred:\n", len(errs))
		for _, err := range errs {
			fmt.Printf("%s\n", err)
		}
		return 1
	}
	return 0
}
