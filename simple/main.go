/*
@author: sk
@date: 2024/12/28
*/
package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

// https://zserge.com/posts/jvm/

type Const struct {
	Tag              byte
	NameIndex        uint16
	ClassIndex       uint16
	NameAndTypeIndex uint16
	StringIndex      uint16
	DescIndex        uint16
	String           string
}

type Field struct {
	Flags      uint16
	Name       string
	Descriptor string
	Attributes map[string]*Attribute
}

type Attribute struct {
	Name string
	Data []byte
}

func main() {
	file, err := os.Open("simple/Add.class")
	HandleErr(err)
	bs := ReadBytes(file, 4)
	major, minor := ReadU16(file), ReadU16(file)
	fmt.Println(bs, major, minor)

	consts := ParseConst(file)
	fmt.Println(consts)

	flags := ReadU16(file)
	name, super := Resolve(consts, ReadU16(file)), Resolve(consts, ReadU16(file))
	fmt.Println(flags, name, super)

	interfaces := ParseInterface(file, consts)
	fmt.Println(interfaces)

	fields := ParseField(file, consts)
	fmt.Println(fields)

	methods := ParseField(file, consts)
	fmt.Println(methods)

	attributes := ParseAttribute(file, consts)
	fmt.Println(attributes)

	code, localCnt := ParseCodeAndLocal(methods["add"])
	fmt.Println(ExecCode(code, localCnt, 22, 33))
}

func ExecCode(code []byte, cnt int, args ...any) any {
	stack := make([]any, 0)
	for _, op := range code {
		switch op {
		case 26: // iload_0
			stack = append(stack, args[0])
		case 27: // iload_1
			stack = append(stack, args[1])
		case 96: // iadd
			num1 := stack[0].(int)
			num2 := stack[1].(int)
			stack = make([]any, 0)
			stack = append(stack, num1+num2)
		case 172: // ireturn
			res := stack[0]
			stack = make([]any, 0)
			return res
		}
	}
	return nil
}

func ParseCodeAndLocal(field *Field) ([]byte, int) {
	data := field.Attributes["Code"].Data
	localCnt := binary.BigEndian.Uint16(data[2:4])
	codeCnt := binary.BigEndian.Uint32(data[4:8])
	return data[8 : 8+codeCnt], int(localCnt)
}

func ParseField(file *os.File, consts []*Const) map[string]*Field {
	fieldsCount := ReadU16(file)
	fields := make(map[string]*Field)
	for i := uint16(0); i < fieldsCount; i++ {
		field := &Field{
			Flags:      ReadU16(file),
			Name:       Resolve(consts, ReadU16(file)),
			Descriptor: Resolve(consts, ReadU16(file)),
			Attributes: ParseAttribute(file, consts),
		}
		fields[field.Name] = field
	}
	return fields
}

func ParseAttribute(file *os.File, consts []*Const) map[string]*Attribute {
	attributesCount := ReadU16(file)
	attributes := make(map[string]*Attribute)
	for i := 0; i < int(attributesCount); i++ {
		attribute := &Attribute{
			Name: Resolve(consts, ReadU16(file)),
			Data: ReadBytes(file, int(ReadU32(file))),
		}
		attributes[attribute.Name] = attribute
	}
	return attributes
}

func ParseInterface(file *os.File, consts []*Const) []string {
	interfaceCount := ReadU16(file)
	interfaces := make([]string, 0)
	for i := 0; i < int(interfaceCount); i++ {
		interfaces = append(interfaces, Resolve(consts, ReadU16(file)))
	}
	return interfaces
}

func ParseConst(file *os.File) []*Const {
	constCount := ReadU16(file)
	consts := make([]*Const, 0)
	for i := 1; i < int(constCount); i++ { // 从 1开始的
		item := &Const{
			Tag: ReadU8(file),
		}
		switch item.Tag {
		case 0x01:
			l := ReadU16(file)
			item.String = string(ReadBytes(file, int(l)))
		case 0x07:
			item.NameIndex = ReadU16(file)
		case 0x08:
			item.StringIndex = ReadU16(file)
		case 0x09, 0x0A:
			item.ClassIndex = ReadU16(file)
			item.NameAndTypeIndex = ReadU16(file)
		case 0x0C:
			item.NameIndex = ReadU16(file)
			item.DescIndex = ReadU16(file)
		default:
			panic(fmt.Errorf("unknown tag 0x%02x", item.Tag))
		}
		consts = append(consts, item)
	}
	return consts
}

func Resolve(consts []*Const, index uint16) string {
	if consts[index-1].Tag == 0x01 {
		return consts[index-1].String
	}
	return ""
}

func ReadU64(file *os.File) uint64 {
	bs := ReadBytes(file, 8)
	return binary.BigEndian.Uint64(bs)
}

func ReadU32(file *os.File) uint32 {
	bs := ReadBytes(file, 4)
	return binary.BigEndian.Uint32(bs)
}

func ReadU16(file *os.File) uint16 {
	bs := ReadBytes(file, 2)
	return binary.BigEndian.Uint16(bs)
}

func ReadU8(file *os.File) uint8 {
	return ReadBytes(file, 1)[0]
}

func ReadBytes(file *os.File, count int) []byte {
	bs := make([]byte, count)
	_, err := file.Read(bs)
	HandleErr(err)
	return bs
}

func HandleErr(err error) {
	if err != nil {
		panic(err)
	}
}
