#10 结构struct

[地址](http://edu.51cto.com/lesson/id-32317.html "地址")

###笔记

---

**结构体**

* Go 中的`struct`与 C 中的`struct`非常相似, 但 GO 没有`class`
* `type <Name> struct{}`

		type person struct{
			Name string
			Age int
		}
		
		func main(){
			a := person{}
			a.Name = "joe"
			a.Age = 19
			
			//等同于:
			//a := person{
			//	Name: "joe",
			//	Age: 19,
			//}
			//
			fmt.Println(a)
		}
		
		输出:
		{joe 19}
		
* 支持指向自身的`指针`类型成员
* 支持`匿名`结构, 可用作成员或定义成员变量

		a := struct {
			Name string
			Age int
		}{
			Name: "joe",
			Age: 19,
		}
		fmt.Println(a)

* 匿名结构也可以用于`map`的值
* 可以使用字面值对结构进行初始化
* 允许直接通过指针来读写结构成员
* 相同类型的成员可进行直接拷贝赋值

		type person struct{
			Name string
			Age int
		}
		
		func main(){
			a := person{Name:"joe", Age:19}
			var b person
			b = a
		}

* 相同类型支持 `==` 与 `!=` 比较运算符, 但不支持 `>` 或 `<`

		type person struct{
			Name string
			Age int
		}
		
		type person1 struct{
			Name string
			Age int
		}
		
		func main(){
			a := person{Name:"joe", Age:19}
			b := person1{Name:"joe", Age:19}				fmt.Prontln(b == a)//报错,person 和 person1 不同类型
			c : = person{Name:"joe", Age:19}
			fmt.Prontln(c == a)//输出 true
		}

* 支持匿名字段, 本质是定义了某个类型名为名称的字段

		type person struct{
			Name string
			Age int
		}
		
		func main(){
			//必须和声明类型一致
			a := person{"joe",19}
			//a := person{19,"joe"} 会报错
			fmt.Println(a)			
		}

* 嵌入结构作为匿名字段看起来像`继承`,但`不是继承`(内部的字段默认给了外部的字段)

		type human struct {
			Sex int
		}
		
		type teacher struct {
			human //嵌入结构作为匿名字段
			Name string
			Age int
		}
		
		func main(){
			a := teacher{Name: "joe", Age: 19,human: human{Sex:0}}
			a.Name = "xxx"//操作Name
			a.Age = 13//操作Age
			a.human.Sex = 1;//操作sex
			a.Sex = 1;//操作sex
		}

* 可以使用`匿名字段指针`

**结构是值拷贝类型**

		type person struct{
			Name string
			Age int
		}
		
		func main(){
			a := person{
				Name: "joe",
				Age: 19,
			}
			fmt.Println(a)
			A(a)
			fmt.Println(a)
		}
		
		func A(per person){
			per.Age = 13
			fmt.Println("A",per)
		}
		
		输出:
		{joe 19}
		A {joe 13}
		{joe 19}
		
`修改为值引用传递`(用指针):

		...
		A(&a)
		...
		...
		func A(per *person){
			...
		}
		...		
		
		或者 初始化为指针:
		
		...
		//a 为指向结构的指针
		a := &person{
				Name: "joe",
				Age: 19,
		}
		a.Name = "xxx" //不用加 * 
		...
		A(a)//不用加取地址符&
		...
		func A(per *person){
			...
		}

**结构体嵌套**

		type person struct {
			Name string
			Age  int
			Contact struct {
				Phone, City string
			}
		}
		
		func main(){
			a := person{Name:"joe", Age:19}
			a.Contact.city = "beijing"
		}
		
###整理知识点

---