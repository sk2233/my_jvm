/*
@author: sk
@date: 2024/12/29
*/
package main

const (
	AccessPublic       = 0x0001 // class field method
	AccessPrivate      = 0x0002 //       field method
	AccessProtected    = 0x0004 //       field method
	AccessStatic       = 0x0008 //       field method
	AccessFinal        = 0x0010 // class field method
	AccessSuper        = 0x0020 // class
	AccessSynchronized = 0x0020 //             method
	AccessVolatile     = 0x0040 //       field
	AccessBridge       = 0x0040 //             method
	AccessTransient    = 0x0080 //       field
	AccessVarargs      = 0x0080 //             method
	AccessNative       = 0x0100 //             method
	AccessInterface    = 0x0200 // class
	AccessAbstract     = 0x0400 // class       method
	AccessStrict       = 0x0800 //             method
	AccessSynthetic    = 0x1000 // class field method
	AccessAnnotation   = 0x2000 // class
	AccessEnum         = 0x4000 // class field
)

type Class struct {
	// 静态读取出来的
	Magic        uint32
	Major, Minor uint16   // 主次版本号
	Consts       []*Const // 常量池
	Access       uint16   // 访问权限
	ThisIndex    uint16
	SupperIndex  uint16
	Interfaces   []uint16 // 实现的接口索引表
	Fields       []*Field
	Methods      []*Field
	Attributes   []*Attribute
	// 动态后来添加的
	InstSlotCount   int
	StaticSlotCount int
	StaticValues    []*Value
}

// 还没有考虑继承
func (c *Class) GetMethod(name string, desc string) *Field {
	for _, method := range c.Methods {
		if c.GetString(method.NameIndex) == name && c.GetString(method.DescIndex) == desc {
			return method
		}
	}
	return nil
}

// 还没有考虑继承
func (c *Class) GetField(name string, desc string) *Field {
	for _, field := range c.Fields {
		if c.GetString(field.NameIndex) == name && c.GetString(field.DescIndex) == desc {
			return field
		}
	}
	return nil
}

func (c *Class) GetString(index uint16) string {
	temp := c.Consts[index]
	if temp.Type == ConstClass || temp.Type == ConstString {
		return c.GetString(temp.Index)
	}
	return temp.String
}

type Field struct {
	Access     uint16
	NameIndex  uint16
	DescIndex  uint16
	Attributes []*Attribute
	// 后面添加的非 class 文件中
	Class  *Class
	SlotID int
}

func (f *Field) GetCodeAttribute() *Code {
	for _, attr := range f.Attributes {
		if attr.Name == AttributeCode {
			return attr.Code
		}
	}
	return nil
}

func (f *Field) GetConstantValueAttribute() uint16 {
	for _, attr := range f.Attributes {
		if attr.Name == AttributeConstantValue {
			return attr.ConstantValueIndex
		}
	}
	return 0
}

func (f *Field) IsTwoSlot() bool {
	desc := f.Class.GetString(f.DescIndex)
	return desc == "J" || desc == "D"
}

func IsStatic(access uint16) bool {
	return access&AccessStatic > 0
}

func IsFinal(access uint16) bool {
	return access&AccessFinal > 0
}

func IsInterface(access uint16) bool {
	return access&AccessInterface > 0
}

func IsAbstract(access uint16) bool {
	return access&AccessAbstract > 0
}

func IsNative(access uint16) bool {
	return access&AccessNative > 0
}

const (
	AttributeCode            = "Code"
	AttributeSourceFile      = "SourceFile"
	AttributeExceptions      = "Exceptions"
	AttributeLineNumberTable = "LineNumberTable"
	AttributeConstantValue   = "ConstantValue"
)

type Attribute struct {
	Name               string
	Data               []byte
	Code               *Code
	SourceFileIndex    uint16
	ExceptionIndexes   []uint16
	LineNumbers        []*LineNumber
	ConstantValueIndex uint16
}

type Code struct { // 解析出来最好不要是裸信息，还是尽可能转换为其包装信息为好
	MaxStack   uint16
	MaxLocal   uint16
	Code       []byte
	Exceptions []*Exception
	Attributes []*Attribute
}

func (c *Code) FindException(class *Class, pc uint16, obj *Object) *Exception {
	for _, item := range c.Exceptions {
		if item.Start > pc || item.End < pc { // 处于范围内
			continue
		} // 这里需要判定处理的类型兼容 允许继承 这里简单判断必须相等
		if obj.Class.GetString(obj.Class.ThisIndex) == class.GetString(item.CatchType) {
			return item
		}
	}
	return nil
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
	// 管理地址
	Start uint16
	End   uint16
	// 处理地址
	Handler uint16
	// 处理类型 指定常量池 ClassRef    0的话 catch all
	CatchType uint16
}

const (
	// 不完整，缺了再说  暂时不支持 long double
	ConstUtf8            = 1 // 直接文本
	ConstInteger         = 3 // 直接数字
	ConstFloat           = 4 // 直接数字
	ConstLong            = 5
	ConstDouble          = 6
	ConstClass           = 7  // 引用类名
	ConstString          = 8  // 引用字符串
	ConstField           = 9  // 类 + 字段(ConstNameType)
	ConstMethod          = 10 //
	ConstInterfaceMethod = 11 //
	ConstNameType        = 12 // Name + Type
	ConstMethodHandle    = 15
	ConstMethodType      = 16
	ConstInvokeDynamic   = 18
)

type Const struct {
	Type uint8
	// ConstClass, ConstString
	Index uint16
	// ConstUtf8
	String string
	// ConstInteger
	Integer int32
	// ConstFloat
	Float float32
	// ConstLong
	Long int64
	// ConstDouble
	Double float64
	// ConstField, ConstMethod, ConstInterfaceMethod
	ClassIndex    uint16
	NameTypeIndex uint16
	// ConstNameType
	NameIndex uint16
	DescIndex uint16
}
