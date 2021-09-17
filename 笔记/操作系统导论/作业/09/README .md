# README 

`lottery.py`这个程序允许你查看彩票调度程序的工作原理。详情参阅`README`文件。

### 1. 问题

计算3个工作在随机种子为`1、2`和`3`时的模拟解。

```
python lottery.py -s 1 -j 3 -c

-s 表示随机种子
-j 表示任务数量
-c 表示展示答案

RG jlist
ARG jobs 3
ARG maxlen 10
ARG maxticket 100
ARG quantum 1
ARG seed 1

Here is the job list, with the run time of each job:
  Job 0 ( length = 1, tickets = 84 )
  Job 1 ( length = 7, tickets = 25 )
  Job 2 ( length = 4, tickets = 44 )


** Solutions **

Random 651593 -> Winning ticket 119 (of 153) -> Run 2
  Jobs:
 (  job:0 timeleft:1 tix:84 )  (  job:1 timeleft:7 tix:25 )  (* job:2 timeleft:4 tix:44 )
Random 788724 -> Winning ticket 9 (of 153) -> Run 0
  Jobs:
 (* job:0 timeleft:1 tix:84 )  (  job:1 timeleft:7 tix:25 )  (  job:2 timeleft:3 tix:44 )
--> JOB 0 DONE at time 2
Random 93859 -> Winning ticket 19 (of 69) -> Run 1
  Jobs:
 (  job:0 timeleft:0 tix:--- )  (* job:1 timeleft:7 tix:25 )  (  job:2 timeleft:3 tix:44 )
Random 28347 -> Winning ticket 57 (of 69) -> Run 2
  Jobs:
 (  job:0 timeleft:0 tix:--- )  (  job:1 timeleft:6 tix:25 )  (* job:2 timeleft:3 tix:44 )
Random 835765 -> Winning ticket 37 (of 69) -> Run 2
  Jobs:
 (  job:0 timeleft:0 tix:--- )  (  job:1 timeleft:6 tix:25 )  (* job:2 timeleft:2 tix:44 )
 ...
```

第一行`119`, 通过随机数对`153`去余得到。总彩票数为三个任务之和`153`, `119`则对应定位到`job2`

* `0-83`: `job0`
* `84-109`: `job1`
* `109-152`: `job2`

### 2. 问题

现在运行两个具体的工作：每个长度为`10`，但是一个（工作0）只有一张彩票，另一个（工作1）有100张（-l 10:1, 10:100）。彩票数量如此不平衡时会发生什么？在工作1完成之前，工作0是否会运行? 多久? 一般来说，这种彩票不平衡对彩票调度的行为有什么影响？

* 在工作1完成之前，工作0会运行
* 工作0运行机会比较小

### 3. 问题

如果运行两个长度为`100`的工作，都有`100`张彩票（-1 100:100, 100:100），调度程序有多不公平？运行一些不同的随机种子来确定（概率上的）答案。不公平性取决于一项工作比另一项工作早完成多少。

`50%`

### 4. 问题

随着量子规模(-q)变大，你对上一个问题的答案如何改变？

量子规模变大，相当于工作长度缩短了。不公平性增加（公平性取决于一项工作比另一项工作早完成多少）。

### 5. 你可以制作类似本章中的图表吗？

以量子规模（即时间片的长度）为x轴，公平度为y轴画一个曲线图。很明显量子规模越大，公平度越低。

