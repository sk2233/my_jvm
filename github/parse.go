/*
@author: sk
@date: 2024/12/28
*/
package main

import (
	"bytes"
	"fmt"
	"os"
)

func ParseClass(path string) *Class {
	file := OpenFile(path)
	class := &Class{}
	class.Magic = ReadU32(file)
	class.Minor = ReadU16(file) // 先次后主
	class.Major = ReadU16(file)
	class.Consts = ParseConsts(file)
	class.Access = ReadU16(file)
	class.Name = ParseString(class.Consts, ReadU16(file))
	class.Supper = ParseString(class.Consts, ReadU16(file))
	class.Interfaces = ParseInterfaces(class.Consts, file)
	class.Fields = ParseFields(class.Consts, file)
	class.Methods = ParseFields(class.Consts, file) // 可以复用 Fields 的解析方式
	class.Attributes = ParseAttributes(class.Consts, file)
	return class
}

func ParseFields(consts []*Const, file *os.File) []*Field {
	count := ReadU16(file)
	res := make([]*Field, 0)
	for i := 0; i < int(count); i++ {
		res = append(res, &Field{
			Access:     ReadU16(file),
			Name:       ParseString(consts, ReadU16(file)),
			Desc:       ParseString(consts, ReadU16(file)),
			Attributes: ParseAttributes(consts, file),
		})
	}
	return res
}

func ParseAttributes(consts []*Const, file *os.File) []*Attribute {
	count := ReadU16(file)
	res := make([]*Attribute, 0)
	for i := 0; i < int(count); i++ {
		attr := &Attribute{
			Name: ParseString(consts, ReadU16(file)),
		}
		switch attr.Name { // 先对 code 特殊解析，其他先不管直接捞取所有 byte 到 data
		case AttributeCode:
			attr.Code = ParseCode(consts, file)
		case AttributeSourceFile:
			attr.SourceFile = ParseSourceFile(consts, file)
		default:
			attr.Data = ReadBytes(file, int(ReadU32(file)))
		}
		res = append(res, attr)
	}
	return res
}

func ParseSourceFile(consts []*Const, file *os.File) string {
	reader := bytes.NewReader(ReadBytes(file, int(ReadU32(file))))
	return ParseString(consts, ReadU16(reader))
}

func ParseCode(consts []*Const, file *os.File) *Code {
	// 先内存暂存
	reader := bytes.NewReader(ReadBytes(file, int(ReadU32(file))))
	return &Code{
		MaxStackDepth: ReadU16(reader),
		MaxLocals:     ReadU16(reader),
		Code:          ReadBytes(reader, int(ReadU32(reader))),
		Exceptions:    ParseExceptions(consts, reader),
		ExtAttrs:      ParseExtAttrs(consts, reader),
	}
}

func ParseExtAttrs(consts []*Const, reader *bytes.Reader) []*ExtAttr {
	count := ReadU16(reader)
	res := make([]*ExtAttr, 0)
	for i := 0; i < int(count); i++ {
		attr := &ExtAttr{
			Name: ParseString(consts, ReadU16(reader)),
		}
		switch attr.Name {
		case ExtAttrLineNumberTable:
			attr.LineNumbers = ParseLineNumbers(consts, reader)
		case ExtAttrLocalVariableTable:
			attr.LocalVariables = ParseLocalVariables(consts, reader)
		default:
			attr.Data = ReadBytes(reader, int(ReadU32(reader)))
		}
		res = append(res, attr)
	}
	return res
}

func ParseLocalVariables(consts []*Const, temp *bytes.Reader) []*LocalVariable {
	reader := bytes.NewReader(ReadBytes(temp, int(ReadU32(temp))))
	count := ReadU16(reader)
	res := make([]*LocalVariable, 0)
	for i := 0; i < int(count); i++ {
		res = append(res, &LocalVariable{
			Start: ReadU16(reader),
			Len:   ReadU16(reader),
			Name:  ParseString(consts, ReadU16(reader)),
			Type:  ReadU16(reader),
			Index: ReadU16(reader),
		})
	}
	return res
}

func ParseLineNumbers(consts []*Const, temp *bytes.Reader) []*LineNumber {
	reader := bytes.NewReader(ReadBytes(temp, int(ReadU32(temp))))
	count := ReadU16(reader)
	res := make([]*LineNumber, 0)
	for i := 0; i < int(count); i++ {
		res = append(res, &LineNumber{
			Start: ReadU16(reader),
			Line:  ReadU16(reader),
		})
	}
	return res
}

func ParseExceptions(consts []*Const, reader *bytes.Reader) []*Exception {
	count := ReadU16(reader)
	res := make([]*Exception, 0)
	for i := 0; i < int(count); i++ {
		res = append(res, &Exception{
			Start:     ReadU16(reader),
			End:       ReadU16(reader),
			Handler:   ReadU16(reader),
			CatchType: ReadU16(reader), // 指向常量池?
		})
	}
	return res
}

func ParseInterfaces(consts []*Const, file *os.File) []*Interface {
	count := ReadU16(file)
	res := make([]*Interface, 0)
	for i := 0; i < int(count); i++ {
		res = append(res, &Interface{
			Name: ParseString(consts, ReadU16(file)),
		})
	}
	return res
}

func ParseString(consts []*Const, index uint16) string {
	item := consts[index-1]
	switch item.Type {
	case ConstUtf8:
		return item.String
	case ConstClass, ConstString:
		return ParseString(consts, item.Index)
	}
	return ""
}

func ParseConsts(file *os.File) []*Const {
	count := ReadU16(file) - 1 // 比读到的长度是 - 1 的
	consts := make([]*Const, 0)
	for i := 0; i < int(count); i++ {
		item := &Const{
			Type: ConstType(ReadU8(file)),
		}
		switch item.Type {
		case ConstUtf8:
			l := ReadU16(file)
			item.String = string(ReadBytes(file, int(l)))
		case ConstInt, ConstFloat:
			item.Num = ReadU32(file)
		case ConstClass, ConstString:
			item.Index = ReadU16(file)
		case ConstField, ConstMethod, ConstInterfaceMethod, ConstNameType:
			item.Index = ReadU16(file)
			item.ExtIndex = ReadU16(file)
		default:
			panic(fmt.Errorf("unknown constant type: %v", item.Type))
		}
		consts = append(consts, item)
	}
	return consts
}
