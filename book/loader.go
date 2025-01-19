/*
@author: sk
@date: 2024/12/29
*/
package main

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Loader struct {
	Paths   []string
	Classes map[string]*Class
}

func (l *Loader) LoadClass(className string) *Class { // 最好静态与运行时分开
	if _, ok := l.Classes[className]; !ok {
		if className[0] == '[' { // 加载数组
			l.Classes[className] = &Class{ // 构造数组 class
				Consts: []*Const{{}, // 第一个空着
					{Type: 7, Index: 2}, {Type: 1, String: className},
					{Type: 7, Index: 4}, {Type: 1, String: "java/lang/Object"},
					{Type: 7, Index: 6}, {Type: 1, String: "java/lang/Cloneable"}, {Type: 7, Index: 8}, {Type: 1, String: "java/io/Serializable"}},
				Access:      AccessPublic,
				ThisIndex:   1,
				SupperIndex: 3,
				Interfaces:  []uint16{5, 7}, // 因该实现序列化接口啥的
			}
		} else { // 加载普通类
			// 加载解析 class
			bs := l.LoadData(className)
			class := NewParser(bs).ParseClass()
			// 定义 class
			l.DefineClass(class)
			// 链接 class
			l.LinkClass(class)
		}
	}
	return l.Classes[className]
}

func (l *Loader) LinkClass(class *Class) {
	// TODO 校验类
	// 计算实例字段下标
	l.calcuInstSlotID(class)
	// 计算静态字段下标
	l.calcuStaticSlotID(class)
	// 为静态变量分配内存与初始化
	l.initStaticFinalField(class)
}

func (l *Loader) initStaticFinalField(class *Class) {
	class.StaticValues = make([]*Value, class.StaticSlotCount)
	for _, field := range class.Fields { // final 值直接存储在常量池 中
		if IsStatic(field.Access) && IsFinal(field.Access) {
			constantValueIndex := field.GetConstantValueAttribute()
			if constantValueIndex == 0 {
				continue
			} // 对于常量值直接赋值
			constantValue := class.Consts[constantValueIndex]
			desc := class.GetString(field.DescIndex)

			switch desc {
			case "Z", "B", "C", "S", "I":
				class.StaticValues[field.SlotID] = NewInteger(constantValue.Integer)
			case "J":
				class.StaticValues[field.SlotID] = NewLong(constantValue.Long)
			case "F":
				class.StaticValues[field.SlotID] = NewFloat(constantValue.Float)
			case "D":
				class.StaticValues[field.SlotID] = NewDouble(constantValue.Double)
			case "Ljava/lang/String;": // 字符串常量
				//class.StaticValues[field.SlotID] = NewString()
			default:
				panic(fmt.Sprintf("unknown field desc %s", desc))
			}
		}
	}
}

func (l *Loader) calcuStaticSlotID(class *Class) {
	slotID := 0
	for _, field := range class.Fields {
		if IsStatic(field.Access) {
			field.SlotID = slotID
			slotID++
			if field.IsTwoSlot() {
				slotID++
			}
		}
	}
	class.StaticSlotCount = slotID
}

func (l *Loader) calcuInstSlotID(class *Class) {
	slotID := 0
	if class.SupperIndex > 0 { // 0 是无效的  有父类实例字段下标要进行累加
		supperClass := class.GetString(class.SupperIndex)
		slotID = l.LoadClass(supperClass).InstSlotCount
	}
	for _, field := range class.Fields {
		if !IsStatic(field.Access) {
			field.SlotID = slotID
			slotID++
			if field.IsTwoSlot() {
				slotID++
			}
		}
	}
	class.InstSlotCount = slotID
}

func (l *Loader) DefineClass(class *Class) {
	className := class.GetString(class.ThisIndex)
	// 先加载父类
	if className != "java/lang/Object" {
		supperClass := class.GetString(class.SupperIndex)
		l.LoadClass(supperClass)
	}
	// 再加载接口
	for _, tempIndex := range class.Interfaces {
		tempClass := class.GetString(tempIndex)
		l.LoadClass(tempClass)
	}
	// 最后定义自己
	l.Classes[className] = class
}

func (l *Loader) LoadData(class string) []byte {
	class = class + ".class" // 转换为路径
	for _, path := range l.Paths {
		if strings.HasSuffix(path, ".jar") { // 两种加载方式
			if bs := l.loadJarData(path, class); bs != nil {
				return bs
			}
		} else {
			if bs := l.loadDirData(path, class); bs != nil {
				return bs
			}
		}
	}
	panic(fmt.Sprintf("class %s not found", class))
}

func (l *Loader) loadJarData(path string, class string) []byte {
	reader, err := zip.OpenReader(path)
	HandleErr(err)
	defer reader.Close()

	file, err := reader.Open(class)
	if err != nil { // 没有找到
		return nil
	}
	defer file.Close()
	return ReadAll(file) // 找到了
}

func (l *Loader) loadDirData(dirPath string, class string) []byte {
	file, err := os.Open(filepath.Join(dirPath, class))
	if err != nil { // 没有找到
		return nil
	}
	defer file.Close()
	return ReadAll(file) // 找到了
}

func NewLoader(path string) *Loader {
	paths := make([]string, 0)
	// 先添加基本搜索路径
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".jar") {
			paths = append(paths, path)
		}
		return err
	})
	HandleErr(err)
	// 再添加用户搜索路径
	wd, err := os.Getwd()
	HandleErr(err)
	paths = append(paths, wd)
	return &Loader{Paths: paths, Classes: make(map[string]*Class)}
}
