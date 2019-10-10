# InnoDB 同步机制

## 概述

* `mutex`互斥
	* 基于`test-and-set`机制实现, 在其基础上做了优化, 流程为
		* 线程调用`test-and-set`返回 1, 说明其他线程已经持有了这把锁, 此时进行自旋. 自旋时间大约为 `20us`
		* 再次获取`mutex`, 如果还是不能获取到就放入`wait array`中, 等待被唤醒.
* `rw-lock`, 可以给临街资源加上`s-latch`或者`x-latch`.
	* `s-latch`允许并发的读取操作
	* `x-latch`完全的互斥操作