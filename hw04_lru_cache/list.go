package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v any) *ListItem
	PushBack(v any) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
	Clear()
}

type ListItem struct {
	Value any
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len       int
	frontItem *ListItem
	backItem  *ListItem
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.frontItem
}

func (l *list) Back() *ListItem {
	return l.backItem
}

func (l *list) PushFront(v any) *ListItem {
	if l.frontItem == nil {
		l.frontItem = new(ListItem)
		l.frontItem.Value = v
		l.backItem = l.frontItem
	} else {
		newFrontItem := new(ListItem)
		newFrontItem.Value = v
		newFrontItem.Next = l.frontItem
		l.frontItem.Prev = newFrontItem
		l.frontItem = newFrontItem
	}
	l.len++
	return l.frontItem
}

func (l *list) PushBack(v any) *ListItem {
	if l.backItem == nil {
		l.backItem = new(ListItem)
		l.backItem.Value = v
		l.frontItem = l.backItem
	} else {
		newBackItem := new(ListItem)
		newBackItem.Value = v
		newBackItem.Prev = l.backItem
		l.backItem.Next = newBackItem
		l.backItem = newBackItem
	}
	l.len++
	return l.frontItem
}

func (l *list) Remove(i *ListItem) {
	if i != nil {
		if l.len > 0 {
			switch {
			case (i.Prev == nil) && (i.Next == nil):
				l.backItem = nil
				l.frontItem = nil
				l.len = 0
			case (i.Prev != nil) && (i.Next != nil):
				i.Prev.Next = i.Next
				i.Next.Prev = i.Prev
				l.len--
			case i.Prev == nil:
				l.frontItem = i.Next
				i.Next.Prev = nil
				l.len--
			case i.Next == nil:
				l.backItem = i.Prev
				i.Prev.Next = nil
				l.len--
			}
		}
		i.Prev = nil
		i.Next = nil
	}
}

func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.PushFront(i.Value)
}

func (l *list) Clear() {
	l.frontItem = nil
	l.backItem = nil
	l.len = 0
}

func NewList() List {
	return new(list)
}
