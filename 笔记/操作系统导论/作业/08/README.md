# README

程序`mlfq.py`允许你查看本章介绍的`MLFQ`调度程序的行为.

### 1. 问题

只用两个工作和两个队列运行几个随机生成的问题。针对每个工作计算`MLFQ`的执行记录。限制每项作业的长度并关闭`I/O`


* -n numqueues, 队列数量
* -j numjobs, 任务数量
* -m maxlen, 任务的最大时长
* -M maxio, io的最大次数
* -s randomseed, 随机种子，影响随机值
* -c 查看答案
* -l jobList, 表示指定任务队列，格式为x1,y1,z1:x2,y2,z2
	* x表示到达时间
	* y表示任务时长
	* z is how often the job issues an I/O request
* -B BOOST, --boost=BOOST
    * how often to boost the priority of all jobs back to high priority (0 means never)

### 2. 问题

如何运行调度程序来重现本周的每个实例

**8.2**

```
python mlfq.py -n 3 -l 0,200,0 -c
```

一个场务，运行三个队列(`-n 3`), 到达时间为`0`, `io`为`0`.

* 第一个队列执行`10ms`片段
* 第二个队列执行`10ms`片段
* 第三个队列执行`180ms`片段

**8.3**

```
python mlfq.py -n 3 -l 0,180,0:100,20,0 -c -q 10
```

* 第2个任务在`100ms`插入，运行`20ms`

**8.4**

```
python mlfq.py -n 3 -l 0,180,0:0,20,1  -i 8 -c 
```

* `io`时间为`8ms`，需要减去`io`开始和`io`结束的`1ms`

### 3. 问题

将如何配置调度程序参数，像轮转调度程序那样工作？

```
python mlfq.py -n 1 -l 0,200,0:0,200,0 -c
```

设定一个队列

### 4. 问题

设计两个工作的负载和调度程序参数，以便一个工作利用较早的规则**4a**和**4b**来“愚弄”调度程序(用`-S`标志打开)，在特定的时间间隔内获得`99%`的CPU。

```
python mlfq.py -n 3 -l 0,100,0:0,100,9 -c -q 10 -S -i 1
```

`job1`每隔`9s`就进行`io`放弃CPU, 获得高优先级

### 5. 问题

给定一个系统，其最高队列中的时间片长度为`10ms`, 你需要如何频繁地将工作推回到最高优先级别（带有`-B`标志），以保证一个长时间运行（并可能饥饿）的工作得到至少`5%`的CPU。

总运行时间`200ms x 5% = 10ms`

`-B 200ms`重置优先级, 可以得到`10ms`的运行时间

### 6. 问题

调度中有一个问题，即刚完成`I/O`的作业添加在队列的哪一端。`-I`标志改变了这个调度模拟器的这方面的行为。尝试一些工作负载，看看你是否能看到这个标志的效果。

`-i IOTIME, --iotime=IOTIME`: how long an I/O should last (fixed constant)。

`io`运行时间。