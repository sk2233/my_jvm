/*
@author: sk
@date: 2024/12/29
*/
package main

import "strings"

// 对于 panic 的 OpCode 可以直接在这里搜代码
// https://docs.oracle.com/javase/specs/jvms/se16/html/jvms-6.html

func main() {
	//if len(os.Args) < 2 {
	//	fmt.Println("usage: book <class> <args...>")
	//	return
	//}
	//class := os.Args[1]
	//args := os.Args[2:]
	//Run(class, args...)
	// 临时测试使用
	Run("ExceptionTest")
}

// 默认使用当前路径作为类搜索路径
func Run(className string, args ...string) {
	className = strings.ReplaceAll(className, ".", "/")
	loader := NewLoader("/Users/bytedance/Library/Java/JavaVirtualMachines/corretto-1.8.0_352/Contents/Home/jre/lib")
	class0 := loader.LoadClass(className) // 静态方法没有调用，这里拿不到 thread
	InitInstruction()
	InitNativeFunc()
	RunMain(class0, loader, args)
}
