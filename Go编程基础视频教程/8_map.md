#8 map

[地址](http://edu.51cto.com/lesson/id-32313.html "地址")

###笔记

---

**map**

* 类似其它语言中的哈希表或者字典,以`key-value`形式存储数据
* `key`必须是支持 `==` 或者 `!=` 比较运算的类型, 不可以是`函数`, `map` 或者 `slice`,所有类型都可以作为 `value` 类型.
* `map` 查找比线性搜索快很多, 但比使用索引访问数据的类型慢100倍
* `map` 使用 `make` 创建, 支持 `:=`. 

		var m map[int]string //声明
		m = map[int]string{} //初始化
		
		//使用make声明
		var m map[int]string = make(map[int]string)
		
		//:= 声明
		m := make(map[int]string)

* `make([ketType]valueType, cap)`, `cap`表示容量, 可省略(make只可以为 `slice` , `map`, `channel` 三种类型进行内存的分配)
* 超出容量时会自动扩容, 但尽量提供一个合理的初始值
* 使用`len()`获取元素个数

* 键值对不存在时自动添加, 使用`delte()`删除某键值对

		m := make(map[int]string)
		m[1] = "ok"
		delete(m,1)//删除key值1

* 使用 `for range` 对 `map` 和 `slice` 进行迭代操作
	
		for i,v:=range slice {
			//i,索引
			//v,值
		}
		
		for i,v:=range map {
			//i,key
			//v,值
		}

**示例**

		sm := make([]map[int]string, 5)
		
		for _, v := range sm {
			v := make(map[int]string,1)
			v[1] = "ok"
			fmt.Println(v)
		}
		fmt.Println(sm)
		
		输出:
		map[1:ok]
		map[1:ok]
		map[1:ok]
		map[1:ok]
		map[1:ok]
		
		[map[] map[] map[] map[] map[]]//最后还是空的,因为v是copy,所以对v操作不会影响map本身
		
		for _, v := range sm {
			sm[i] := make(map[int]string,1)
			sm[i][1] = "ok"
			fmt.Println(v)
		}
		
		输出:
		map[1:ok]
		map[1:ok]
		map[1:ok]
		map[1:ok]
		map[1:ok]
		
		[map[1:ok] map[1:ok] map[1:ok] map[1:ok] map[1:ok]]//会影响本身
		
**map的排序**

		m :=map[int]string(1:"a",2:"b",3:"c",4:"d",5:"e")
		s := make([]int,len(m))
		i := 0
		for k,_:=range m{
			s[i] = k
			i++
		}
		
		fmt.Println(s)
		
		输出:
		//无序
		[2 1 3 5 4]
		[3 5 4 2 1]
		
		//修改成有序的:
		1.引入包 sort
		
		...
		sort.Ints(s)
		...
		fmt.Println(s)
		输出:
		[1 2 3 4 5]
		[1 2 3 4 5]
		
###整理知识点

---

**map的无序性**

多次通过range循环来迭代访问map中元素时,尽管您访问的是同一个map,但是访问元素的顺序在前后两次range中是不会完全相同的.

Go在range遍历Map中元素的时候,从随机的一个位置开始迭代.