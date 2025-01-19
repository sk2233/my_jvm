/*
@author: sk
@date: 2024/12/29
*/
package main

type Stack[T any] struct {
	Data  []T
	Index int
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{}
}
