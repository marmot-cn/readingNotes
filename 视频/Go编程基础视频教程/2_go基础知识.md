#2 Go基础知识

[地址](http://edu.51cto.com/lesson/id-32302.html "地址")

###笔记

---

**Go程序结构**

* Go 程序是通过 `package` 来组织的.
* 只有 `package` 名称为 `main` 的包可以包含 `main` 函数. (main函数作为程序入口点启动)
* 一个可执行程序`有且仅有`一个 `main` 包.

* 通过 `import` 关键字老导入其他非 `main` 包 (有且仅有一个main函数).
* 通过 `const` 关键字来进行`常量`定义.
* 通过在函数体外部使用 `var` 关键字来进行 `全局变量`(在整个包中使用) 的声明和赋值.
* 通过 `type` 关键字来进行 结构(`struct`) 或 接口(`interface`)的声明
* 通过 `func` 关键字来进行函数的声明.

**Go导入package的格式**

		import "fmt"
		import "os"
		import "io"
		import "time"
		
		等同于:
		
		import(
			"fmt"
			"os"
			"io"
			"time"
		)

* 导入之后, 可以使用 <PackageName>.<FuncName> 来对包中的函数进行调用.
* 如果导入包之后, `未调用` 会编译错误.

**package别名**

		import io "fmt"
		
		import(
			io "fmt"
		)
		
		调用:
		
		io.Println("Hello world")

**package省略调用**

		import . "fmt"
		
		调用:
		
		Println("Hello world")


**可见性规则**

使用`大小写`来决定该`常量`,`变量`,`类型`,`接口`,`结构`或`函数` 是否可以被外部包所调用.


`函数名`首字母`小写`等于`private`.(外部包不能调用, 内部可以调用)

		func getField(...)
		
`函数名`首字母`大写`等于`public`

		func Printf(...)

###整理知识点

---
