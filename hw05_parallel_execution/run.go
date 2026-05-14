package hw05parallelexecution

import (
	"errors"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

var ErrNoExeThreadsAreGiven = errors.New("no exe threads are given")

var ErrExeFailure = errors.New("error getting tasks exe results")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	// Place your code here.
	if n < 1 {
		return ErrNoExeThreadsAreGiven
	}

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
		go func() {
			chIn <- tasks[i]()
		}()
	}

	for err := range chOut {
		finishedTasksCount++
		if err != nil {
			errTasksCount++
		}
		if errTasksCount >= m {
			if startedTasksCount == finishedTasksCount {
				return ErrErrorsLimitExceeded
			}
		} else {
			if (startedTasksCount == finishedTasksCount) && (finishedTasksCount == tasksCount) {
				return nil
			}
			if startedTasksCount < tasksCount {
				startedTasksCount++
				go func() {
					chIn <- tasks[startedTasksCount-1]()
				}()
			}
		}
	}

	if (startedTasksCount == finishedTasksCount) && (finishedTasksCount == tasksCount) {
		return nil
	}
	return ErrExeFailure
}
