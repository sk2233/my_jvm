/*
@author: sk
@date: 2024/12/28
*/
package main

// https://github.com/ArosyW/JVM
// 对于 panic 的 OpCode 可以直接在这里搜代码
// https://docs.oracle.com/javase/specs/jvms/se16/html/jvms-6.html

func main() {
	jvm := NewJVM()
	class := jvm.LoadClass("github/HelloJVM.class")
	method := class.GetMethod("main", "([Ljava/lang/String;)V")
	jvm.CallStaticMethod(class, method)
}
