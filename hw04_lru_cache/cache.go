package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value any) bool
	Get(key Key) (any, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	mu       sync.Mutex
}

type keyValue struct {
	key   Key
	value any
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (lc *lruCache) Set(key Key, value any) bool {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	itm, ok := lc.items[key]
	if ok {
		lc.queue.MoveToFront(itm)
		itm.Value = keyValue{key: key, value: value}
		return true
	}
	lc.items[key] = lc.queue.PushFront(keyValue{key: key, value: value})
	for lc.queue.Len() > lc.capacity {
		litm := lc.queue.Back().Value.(keyValue)
		lc.queue.Remove(lc.queue.Back())
		delete(lc.items, litm.key)
	}
	return false
}

func (lc *lruCache) Get(key Key) (any, bool) {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	itm, ok := lc.items[key]
	if ok {
		litm := itm.Value.(keyValue)
		lc.queue.MoveToFront(itm)
		return litm.value, true
	}
	return nil, false
}

func (lc *lruCache) Clear() {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	lc.queue.Clear()
	clear(lc.items)
}
