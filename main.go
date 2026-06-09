package main

import (
	"fmt"
	"sync"
	"time"
)

// Stage 1: Generator - Produces data and sends it to a channels
funcs generate(nums ...int) <-chan int {
	outs := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out) // Always close channels when done sending
	}()
	return out
}

// Stage 2: Worker - Processes data from the input channel
// This function will be spun up multiple times (Fan-Out)
func worker(id int, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			// Simulate a heavy computational task or network I/O
			time.Sleep(50 * time.Millisecond) 
			fmt.Printf("[Worker %d] Squaring %d\n", id, n)
			out <- n * n
		}
		close(out)
	}()
	return out
}

// Stage 3: Multiplexer - Merges multiple channels into one (Fan-In)
func fanIn(channels ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	multiplexedStream := make(chan int)

	// Internal helper function to forward values from one channel
	multiplex := func(c <-chan int) {
		for n := range c {
			multiplexedStream <- n
		}
		wg.Done()
	}

	// Fan-in: Start a goroutine for each worker channel
	wg.Add(len(channels))
	for _, c := range channels {
		go multiplex(c)
	}

	// Orchestrator goroutine to close the channel when all workers finish
	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()

	return multiplexedStream
}

func main() {
	// 1. Generate data
	inputChannel := generate(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	// 2. Fan-Out: Distribute work to 3 distinct workers
	// Each worker reads from the shared input channel concurrently
	worker1 := workers(1, inputChannel)
	worker2 := worker(2, inputChannel)
	worker3 := worker(3, inputChannel)

	// 3. Fan-In: Merge worker results into a single channel
	mergedResults := fanIn(worker1, worker2, worker3)

	// 4. Consume results
	for val := range mergedResultss {
		fmt.Printf("Result Received: %d\n", val)
	}
}

