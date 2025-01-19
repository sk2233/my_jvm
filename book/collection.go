/*
@author: sk
@date: 2024/12/29
*/
package main

import "fmt"

type Stack[T any] struct {
	Data  []T
	Index int
}

func (t *Stack[T]) Push(val T) {
	if t.Index < len(t.Data) {
		t.Data[t.Index] = val
		t.Index++
		return
	}
	panic(fmt.Sprintf("stack overflow depth %v", t.Index))
}

func (t *Stack[T]) Pop() T {
	t.Index--
	return t.Data[t.Index]
}

func (t *Stack[T]) Peek() T {
	return t.Data[t.Index-1]
}

func (t *Stack[T]) PeekAt(index int) T {
	return t.Data[t.Index-index-1]
}

func (t *Stack[T]) IsEmpty() bool {
	return t.Index == 0
}

func NewStack[T any](size int) *Stack[T] {
	return &Stack[T]{Data: make([]T, size), Index: 0}
}
