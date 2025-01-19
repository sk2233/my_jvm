/*
@author: sk
@date: 2025/1/4
*/
package main

import (
	"fmt"
	"reflect"
)

type NativeFunc func(thread *Thread)

var (
	nativeFuncs = make(map[string]NativeFunc)
)

func RegisterNativeFunc(class string, name string, desc string, func0 NativeFunc) {
	nativeFuncs[fmt.Sprintf("%s-%s-%s", class, name, desc)] = func0
}

func GetNativeFunc(class string, name string, desc string) NativeFunc {
	return nativeFuncs[fmt.Sprintf("%s-%s-%s", class, name, desc)]
}

func InitNativeFunc() {
	RegisterNativeFunc("HelloWorld", "max", "(II)I", func(thread *Thread) {
		frame := thread.Peek()
		val1 := frame.Pop()
		val2 := frame.Pop()
		frame.Push(NewInteger(max(val1.Integer, val2.Integer)))
	})
	RegisterNativeFunc("java/lang/Object", "getClass", "()Ljava/lang/Class;", func(thread *Thread) {
		frame := thread.Peek()
		class := frame.Pop().Object.Class
		name := class.GetString(class.ThisIndex)
		frame.Push(NewObject(makeClassObject(thread, name)))
	})
	RegisterNativeFunc("java/lang/Object", "hashCode", "()I", func(thread *Thread) {
		frame := thread.Peek()
		obj := frame.Pop().Object
		frame.Push(NewInteger(int32(reflect.ValueOf(obj).Pointer())))
	})
}
