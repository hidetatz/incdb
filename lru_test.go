package main

import (
	"reflect"
	"testing"
)

func TestLRU(t *testing.T) {
	lru := NewLRU(3)
	assert := func(t *testing.T, expected []string) {
		if !reflect.DeepEqual(lru.Keys(), expected) {
			t.Fatalf("unexpected lru: got: %v, expected: %v", lru.Keys(), expected)
		}
	}

	lru.Put("1", &Record{Vals: []string{"1"}})
	assert(t, []string{"1"})

	lru.Put("2", &Record{Vals: []string{"2"}})
	assert(t, []string{"2", "1"})

	lru.Put("3", &Record{Vals: []string{"3"}})
	assert(t, []string{"3", "2", "1"})

	lru.Put("4", &Record{Vals: []string{"4"}})
	assert(t, []string{"4", "3", "2"})

	lru.Get("2")
	assert(t, []string{"2", "4", "3"})

	lru.Put("5", &Record{Vals: []string{"5"}})
	assert(t, []string{"5", "2", "4"})

	got, _ := lru.Get("2")
	assert(t, []string{"2", "5", "4"})
	if !reflect.DeepEqual(got, &Record{Vals: []string{"2"}}) {
		t.Fatalf("unexpected lru: got: %v", got)
	}
}
