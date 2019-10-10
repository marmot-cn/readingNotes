#javaScript的匿名函数及闭包



###匿名函数

---

`匿名函数`: 就是没有函数名的函数.

**函数的定义**

1. 常规:


		function double(x){
			return 2 * x;
		}

2. Function构造函数, 把参数列表和函数体都作为字符串:

		var double = new Function('x', 'return 2 * x;');
		
3. 匿名函数的创建:

		var double = function(x) { return 2*x; }
		
**创建匿名函数**

1. 第一种就是上面的创建方法

2. 第二种:

		(function(x,y)){
		
			alert(x + y);
		
		})(2,3);
		
###闭包

---

闭包`closure`: 函数的嵌套, `内层的函数`可以`使用外层函数的所有变量`, 即时外层函数已经执行完毕.

**示例一**

		function checkClosure(){
    		var str = 'rain-man';
    		setTimeout(
        		function(){ alert(str); } //这是一个匿名函数
    		, 2000);
		}
		checkClosure();
		
1. checkClosure函数的执行是瞬间的.
2. 在checkClosure的函数体内创建了一个变量str.
3. `在checkClosure执行完毕之后str并没有被释放`.
4. 这是因为setTimeout内的`匿名函数存在这对str的引用`.
5. 待到2秒后函数体内的匿名函数被执行完毕,str才被释放.	

**示例二,优化代码**

		function forTimeout(x, y){
    		alert(x + y);
		}
		function delay(x , y  , time){
    		setTimeout('forTimeout(' +  x + ',' +  y + ')' , time);    
		}

		上面的delay函数十分难以阅读，也不容易编写，但如果使用闭包就可以让代码更加清晰:
		
		function delay(x , y , time){
			setTimeout(
          		function(){
              		forTimeout(x , y) 
          		}          
      		, time);   
  		}
	
###匿名函数和闭包

---

* 匿名函数最大的用途是创建闭包
* 构建命名空间,以减少全局变量的使用

**示例**

		var oEvent = {};
		(function(){ 
    		var addEvent = function(){ /*代码的实现省略了*/ };
    		function removeEvent(){}

    		oEvent.addEvent = addEvent;
    		oEvent.removeEvent = removeEvent;
		})();
		
在这段代码中函数addEvent和removeEvent都是局部变量,但我们可以通过全局变量oEvent使用它,这就大大减少了全局变量的使用,增强了网页的安全性.

使用代码:

		oEvent.addEvent(document.getElementById('box') , 'click' , function(){});
		
**示例**

建了一个变量rainman,并通过`直接调用匿名函数`初始化为5:

		var rainman = (function(x , y){
    		return x + y;
		})(2 , 3);
		
		也可以写成下面的形式,因为第一个括号只是帮助我们阅读,但是不推荐使用下面这种书写格式.
		
		var rainman = function(x , y){
				return x + y;
		 }(2 , 3);
		 
**示例**

这段代码中的变量one是一个局部变量(因为它被定义在一个函数之内),因此外部是不可以访问的.但是这里我们创建了inner函数,inner函数是可以访问变量one的;又将全局变量outer引用了inner,所以三次调用outer会弹出递增的结果.

		var outer = null;

		(function(){
    		var one = 1;
    		function inner (){
        		one += 1;
        		alert(one);
    		}
    		outer = inner;
		})();

		outer();    //2
		outer();    //3
		outer();    //4
		
###一些注意事项

---

**闭包允许内层函数引用父函数中的变量,但是该变量是最终值**

		<body>
  		<ul>
 	    	<li>one</li>
      		<li>two</li>
 		    <li>three</li>
      		<li>one</li>
 		 </ul>
 		 
 		 
 		 var lists = document.getElementsByTagName('li');
		 for(var i = 0 , len = lists.length ; i < len ; i++){
    		lists[ i ].onmouseover = function(){
        		alert(i);    
    		};
		 }
		 
当鼠标移过每一个`<li>`元素时,总是弹出`4`,而`不是`我们期待的`元素下标`.

`因为`: 

1. 当mouseover事件调用监听函数时,首先在匿名函数(function(){ alert(i); })内部查找是否定义了 `i`, 结果是`没有定义`. 
2. 因此它会`向上查找`,查找`结果是已经定义了`,并且`i`的值是`4`(`循环后的i值`);
3. 所以,最终每次弹出的都是4.

`解决方法`:

		var lists = document.getElementsByTagName('li');
		for(var i = 0 , len = lists.length ; i < len ; i++){
    		(function(index){
        		lists[ index ].onmouseover = function(){
            		alert(index);    
        		};                    
    		})(i);
		}

**内存泄露**

闭包十分容易造成浏览器的内存泄露,严重情况下会是浏览器挂死.
