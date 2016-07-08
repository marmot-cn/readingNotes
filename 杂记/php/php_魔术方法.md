#PHP 魔术方法

---

###魔术方法

* `__construct()`
* `__destruct()`
* `__call()`
* `__callStatic()`
* `__get()`
* `__set()`
* `__isset()`
* `__unset()`
* `__sleep()`: serialize时调用,必须返回一个包含对象中所有应被序列化的变量名称的数组.
* `__wakeup()`: unserialize时调用.
* `__toString()`: 一个类被当成字符串时应怎样回应.`echo $obj;`这样的情况.
* `__invoke()`: 以调用函数的方式调用一个对象时
* `__set_state()`: `var_export`时调用该静态方法
* `__clone()`
* `__debugInfo()`

####`__sleep()`和`__wakeup()`

`serialize()` 函数会检查类中是否存在一个魔术方法 `__sleep()`. 如果存在,该方法会先被调用,然后才执行序列化操作,并返回一个包含对象中所有应被序列化的变量名称的数组.如果该方法未返回任何内容,则 `NULL` 被序列化,并产生一个 E_NOTICE 级别的错误.

`__sleep()`方法必须返回一个数组,包含需要串行化的属性,PHP会抛弃其它属性的值.如果没有`__sleep()`方法,PHP将保存所有属性.

`unserialize()` 会检查是否存在一个 `__wakeup()` 方法.如果存在,则会先调用 `__wakeup` 方法,预先准备对象需要的资源.

**`__sleep()`示例**

		class A 
		{
			private $a = 1;
			private $b = 2;
		}

		$obj = new A();
		
		var_dump(serialize($obj));
		
		输出:string(42) "O:1:"A":2:{s:4:"Aa";i:1;s:4:"Ab";i:2;}"
		
我们使用__sleep,不期望序列化$b的内容
		
		class A 
		{
			private $a = 1;
			private $b = 2;
			
			public function __sleep()
		    {
		        return array('a');
		    }
		}

		$obj = new A();
		
		var_dump(serialize($obj));
		
		输出:string(27) "O:1:"A":1:{s:4:"Aa";i:1;}"
		
我们用`__sleep()`,不返回任何内容.

		
		class A 
		{
			private $a = 1;
			private $b = 2;
			
			public function __sleep()
		    {
		        return ;
		    }
		}

		$obj = new A();
		
		var_dump(serialize($obj));	
		
		输出(如上所述 NULL 被序列化,并产生一个 E_NOTICE 级别的错误):
		PHP Notice:  serialize(): __sleep should return an array only containing the names of instance-variables to serialize in /Users/chloroplast1983/Sites/test/MagicMthods/sleepAndWakeup.php on line 21
		
		Notice: serialize(): __sleep should return an array only containing the names of instance-variables to serialize in /Users/chloroplast1983/Sites/test/MagicMthods/sleepAndWakeup.php on line 21
		string(2) "N;"	
		
**`__wakeup()`示例**

		class A 
		{
			private $a = 1;
			private $b = 2;
		    
		    public function __wakeup()
		    {
		        echo 'wake up', PHP_EOL;
		    }
		}
		
		$obj = new A();
		$serializeObj = serialize($obj);
		var_dump(unserialize($serializeObj));
		
		输出(第一行输出wakeup):
		wake up
		object(A)#2 (2) {
		  ["a":"A":private]=>
		  int(1)
		  ["b":"A":private]=>
		  int(2)
		}
		
####`__toString()`

`__toString()` 方法用于一个类被当成字符串时应怎样回应.例如 `echo $obj;` 应该显示些什么.
此方法必须**返回**一个字符串,否则将发出一条 `E_RECOVERABLE_ERROR` 级别的致命错误.

**没有`__toString()`示例**

		class A 
		{
			private $a = 1;
			private $b = 2;
		}

		$obj = new A();
		
		echo $obj;
		
		输出:
		PHP Catchable fatal error:  Object of class A could not be converted to string in /Users/chloroplast1983/Sites/test/MagicMthods/toString.php on line 16

		Catchable fatal error: Object of class A could not be converted to string in /Users/chloroplast1983/Sites/test/MagicMthods/toString.php on line 16

**`__toString()`示例**

		class A 
		{
			private $a = 1;
			private $b = 2;
			
		    public function __toString()
		    {
		       return "I am Class A".PHP_EOL;
		    }
		}

		$obj = new A();
		
		echo $obj;
		
		输出:
		I am Class A

####`__invoke()`

当尝试**以调用函数的方式调用一个对象时**,`__invoke()` 方法会被自动调用.

**示例**

		class A 
		{
			function __invoke($x) {
		        var_dump($x);
		    }
		}
		
		$obj = new A();
		
		$obj(5);

####`__set_state()`

`PHP 5.1.0` 起当调用 `var_export()` 导出类时,此静态方法会被调用.

本方法的唯一参数是一个数组,其中包含按 array('property' => value, ...)格式排列的类属性.

就是`var_export()`的`回调函数`.

**`var_export`**

输出或返回一个变量的字符串表示.

`mixed var_export ( mixed $expression [, bool $return ] )`

此函数返回关于传递给该函数的变量的结构信息,它和 `var_dump()` 类似, 不同的是其返回的表示是合法的 PHP 代码.

也就是说,`var_export()`返回的代码,可以直接当作php代码赋值个一个变量.而这个变量就会取得和被var_export一样的类型的值.但是,当变量类型为resource的时候,是无法简单copy复制的,所以,当var_export的变量是resource类型时,var_export会返回NULL.
		
示例:

		$a = array (1, 2, array ("a", "b", "c"));
		var_export ($a);
		输出:
		array (
		  0 => 1,
		  1 => 2,
		  2 =>
		  array (
		    0 => 'a',
		    1 => 'b',
		    2 => 'c',
		  ),
		)
		
		$a = array (1, 2, array ("a", "b", "c"));
		$b = var_export ($a,true);
		
		echo $b;
		输出:
		array (
		  0 => 1,
		  1 => 2,
		  2 =>
		  array (
		    0 => 'a',
		    1 => 'b',
		    2 => 'c',
		  ),
		)

**示例`__set_state()`**

我们创建一个类然后执行`var_export()`观察其输出:

		class A
		{
		    public $var1 = 1;
		    public $var2 = 2;
		}

		$a = new A();
		$b = var_export($a,true);
		
		var_dump($b);
		
输出是:

		string(56) "A::__set_state(array(
		   'var1' => 1,
		   'var2' => 2,
		))"
		
这是一段合法的php代码,它代表一个 static __set_state 方法在我们的 Class A.如果我们要执行它的话就要使用`eval`函数.

		class A
		{
		    public $var1 = 1;
		    public $var2 = 2;
		}

		$a = new A();
		eval('$b = '.var_export($a,true).';');
		
		var_dump($b);
		
输出是:

		PHP Fatal error:  Call to undefined method A::__set_state()...
		
我们不能执行这段PHP代码,因为我们的PHP类没有`__set_state()`,我们必须创建一个.

		class A
		{
		    public $var1 = 1;
		    public $var2 = 2;
		
		    public static function __set_state($array) // As of PHP 5.1.0
		    {
		       $tmp = new A();
		       $tmp->var1 = 3;
		       $tmp->var2 = 4;
		       return $tmp;
		    }
		}
		
		$a = new A();
		eval('$b = '.var_export($a,true).';');
		var_dump($b);
		
输出是:

		object(A)#2 (2) {
		  ["var1"]=>
		  int(3)
		  ["var2"]=>
		  int(4)
		}
		
我们可以看见我们`__set_state`内部的代码已经被执行了.	
		
####`__debugInfo()`

####`__construct()`

####`__destruct()`

####`__call()`

####`__callStatic()`

####`__get()`

####`__set()`

####`__isset()`

####`__unset()`

####`__clone()`