# FLUSH TABLE WITH READ LOCK

## 简介

关闭所有打开的表, 同时对于所有数据库中的表都加一个读锁, 直到显示地执行`unlock tables`.

将所有的脏页都要刷新到磁盘, 然后对所有的表加上了读锁, 于是这时候直接拷贝数据文件也就是安全的.

## flush tables with read lock的处理

### 请求锁

请求全局`read lock`.

### 等待锁

在`flush tables with read lock`成功获得锁之前, 必须等待所有语句执行完成(包括SELECT). 所以如果有个慢查询在执行, 或者一个打开的事务, 或者其他进程拿着表锁, `flush tables with read lock`就会被阻塞, 直到所有的锁被释放. 

### 刷新表

当`flush tables with read lock`拿到锁后, 必定`flush data`.

### 持有锁

我们可以使用`unlock tables`或者其它命令来释放锁.