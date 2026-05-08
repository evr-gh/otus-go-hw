package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)

		c.Clear()

		val, ok = c.Get("aaa")
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get("bbb")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("simple purge logic", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("a", 1)
		require.False(t, wasInCache)

		wasInCache = c.Set("b", 2)
		require.False(t, wasInCache)

		wasInCache = c.Set("c", 3)
		require.False(t, wasInCache)

		wasInCache = c.Set("d", 4)
		require.False(t, wasInCache)

		wasInCache = c.Set("e", 5)
		require.False(t, wasInCache)

		wasInCache = c.Set("f", 6)
		require.False(t, wasInCache)

		val, ok := c.Get("a")
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get("b")
		require.True(t, ok)
		require.Equal(t, 2, val)

		val, ok = c.Get("c")
		require.True(t, ok)
		require.Equal(t, 3, val)

		val, ok = c.Get("d")
		require.True(t, ok)
		require.Equal(t, 4, val)

		val, ok = c.Get("e")
		require.True(t, ok)
		require.Equal(t, 5, val)

		val, ok = c.Get("f")
		require.True(t, ok)
		require.Equal(t, 6, val)
	})

	t.Run("comlex purge logic", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("a", 1)
		require.False(t, wasInCache)

		wasInCache = c.Set("b", 2)
		require.False(t, wasInCache)

		wasInCache = c.Set("c", 3)
		require.False(t, wasInCache)

		wasInCache = c.Set("d", 4)
		require.False(t, wasInCache)

		wasInCache = c.Set("e", 5)
		require.False(t, wasInCache)

		val, ok := c.Get("b")
		require.True(t, ok)
		require.Equal(t, 2, val)

		wasInCache = c.Set("f", 6)
		require.False(t, wasInCache)

		wasInCache = c.Set("a", 11)
		require.False(t, wasInCache)

		wasInCache = c.Set("e", 15)
		require.True(t, wasInCache)

		val, ok = c.Get("a")
		require.True(t, ok)
		require.Equal(t, 11, val)

		val, ok = c.Get("b")
		require.True(t, ok)
		require.Equal(t, 2, val)

		val, ok = c.Get("c")
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get("d")
		require.True(t, ok)
		require.Equal(t, 4, val)

		val, ok = c.Get("e")
		require.True(t, ok)
		require.Equal(t, 15, val)

		val, ok = c.Get("f")
		require.True(t, ok)
		require.Equal(t, 6, val)
	})
}

func TestCacheMultithreading(t *testing.T) {
	t.Run("test 1", func(t *testing.T) {
		c := NewCache(10)
		wg := &sync.WaitGroup{}
		wg.Add(2)

		go func() {
			defer wg.Done()
			for i := range 1_000_000 {
				c.Set(Key(strconv.Itoa(i)), i)
			}
		}()

		go func() {
			defer wg.Done()
			for range 1_000_000 {
				c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
			}
		}()

		wg.Wait()

		for i := 999999; i >= 999990; i-- {
			val, ok := c.Get(Key(strconv.Itoa(i)))
			require.Equal(t, i, val)
			require.True(t, ok)
		}
	})

	t.Run("test 2", func(t *testing.T) {
		c := NewCache(10)
		wg := &sync.WaitGroup{}
		wg.Add(2)

		go func() {
			defer wg.Done()
			for i := range 1_000_000 {
				if (i % 20) == 0 {
					c.Clear()
				}
				c.Set(Key(strconv.Itoa(i)), i)
			}
		}()

		go func() {
			defer wg.Done()
			for i := range 1_000_000 {
				c.Get(Key(strconv.Itoa(rand.Intn(999990 + (i % 10)))))
			}
		}()

		wg.Wait()

		for i := 999999; i >= 999990; i-- {
			val, ok := c.Get(Key(strconv.Itoa(i)))
			require.Equal(t, i, val)
			require.True(t, ok)
		}
	})

	t.Run("test 3", func(t *testing.T) {
		c := NewCache(10)
		wg := &sync.WaitGroup{}
		wg.Add(2)

		c.Set(Key(strconv.Itoa(100)), 100)

		go func() {
			defer wg.Done()
			for i := range 1_000_000 {
				c.Set(Key(strconv.Itoa(i%5)), i%5)
			}
		}()

		go func() {
			defer wg.Done()
			for i := range 1_000_000 {
				c.Set(Key(strconv.Itoa(5+i%5)), 5+i%5)
			}
		}()

		wg.Wait()

		for i := range 10 {
			val, ok := c.Get(Key(strconv.Itoa(i)))
			require.Equal(t, i, val)
			require.True(t, ok)
		}

		c.Clear()

		for i := range 10 {
			val, ok := c.Get(Key(strconv.Itoa(i)))
			require.Nil(t, val)
			require.False(t, ok)
		}
	})
}
