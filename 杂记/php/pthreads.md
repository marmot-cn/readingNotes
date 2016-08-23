#pthreads

---

####`set_time_limit`

		void set_time_limit ( int $seconds )

设置允许脚本运行的时间,单位为秒.如果超过了此设置,脚本返回一个致命的错误.默认值为`30`秒,或者是在`php.ini`的`max_execution_time`被定义的值,如果此值存在.

当此函数被调用时,`set_time_limit()`会从零开始重新启动超时计数器.换句话说,如果超时默认是30秒,在脚本运行了了25秒时调用`set_time_limit(20)`,那么,脚本在超时之前可运行总时间为45秒.

**`cli`模式**

在`cli`模式下,这个值会被硬编码为0.

		$max_time = ini_get("max_execution_time");
		echo $max_time, PHP_EOL;
		
		root@3b4f72d7fb9c:/var/www/html# php test.php
		0
		
####pthreads

**不用线程执行任务**

我们先用传统的方式执行一个任务,假设这个任务需要执行一秒钟

		//我们假设这个任务执行需要1秒,我们需要执行5次,则需要一共5秒
		for ($i = 1; $i <= 5; $i++) {
		    sleep(1);
		    echo "Hello World", PHP_EOL;
		}
		
我们执行然后查看时间:

		time php test.php

		Hello World
		Hello World
		Hello World
		Hello World
		Hello World		
		
		real	0m5.063s
		user	0m0.040s
		sys	0m0.010s
		
**使用线程执行任务**

		class SomeThreadedClass extends Thread {
		
		    public function run() {
		        sleep(1);
		        echo "Hello World", PHP_EOL;
		    }
		};
		
		$threads = [];
		
		for ($i = 1; $i <= 5; $i++) {
		    $threads[$i] = new SomeThreadedClass(); 
		    $threads[$i]->start();
		}
		
		for ($i = 1; $i <= 5; $i++) {
		    $threads[$i]->join();
		}
		
我们执行然后查看时间:		
		
		Hello World
		Hello World
		Hello World
		Hello World
		Hello World
		
		real	0m1.142s
		user	0m0.080s
		sys	0m0.060s		
		
我们几乎只使用了`1/5`的时间.		
		
**`Stackable`**

`Stackable` 是 `Threaded` 的一个别称,直到 `pthreads v.2.0.0`