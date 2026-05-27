package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	errCnt := 0
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	if m < 0 {
		m = len(tasks) * 100
	}
	if n > len(tasks) {
		n = len(tasks)
	}
	for iTask := 0; iTask < len(tasks); iTask++ {
		fmt.Printf("start iTask=%d\n", iTask)
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func(taskNum int) {
				if taskNum < len(tasks) {
					err := tasks[taskNum]()
					if err != nil {
						mu.Lock()
						errCnt++
						mu.Unlock()
					}
				}
				wg.Done()
			}(iTask)
			mu.Lock()
			iTask++
			mu.Unlock()
		}
		wg.Wait()
		if errCnt >= m {
			return ErrErrorsLimitExceeded
		}
		iTask--
	}
	return nil
}
