package main

import (
	"container/list"
	"sync"
)

type LRU struct {
	cap  int
	list *list.List
	elem map[string]*list.Element
	m    sync.Mutex
}

type element struct {
	key   string
	value *Record
}

func NewLRU(cap int) *LRU {
	return &LRU{
		cap:  cap,
		list: list.New(),
		elem: make(map[string]*list.Element),
	}
}

func (l *LRU) Get(key string) (*Record, bool) {
	l.m.Lock()
	defer l.m.Unlock()

	if e, ok := l.elem[key]; ok {
		v := e.Value
		l.list.MoveToFront(e)
		return v.(*element).value, true
	}

	return nil, false
}

func (l *LRU) Put(key string, r *Record) {
	l.m.Lock()
	defer l.m.Unlock()

	e := &element{key: key, value: r}
	l.list.PushFront(e)
	l.elem[key] = l.list.Front()
	if l.list.Len() > l.cap {
		evicted := l.list.Back()
		delete(l.elem, evicted.Value.(*element).key)
		l.list.Remove(evicted)
	}
}

func (l *LRU) Keys() []string {
	l.m.Lock()
	defer l.m.Unlock()

	keys := make([]string, l.list.Len())
	i := 0
	for e := l.list.Front(); e != nil; e = e.Next() {
		keys[i] = e.Value.(*element).key
		i++
	}
	return keys
}
