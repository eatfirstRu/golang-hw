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
	// игнорировать ошибки в принципе.
	if m < 0 {
		m = len(tasks) * 100
	}

	if n > len(tasks) {
		n = len(tasks)
	}

	for iTask := 0; iTask < len(tasks); iTask++ {
		fmt.Printf("start iTask=%d\n", iTask)
		for i := 0; i < n; i++ {
			// fmt.Printf("start in gorutine n=%d iTask=%d\n", i, iTask)
			wg.Add(1)
			// fmt.Printf("ready to run iTask=%d,i=%d\n", iTask, i)
			go func(taskNum int) {
				if taskNum < len(tasks) {
					err := tasks[taskNum]()
					// fmt.Printf("ran iTask=%d\n", taskNum)
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
			// fmt.Printf("ErrErrorsLimitExceeded in end iTask=%d\n", iTask)
			return ErrErrorsLimitExceeded
		}
		iTask--
	}
	// fmt.Printf("received errCnt=%d errors limit=%d\n", errCnt, m)
	return nil
}
