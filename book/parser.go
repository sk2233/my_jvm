/*
@author: sk
@date: 2024/12/29
*/
package main

import (
	"encoding/binary"
	"fmt"
	"math"
)

type Parser struct {
	Data  []byte
	Index int
}

func (p *Parser) ParseClass() *Class {
	class := &Class{}
	class.Magic = p.ReadU32()
	class.Minor = p.ReadU16()
	class.Major = p.ReadU16()
	class.Consts = p.ParseConsts()
	class.Access = p.ReadU16()
	class.ThisIndex = p.ReadU16()
	class.SupperIndex = p.ReadU16()
	class.Interfaces = p.ReadU16s()
	class.Fields = p.ParseFields(class)
	class.Methods = p.ParseFields(class)
	class.Attributes = p.ParseAttributes(class.Consts)
	return class
}

func (p *Parser) ReadBytes(count int) []byte {
	index := p.Index
	p.Index += count
	return p.Data[index:p.Index]
}

func (p *Parser) ReadU8() uint8 {
	return p.ReadBytes(1)[0]
}

func (p *Parser) ReadU16() uint16 {
	bs := p.ReadBytes(2)
	return binary.BigEndian.Uint16(bs)
}

func (p *Parser) ReadU64() uint64 {
	bs := p.ReadBytes(8)
	return binary.BigEndian.Uint64(bs)
}

func (p *Parser) ReadU32() uint32 {
	bs := p.ReadBytes(4)
	return binary.BigEndian.Uint32(bs)
}

func (p *Parser) ParseConsts() []*Const {
	count := p.ReadU16()
	consts := make([]*Const, 0)
	consts = append(consts, &Const{}) // 第一个位置不使用
	for i := 1; i < int(count); i++ { // 从 1 开始
		item := &Const{
			Type: p.ReadU8(),
		}
		switch item.Type {
		case ConstUtf8:
			l := p.ReadU16()
			item.String = string(p.ReadBytes(int(l)))
		case ConstInteger: // bool byte char 也都是这个
			item.Integer = int32(p.ReadU32())
		case ConstFloat:
			item.Float = math.Float32frombits(p.ReadU32())
		case ConstLong:
			item.Long = int64(p.ReadU64())
		case ConstDouble:
			item.Double = math.Float64frombits(p.ReadU64())
		case ConstClass, ConstString, ConstMethodType:
			item.Index = p.ReadU16()
		case ConstField, ConstMethod, ConstInterfaceMethod, ConstInvokeDynamic:
			item.ClassIndex = p.ReadU16()
			item.NameTypeIndex = p.ReadU16()
		case ConstNameType:
			item.NameIndex = p.ReadU16()
			item.DescIndex = p.ReadU16()
		case ConstMethodHandle:
			p.ReadU8()
			p.ReadU16()
		default:
			panic(fmt.Errorf("unknown constant type: %v", item.Type))
		}
		consts = append(consts, item)
		// https://docs.oracle.com/javase/specs/jvms/se8/html/jvms-4.html#jvms-4.4.5
		if item.Type == ConstLong || item.Type == ConstDouble {
			consts = append(consts, &Const{}) // Long Double 占用两个位置，官方都承认这是一个糟糕的设计
			i++                               // 占用两个位置
		}
	}
	return consts
}

func (p *Parser) ReadU16s() []uint16 {
	count := p.ReadU16()
	res := make([]uint16, 0)
	for i := 0; i < int(count); i++ {
		res = append(res, p.ReadU16())
	}
	return res
}

func (p *Parser) ParseFields(class *Class) []*Field {
	count := p.ReadU16()
	fields := make([]*Field, 0)
	for i := 0; i < int(count); i++ {
		fields = append(fields, &Field{
			Access:     p.ReadU16(),
			NameIndex:  p.ReadU16(),
			DescIndex:  p.ReadU16(),
			Attributes: p.ParseAttributes(class.Consts),
			Class:      class,
		})
	}
	return fields
}

func (p *Parser) ParseAttributes(consts []*Const) []*Attribute {
	count := p.ReadU16()
	attrs := make([]*Attribute, 0)
	for i := 0; i < int(count); i++ {
		attr := &Attribute{
			Name: ParseString(consts, int(p.ReadU16())),
		} // 确保耗尽
		temp := NewParser(p.ReadBytes(int(p.ReadU32())))
		switch attr.Name {
		case AttributeCode:
			attr.Code = temp.ParseCode(consts)
		case AttributeSourceFile:
			attr.SourceFileIndex = temp.ReadU16()
		case AttributeExceptions:
			attr.ExceptionIndexes = temp.ReadU16s()
		case AttributeLineNumberTable:
			attr.LineNumbers = temp.ParseLineNumbers()
		case AttributeConstantValue:
			attr.ConstantValueIndex = temp.ReadU16()
		default:
			fmt.Println("unknown attribute:", attr.Name)
			attr.Data = temp.ReadAll()
		}
		attrs = append(attrs, attr)
	}
	return attrs
}

func (p *Parser) ParseCode(consts []*Const) *Code {
	return &Code{
		MaxStack:   p.ReadU16(),
		MaxLocal:   p.ReadU16(),
		Code:       p.ReadBytes(int(p.ReadU32())),
		Exceptions: p.ParseExceptions(),
		Attributes: p.ParseAttributes(consts),
	}
}

func (p *Parser) ParseExceptions() []*Exception {
	count := p.ReadU16()
	exceptions := make([]*Exception, 0)
	for i := 0; i < int(count); i++ {
		exceptions = append(exceptions, &Exception{
			Start:     p.ReadU16(),
			End:       p.ReadU16(),
			Handler:   p.ReadU16(),
			CatchType: p.ReadU16(),
		})
	}
	return exceptions
}

func (p *Parser) ReadAll() []byte {
	index := p.Index
	p.Index = len(p.Data)
	return p.Data[index:]
}

func (p *Parser) ParseLineNumbers() []*LineNumber {
	count := p.ReadU16()
	res := make([]*LineNumber, 0)
	for i := 0; i < int(count); i++ {
		res = append(res, &LineNumber{
			Start: p.ReadU16(),
			Line:  p.ReadU16(),
		})
	}
	return res
}

func ParseString(consts []*Const, index int) string {
	return consts[index].String
}

func NewParser(data []byte) *Parser {
	return &Parser{Data: data, Index: 0}
}
