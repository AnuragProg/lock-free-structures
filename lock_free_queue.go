package main

import "sync/atomic"

type LockFreeNode[T any] struct {
	val  T
	next *atomic.Pointer[LockFreeNode[T]]
}

type LockFreeQueue[T any] struct {
	head, tail *atomic.Pointer[LockFreeNode[T]]
}

func NewLockFreeQueue[T any]() LockFreeQueue[T] {
	dummyNode := &LockFreeNode[T]{
		next: &atomic.Pointer[LockFreeNode[T]]{},
	}
	head := atomic.Pointer[LockFreeNode[T]]{}
	head.Store(dummyNode)
	tail := atomic.Pointer[LockFreeNode[T]]{}
	tail.Store(dummyNode)

	return LockFreeQueue[T]{
		head: &head,
		tail: &tail,
	}
}

func (lfq *LockFreeQueue[T]) Enqueue(val T) {
	newNode := LockFreeNode[T]{
		val: val,
		next: &atomic.Pointer[LockFreeNode[T]]{},
	}

	for {
		tail := lfq.tail.Load()

		if lfq.tail.CompareAndSwap(tail, &newNode) { // swap was successfull
			tail.next.Store(&newNode) // now previous tail is immutable 
			break
		}
	}
}


func (lfq *LockFreeQueue[T]) Slice() []T {
	slice := []T{}
	
	current := lfq.head.Load().next.Load() // head will never be nil

	for current != nil {
		slice = append(slice, current.val)
		current = current.next.Load()
	}

	return slice
}
