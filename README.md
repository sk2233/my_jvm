# 使用Go实现一个简单的JVM
> 注意：仅用于学习，按需实现 jvm 指令，对于没有实现的指令会报错
## 支持功能
- class 文件解析
- jvm 虚拟机运行时
- jvm 常见指令
- 类加载器与类初始化
- 对象创建与初始化
- 方法与本地方法调用
- 数组与字符串常量池
- 异常捕获与处理
## 参考资料
jvm 指令集：https://docs.oracle.com/javase/specs/jvms/se16/html/jvms-6.html<br>
https://zserge.com/posts/jvm/<br>
https://github.com/ArosyW/JVM<br>
[自己动手写Java虚拟机.pdf](book%2Fnote%2F%E8%87%AA%E5%B7%B1%E5%8A%A8%E6%89%8B%E5%86%99Java%E8%99%9A%E6%8B%9F%E6%9C%BA.pdf)
## 示例
[ArrayDemo.java](ArrayDemo.java)<br>
[BubbleSortTest.java](BubbleSortTest.java)<br>
[ExceptionTest.java](ExceptionTest.java)<br>
[FibonacciTest.java](FibonacciTest.java)<br>
[GaussTest.java](GaussTest.java)<br>
[GetClassTest.java](GetClassTest.java)<br>
[HelloWorld.java](HelloWorld.java)<br>
[InvokeDemo.java](InvokeDemo.java)<br>
[MyObject.java](MyObject.java)<br>
[ObjectTest.java](ObjectTest.java)<br>
[StringTest.java](StringTest.java)<br>
```shell
pc:0 opcode:3 IConst0 ExceptionTest.main:4
pc:1 opcode:b8 InvokeStatic ExceptionTest.main:5
pc:0 opcode:1a Load0 ExceptionTest.test:12
pc:1 opcode:9a IfNe ExceptionTest.test:13
pc:4 opcode:bb New ExceptionTest.test:13
pc:7 opcode:59 Dup ExceptionTest.test:15
pc:8 opcode:12 Ldc ExceptionTest.test:15
pc:10 opcode:b7 InvokeSpecial ExceptionTest.test:15
pc:0 opcode:2a Load0 java/lang/IllegalArgumentException.<init>:52
pc:1 opcode:2b Load1 java/lang/IllegalArgumentException.<init>:53
pc:2 opcode:b7 InvokeSpecial java/lang/IllegalArgumentException.<init>:53
pc:0 opcode:2a Load0 java/lang/RuntimeException.<init>:62
pc:1 opcode:2b Load1 java/lang/RuntimeException.<init>:63
pc:2 opcode:b7 InvokeSpecial java/lang/RuntimeException.<init>:63
pc:0 opcode:2a Load0 java/lang/Exception.<init>:66
pc:1 opcode:2b Load1 java/lang/Exception.<init>:67
pc:2 opcode:b7 InvokeSpecial java/lang/Exception.<init>:67
pc:0 opcode:2a Load0 java/lang/Throwable.<init>:265
pc:1 opcode:b7 InvokeSpecial java/lang/Throwable.<init>:198
pc:0 opcode:b1 Return java/lang/Object.<init>:37
pc:4 opcode:2a Load0 java/lang/Throwable.<init>:198
pc:5 opcode:2a Load0 java/lang/Throwable.<init>:211
pc:6 opcode:b5 PutField java/lang/Throwable.<init>:211
pc:9 opcode:2a Load0 java/lang/Throwable.<init>:211
pc:10 opcode:b2 GetStatic java/lang/Throwable.<init>:228
pc:13 opcode:b5 PutField java/lang/Throwable.<init>:228
pc:16 opcode:2a Load0 java/lang/Throwable.<init>:228
pc:17 opcode:b2 GetStatic java/lang/Throwable.<init>:266
pc:20 opcode:b5 PutField java/lang/Throwable.<init>:266
pc:23 opcode:2a Load0 java/lang/Throwable.<init>:266
pc:24 opcode:b6 InvokeVirtual java/lang/Throwable.<init>:267
pc:0 opcode:2a Load0 java/lang/Throwable.fillInStackTrace:782
pc:1 opcode:b4 GetField java/lang/Throwable.fillInStackTrace:784
pc:4 opcode:c7 IfNonNull java/lang/Throwable.fillInStackTrace:784
pc:7 opcode:2a Load0 java/lang/Throwable.fillInStackTrace:784
pc:8 opcode:b4 GetField java/lang/Throwable.fillInStackTrace:784
pc:11 opcode:c6 IfNull java/lang/Throwable.fillInStackTrace:784
pc:27 opcode:2a Load0 java/lang/Throwable.fillInStackTrace:787
pc:28 opcode:b0 Return1 java/lang/Throwable.fillInStackTrace:0
pc:27 opcode:57 Pop java/lang/Throwable.<init>:267
pc:28 opcode:2a Load0 java/lang/Throwable.<init>:267
pc:29 opcode:2b Load1 java/lang/Throwable.<init>:268
pc:30 opcode:b5 PutField java/lang/Throwable.<init>:268
pc:33 opcode:b1 Return java/lang/Throwable.<init>:268
pc:5 opcode:b1 Return java/lang/Exception.<init>:67
pc:5 opcode:b1 Return java/lang/RuntimeException.<init>:63
pc:5 opcode:b1 Return java/lang/IllegalArgumentException.<init>:53
pc:13 opcode:bf AThrow ExceptionTest.test:15
pc:54 opcode:4c Store1 ExceptionTest.test:21
pc:55 opcode:b2 GetStatic ExceptionTest.test:22
pc:58 opcode:12 Ldc ExceptionTest.test:28
pc:60 opcode:b6 InvokeVirtual ExceptionTest.test:28
catch IllegalArgumentException
pc:63 opcode:b2 GetStatic ExceptionTest.test:28
pc:66 opcode:1a Load0 ExceptionTest.test:29
pc:67 opcode:b6 InvokeVirtual ExceptionTest.test:29
0
pc:70 opcode:a7 GoTo ExceptionTest.test:29
pc:121 opcode:b1 Return ExceptionTest.test:30
pc:4 opcode:4 IConst1 ExceptionTest.main:5
pc:5 opcode:b8 InvokeStatic ExceptionTest.main:6
pc:0 opcode:1a Load0 ExceptionTest.test:12
pc:1 opcode:9a IfNe ExceptionTest.test:13
pc:14 opcode:1a Load0 ExceptionTest.test:15
pc:15 opcode:4 IConst1 ExceptionTest.test:16
pc:16 opcode:a0 IfICmpNe ExceptionTest.test:16
pc:19 opcode:bb New ExceptionTest.test:16
pc:22 opcode:59 Dup ExceptionTest.test:18
pc:23 opcode:12 Ldc ExceptionTest.test:18
pc:25 opcode:b7 InvokeSpecial ExceptionTest.test:18
pc:0 opcode:2a Load0 java/lang/RuntimeException.<init>:62
pc:1 opcode:2b Load1 java/lang/RuntimeException.<init>:63
pc:2 opcode:b7 InvokeSpecial java/lang/RuntimeException.<init>:63
pc:0 opcode:2a Load0 java/lang/Exception.<init>:66
pc:1 opcode:2b Load1 java/lang/Exception.<init>:67
pc:2 opcode:b7 InvokeSpecial java/lang/Exception.<init>:67
pc:0 opcode:2a Load0 java/lang/Throwable.<init>:265
pc:1 opcode:b7 InvokeSpecial java/lang/Throwable.<init>:198
pc:0 opcode:b1 Return java/lang/Object.<init>:37
pc:4 opcode:2a Load0 java/lang/Throwable.<init>:198
pc:5 opcode:2a Load0 java/lang/Throwable.<init>:211
pc:6 opcode:b5 PutField java/lang/Throwable.<init>:211
pc:9 opcode:2a Load0 java/lang/Throwable.<init>:211
pc:10 opcode:b2 GetStatic java/lang/Throwable.<init>:228
pc:13 opcode:b5 PutField java/lang/Throwable.<init>:228
pc:16 opcode:2a Load0 java/lang/Throwable.<init>:228
pc:17 opcode:b2 GetStatic java/lang/Throwable.<init>:266
pc:20 opcode:b5 PutField java/lang/Throwable.<init>:266
pc:23 opcode:2a Load0 java/lang/Throwable.<init>:266
pc:24 opcode:b6 InvokeVirtual java/lang/Throwable.<init>:267
pc:0 opcode:2a Load0 java/lang/Throwable.fillInStackTrace:782
pc:1 opcode:b4 GetField java/lang/Throwable.fillInStackTrace:784
pc:4 opcode:c7 IfNonNull java/lang/Throwable.fillInStackTrace:784
pc:7 opcode:2a Load0 java/lang/Throwable.fillInStackTrace:784
pc:8 opcode:b4 GetField java/lang/Throwable.fillInStackTrace:784
pc:11 opcode:c6 IfNull java/lang/Throwable.fillInStackTrace:784
pc:27 opcode:2a Load0 java/lang/Throwable.fillInStackTrace:787
pc:28 opcode:b0 Return1 java/lang/Throwable.fillInStackTrace:0
pc:27 opcode:57 Pop java/lang/Throwable.<init>:267
pc:28 opcode:2a Load0 java/lang/Throwable.<init>:267
pc:29 opcode:2b Load1 java/lang/Throwable.<init>:268
pc:30 opcode:b5 PutField java/lang/Throwable.<init>:268
pc:33 opcode:b1 Return java/lang/Throwable.<init>:268
pc:5 opcode:b1 Return java/lang/Exception.<init>:67
pc:5 opcode:b1 Return java/lang/RuntimeException.<init>:63
pc:28 opcode:bf AThrow ExceptionTest.test:18
pc:73 opcode:4c Store1 ExceptionTest.test:23
pc:74 opcode:b2 GetStatic ExceptionTest.test:24
pc:77 opcode:12 Ldc ExceptionTest.test:28
pc:79 opcode:b6 InvokeVirtual ExceptionTest.test:28
catch RuntimeException
pc:82 opcode:b2 GetStatic ExceptionTest.test:28
pc:85 opcode:1a Load0 ExceptionTest.test:29
pc:86 opcode:b6 InvokeVirtual ExceptionTest.test:29
1
pc:89 opcode:a7 GoTo ExceptionTest.test:29
pc:121 opcode:b1 Return ExceptionTest.test:30
pc:8 opcode:5 IConst2 ExceptionTest.main:6
pc:9 opcode:b8 InvokeStatic ExceptionTest.main:7
pc:0 opcode:1a Load0 ExceptionTest.test:12
pc:1 opcode:9a IfNe ExceptionTest.test:13
pc:14 opcode:1a Load0 ExceptionTest.test:15
pc:15 opcode:4 IConst1 ExceptionTest.test:16
pc:16 opcode:a0 IfICmpNe ExceptionTest.test:16
pc:29 opcode:1a Load0 ExceptionTest.test:18
pc:30 opcode:5 IConst2 ExceptionTest.test:19
pc:31 opcode:a0 IfICmpNe ExceptionTest.test:19
pc:34 opcode:bb New ExceptionTest.test:19
pc:37 opcode:59 Dup ExceptionTest.test:28
pc:38 opcode:12 Ldc ExceptionTest.test:28
pc:40 opcode:b7 InvokeSpecial ExceptionTest.test:28
pc:0 opcode:2a Load0 java/lang/Exception.<init>:66
pc:1 opcode:2b Load1 java/lang/Exception.<init>:67
pc:2 opcode:b7 InvokeSpecial java/lang/Exception.<init>:67
pc:0 opcode:2a Load0 java/lang/Throwable.<init>:265
pc:1 opcode:b7 InvokeSpecial java/lang/Throwable.<init>:198
pc:0 opcode:b1 Return java/lang/Object.<init>:37
pc:4 opcode:2a Load0 java/lang/Throwable.<init>:198
pc:5 opcode:2a Load0 java/lang/Throwable.<init>:211
pc:6 opcode:b5 PutField java/lang/Throwable.<init>:211
pc:9 opcode:2a Load0 java/lang/Throwable.<init>:211
pc:10 opcode:b2 GetStatic java/lang/Throwable.<init>:228
pc:13 opcode:b5 PutField java/lang/Throwable.<init>:228
pc:16 opcode:2a Load0 java/lang/Throwable.<init>:228
pc:17 opcode:b2 GetStatic java/lang/Throwable.<init>:266
pc:20 opcode:b5 PutField java/lang/Throwable.<init>:266
pc:23 opcode:2a Load0 java/lang/Throwable.<init>:266
pc:24 opcode:b6 InvokeVirtual java/lang/Throwable.<init>:267
pc:0 opcode:2a Load0 java/lang/Throwable.fillInStackTrace:782
pc:1 opcode:b4 GetField java/lang/Throwable.fillInStackTrace:784
pc:4 opcode:c7 IfNonNull java/lang/Throwable.fillInStackTrace:784
pc:7 opcode:2a Load0 java/lang/Throwable.fillInStackTrace:784
pc:8 opcode:b4 GetField java/lang/Throwable.fillInStackTrace:784
pc:11 opcode:c6 IfNull java/lang/Throwable.fillInStackTrace:784
pc:27 opcode:2a Load0 java/lang/Throwable.fillInStackTrace:787
pc:28 opcode:b0 Return1 java/lang/Throwable.fillInStackTrace:0
pc:27 opcode:57 Pop java/lang/Throwable.<init>:267
pc:28 opcode:2a Load0 java/lang/Throwable.<init>:267
pc:29 opcode:2b Load1 java/lang/Throwable.<init>:268
pc:30 opcode:b5 PutField java/lang/Throwable.<init>:268
pc:33 opcode:b1 Return java/lang/Throwable.<init>:268
pc:5 opcode:b1 Return java/lang/Exception.<init>:67
pc:43 opcode:bf AThrow ExceptionTest.test:28
pc:92 opcode:4c Store1 ExceptionTest.test:25
pc:93 opcode:b2 GetStatic ExceptionTest.test:26
pc:96 opcode:12 Ldc ExceptionTest.test:28
pc:98 opcode:b6 InvokeVirtual ExceptionTest.test:28
catch Exception
pc:101 opcode:b2 GetStatic ExceptionTest.test:28
pc:104 opcode:1a Load0 ExceptionTest.test:29
pc:105 opcode:b6 InvokeVirtual ExceptionTest.test:29
2
pc:108 opcode:a7 GoTo ExceptionTest.test:29
pc:121 opcode:b1 Return ExceptionTest.test:30
pc:12 opcode:6 IConst3 ExceptionTest.main:7
pc:13 opcode:b8 InvokeStatic ExceptionTest.main:8
pc:0 opcode:1a Load0 ExceptionTest.test:12
pc:1 opcode:9a IfNe ExceptionTest.test:13
pc:14 opcode:1a Load0 ExceptionTest.test:15
pc:15 opcode:4 IConst1 ExceptionTest.test:16
pc:16 opcode:a0 IfICmpNe ExceptionTest.test:16
pc:29 opcode:1a Load0 ExceptionTest.test:18
pc:30 opcode:5 IConst2 ExceptionTest.test:19
pc:31 opcode:a0 IfICmpNe ExceptionTest.test:19
pc:44 opcode:b2 GetStatic ExceptionTest.test:28
pc:47 opcode:1a Load0 ExceptionTest.test:29
pc:48 opcode:b6 InvokeVirtual ExceptionTest.test:29
3
pc:51 opcode:a7 GoTo ExceptionTest.test:29
pc:121 opcode:b1 Return ExceptionTest.test:30
pc:16 opcode:b1 Return ExceptionTest.main:8
```