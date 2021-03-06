#Stack & Heap

---

堆(heap)栈(stack)是两种数据结构.堆栈都是一种数据项按序排列的数据结构.

* 栈 FILO, 先进后出
* 堆 FIFO, 先进先出

###空间分配

* 栈:由`操作系统`自动分配释放,存放`函数的参数值`,`局部变量的值`.
* 堆:一般由`程序员`分配释放,若程序员不释放,程序结束时可能有OS回收.

###缓存方式

* 栈: 使用的是`一级缓存`,他们通常都是被调用时处于存储空间,调用完毕激励释放.
* 堆: 存放在`二级缓存`中,声明周期由虚拟机的垃圾回收算法来决定.所以调用这些对象的速度相对来的低一些.

###申请方式
* 栈:由`系统`自动分配.
		
		例如:声明在函数中的一个局部变量 int b, 系统自动在栈中为b开辟空间.
		
* 堆:需要`程序员`自动申请,并指明大小.

###申请后系统的响应
* 栈: 只要栈的剩余空间大于所申请空间,系统将为程序提供内存,否则将报错提示栈溢出.
* 堆: 操作系统有一个记录空闲内存地址的链表,当系统收到程序的申请时,会遍历该链表,寻找第一个空间大于所申请空间的堆结点,然后将该结点从空闲结点链表中删除,并将该结点
的空间分配给程序.另外,对于大多数系统,会在这块内存空间中的首地址处记录本次分配的大小.
这样,代码中的delete语句才能正确的释放本内存空间.另外,由于找到的堆结点的大小不一定正
好等于申请的大小,系统会自动的将多余的那部分重新放入空闲链表中.

###申请效率
* 栈: 由系统自动分配,速度较快.但是程序员无法控制
* 堆: 由`new`分配的内存,一般速度比较慢,而且容易产生内存碎片,不过用起来最方便.

###存储内容
* 栈: 在函数调用时,在大多数的C编译器中,参数是由右往左入栈的,然后是函数中的局部变量.注意静态变量是不入栈的.当本次函数调用结束后,局部变量先出栈,然后是参数,最后栈顶指针指向函数的返回地址,也就是主函数中的下一条指令的地址,程序由该点继续运行.
* 堆: 一般是在堆的头部用一个字节存放堆的大小.堆中的具体内容由程序员安排.


###JAVA中的区别

**JVM中的功能**

* 栈: 内存指令区
* 堆: 内存数据区

**存储数据**

* 栈: 基本数据类型, 指令代码,常量,对象的引用地址
* 堆: 对象实例

**保存对象实例**

保存对象实例,实际上是保存对象实例的属性值,属性的类型和对象本身的类型标记等,并不保存对象的方法(方法是指令,保存在stack中).

对象实例在heap中分配好以后,需要在stack中保存一个4字节的heap内存地址,用来定位该对象实例在heap中的位置,便于找到该对象实例.

**基本数据类型**

基本数据类型包括byte、int、char、long、float、double、boolean和short.

**总结**

		Java 的堆是一个运行时数据区,类的(对象从中分配空间.这些对象通过new、newarray、anewarray和multianewarray等指令建立,它们不需要程序代码来显式的释放.堆是由垃圾回收来负责的,堆的优势是可以动态地分配内存大小,生存期也不必事先告诉编译器,因为它是在运行时动态分配内存的,Java的垃圾收集器会自动收走这些不再使用的数据.但缺点是,由于要在运行时动态分配内存,存取速度较慢.
		
		栈的优势是,存取速度比堆要快,仅次于寄存器,栈数据可以共享.但缺点是,存在栈中的数据大小与生存期必须是确定的,缺乏灵活性.栈中主要存放一些基本类型的变量（,int, short, long, byte, float, double, boolean, char）和对象句柄.
		
###小结

* 使用栈就象我们去饭馆里吃饭,只管点菜(发出申请),付钱,和吃(使用),吃饱了就走,不必理会切菜,洗菜等准备工作和洗碗,刷锅等扫尾工作,他的好处是快捷,但是自由度小.
* 使用堆就象是自己动手做喜欢吃的菜肴,比较麻烦,但是比较符合自己的口味,而且自由度大.



 

 




