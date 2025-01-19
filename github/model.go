/*
@author: sk
@date: 2024/12/28
*/
package main

const (
	AccessPublic     = 0x0001
	AccessFinal      = 0x0010
	AccessSuper      = 0x0020
	AccessInterface  = 0x0200
	AccessAbstract   = 0x0300
	AccessSynthetic  = 0x1000
	AccessAnnotation = 0x2000
	AccessEnum       = 0x4000
)

type Class struct {
	Magic        uint32
	Major, Minor uint16   // 主次版本号
	Consts       []*Const // 常量池
	Access       uint16   // 访问权限
	Name         string
	Supper       string
	Interfaces   []*Interface
	Fields       []*Field
	Methods      []*Field
	Attributes   []*Attribute
}

func (c *Class) GetMethod(name string, desc string) *Field {
	for _, method := range c.Methods {
		if method.Name == name && method.Desc == desc {
			return method
		}
	}
	return nil
}

type Field struct {
	Access     uint16
	Name       string
	Desc       string
	Attributes []*Attribute
}

func (f *Field) GetCode() *Code {
	for _, attribute := range f.Attributes {
		if attribute.Name == AttributeCode {
			return attribute.Code
		}
	}
	return nil
}

const (
	AttributeCode       = "Code"
	AttributeSourceFile = "SourceFile"
)

type Attribute struct {
	Name       string
	Data       []byte
	Code       *Code
	SourceFile string
}

type Code struct {
	MaxStackDepth uint16
	MaxLocals     uint16
	Code          []byte
	Exceptions    []*Exception
	ExtAttrs      []*ExtAttr
}

const (
	ExtAttrLineNumberTable    = "LineNumberTable"
	ExtAttrLocalVariableTable = "LocalVariableTable"
)

type ExtAttr struct {
	Name           string
	Data           []byte
	LineNumbers    []*LineNumber
	LocalVariables []*LocalVariable
}

type LineNumber struct {
	Start uint16
	Line  uint16
}

type LocalVariable struct {
	Start uint16
	Len   uint16
	Name  string
	Type  uint16 // 常量池索引
	Index uint16
}

type Exception struct {
	Start     uint16
	End       uint16
	Handler   uint16
	CatchType uint16
}

type Interface struct {
	Name string
}

type ConstType int8

const (
	ConstUtf8            ConstType = 1  // 直接文本
	ConstInt             ConstType = 3  // 直接数字
	ConstFloat           ConstType = 4  // 直接数字
	ConstClass           ConstType = 7  // 引用类名
	ConstString          ConstType = 8  // 引用字符串
	ConstField           ConstType = 9  // 类 + 字段(ConstNameType)
	ConstMethod          ConstType = 10 //
	ConstInterfaceMethod ConstType = 11 //
	ConstNameType        ConstType = 12 // Name + Type
)

type Const struct {
	Type     ConstType
	Index    uint16
	ExtIndex uint16
	String   string
	Num      uint32
}
