# 实验1 查看CPU和内存,用机器指令和汇编指令编程

---

因为我使用的是mac,需要使用`debug.exe`,做如下操作

**下载 dosbox**

[dosbox](http://www.dosbox.com/download.php?main=1
 "dosbox")
 
**下载 64位debug.exe**
 
**dos box挂载**

把`debug.exe`放入`~/Downloads/debug`中.

在`Dosbox`中输入

		mount c ~/Downloads/debug
		
		C:
		
		dir
		
		就能看见 debug.exe
		
		再次输入
		DEBUG.exe 就可以运行了

### 常用功能

* `R` 
	* 查看 `寄存器`的内容
	* 改变 `寄存器`的内容
* `D` 查看 `内存中`的内容
* `E` 改写 `内存中`的内容
* `U` 将`内存中`的机器指令翻译成汇编指令
* `T` 执行一条机器指令
* `A` 以汇编指令的格式在内存中写入一条机器指令


		
		



