# 第五章 Ajax

### Ajax 和 JavaScript (前篇)

`Ajax` (Asynchronous JavaScript and XML). 异步 JavaScript 及 XML.

#### 通信及异步页面更新

异步通信是将处理细分成多个处理来交替执行.

#### 技术要素之一: JavaScript

#### 技术要素之二: XML

用`JavaScript`进行异步通信的对象的名字是 ==XMLHttpRequest==.

#### 技术要素之三: DHTML

利用装在在网页中的`JavaScript`,使用`DOM`(文档对象模型)对网页数据进行操作.使用`DOM`可以进行下述处理:

* 取得页面中特定标签中的数据.
* 修改标签的数据(文字,属性等)
* 在页面中添加标签.
* 设定事件处理程序.

#### JavaScript技术基础

`JavaScript`是以对象为基础的语言.

#### 原型模式的面向对象编程语言

以类为中心的传统面向对象编程,是以类为进出生成新对象.类和对象的关系可以类别铸模和铸件的关系.

原型模式的面向对象编程语言没有类这样一个概念.

需要生成新的对象时,只要给对象追加属性.设置函数对象作为属性的话,就成为方法.当访问对象中不存在的属性时,`JavaScript`回去搜索该对象==prototype==属性所指向的对象.

`JavaScript`利用这个功能,使用"委派"而非"继承"来实现面向对象编程.

**实例**

		// 生成Dog...A
		function Dog(){
			this.sig = function() {return "I'm sitting"}
		}
		
		// 从Dog 生成对象dog...B
		var dog = new Dog()
		// dog 是狗, 所以能 sit...C
		alert(dog.sit())
		// 生成新型myDog...D
		function MyDog()
		// 指定委派原型
		MyDog.prototype = new Dog()
		// 从MyDog 生成新对象myDog...E
		var myDog = new MyDog()
		document.write(myDog.sit())
		
`A`: `JavaScript`使用了函数对象来代替类.函数对象起到了对象构造器的作用.执行构造器的处理时,给`this`追加新的属性,从而完成对象的初始化.

`B`: 原型`Dog`调用`new`完成了以下处理:
	
1. 生成对象.
2. 将委派原型的内部属性(__proto__)设置为Dog.prototype.
3. 调用函数`Dog`,参数即为传递给`new`时的参数.
4. 返回新生成的对象.

`E`: 对由新定义原型生成的对象来说,当访问它不知道的属性时,就会为派给对象`dog`.`JavaScript`通过这种方式实现了类模式语言中相当于继承的功能.

### Ajax 和 JavaScript (后篇)

#### 和服务器间的通信

`Ajax` 是利用 `XMLHttpRequest` 对象来进项异步通信的.

#### 命名的重要性

起了一个合适的名字本身就意味着功能设计得正确.反过来,起了不好的名字说明设计者自己也没有完全理解应完成什么样的功能.
