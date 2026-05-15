package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

var ErrProcExeWithErrors = errors.New("procs execution with errors")

var ErrNoExeThreadsAreGiven = errors.New("no exe threads are given")

var ErrExeFailure = errors.New("error getting tasks exe results")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	// Place your code here.
	if n < 1 {
		return ErrNoExeThreadsAreGiven
	}

	wg := sync.WaitGroup{}

	resChan := make(chan error)
	var chIn chan<- error
	var chOut <-chan error
	chIn = resChan
	chOut = resChan

	tasksCount := len(tasks)
	finishedTasksCount := 0
	startedTasksCount := 0
	errTasksCount := 0

	for i := 0; (i < n) && (i < tasksCount); i++ {
		startedTasksCount++
		wg.Go(func() {
			chIn <- tasks[i]()
		})
	}

	for err := range chOut {
		finishedTasksCount++
		if err != nil {
			errTasksCount++
		}
		if (m >= 0) && (errTasksCount >= m) {
			if startedTasksCount == finishedTasksCount {
				wg.Wait()
				return ErrErrorsLimitExceeded
			}
		} else {
			switch {
			case (startedTasksCount == tasksCount) && (finishedTasksCount == tasksCount):
				wg.Wait()
				if errTasksCount > 0 {
					return ErrProcExeWithErrors
				}
				return nil
			case startedTasksCount < tasksCount:
				i := startedTasksCount
				startedTasksCount++

				wg.Go(func() {
					chIn <- tasks[i]()
				})
			}
		}
	}

	wg.Wait()

	if (startedTasksCount == finishedTasksCount) && (finishedTasksCount == tasksCount) {
		return nil
	}
	return ErrExeFailure
}
