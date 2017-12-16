# 物理CPU, 逻辑CPU和CPU核数

---

## 概念

### 物理CPU

实际`Server`中插槽上的CPU个数.

物理CPU数量, 可以数不重复的`phusical id`有几个.

### 逻辑CPU

`/proc/cpuinfo`, 信息内容分别列出了`processor 0-n`的规格.

一般情况, 我们认为一颗cpu可以有多核, 加上intel的超线程技术(HT), 可以在逻辑上再分一倍数量的cpu core出来.

逻辑CPU数量=物理cpu数量 x cpu cores 这个规格值 x 2.

Linux下top查看的CPU也是逻辑CPU个数.

### CPU 核数

一块CPU上面能处理数据的芯片组的数量.

一般来说, 物理CPU个数×每颗核数就应该等于逻辑CPU的个数. 如果不相等的话,则表示服务器的CPU支持超线程技术.

## 查看CPU信息   

* `vendor id`: 厂商id. 如果处理器为英特尔处理器,则字符串是`GenuineIntel`.
* `processor`: 包括这一逻辑处理器的唯一标识符.
* `physical id`: 包括每个物理封装的唯一标识符.
* `core id`: 保存每个内核的唯一标识符.
* `siblings`: 列出了位于相同物理封装中的逻辑处理器的数量.
* `cpu cores`:  包含位于相同物理封装中的内核数量.

1. 拥有相同`physical id`的所有逻辑处理器共享同一个物理插座, 每个`physical id`代表一个唯一的物理封装.
2. `siblings`表示位于这一物理封装上的逻辑处理器的数量, 它们可能支持也可能不支持超线程(HT)技术.
3. 每个`core id`均代表一个唯一的处理器内核, 所有带有相同 core id 的逻辑处理器均位于同一个处理器内核上.
4. 如果有一个以上逻辑处理器拥有相同的`core id`和`physical id`, 则说明系统支持超线程(HT)技术.
5. 如果有两个或两个以上的逻辑处理器拥有相同的`physical id`, 但是`core id`不同, 则说明这是一个多内核处理器.
cpu cores条目也可以表示是否支持多内核.

## 示例

### 查看物理CPU个数

```shell
[ansible@production-front-1 ~]$ cat /proc/cpuinfo |grep "physical id"|sort |uniq|wc -l
1
```

### 查看逻辑CPU的个数

```
[ansible@production-front-1 ~]$ cat /proc/cpuinfo |grep "processor"|wc -l
2
```

### 查看CPU是几核

```
cat /proc/cpuinfo |grep "cores"|uniq
cpu cores	: 2
```