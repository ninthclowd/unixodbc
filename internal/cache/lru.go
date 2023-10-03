package cache

import (
	"container/list"
	"sync"
)

type OnEvict[T any] func(key string, value *T) error

type lruNode[T any] struct {
	checkout sync.Mutex
	key      string
	value    *T
}

type LRU[T any] struct {
	capacity      int
	mux           sync.Mutex
	elements      *list.List
	elementForKey map[string]*list.Element
	onEvict       OnEvict[T]
}

func NewLRU[T any](capacity int, onEvict OnEvict[T]) *LRU[T] {
	return &LRU[T]{elementForKey: make(map[string]*list.Element), elements: list.New(), capacity: capacity, onEvict: onEvict}
}

func (l *LRU[T]) evictLRU() error {

	lru := l.elements.Back()
	node := lru.Value.(*lruNode[T])
	l.elements.Remove(lru)
	delete(l.elementForKey, node.key)
	if l.onEvict != nil {
		if err := l.onEvict(node.key, node.value); err != nil {
			return err
		}
	}
	return nil
}

func (l *LRU[T]) Put(key string, value *T) error {
	if l.capacity == 0 {
		if l.onEvict != nil {
			return l.onEvict(key, value)
		}
		return nil
	}
	l.mux.Lock()
	defer l.mux.Unlock()
	if element, ok := l.elementForKey[key]; ok {
		l.elements.MoveToFront(element)
		return nil
	}
	if l.capacity == len(l.elementForKey) {
		if err := l.evictLRU(); err != nil {
			return err
		}
	}
	newNode := &lruNode[T]{key: key, value: value}
	l.elementForKey[key] = l.elements.PushFront(newNode)
	return nil
}

func (l *LRU[T]) Get(key string, removeIfFound bool) *T {
	if l.capacity == 0 {
		return nil
	}
	l.mux.Lock()
	defer l.mux.Unlock()
	if element, ok := l.elementForKey[key]; ok {
		node := element.Value.(*lruNode[T])
		if removeIfFound {
			l.elements.Remove(element)
			delete(l.elementForKey, key)
		} else {
			l.elements.MoveToFront(element)
		}
		return node.value
	}
	return nil
}

func (l *LRU[T]) Purge() error {
	l.mux.Lock()
	defer l.mux.Unlock()

	for key, element := range l.elementForKey {
		l.elements.Remove(element)
		delete(l.elementForKey, key)
		node := element.Value.(*lruNode[T])
		if l.onEvict != nil {
			if err := l.onEvict(node.key, node.value); err != nil {
				return err
			}
		}
	}
	return nil
}
