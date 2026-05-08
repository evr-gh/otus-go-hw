package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]

		l.MoveToFront(l.Back()) // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})
}

func TestAddList(t *testing.T) {
	t.Run("additional 1", func(t *testing.T) {
		l := NewList()

		l.PushFront(1) // [1]
		l.PushBack(2)  // [1, 2]
		l.PushFront(2) // [2, 1, 2]
		l.PushBack(3)  // [2, 1, 2, 3]
		l.PushFront(3) // [3, 2, 1, 2, 3]

		require.Equal(t, 5, l.Len())

		middle := l.Front().Next.Next // 1
		l.Remove(middle)              // [3, 2, 2, 3]
		require.Equal(t, 4, l.Len())

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			v, ok := i.Value.(int)
			if ok {
				elems = append(elems, v)
			} else {
				require.True(t, ok, "Error list value type")
			}
		}
		require.Equal(t, []int{3, 2, 2, 3}, elems)

		for _, v := range [...]int{4, 5, 6} {
			l.PushFront(v)
			l.PushBack(v)
		} // [6,5,4,3,2,2,3,4,5,6]

		require.Equal(t, 10, l.Len())
		require.Equal(t, 6, l.Front().Value)
		require.Equal(t, 6, l.Back().Value)

		l.MoveToFront(l.Front().Next) // [5,6,4,3,2,2,3,4,5,6]
		l.MoveToFront(l.Back().Prev)  // [5,5,6,4,3,2,2,3,4,6]

		elems = make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			v, ok := i.Value.(int)
			if ok {
				elems = append(elems, v)
			} else {
				require.True(t, ok, "Error list value type")
			}
		}
		require.Equal(t, []int{5, 5, 6, 4, 3, 2, 2, 3, 4, 6}, elems)
	})

	t.Run("additional 2", func(t *testing.T) {
		l := NewList()

		l.PushFront("dog") // [("dog"]
		l.PushBack("cat")  // ["dog", 2]
		l.PushFront("cat") // ["cat", "dog", "cat"]
		l.PushBack("rat")  // ["cat", "dog", "cat", "rat"]
		l.PushFront("rat") // ["rat", "cat", "dog", "cat", "rat"]

		require.Equal(t, 5, l.Len())

		middle := l.Front().Next.Next // "dog"
		l.Remove(middle)              // ["rat", "cat",  "cat", "rat"]
		require.Equal(t, 4, l.Len())

		elems := make([]string, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			v, ok := i.Value.(string)
			if ok {
				elems = append(elems, v)
			} else {
				require.True(t, ok, "Error list value type")
			}
		}
		require.Equal(t, []string{"rat", "cat", "cat", "rat"}, elems)

		for _, v := range [...]string{"frog", "caw"} {
			l.PushFront(v)
			l.PushBack(v)
		} // ["caw", "frog",  "rat", "cat",  "cat", "rat",  "frog", "caw"]

		require.Equal(t, 8, l.Len())
		require.Equal(t, "caw", l.Front().Value)
		require.Equal(t, "caw", l.Back().Value)

		l.MoveToFront(l.Front().Next) // ["frog", "caw",  "rat", "cat",  "cat", "rat",  "frog", "caw"]
		l.MoveToFront(l.Back().Prev)  // ["frog", "frog", "caw",   "rat", "cat",  "cat", "rat",   "caw"]

		elems = make([]string, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			v, ok := i.Value.(string)
			if ok {
				elems = append(elems, v)
			} else {
				require.True(t, ok, "Error list value type")
			}
		}
		require.Equal(t, []string{"frog", "frog", "caw", "rat", "cat", "cat", "rat", "caw"}, elems)
	})
}
