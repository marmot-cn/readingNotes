#9 函数function

[地址](http://edu.51cto.com/lesson/id-32316.html "地址")

###笔记

---

**函数function**

* Go 函数 `不支持` `嵌套`(nested),`重载`(overload)和`默认参数`
* `支持`:
	* 无序声明原型
	* 不定长度变参
	* 多返回值
	* 命名返回值参数
	* 匿名函数
	* 闭包

* `func` 关键字定义函数
* 函数可以作为一种类型使用(go 语言一切皆类型)

		func main() {
			a := A
			a() //函数作为类型使用
		}
		
		func A() {
			fmt.Println("Func A")
		}

**示例**

		func A(a int,b string) (int,string){
			
		}
		
		func A(a int,b int,c int){
			...
		}
		等价于
		func A(a,b,c int){
			...
		}
		同理用于返回值:
		func A() (a,b,c int){
			...
		}
		
		function A()(int,int,int){
			a,b,c := 1,2,3
			return a,b,c
		}
		等价于
		function A()(a,b,c int){
			a,b,c = 1,2,3 (注意没有冒号,因为a,b,c在前面已经定义了,已经分配好内存地址了)
			return 
		}
		但是为了可读性,一般都这样写
		...
		return a,b,c
		...
		
`不定长变参`:

		fun A(b string,a ...int){// (a ...int, b string) 这样不行
			//a变为slice,在接受一系列参数后,得到是"值copy"		}
		
`slice传参`

* 传递`地址拷贝`:
		
		fun A(a ...int){
			//a变为slice,在接受一系列参数后,得到是"值copy"		}

* 传递`地址拷贝`:(不是传递指针)

		fun A(s []int){
			//s为值引用
		}

	`slice`,`map`,`channe` 只有这三种类型是拷贝内存地址.

**匿名函数**

		func main(){
			
			a := func(){
				fmt.Println("Func A")
			}
			a()
		}
		
**闭包**

		func main(){
			f := closure(10)//f代表函数
			fmt.Println(f(1))
			fmt.Println(f(2))
		}

		func closure(x int) func(int) int {
			fmt.Printf("%p",&x)//打印x地址
			return func(y int) int {
				fmt.Printf("%p",&x)
				return x + y
			}
		}
	
		输出:
		//3次x地址相同,指向同一个x,不是值的拷贝
		//变量如果不是传参传递,则都是值引用而非值拷贝
		0xc200000038
		0xc200000038
		11
		0xc200000038
		12
	
**defer**

* 执行方式类似其他语言中的`析构`函数, 在函数体执行结束后按照调用顺序的`相反顺序`逐个执行(先进后出, 后出先进)
* 即时函数发生`严重错误`也会执行
* 支持`匿名`函数的调用
* 常用于`资源清理`,`文件关闭`,`解锁`以及`记录时间`等操作.
* 通过与匿名函数配合可在`return`之后`修改`函数计算结果.
* 如果函数体内某个变量作为defer时匿名函数的参数,则在定义defer时即已经获得了拷贝,否则则是引用某个变量的地址(闭包).

* `Go` 没有异常机制, 担忧 `painc` `recover`模式来处理错误.
* `Panic` 可以在任何地方引发, 但 `recover` `只有`在 `defer` 调用的函数中有效.

`示例`:

		func main(){
			fmt.Println("a")
			defer fmt.Println("b")
			defer fmt.Println("c")
		}
		输出:
		a
		c
		b
		
		...
		for i:=0; i < 3; i++ {
			defer fmt.Println(i)//值拷贝,所以输出2,1,0
		}
		输出:
		2
		1
		0
	
`匿名函数调用`:
		
		func main(){
			for i:=0; i < 3; i++ {
				defer func() {
					fmt.Println(i) //传递的是地址
				}()//后面这个括号代表调用这个函数
			}
		}	
		输出:
		3
		3
		3
	
**painc**

		func main() {
			A()
			B() //输出到B,因为painc,程序终止
			C()
		}	
		
		func A(){
			fmt.Println("Func A")
		}
		
		func B(){
			panic("Paninc in B")
		}
		
		func C(){
			fmt.Println("Func C")
		}
		
		输出:
		//输出到函数B终止
		Func A
		paint: Painc in B
		....
		
`使用recover将程序从painc恢复到正常状态`:

		//修改函数B
		func B(){
			//必须在 panic() 之前进行defer -> recover
			defer func(){
				if err:=recover();err!=nil{
					fmt.Println("Recover in B")
				}
			}()
			panic("Paninc in B")
		}
		输出:
		Func A
		Recover in B
		Func C
		
**示例**

		var fs = [4]func(){}

		for i := 0; i < 4; i++ {
			defer fmt.Println("defer i= ",i)
			defer func() {fmt.Println("defer_closure i= ",i)}()
			fs[i] = func() {fmt.Println("clousre i=",i)}
		}

		for _, f := range fs {
			f()
		}
	
		输出:
		//因为闭包i是引用地址,执行完for循环之后i=4,所以输出4
		clousre i= 4
		clousre i= 4
		clousre i= 4
		clousre i= 4
		defer_closure i=  4//闭包,值引用
		defer i=  3//i作为参数,值拷贝
		defer_closure i=  4
		defer i=  2
		defer_closure i=  4
		defer i=  1
		defer_closure i=  4
		defer i=  0
		
###整理知识点

---

**Go函数嵌套**

Go不能在函数内部显式嵌套定义函数,但是可以定义一个`匿名函数`.

**函数声明和定义**

* 函数`声明`, `返回类型 函数名 (参数类型1 参数名1,·····,参数类型n 参数名n)`;

		int fun(int a,int b)
		
		//示例
		double square ( double );  //或 double square ( double x)
		int main(void)
		{
  			printf("%f\n" , square(3.) );
    		return 0;
		}
		double square ( double x) 
		{
   			return x * x ;
		}
		
* 函数定义:

		返回类型  函数名（参数类型1  参数名1，·····，参数类型n  参数名n）{
     		函数体······
		}
		
		如:
		
		int  fun(int a,int b){  
		   int  c;
           c=a+b;
           return c;
        }
        
        //示例
        double square ( double x) 
		{
   			return x * x ;
		}
		int main(void)
		{
  			printf("%f\n" , square(3.) );
  			return 0;
		}	


**函数原型和函数声明**

函数`声明`也称函数`原型`.






