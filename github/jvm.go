/*
@author: sk
@date: 2024/12/29
*/
package main

import (
	"encoding/binary"
	"fmt"
)

type JavaFrame struct {
	Stack  *Stack[any]
	Locals []any
}

type JVM struct {
	Classes map[string]*Class
	Stack   *Stack[*JavaFrame]
}

func (j *JVM) LoadClass(path string) *Class {
	if _, ok := j.Classes[path]; !ok {
		j.Classes[path] = ParseClass(path)
	}
	return j.Classes[path]
}

func (j *JVM) CallStaticMethod(class *Class, method *Field) {
	code := method.GetCode()
	pc := 0
	for pc < len(code.Code) {
		opCode := code.Code[pc]
		pc++
		switch opCode {
		case 0xB2: // getstatic
			index := binary.BigEndian.Uint16(code.Code[pc : pc+2])
			pc += 2
			field := class.Consts[index-1]
			name := ParseString(class.Consts, field.Index)
			nameType := class.Consts[field.ExtIndex-1]
			typeName := ParseString(class.Consts, nameType.Index)
			fmt.Println(name, typeName)
		case 0x03: // iconst0
			fmt.Println(0)
		case 0xB3: // putstatic
			index := binary.BigEndian.Uint16(code.Code[pc : pc+2])
			pc += 2
			field := class.Consts[index-1]
			name := ParseString(class.Consts, field.Index)
			nameType := class.Consts[field.ExtIndex-1]
			typeName := ParseString(class.Consts, nameType.Index)
			fmt.Println(name, typeName)
		case 0x12: // ldc
			index := code.Code[pc]
			pc++
			val := ParseString(class.Consts, uint16(index))
			fmt.Println(val)
		case 0x59: // dup
			fmt.Println("dup")
		case 0x2A: // aload0
			fmt.Println("aload0")
		case 0x2B: // aload1
			fmt.Println("aload1")
		case 0xBB: // new
			index := binary.BigEndian.Uint16(code.Code[pc : pc+2])
			pc += 2
			className := ParseString(class.Consts, index)
			fmt.Println(className)
		case 0xB6: // invokevirtual
			index := binary.BigEndian.Uint16(code.Code[pc : pc+2])
			pc += 2
			method0 := class.Consts[index-1]
			className := ParseString(class.Consts, method0.Index)
			nameType := class.Consts[method0.ExtIndex-1]
			typeName := ParseString(class.Consts, nameType.Index)
			fmt.Println(className, typeName)
		case 0xB7: // invokespecial
			index := binary.BigEndian.Uint16(code.Code[pc : pc+2])
			pc += 2
			fmt.Println(index)
		case 0xB8: // invokestatic
			index := binary.BigEndian.Uint16(code.Code[pc : pc+2])
			pc += 2
			method0 := class.Consts[index-1]
			className := ParseString(class.Consts, method0.Index)
			nameType := class.Consts[method0.ExtIndex-1]
			typeName := ParseString(class.Consts, nameType.Index)
			fmt.Println(className, typeName)
		case 0xB1: // return
			fmt.Println("return")
		default:
			panic(fmt.Sprintf("unknown opCode = %v", opCode))
		}
	}
}

func NewJVM() *JVM {
	return &JVM{Classes: make(map[string]*Class)}
}
