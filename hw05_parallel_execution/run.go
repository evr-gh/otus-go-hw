package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

var ErrProcExeWithErrors = errors.New("procs execution with errors")

var ErrNoExeThreadsAreGiven = errors.New("no exe threads are given")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if n < 1 {
		return ErrNoExeThreadsAreGiven
	}

	tasksCount := len(tasks)
	wg := sync.WaitGroup{}

	taskChan := make(chan Task, tasksCount)
	var chIn chan<- Task
	var chOut <-chan Task
	chIn = taskChan
	chOut = taskChan

	errorCount := int64(0)

	for i := 0; (i < n) && (i < tasksCount); i++ {
		wg.Go(func() {
			for task := range chOut {
				cnt := int64(0)
				if task != nil {
					err := task()
					if err != nil {
						cnt = atomic.AddInt64(&errorCount, 1)
						fmt.Println("Ошибка при выполнении задачи: ", err)
					}
				}
				if (m >= 0) && (cnt >= int64(m)) {
					break
				}
			}
		})
	}

	for _, task := range tasks {
		chIn <- task
	}

	close(chIn)

	wg.Wait()

	if (m >= 0) && (errorCount >= int64(m)) {
		return ErrErrorsLimitExceeded
	}

	if errorCount > 0 {
		return ErrProcExeWithErrors
	}
	return nil
}
