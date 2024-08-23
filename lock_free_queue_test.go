package main

import (
	"math/rand"
	"slices"
	"sync"
	"testing"
)

func TestEnqueue(t *testing.T) {

	n := 10_000
	expected := []int{}
	for range n {
		num := (rand.Int()+100)%1000
		expected = append(expected, num)
	}

	lfq := NewLockFreeQueue[int]()
	var wg sync.WaitGroup
	wg.Add(len(expected))
	for _, num := range expected {
		go func(val int){
			defer wg.Done()
			lfq.Enqueue(val)
		}(num)
	} 
	wg.Wait()

	slices.Sort(expected)
	actual := lfq.Slice()
	slices.Sort(actual)

	if slices.Compare(expected, actual) != 0 {
		t.Fatalf("enqueueing failed")
	}
}
