# 22 | 想成为架构师,你必须知道CAP理论 

## 笔记

### CAP 定理

又称作**布鲁尔(Brewer's theorem)定理**

### CAP 理论

#### 第一版

对于一个分布式计算系统, 不可能同时满足:

* 一致性(Consistence)
* 可用性(Availability)
* 分区容错性(Partition Tolerance)

#### 第二版

在一个分布式系统(互相连接并共享数据的节点的集合)中, 当涉及读写操作时, 只能保证一致性(Consistence), 可用性(Availability), 分区容错性(Partition Tolerance)三者中的两个, 另外一个必须被牺牲.

强调了:

* 互相连接,共享数据
* write/read pair

因为**分布式系统并不一定互联和共享数据**(Memcached集群, 相互之间没有连接和共享数据, 不符合CAP. Mysql集群互联和进行数据复制, 因此是CAP理论探讨的对象).

强调`write/read pair`, **`CAP`关注的是对数据的读写操作**, 而不是分布式系统的所有功能. 如`ZooKeeper`的选举机制就不是`CAP`探讨的对象.

**互联, 共享, 读写操作**

### 一致性(Consistence)

#### 第一版

所有节点在同一时刻都能看到相同的数据.

```
All nodes see the same data at the same time.
```

强调同一时刻拥有相同的数据.

#### 第二版

对某个指定的客户端来说, 读操作保证能够返回最新的写操作结果.

```
A read is guaranteed to return the most recent write for a given client.
```

**没有**强调同一时刻拥有相同的数据, 意味着实际上对于节点来说, 可能同一时刻拥有不同的数据(same time + diffetent data).

对于系统执行事务来说, **在事务执行过程中, 系统其实处于一个不一致的状态, 不同的节点的数据并不完全一致**.

第二版强调`client`读操作能够获取最新的写结果就没有问题, 因为事务在执行过程中, `client`是无法读取到未提交的数据的, 只有在等到事务提交后, `client`才能读取到事务写入的数据, 而如果事务失败则会进行回滚, `clinet`也不会读取到事务中间写入的数据.

### 可用性(Availability)

#### 第一版

```
Every request gets a response on success/failure.
```

每个请求都能得到成功或者失败的响应.

#### 第二版

```
A non-failing node will return a reasonable response within a reasonable amount of time (no error or timeout).
```

非故障的节点在合理的时间内返回合理的响应(不是错误和超时的响应).

**合理, 不是代表是正确的**

比如应该返回`100`, 但实际上返回了`90`, 肯定是不正确的结果, 但可以是一个**合理**的结果.

### 分区容错性(Partition Tolerance)

#### 第一版

```
System continues to work despite message loss or partial failure.
```

出现消息丢失或者分区错误时系统能够继续运行.

#### 第二版

```
The system will continue to function when network partitions occur.
```

当出现网络分区后, 系统能够继续"履行职责".

第二版定义更宽泛, 即发生了**分区现象**, 不管是什么原因, 可能是丢包, 可能是连接中断, 还可能是拥塞, 只要导致了网络分区, 就通通算在里面.

### CAP 应用

定义了三个要素中只能取两个, 如果选择了`CA`放弃了`P`, 那么当发生分区现象时, 为了保证`C`, 系统需要禁止写入, 当有写入请求时, 系统返回`error`(如: 当前系统不可写入), 这又和`A`冲突了, 因为`A`要求返回`no error`和`no timeout`. 因此, 分布式系统理论上不可能选择`CA`架构, 只能宣策`CP`或者`AP`架构.

#### CP - Consistency/Partition Tolerance

![](./img/22_01.png)

为了保证一致性, 当发生分区现象后, `N1`节点上的数据已经更新到`y`, 但由于`N1`和`N2`之间的复制通道中断, 数据`y`无法同步到`N2`, `N2`节点上的数据还是`x`. 这时客户端`C`访问`N2`时, `N2`需要返回`Error`, 提示客户端`C`系统现在发生了错误. 这种处理**违背了可用性(Availability)的要求**.

**返回错误, 违背`A`**

#### AP - Availability/Partition Tolerance

![](./img/22_02.png)

为了保证可用性, 当发生分区现象有, `N1`节点上的数据已经更新到`y`, 但由于`N1`和`N2`之间的复制通道中断, 数据`y`无法同步到`N2`, `N2`节点上的数据还是`x`. 这时客户端`C`访问`N2`时, `N2`将当前自己拥有的数据`X`返回给客户端`C`了, 而实际上当前最新的数据已经是`y`了. 这就不满足一致性（Consistency）的要求了,因此 CAP 三者只能满足 AP.注意：这里 N2 节点返回 x,虽然不是一个“正确”的结果,但是一个“合理”的结果,因为 x 是旧的数据,并不是一个错乱的值,只是不是最新的数据而已.

**返回不正确但是合理的结果, 违背`C`(没有读取到最新的写操作结果)**

## 扩展