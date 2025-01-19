/*
@author: sk
@date: 2024/12/29
*/
package main

import (
	"fmt"
)

const (
	ValueNull    = 1
	ValueDouble  = 2
	ValueFloat   = 3
	ValueInteger = 4
	ValueLong    = 5
	ValueObject  = 6
)

// double long 占用两个其他包含指针等都是占用一个 会采用第二个留空实际信息都存储在第一个的方式来占用两个
type Value struct { // 通用值类型
	Type    uint8
	Double  float64
	Float   float32
	Integer int32
	Long    int64
	Object  *Object
}

func (v *Value) String() string {
	switch v.Type {
	case ValueDouble:
		return fmt.Sprintf("%f", v.Double)
	case ValueFloat:
		return fmt.Sprintf("%f", v.Float)
	case ValueInteger:
		return fmt.Sprintf("%d", v.Integer)
	case ValueLong:
		return fmt.Sprintf("%d", v.Long)
	case ValueObject:
		return v.Object.String()
	case ValueNull:
		return "null"
	default:
		panic(fmt.Sprintf("unknown type %d", v.Type))
	}
}

const (
	ArrayBoolean = 4
	ArrayChar    = 5
	ArrayFloat   = 6
	ArrayDouble  = 7
	ArrayByte    = 8
	ArrayShort   = 9
	ArrayInt     = 10
	ArrayLong    = 11
)

type Object struct {
	Class  *Class
	Fields []*Value // 数组不使用
	// 数组专用，暂时只支持  int 与 Object
	ArrayType uint8
	ArrayData []*Value // 支持多种数据
}

func (o *Object) String() string {
	return fmt.Sprintf("<%s inst>", o.Class.GetString(o.Class.ThisIndex))
}

func NewObject(object *Object) *Value {
	return &Value{Object: object, Type: ValueObject}
}

func NewLong(long int64) *Value {
	return &Value{Long: long, Type: ValueLong}
}

func NewInteger(integer int32) *Value {
	return &Value{Integer: integer, Type: ValueInteger}
}

func NewFloat(float float32) *Value {
	return &Value{Float: float, Type: ValueFloat}
}

func NewDouble(double float64) *Value {
	return &Value{Double: double, Type: ValueDouble}
}

func NewNull() *Value {
	return &Value{Type: ValueNull}
}

type Frame struct {
	Method *Field
	Local  []*Value       // double long 占用两个其他包含指针等都是占用一个
	Stack  *Stack[*Value] // double long 占用两个其他包含指针等都是占用一个
}

func (f *Frame) Push(val *Value) {
	f.Stack.Push(val)
}

func (f *Frame) Push2(val *Value) { // 用于double  long
	f.Stack.Push(val)
	f.Stack.Push(val)
}

func (f *Frame) Pop() *Value {
	return f.Stack.Pop()
}

func (f *Frame) Pop2() *Value { // 用于double  long
	f.Stack.Pop()
	return f.Stack.Pop()
}

func (f *Frame) Peek() *Value {
	return f.Stack.Peek()
}

func (f *Frame) PeekAt(index int) *Value {
	return f.Stack.PeekAt(index)
}

func (f *Frame) Get(index int) *Value { // 不需要 Get2 只取对应的位置
	return f.Local[index]
}

func (f *Frame) Set(val *Value, index int) { // 不需要 Set2 只取对应的位置
	f.Local[index] = val
}

func (f *Frame) Clear() {
	for !f.Stack.IsEmpty() {
		f.Stack.Pop()
	}
}

func NewFrame(method *Field, maxLocal int, maxStack int, args []*Value) *Frame {
	local := make([]*Value, maxLocal)
	for i, arg := range args { // 接收初始化参数
		local[i] = arg
	}
	return &Frame{Method: method, Local: local, Stack: NewStack[*Value](maxStack)}
}

type Thread struct {
	Pc     int
	Stack  *Stack[*Frame]
	Loader *Loader
}

func (t *Thread) Push(frame *Frame) {
	t.Stack.Push(frame)
}

func (t *Thread) Peek() *Frame {
	return t.Stack.Peek()
}

func (t *Thread) Pop() *Frame {
	return t.Stack.Pop()
}

func (t *Thread) IsEmpty() bool {
	return t.Stack.IsEmpty()
}

func NewThread(loader *Loader) *Thread {
	return &Thread{Pc: 0, Stack: NewStack[*Frame](MaxStackDepth), Loader: loader}
}

func RunMain(class *Class, loader *Loader, args []string) {
	method := class.GetMethod("main", "([Ljava/lang/String;)V")
	thread := NewThread(loader)
	// 构造参数
	argsClass := loader.LoadClass("[java/lang/String")
	data := make([]*Value, 0)
	for _, arg := range args {
		data = append(data, NewString(thread, arg))
	}
	argVal := NewObject(&Object{Class: argsClass, ArrayData: data})
	RunMethod(thread, method, []*Value{argVal})
}

var (
	Instructions = make(map[byte]Instruction)
)

func InitInstruction() {
	Instructions = map[byte]Instruction{
		// constants
		0x00: InstructionNop,
		0x01: InstructionAConstNull,
		0x03: InstructionIConst0,
		0x04: InstructionIConst1,
		0x05: InstructionIConst2,
		0x06: InstructionIConst3,
		0x07: InstructionIConst4,
		0x08: InstructionIConst5,
		0x09: InstructionLConst0,
		0x0A: InstructionLConst1,
		0x0B: InstructionFConst0,
		0x0E: InstructionDConst0,
		0x10: InstructionBIPush,
		0x11: InstructionSIPush,
		0x12: InstructionLdc,
		0x13: InstructionLdcW,
		0x14: InstructionLdcW,
		// Integer
		0x15: InstructionLoad,
		0x1A: InstructionLoad0,
		0x1B: InstructionLoad1,
		0x1C: InstructionLoad2,
		0x1D: InstructionLoad3,
		// Long
		0x16: Instruction2Load,
		0x1E: Instruction2Load0,
		0x1F: Instruction2Load1,
		0x20: Instruction2Load2,
		0x21: Instruction2Load3,
		// Float
		0x17: InstructionLoad,
		0x22: InstructionLoad0,
		0x23: InstructionLoad1,
		0x24: InstructionLoad2,
		0x25: InstructionLoad3,
		// Double
		0x18: Instruction2Load,
		0x26: Instruction2Load0,
		0x27: Instruction2Load1,
		0x28: Instruction2Load2,
		0x29: Instruction2Load3,
		// 对象
		0x19: InstructionLoad,
		0x2A: InstructionLoad0,
		0x2B: InstructionLoad1,
		0x2C: InstructionLoad2,
		0x2D: InstructionLoad3,
		// 数组
		0x2E: InstructionALoad,
		0x2F: Instruction2ALoad,
		0x30: InstructionALoad,
		0x31: Instruction2ALoad,
		0x32: InstructionALoad,
		0x33: InstructionALoad,
		0x34: InstructionALoad,
		0x35: InstructionALoad,
		// Integer
		0x36: InstructionStore,
		0x3B: InstructionStore0,
		0x3C: InstructionStore1,
		0x3D: InstructionStore2,
		0x3E: InstructionStore3,
		// Long
		0x37: Instruction2Store,
		0x3F: Instruction2Store0,
		0x40: Instruction2Store1,
		0x41: Instruction2Store2,
		0x42: Instruction2Store3,
		// Float
		0x38: InstructionStore,
		0x43: InstructionStore0,
		0x44: InstructionStore1,
		0x45: InstructionStore2,
		0x46: InstructionStore3,
		// Double
		0x39: Instruction2Store,
		0x47: Instruction2Store0,
		0x48: Instruction2Store1,
		0x49: Instruction2Store2,
		0x4A: Instruction2Store3,
		// 对象
		0x3A: InstructionStore,
		0x4B: InstructionStore0,
		0x4C: InstructionStore1,
		0x4D: InstructionStore2,
		0x4E: InstructionStore3,
		// array
		0x4F: InstructionAStore,
		0x50: Instruction2AStore,
		0x51: InstructionAStore,
		0x52: Instruction2AStore,
		0x53: InstructionAStore,
		0x54: InstructionAStore,
		0x55: InstructionAStore,
		0x56: InstructionAStore,
		// stack
		0x57: InstructionPop,
		0x58: InstructionPop2,
		0x59: InstructionDup,
		0x5A: InstructionDupX1,
		0x5B: InstructionDupX2,
		0x5C: InstructionDup2,
		0x5D: InstructionDup2X1,
		0x5E: InstructionDup2X2,
		0x5F: InstructionSwap,
		// math
		0x60: InstructionIAdd,
		0x61: InstructionLAdd,
		0x62: InstructionFAdd,
		0x63: InstructionDAdd,
		0x64: InstructionISub,
		0x65: InstructionLSub,
		0x66: InstructionFSub,
		0x67: InstructionDSub,
		0x68: InstructionIMul,
		0x69: InstructionLMul,
		0x6A: InstructionFMul,
		0x6B: InstructionDMul,
		0x6C: InstructionIDiv,
		0x6D: InstructionLDiv,
		0x6E: InstructionFDiv,
		0x6F: InstructionDDiv,
		0x70: InstructionIMod,
		0x71: InstructionLMod,
		0x72: InstructionFMod,
		0x73: InstructionDMod,
		0x74: InstructionINeg,
		0x75: InstructionLNeg,
		0x76: InstructionFNeg,
		0x77: InstructionDNeg,
		0x7e: InstructionIAnd,
		0x7F: InstructionLAnd,
		0x80: InstructionIOr,
		0x81: InstructionLOr,
		0x82: InstructionIXor,
		0x83: InstructionLXor,
		0x84: InstructionIInc,
		// conversions
		0x91: InstructionI2B,
		// comparisons
		0x94: InstructionLCmp,
		0x99: InstructionIfEq,
		0x9A: InstructionIfNe,
		0x9D: InstructionIfGt,
		0xA0: InstructionIfICmpNe,
		0xA2: InstructionIfICmpGe,
		0xA3: InstructionIfICmpGt,
		0xA4: InstructionIfICmpLe,
		0xA6: InstructionIfACmpNe,
		0xA7: InstructionGoTo,
		// control
		0xAC: InstructionReturn1,
		0xAD: InstructionReturn2,
		0xAE: InstructionReturn1,
		0xAF: InstructionReturn2,
		0xB0: InstructionReturn1,
		0xB1: InstructionReturn,
		// references
		0xB2: InstructionGetStatic,
		0xB3: InstructionPutStatic,
		0xB4: InstructionGetField,
		0xB5: InstructionPutField,
		0xB6: InstructionInvokeVirtual,
		0xB7: InstructionInvokeSpecial,
		0xB8: InstructionInvokeStatic,
		0xB9: InstructionInvokeInterface,
		0xBB: InstructionNew,
		0xBC: InstructionNewArray,
		0xBD: InstructionObjArray,
		0xBE: InstructionArrayLen,
		0xBF: InstructionAThrow,
		0xC0: InstructionCheckCast,
		0xC1: InstructionInstanceOf,
		0xC5: InstructionMultiArray,
		// extended
		0xC6: InstructionIfNull,
		0xC7: InstructionIfNonNull,
	}
}

var (
	internStrings = make(map[string]*Value)
)

func NewString(thread *Thread, val string) *Value {
	if _, ok := internStrings[val]; !ok {
		// string 对象
		class := thread.Loader.LoadClass("java/lang/String")
		res := NewObject(&Object{Class: class, Fields: make([]*Value, class.InstSlotCount)})
		// char[] 对象
		fieldClass := thread.Loader.LoadClass("[C")
		data := make([]*Value, 0)
		for i := 0; i < len(val); i++ { // 这里使用的 utf-8 编码 非  utf-16 编码
			data = append(data, NewInteger(int32(val[i])))
		}
		value := NewObject(&Object{Class: fieldClass, ArrayType: ArrayChar, ArrayData: data})
		// 设置值
		field := class.GetField("value", "[C")
		res.Object.Fields[field.SlotID] = value
		internStrings[val] = res
	}
	return internStrings[val]
}

func RunMethod(thread *Thread, method *Field, args []*Value) {
	code := method.GetCodeAttribute()
	thread.Push(NewFrame(method, int(code.MaxLocal), int(code.MaxStack), args))
	pc := 0
	for pc < len(code.Code) {
		opCode := code.Code[pc]
		fmt.Printf("pc:%d inst:%T %x\n", pc, Instructions[opCode], opCode)
		pc++
		if instruction, ok := Instructions[opCode]; ok {
			pc = instruction(thread, method.Class, code, pc)
		} else {
			panic(fmt.Sprintf("opcode %x not found", opCode))
		}
	}
}

type MethodDesc struct {
	ArgTypes []string
	RetType  string
}

type MethodDescParser struct {
	Desc  string
	Index int
}

func (p *MethodDescParser) Parse() *MethodDesc {
	desc := &MethodDesc{}
	p.Must('(')
	for !p.Match(')') {
		desc.ArgTypes = append(desc.ArgTypes, p.ParseType())
	}
	desc.RetType = p.ParseType()
	return desc
}

func (p *MethodDescParser) ParseType() string {
	switch p.Desc[p.Index] {
	case 'I', 'V', 'Z', 'J':
		p.Index++
		return p.Desc[p.Index-1 : p.Index]
	case 'L': // 引用类型
		p.Index++
		index := p.Index
		for p.Desc[p.Index] != ';' {
			p.Index++
		}
		res := p.Desc[index:p.Index]
		p.Index++
		return res
	case '[': // 数组类型
		p.Index++
		return "[" + p.ParseType()
	default:
		panic(fmt.Sprintf("unknown token %v", p.Desc[p.Index]))
	}
}

func (p *MethodDescParser) Must(token uint8) {
	if !p.Match(token) {
		panic(fmt.Sprintf("token %v not match", token))
	}
}

func (p *MethodDescParser) Match(token uint8) bool {
	if p.Desc[p.Index] != token {
		return false
	}
	p.Index++
	return true
}

func NewMethodDescParser(desc string) *MethodDescParser {
	return &MethodDescParser{Desc: desc, Index: 0}
}
