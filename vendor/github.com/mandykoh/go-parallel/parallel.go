package parallel

import (
	"sync"
)

// RunWorkers executes the specified worker function using n goroutines, passing
// each a workerNum from 0-n and a workerCount of n. This function returns after
// all workers have run to completion.
func RunWorkers(n int, worker func(workerNum, workerCount int)) {
	allDone := sync.WaitGroup{}
	allDone.Add(n)

	for workerNum := 0; workerNum < n; workerNum++ {
		go func(workerNum int) {
			defer allDone.Done()
			worker(workerNum, n)
		}(workerNum)
	}

	allDone.Wait()
}
