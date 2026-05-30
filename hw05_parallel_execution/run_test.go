package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, int32(tasksCount), runTasksCount, "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})
}

func tc1(t *testing.T) {
	tasksCount := int32(50)
	tasks := make([]Task, 0, tasksCount)

	var m sync.Map

	var runTasksCount int32

	for i := int32(0); i < tasksCount; i++ {
		err := fmt.Errorf("error from task %d", i)
		taskSleep := time.Millisecond * time.Duration(rand.Intn(100)+1)
		tasks = append(tasks, func() error {
			require.Eventually(t, func() bool { return true }, taskSleep*10, taskSleep)
			atomic.AddInt32(&runTasksCount, 1)
			m.Store(i, err)
			return err
		})
	}

	workersCount := 10
	maxErrorsCount := 23
	err := Run(tasks, workersCount, maxErrorsCount)

	require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
	require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")

	for i := int32(0); i < runTasksCount; i++ {
		err, ok := m.Load(i)
		require.Truef(t, ok, "i=%v", i)
		expErr := fmt.Errorf("error from task %d", i)
		require.Equal(t, err, expErr)
	}
}

func tc2(t *testing.T) {
	tasksCount := 50
	tasks := make([]Task, 0, tasksCount)

	var m sync.Map

	var runTasksCount int32
	var sumTime time.Duration

	for i := 0; i < tasksCount; i++ {
		taskSleep := time.Millisecond * time.Duration(rand.Intn(100)+1)
		sumTime += taskSleep

		tasks = append(tasks, func() error {
			time.Sleep(taskSleep)

			require.Eventually(t, func() bool { return true }, taskSleep*10, taskSleep)

			atomic.AddInt32(&runTasksCount, 1)
			m.Store(i, "")
			return nil
		})
	}

	workersCount := 5
	maxErrorsCount := 1

	start := time.Now()
	err := Run(tasks, workersCount, maxErrorsCount)
	elapsedTime := time.Since(start)
	require.NoError(t, err)

	require.Equal(t, int32(tasksCount), runTasksCount, "not all tasks were completed")
	require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")

	for i := 0; i < tasksCount; i++ {
		err, ok := m.Load(i)
		require.True(t, ok)
		require.Empty(t, err)
	}
}

func tc3(t *testing.T) {
	tasksCount := int32(50)
	tasks := make([]Task, 0, tasksCount)

	var m sync.Map

	var runTasksCount int32

	for i := int32(0); i < tasksCount; i++ {
		var err error
		taskSleep := time.Millisecond * time.Duration(rand.Intn(100)+1)
		if i%2 == 0 {
			err = fmt.Errorf("error from task %d", i)
			taskSleep = time.Millisecond * time.Duration(rand.Intn(100)+100)
		}
		tasks = append(tasks, func() error {
			require.Eventually(t, func() bool { return true }, taskSleep*10, taskSleep)
			atomic.AddInt32(&runTasksCount, 1)
			m.Store(i, err)
			return err
		})
	}

	workersCount := 10
	maxErrorsCount := 26
	err := Run(tasks, workersCount, maxErrorsCount)

	require.Equal(t, tasksCount, runTasksCount, "not all tasks were completed")
	require.Truef(t, errors.Is(err, ErrProcExeWithErrors), "actual err - %v", err)

	for i := int32(0); i < tasksCount; i++ {
		err, ok := m.Load(i)
		require.Truef(t, ok, "i=%v", i)
		if i%2 == 0 {
			expErr := fmt.Errorf("error from task %d", i)
			require.Equal(t, err, expErr)
		} else {
			require.Empty(t, err)
		}
	}
}

func tc4(t *testing.T) {
	tasksCount := int32(50)
	tasks := make([]Task, 0, tasksCount)

	var m sync.Map

	var runTasksCount int32

	for i := int32(0); i < tasksCount; i++ {
		var err error
		taskSleep := time.Millisecond * time.Duration(rand.Intn(100)+1)
		if i%2 == 0 {
			err = fmt.Errorf("error from task %d", i)
		}
		tasks = append(tasks, func() error {
			require.Eventually(t, func() bool { return true }, taskSleep*10, taskSleep)
			atomic.AddInt32(&runTasksCount, 1)
			m.Store(i, err)
			return err
		})
	}

	workersCount := 10
	maxErrorsCount := -1
	err := Run(tasks, workersCount, maxErrorsCount)

	require.Equal(t, tasksCount, runTasksCount, "not all tasks were completed")
	require.Truef(t, errors.Is(err, ErrProcExeWithErrors), "actual err - %v", err)

	for i := int32(0); i < tasksCount; i++ {
		err, ok := m.Load(i)
		require.Truef(t, ok, "i=%v", i)
		if i%2 == 0 {
			expErr := fmt.Errorf("error from task %d", i)
			require.Equal(t, err, expErr)
		} else {
			require.Empty(t, err)
		}
	}
}

func tc5(t *testing.T) {
	tasksCount := int32(50)
	tasks := make([]Task, 0, tasksCount)

	var m sync.Map

	var runTasksCount int32

	for i := int32(0); i < tasksCount; i++ {
		var err error
		taskSleep := time.Millisecond * time.Duration(rand.Intn(100)+1)
		if i%2 == 0 {
			err = fmt.Errorf("error from task %d", i)
		}
		tasks = append(tasks, func() error {
			require.Eventually(t, func() bool { return true }, taskSleep*10, taskSleep)
			atomic.AddInt32(&runTasksCount, 1)
			m.Store(i, err)
			return err
		})
	}

	workersCount := 10
	maxErrorsCount := 25
	err := Run(tasks, workersCount, maxErrorsCount)

	require.Equal(t, tasksCount, runTasksCount, "not all tasks were completed")
	require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)

	for i := int32(0); i < runTasksCount; i++ {
		err, ok := m.Load(i)
		require.Truef(t, ok, "i=%v", i)
		if i%2 == 0 {
			expErr := fmt.Errorf("error from task %d", i)
			require.Equal(t, err, expErr)
		} else {
			require.Empty(t, err)
		}
	}
}

func TestMoreRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", tc1)

	t.Run("tasks without errors", tc2)

	t.Run("half tasks with errors, errors limit didn't exceed", tc3)

	t.Run("m<0: errors are ignored", tc4)

	t.Run("half tasks with errors, errors limit exceeded", tc5)
}
