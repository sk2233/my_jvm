/*
@author: sk
@date: 2024/12/29
*/
package main

import (
	"fmt"
	"math"
	"strings"
)

type Instruction func(*Thread, *Class, *Code, int) int

//========================constants===========================

func InstructionNop(thread *Thread, class *Class, code *Code, pc int) int {
	return pc
}

func InstructionAConstNull(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	frame.Push(NewNull())
	return pc
}

func InstructionDConst0(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	frame.Push2(NewDouble(0))
	return pc
}

func InstructionFConst0(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	frame.Push(NewFloat(0))
	return pc
}

func InstructionIConst0(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	frame.Push(NewInteger(0))
	return pc
}

func InstructionIConst1(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	frame.Push(NewInteger(1))
	return pc
}

func InstructionIConst2(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	frame.Push(NewInteger(2))
	return pc
}

func InstructionIConst3(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	frame.Push(NewInteger(3))
	return pc
}

func InstructionIConst4(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	frame.Push(NewInteger(4))
	return pc
}

func InstructionIConst5(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	frame.Push(NewInteger(5))
	return pc
}

func InstructionLConst0(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	frame.Push2(NewLong(0))
	return pc
}

func InstructionLConst1(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	frame.Push2(NewLong(1))
	return pc
}

func InstructionBIPush(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	frame.Push(NewInteger(int32(ParseU8(code.Code, pc))))
	return pc + 1
}

func InstructionSIPush(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	frame.Push(NewInteger(int32(ParseU16(code.Code, pc))))
	return pc + 2
}

func InstructionLdc(thread *Thread, class *Class, code *Code, pc int) int {
	index := ParseU8(code.Code, pc)
	ldc(thread, class, int(index))
	return pc + 1
}

func InstructionLdcW(thread *Thread, class *Class, code *Code, pc int) int {
	index := ParseU16(code.Code, pc)
	ldc(thread, class, int(index))
	return pc + 2
}

func toJavaName(name string) string {
	return strings.ReplaceAll(name, "/", ".")
}

func makeClassObject(thread *Thread, name string) *Object {
	class := thread.Loader.LoadClass("java/lang/Class")
	field := class.GetField("name", "Ljava/lang/String;")
	fields := make([]*Value, class.InstSlotCount)
	fields[field.SlotID] = NewString(thread, name)
	return &Object{Class: class, Fields: fields}
}

func ldc(thread *Thread, class *Class, index int) {
	temp := class.Consts[index]
	frame := thread.Peek()
	switch temp.Type {
	case ConstInteger:
		frame.Push(NewInteger(temp.Integer))
	case ConstFloat:
		frame.Push(NewFloat(temp.Float))
	case ConstLong:
		frame.Push2(NewLong(temp.Long))
	case ConstDouble:
		frame.Push2(NewDouble(temp.Double))
	case ConstString:
		value := class.GetString(temp.Index)
		frame.Push(NewString(thread, value))
	case ConstClass:
		// 获取类名称
		name := toJavaName(class.GetString(temp.Index))
		// 构造类对象
		frame.Push(NewObject(makeClassObject(thread, name)))
	default:
		panic(fmt.Sprintf("unknown type: %v", temp.Type))
	}
}

//========================loads===========================

func instructionLoad(thread *Thread, index int) {
	frame := thread.Peek()
	frame.Push(frame.Get(index))
}

func InstructionLoad(thread *Thread, class *Class, code *Code, pc int) int {
	index := ParseU8(code.Code, pc)
	instructionLoad(thread, int(index))
	return pc + 1
}

func InstructionLoad0(thread *Thread, class *Class, code *Code, pc int) int {
	instructionLoad(thread, 0)
	return pc
}

func InstructionLoad1(thread *Thread, class *Class, code *Code, pc int) int {
	instructionLoad(thread, 1)
	return pc
}

func InstructionLoad2(thread *Thread, class *Class, code *Code, pc int) int {
	instructionLoad(thread, 2)
	return pc
}

func InstructionLoad3(thread *Thread, class *Class, code *Code, pc int) int {
	instructionLoad(thread, 3)
	return pc
}

func instruction2Load(thread *Thread, index int) { // 用于 Double Long
	frame := thread.Peek()
	frame.Push2(frame.Get(index))
}

func Instruction2Load(thread *Thread, class *Class, code *Code, pc int) int {
	index := ParseU8(code.Code, pc)
	instruction2Load(thread, int(index))
	return pc + 1
}

func Instruction2Load0(thread *Thread, class *Class, code *Code, pc int) int {
	instruction2Load(thread, 0)
	return pc
}

func Instruction2Load1(thread *Thread, class *Class, code *Code, pc int) int {
	instruction2Load(thread, 1)
	return pc
}

func Instruction2Load2(thread *Thread, class *Class, code *Code, pc int) int {
	instruction2Load(thread, 2)
	return pc
}

func Instruction2Load3(thread *Thread, class *Class, code *Code, pc int) int {
	instruction2Load(thread, 3)
	return pc
}

func InstructionALoad(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	index := frame.Pop().Integer
	arr := frame.Pop().Object
	frame.Push(arr.ArrayData[index])
	return pc
}

func Instruction2ALoad(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	index := frame.Pop().Integer
	arr := frame.Pop().Object
	frame.Push2(arr.ArrayData[index])
	return pc
}

//========================stores===========================

func instructionStore(thread *Thread, index int) {
	frame := thread.Peek()
	frame.Set(frame.Pop(), index)
}

func InstructionStore(thread *Thread, class *Class, code *Code, pc int) int {
	index := ParseU8(code.Code, pc)
	instructionStore(thread, int(index))
	return pc + 1
}

func InstructionStore0(thread *Thread, class *Class, code *Code, pc int) int {
	instructionStore(thread, 0)
	return pc
}

func InstructionStore1(thread *Thread, class *Class, code *Code, pc int) int {
	instructionStore(thread, 1)
	return pc
}

func InstructionStore2(thread *Thread, class *Class, code *Code, pc int) int {
	instructionStore(thread, 2)
	return pc
}

func InstructionStore3(thread *Thread, class *Class, code *Code, pc int) int {
	instructionStore(thread, 3)
	return pc
}

func instruction2Store(thread *Thread, index int) {
	frame := thread.Peek()
	frame.Set(frame.Pop2(), index)
}

func Instruction2Store(thread *Thread, class *Class, code *Code, pc int) int {
	index := ParseU8(code.Code, pc)
	instruction2Store(thread, int(index))
	return pc + 1
}

func Instruction2Store0(thread *Thread, class *Class, code *Code, pc int) int {
	instruction2Store(thread, 0)
	return pc
}

func Instruction2Store1(thread *Thread, class *Class, code *Code, pc int) int {
	instruction2Store(thread, 1)
	return pc
}

func Instruction2Store2(thread *Thread, class *Class, code *Code, pc int) int {
	instruction2Store(thread, 2)
	return pc
}

func Instruction2Store3(thread *Thread, class *Class, code *Code, pc int) int {
	instruction2Store(thread, 3)
	return pc
}

func InstructionAStore(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	val := frame.Pop()
	index := frame.Pop().Integer
	arr := frame.Pop().Object
	arr.ArrayData[index] = val
	return pc
}

func Instruction2AStore(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	val := frame.Pop2()
	index := frame.Pop().Integer
	arr := frame.Pop().Object
	arr.ArrayData[index] = val
	return pc
}

//========================stack=======================

func InstructionPop(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	frame.Pop()
	return pc
}

func InstructionPop2(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	frame.Pop()
	return pc
}

func InstructionDup(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	val := frame.Pop()
	frame.Push(val)
	frame.Push(val)
	return pc
}

func InstructionDupX1(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	val1 := frame.Pop()
	val2 := frame.Pop()
	frame.Push(val1)
	frame.Push(val2)
	frame.Push(val1)
	return pc
}

func InstructionDupX2(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	val1 := frame.Pop()
	val2 := frame.Pop()
	val3 := frame.Pop()
	frame.Push(val1)
	frame.Push(val3)
	frame.Push(val2)
	frame.Push(val1)
	return pc
}

func InstructionDup2(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	val1 := frame.Pop()
	val2 := frame.Pop()
	frame.Push(val2)
	frame.Push(val1)
	frame.Push(val2)
	frame.Push(val1)
	return pc
}

func InstructionDup2X1(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	val1 := frame.Pop()
	val2 := frame.Pop()
	val3 := frame.Pop()
	frame.Push(val2)
	frame.Push(val1)
	frame.Push(val3)
	frame.Push(val2)
	frame.Push(val1)
	return pc
}

func InstructionDup2X2(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	val1 := frame.Pop()
	val2 := frame.Pop()
	val3 := frame.Pop()
	val4 := frame.Pop()
	frame.Push(val2)
	frame.Push(val1)
	frame.Push(val4)
	frame.Push(val3)
	frame.Push(val2)
	frame.Push(val1)
	return pc
}

func InstructionSwap(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	val1 := frame.Pop()
	val2 := frame.Pop()
	frame.Push(val1)
	frame.Push(val2)
	return pc
}

//=====================math=====================

func InstructionDAdd(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	val1 := frame.Pop2()
	val2 := frame.Pop2()
	frame.Push2(NewDouble(val1.Double + val2.Double))
	return pc
}

func InstructionFAdd(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	val1 := frame.Pop()
	val2 := frame.Pop()
	frame.Push(NewFloat(val1.Float + val2.Float))
	return pc
}

func InstructionLAdd(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	val1 := frame.Pop2()
	val2 := frame.Pop2()
	frame.Push2(NewLong(val1.Long + val2.Long))
	return pc
}

func InstructionIAdd(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	val1 := frame.Pop()
	val2 := frame.Pop()
	frame.Push(NewInteger(val1.Integer + val2.Integer))
	return pc
}

func InstructionLAnd(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	val1 := frame.Pop2()
	val2 := frame.Pop2()
	frame.Push2(NewLong(val1.Long & val2.Long))
	return pc
}

func InstructionIAnd(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	val1 := frame.Pop()
	val2 := frame.Pop()
	frame.Push(NewInteger(val1.Integer & val2.Integer))
	return pc
}

func InstructionLOr(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	val1 := frame.Pop2()
	val2 := frame.Pop2()
	frame.Push2(NewLong(val1.Long | val2.Long))
	return pc
}

func InstructionIOr(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	val1 := frame.Pop()
	val2 := frame.Pop()
	frame.Push(NewInteger(val1.Integer | val2.Integer))
	return pc
}

func InstructionDDiv(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	val1 := frame.Pop2()
	val2 := frame.Pop2()
	frame.Push2(NewDouble(val2.Double / val1.Double))
	return pc
}

func InstructionFDiv(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	val1 := frame.Pop()
	val2 := frame.Pop()
	frame.Push(NewFloat(val2.Float / val1.Float))
	return pc
}

func InstructionLDiv(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	val1 := frame.Pop2()
	val2 := frame.Pop2()
	frame.Push2(NewLong(val2.Long / val1.Long))
	return pc
}

func InstructionIDiv(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	val1 := frame.Pop()
	val2 := frame.Pop()
	frame.Push(NewInteger(val2.Integer / val1.Integer))
	return pc
}

func InstructionIInc(thread *Thread, class *Class, code *Code, pc int) int {
	index := int(ParseU8(code.Code, pc))
	change := ParseU8(code.Code, pc+1)
	frame := thread.Peek()
	frame.Set(NewInteger(frame.Get(index).Integer+int32(change)), index)
	return pc + 2
}

func InstructionDMul(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	val1 := frame.Pop2()
	val2 := frame.Pop2()
	frame.Push2(NewDouble(val2.Double * val1.Double))
	return pc
}

func InstructionFMul(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	val1 := frame.Pop()
	val2 := frame.Pop()
	frame.Push(NewFloat(val2.Float * val1.Float))
	return pc
}

func InstructionLMul(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	val1 := frame.Pop2()
	val2 := frame.Pop2()
	frame.Push2(NewLong(val2.Long * val1.Long))
	return pc
}

func InstructionIMul(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	val1 := frame.Pop()
	val2 := frame.Pop()
	frame.Push(NewInteger(val2.Integer * val1.Integer))
	return pc
}

func InstructionDNeg(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	val := frame.Pop2()
	frame.Push2(NewDouble(-val.Double))
	return pc
}

func InstructionFNeg(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	val := frame.Pop()
	frame.Push(NewFloat(-val.Float))
	return pc
}

func InstructionLNeg(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	val := frame.Pop2()
	frame.Push2(NewLong(-val.Long))
	return pc
}

func InstructionINeg(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	val := frame.Pop()
	frame.Push(NewInteger(-val.Integer))
	return pc
}

func InstructionDMod(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	val1 := frame.Pop2()
	val2 := frame.Pop2()
	frame.Push2(NewDouble(math.Mod(val2.Double, val1.Double)))
	return pc
}

func InstructionFMod(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	val1 := frame.Pop()
	val2 := frame.Pop()
	frame.Push(NewFloat(float32(math.Mod(float64(val2.Float), float64(val1.Float)))))
	return pc
}

func InstructionLMod(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	val1 := frame.Pop2()
	val2 := frame.Pop2()
	frame.Push2(NewLong(val2.Long % val1.Long))
	return pc
}

func InstructionIMod(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	val1 := frame.Pop()
	val2 := frame.Pop()
	frame.Push(NewInteger(val2.Integer % val1.Integer))
	return pc
}

func InstructionDSub(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	val1 := frame.Pop2()
	val2 := frame.Pop2()
	frame.Push2(NewDouble(val2.Double - val1.Double))
	return pc
}

func InstructionFSub(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	val1 := frame.Pop()
	val2 := frame.Pop()
	frame.Push(NewFloat(val2.Float - val1.Float))
	return pc
}

func InstructionLSub(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	val1 := frame.Pop2()
	val2 := frame.Pop2()
	frame.Push2(NewLong(val2.Long - val1.Long))
	return pc
}

func InstructionISub(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	val1 := frame.Pop()
	val2 := frame.Pop()
	frame.Push(NewInteger(val2.Integer - val1.Integer))
	return pc
}

func InstructionLXor(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	val1 := frame.Pop2()
	val2 := frame.Pop2()
	frame.Push2(NewLong(val2.Long ^ val1.Long))
	return pc
}

func InstructionIXor(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	val1 := frame.Pop()
	val2 := frame.Pop()
	frame.Push(NewInteger(val2.Integer ^ val1.Integer))
	return pc
}

//======================conversions=======================

func InstructionI2B(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	val := frame.Pop()
	val.Integer &= 1
	frame.Push(val)
	return pc
}

//=====================comparisons=========================

func InstructionIfEq(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	offset := ParseI16(code.Code, pc)
	if frame.Pop().Integer == 0 {
		return pc + int(offset) - 1
	}
	return pc + 2
}

func InstructionIfNe(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	offset := ParseI16(code.Code, pc)
	if frame.Pop().Integer != 0 {
		return pc + int(offset) - 1
	}
	return pc + 2
}

func InstructionIfGt(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	offset := ParseI16(code.Code, pc)
	if frame.Pop().Integer > 0 {
		return pc + int(offset) - 1
	}
	return pc + 2
}

func InstructionLCmp(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	val1 := frame.Pop2()
	val2 := frame.Pop2()
	if val2.Long > val1.Long {
		frame.Push(NewInteger(1))
	} else if val2.Long < val1.Long {
		frame.Push(NewInteger(-1))
	} else {
		frame.Push(NewInteger(0))
	}
	return pc
}

func InstructionIfICmpNe(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	val1 := frame.Pop()
	val2 := frame.Pop()
	offset := ParseI16(code.Code, pc)
	if val2.Integer != val1.Integer {
		return pc + int(offset) - 1
	}
	return pc + 2
}

func InstructionIfICmpGt(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	val1 := frame.Pop()
	val2 := frame.Pop()
	offset := ParseI16(code.Code, pc)
	if val2.Integer > val1.Integer {
		return pc + int(offset) - 1
	}
	return pc + 2
}

func InstructionIfICmpGe(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	val1 := frame.Pop()
	val2 := frame.Pop()
	offset := ParseI16(code.Code, pc)
	if val2.Integer >= val1.Integer {
		return pc + int(offset) - 1
	}
	return pc + 2
}

func InstructionIfICmpLe(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	val1 := frame.Pop()
	val2 := frame.Pop()
	offset := ParseI16(code.Code, pc)
	if val2.Integer <= val1.Integer {
		return pc + int(offset) - 1
	}
	return pc + 2
}

func InstructionGoTo(thread *Thread, class *Class, code *Code, pc int) int {
	offset := ParseI16(code.Code, pc)
	return pc - 1 + int(offset)
}

func InstructionIfACmpNe(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	val1 := frame.Pop()
	val2 := frame.Pop()
	offset := ParseI16(code.Code, pc)
	if val2.Object != val1.Object {
		return pc + int(offset) - 1
	}
	return pc + 2
}

//====================control======================

func InstructionReturn(thread *Thread, class *Class, code *Code, pc int) int {
	thread.Pop()
	return 0xFFFFFFFF // return 强制退出方法
}

func InstructionReturn1(thread *Thread, class *Class, code *Code, pc int) int {
	oldFrame := thread.Pop()
	frame := thread.Peek()
	frame.Push(oldFrame.Pop())
	return 0xFFFFFFFF // return 强制退出方法
}

func InstructionReturn2(thread *Thread, class *Class, code *Code, pc int) int {
	oldFrame := thread.Pop()
	frame := thread.Peek()
	frame.Push2(oldFrame.Pop2())
	return 0xFFFFFFFF // return 强制退出方法
}

//=====================references======================

func InstructionNew(thread *Thread, class *Class, code *Code, pc int) int {
	index := ParseU16(code.Code, pc)
	className := class.GetString(index)
	newClass := thread.Loader.LoadClass(className)
	if IsInterface(newClass.Access) || IsAbstract(newClass.Access) {
		panic(fmt.Sprintf("interface or abstract class %s", className))
	}
	frame := thread.Peek()
	frame.Push(NewObject(&Object{Class: newClass, Fields: make([]*Value, newClass.InstSlotCount)}))
	return pc + 2
}

func loadClassAndField(thread *Thread, class *Class, index int) (*Class, *Field) {
	// ConstField
	fieldIndex := class.Consts[index]
	// 静态变量的目标 class
	className := class.GetString(fieldIndex.ClassIndex)
	resClass := thread.Loader.LoadClass(className)
	// 静态变量的目标 field
	nameType := class.Consts[fieldIndex.NameTypeIndex]
	name := class.GetString(nameType.NameIndex)
	desc := class.GetString(nameType.DescIndex)
	resField := resClass.GetField(name, desc)
	return resClass, resField
}

func InstructionPutStatic(thread *Thread, class *Class, code *Code, pc int) int {
	index := ParseU16(code.Code, pc)
	targetClass, targetField := loadClassAndField(thread, class, int(index))
	if !IsStatic(targetField.Access) {
		panic(fmt.Sprintf("%s is not static", targetClass.GetString(targetField.NameIndex)))
	}
	// 设置值
	frame := thread.Peek()
	slotID := targetField.SlotID
	if targetField.IsTwoSlot() {
		targetClass.StaticValues[slotID] = frame.Pop2()
	} else {
		targetClass.StaticValues[slotID] = frame.Pop()
	}
	return pc + 2
}

func InstructionGetStatic(thread *Thread, class *Class, code *Code, pc int) int {
	index := ParseU16(code.Code, pc)
	targetClass, targetField := loadClassAndField(thread, class, int(index))
	if !IsStatic(targetField.Access) {
		panic(fmt.Sprintf("%s is not static", targetClass.GetString(targetField.NameIndex)))
	}
	// 获取值
	frame := thread.Peek()
	slotID := targetField.SlotID
	if targetField.IsTwoSlot() {
		frame.Push2(targetClass.StaticValues[slotID])
	} else {
		frame.Push(targetClass.StaticValues[slotID])
	}
	return pc + 2
}

func InstructionPutField(thread *Thread, class *Class, code *Code, pc int) int {
	index := ParseU16(code.Code, pc)
	targetClass, targetField := loadClassAndField(thread, class, int(index))
	if IsStatic(targetField.Access) {
		panic(fmt.Sprintf("%s is static", targetClass.GetString(targetField.NameIndex)))
	}
	// 设置值
	frame := thread.Peek()
	slotID := targetField.SlotID
	if targetField.IsTwoSlot() {
		val := frame.Pop2()
		inst := frame.Pop()
		inst.Object.Fields[slotID] = val
	} else {
		val := frame.Pop()
		inst := frame.Pop()
		inst.Object.Fields[slotID] = val
	}
	return pc + 2
}

func InstructionGetField(thread *Thread, class *Class, code *Code, pc int) int {
	index := ParseU16(code.Code, pc)
	targetClass, targetField := loadClassAndField(thread, class, int(index))
	if IsStatic(targetField.Access) {
		panic(fmt.Sprintf("%s is static", targetClass.GetString(targetField.NameIndex)))
	}
	// 设置值
	frame := thread.Peek()
	slotID := targetField.SlotID
	if targetField.IsTwoSlot() {
		inst := frame.Pop()
		frame.Push2(inst.Object.Fields[slotID])
	} else {
		inst := frame.Pop()
		frame.Push(inst.Object.Fields[slotID])
	}
	return pc + 2
}

func InstructionInstanceOf(thread *Thread, class *Class, code *Code, pc int) int {
	// 获取目标Class
	index := ParseU16(code.Code, pc)
	className := class.GetString(index)
	targetClass := thread.Loader.LoadClass(className)
	// 目标实例
	frame := thread.Peek()
	inst := frame.Pop().Object
	// 判断是否符合
	if instanceOf(thread, inst.Class, targetClass) {
		frame.Push(NewInteger(1))
	} else {
		frame.Push(NewInteger(0))
	}
	return pc + 2
}

func InstructionCheckCast(thread *Thread, class *Class, code *Code, pc int) int {
	// 获取目标Class
	index := ParseU16(code.Code, pc)
	className := class.GetString(index)
	targetClass := thread.Loader.LoadClass(className)
	// 目标实例
	frame := thread.Peek()
	inst := frame.Peek().Object // 不要弹出对象，仅检查
	// 判断是否符合
	if !instanceOf(thread, inst.Class, targetClass) {
		panic(fmt.Sprintf("inst not a %s", className))
	}
	return pc + 2
}

func instanceOf(thread *Thread, subClass *Class, class *Class) bool {
	for subClass != class { // 先只看继承
		if subClass.SupperIndex == 0 {
			return false // 到头了
		}
		supperName := subClass.GetString(subClass.SupperIndex)
		subClass = thread.Loader.LoadClass(supperName)
	}
	return true
}

func loadClassAndMethod(thread *Thread, class *Class, index int) (*Class, *Field) {
	// ConstMethod
	methodIndex := class.Consts[index]
	// 变量的目标 class
	className := class.GetString(methodIndex.ClassIndex)
	resClass := thread.Loader.LoadClass(className)
	// 变量的目标 field
	nameType := class.Consts[methodIndex.NameTypeIndex]
	name := class.GetString(nameType.NameIndex)
	desc := class.GetString(nameType.DescIndex)
	resMethod := resClass.GetMethod(name, desc)
	return resClass, resMethod
}

func parseArgCount(class *Class, method *Field) int {
	count := 0
	desc := class.GetString(method.DescIndex)
	methodDesc := NewMethodDescParser(desc).Parse()
	for _, argType := range methodDesc.ArgTypes {
		count++
		if argType == "D" || argType == "J" {
			count++ // 占用两位
		}
	}
	if !IsStatic(method.Access) { // 非静态方法 this 传递
		count++
	}
	return count
}

func invokeMethod(thread *Thread, targetClass *Class, targetMethod *Field) {
	if IsNative(targetMethod.Access) { // 本地方法调用
		class := targetClass.GetString(targetClass.ThisIndex)
		name := targetClass.GetString(targetMethod.NameIndex)
		desc := targetClass.GetString(targetMethod.DescIndex)
		nativeFunc := GetNativeFunc(class, name, desc)
		nativeFunc(thread)
	} else { // 正常方法调用
		frame := thread.Peek()
		argCount := parseArgCount(targetClass, targetMethod)
		args := make([]*Value, argCount)
		for i := argCount - 1; i >= 0; i-- {
			args[i] = frame.Pop()
		}
		RunMethod(thread, targetMethod, args)
	}
}

// 静态方法
func InstructionInvokeStatic(thread *Thread, class *Class, code *Code, pc int) int {
	index := ParseU16(code.Code, pc)
	targetClass, targetMethod := loadClassAndMethod(thread, class, int(index))
	// 校验方法
	if !IsStatic(targetMethod.Access) {
		panic(fmt.Sprintf("%s not is static", targetClass.GetString(targetMethod.NameIndex)))
	}
	invokeMethod(thread, targetClass, targetMethod)
	return pc + 2
}

// 调用私有方法，静态方法等不需要动态绑定的方法没有搜索过程加快速度
func InstructionInvokeSpecial(thread *Thread, class *Class, code *Code, pc int) int {
	index := ParseU16(code.Code, pc)
	targetClass, targetMethod := loadClassAndMethod(thread, class, int(index))
	if IsStatic(targetMethod.Access) {
		panic(fmt.Sprintf("%s is static", targetClass.GetString(targetMethod.NameIndex)))
	}
	// TODO 捞取 inst 校验方法调用的合法性
	invokeMethod(thread, targetClass, targetMethod)
	return pc + 2
}

// 需要动态绑定的方法
func InstructionInvokeVirtual(thread *Thread, class *Class, code *Code, pc int) int {
	index := ParseU16(code.Code, pc)
	targetClass, targetMethod := loadClassAndMethod(thread, class, int(index))
	if IsStatic(targetMethod.Access) {
		panic(fmt.Sprintf("%s is static", targetClass.GetString(targetMethod.NameIndex)))
	}

	// MOCK
	name := targetClass.GetString(targetMethod.NameIndex)
	if name == "println" {
		frame := thread.Peek()
		desc := targetClass.GetString(targetMethod.DescIndex)
		switch desc {
		case "(Z)V", "(C)V", "(I)V", "(B)V", "(S)V", "(F)V":
			fmt.Println(frame.Pop())
		case "(J)V", "(D)V":
			fmt.Println(frame.Pop2())
		case "(Ljava/lang/String;)V":
			obj := frame.Pop().Object
			field := obj.Class.GetField("value", "[C")
			values := obj.Fields[field.SlotID]
			bs := make([]byte, 0)
			for _, item := range values.Object.ArrayData {
				bs = append(bs, byte(item.Integer))
			}
			fmt.Println(string(bs))
		default:
			panic(fmt.Sprintf("unknown %s", desc))
		}
		frame.Pop()
		return pc + 2
	}
	invokeMethod(thread, targetClass, targetMethod)
	return pc + 2
}

// 接口方法 使用 InvokeVirtual 也行，不过为了减小搜索范围
func InstructionInvokeInterface(thread *Thread, class *Class, code *Code, pc int) int {
	index := ParseU16(code.Code, pc)
	targetClass, targetMethod := loadClassAndMethod(thread, class, int(index)) // 接口中的定义
	if !IsInterface(targetClass.Access) {
		panic(fmt.Sprintf("%s not is interface", targetClass.GetString(targetClass.ThisIndex)))
	}
	if IsStatic(targetMethod.Access) {
		panic(fmt.Sprintf("%s is static", targetClass.GetString(targetMethod.NameIndex)))
	}

	frame := thread.Peek()
	argCount := parseArgCount(targetClass, targetMethod)
	inst := frame.PeekAt(argCount - 1)
	// TODO 捞取 inst 校验是否实现了对应接口
	// 转换为具体实现
	name := targetClass.GetString(targetMethod.NameIndex)
	desc := targetClass.GetString(targetMethod.DescIndex)
	targetClass = inst.Object.Class
	targetMethod = targetClass.GetMethod(name, desc)
	invokeMethod(thread, targetClass, targetMethod)
	return pc + 4 // 还有 2 byte 历史遗留不用管
}

var (
	arrayTypes = map[uint8]string{
		ArrayInt:  "[I",
		ArrayChar: "[C",
	}
)

func InstructionNewArray(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	count := frame.Pop().Integer
	if count < 0 {
		panic(fmt.Sprintf("invalid array size: %d", count))
	}

	arrayType := ParseU8(code.Code, pc)
	newClass := thread.Loader.LoadClass(arrayTypes[arrayType])
	frame.Push(NewObject(&Object{Class: newClass, ArrayType: arrayType, ArrayData: make([]*Value, count)}))
	return pc + 1
}

func InstructionObjArray(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	count := frame.Pop().Integer
	if count < 0 {
		panic(fmt.Sprintf("invalid array size: %d", count))
	}

	index := ParseU16(code.Code, pc)
	className := class.GetString(index) // 是基本元素的类型
	newClass := thread.Loader.LoadClass("[" + className)
	frame.Push(NewObject(&Object{Class: newClass, ArrayData: make([]*Value, count)}))
	return pc + 2
}

func InstructionArrayLen(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	obj := frame.Pop().Object
	if obj == nil {
		panic(fmt.Sprintf("obj is null"))
	}
	frame.Push(NewInteger(int32(len(obj.ArrayData))))
	return pc
}

func InstructionMultiArray(thread *Thread, class *Class, code *Code, pc int) int {
	index := ParseU16(code.Code, pc)
	className := class.GetString(index)
	dimension := ParseU8(code.Code, pc+2)
	counts := make([]int32, dimension)
	frame := thread.Peek()
	for i := len(counts) - 1; i >= 0; i-- {
		counts[i] = frame.Pop().Integer
	}

	frame.Push(makeMultiArray(thread, className, counts))
	return pc + 3
}

func makeMultiArray(thread *Thread, className string, counts []int32) *Value {
	data := make([]*Value, counts[0])
	if len(counts) > 1 {
		for j := 0; j < len(data); j++ { // 逐渐加载
			data[j] = makeMultiArray(thread, className[1:], counts[1:])
		}
	}
	class := thread.Loader.LoadClass(className)
	return NewObject(&Object{Class: class, ArrayData: data})
}

func InstructionAThrow(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	obj := frame.Pop().Object // 获取异常对象
	if obj == nil {
		panic(fmt.Sprintf("obj is null"))
	}
	for !thread.IsEmpty() { // 寻找异常处理函数
		frame = thread.Peek()
		code = frame.Method.GetCodeAttribute()
		exception := code.FindException(class, uint16(pc-1), obj)
		if exception != nil {
			frame.Clear() // 压入异常对象，跳转到异常处理函数
			frame.Push(NewObject(obj))
			return int(exception.Handler)
		} // 当前实际不支持上层捕获下层的异常 因为外层 for 循环没有弹出方法栈 暂时只支持同层捕获异常
		thread.Pop() // 没有找到，继续查找调用栈
	}
	panic(fmt.Sprintf("un handle exception %s", obj.String()))
}

//===================extended===================

func InstructionIfNonNull(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	val := frame.Pop()
	offset := ParseI16(code.Code, pc)
	if val != nil && val.Object != nil {
		return pc + int(offset) - 1
	}
	return pc + 2
}

func InstructionIfNull(thread *Thread, class *Class, code *Code, pc int) int {
	frame := thread.Peek()
	val := frame.Pop()
	offset := ParseI16(code.Code, pc)
	if val == nil || val.Object == nil {
		return pc + int(offset) - 1
	}
	return pc + 2
}
