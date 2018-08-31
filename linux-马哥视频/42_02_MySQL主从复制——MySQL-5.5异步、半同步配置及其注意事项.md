# 42_02_MySQL主从复制——MySQL-5.5异步、半同步配置及其注意事项

---

## 笔记

### 主从

一个从只能属于一个主服务器.

MySQL 5.6, 在复制上引入了`gtid`(全局事务号). 引入多线程复制, `multi-thread replication`多线程复制.

### 配置MySQL复制基本步骤

#### master

**1.** 启用二进制日志

```
log-bin = master-bin

# 索引文件名称
log-bin-index = master-bin.index
```

**2.** 选择一个唯一`server-id`

```
server-id = {0-2^32}
```

**3.** 创建具有复制权限的用户

```
# 从主服务器的二进制文件中复制事件的权限
REPLICATION SLAVE

# 链接主服务器获取相关信息权限
REPLICATION CLIENT
```

授予复制账号`REPLICATION CLIENT`权限, **复制用户**可以使用`SHOW MASTER STATUS`,`SHOW SLAVE STATUS`和`SHOW BINARY LOGS`来确定复制状态.

#### slave

**1.** 启用中继日志

```
relay-log = relay-log
relay-log-index = relay-log.index
```

**2.** 关闭二进制日志(如果不需要作为其他服务器的主)

**3.** 选择一个唯一的`server-id`

```
server-id = {0 - 2^32-1}
```

**4.** 链接至主服务器, 并开始复制数据

`mysql`命令行下执行

```
mysql> CHANGE MASTER TO 
MASTER_HOST='', MASTER_PORT='',MASTER_LOG_FILE='',MASTER_LOG_POS=xxx,MASTER_USER='',MASTER_PASSWORD='';
```

不会立即启动从服务器. 需要启动从服务器线程.

```
mysql> START SLAVE;
```

**5.** 登录主服务器查看主服务器日志, 以及位置

```
mysql> SHOW MATSER STATUS;
```

**6.** 让从服务器变只读

修改配置文件

```
[mysqld]
read-only = ON
```

对具有`SUPER`权限的用户不生效.

### 复制线程

`mysql`主服务器写数据, 既要写到数据库中一份, 又要写到二进制日志中. 要为每一个从服务器, 启动一个线程`dump`(线程）, 将数据发给从服务器.

每个从服务器都会对应一个`dump`线程.

**多线程复制**, 但是一个库对应一个线程. 

主从线程

* master: dump
* salve: IO_Thread(连到主服务器`dump`线程), SQL_Thread

```
mysql> START SLAVE;
```

会启动两个线程.

```
mysql> START SLAVE IO_Thread;
mysql> START SLAVE SQL_Thread;
```

### 半同步超时

如果所有从服务器都联系不上, 主服务器在超过超时时间后, 主从可以断开, 而后降级为异步模式继续工作.

### 从服务器开机自动启动

依赖如下2个文件

#### `master.info`

记录主从同步中, 主的相关信息.

* 用户名
* 密码
* 同步文件
* 同步位置
* ...

#### `relay-log.info`

* 当前使用`relay-log`文件
* 当前使用`relay-log`位置
* 当前使用主二进制文件
* 当时使用主二进制文件位置

### 半同步复制

**现在是 mysql5.5**

#### 插件

* 主: `semisync_master.so`
* 从: `semisync_slave.so`

#### 安装插件

主

```
mysql> INSTALL PLUGIN rpl_semi_sync_master SONAME 'semisync_master.so';
mysql> SET GLOBAL rpl_semi_sync_master_enabled = 1;
mysql> SET GLOBAL rpl_semi_sync_master_timeout = 1000;
```

`rpl_semi_sync_master_enabled`和`rpl_semi_sync_master_timeout`也可写在`[mysqld]`配置文件中.


从

```
mysql> INSTALL PLUGIN rpl_semi_sync_master SONAME 'semisync_slave.so';
mysql> SET GLOBAL rpl_semi_sync_salve_enabled = 1;
mysql> STOP SALVE IO_THREAD; START SLAVE IO_THREAD;
```

`rpl_semi_sync_salve_enabled`也可写在`[mysqld]`配置文件中.

## 整理知识点

---