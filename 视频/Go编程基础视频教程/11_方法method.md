#11 方法method

[地址](http://edu.51cto.com/lesson/id-32318.html "地址")

###笔记

---

**method**

* GO 中没有class, 但是有 `method`
* 通过显示说明`receiver`来实现与某个类型的组合,编译器根据`接受者的类型`来判断是哪个结构的方法
		
		type A struct {
			Name string
		}
		
		func main() {
			a := A{}
			a.print()
		}
		
		func (a A)print(){
			fmt.Println("A")
		}	
		
* 只能为同一个包中的类型定义方法
* `receiver`可以是类型的值或者指针
		
		//值拷贝传递
		func (a A)print(){
			fmt.Println("A")
		}
		或者
		//值引用类型传递
		func (a *A)print(){
			fmt.Println("A")
		}		

* 不存在方法重载
* 可以使用值或指针来调用方法,编译器会自动完成转换
* 方法是函数的语法糖, 因为`receiver`其实就是方法锁接受的第一个参数 (Method Value vs Method Expression)

		Method Value:
		...
		a.Print()
		...
		
		Method Expression:
		...
		(*A).Print(&a)
		
* 如果外部结构和嵌入结构存在同名方法, 则优先调用外部结构的方法
* 类型别名不会拥有`底层`类型所附带的方法
* 方法可以调用结构中的非公开(私有字段)字段(小写:私有. 大写:共有),访问权限较高

**类型断言**

* 通过类型断言的ok pattern可以判断接口中的数据类型
* 使用type switch则可针对空接口进行比较全面的类型判断

**接口转换**

可以将拥有超集的接口转换为子集的接口


###整理知识点

---

**私有字段,共有字段**

* `私有字段`: 包内访问
* `共有字段`: 包外访问
