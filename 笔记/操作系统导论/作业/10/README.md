# README

## Homework (Simulation)

In this homework, we’ll use multi.py to simulate a multi-processor
CPU scheduler, and learn about some of its details. Read the related
README for more information about the simulator and its options.

## Questions

### 1

To start things off, let’s learn how to use the simulator to study how
to build an effective multi-processor scheduler. The first simulation
will run just one job, which has a run-time of 30, and a working-set
size of 200. Run this job (called job ’a’ here) on one simulated CPU
as follows: ./multi.py -n 1 -L a:30:200. How long will it
take to complete? Turn on the -c flag to see a final answer, and the
-t flag to see a tick-by-tick trace of the job and how it is scheduled.

```
python multi.py -n 1 -L a:30:200 -c -t

Scheduler central queue: ['a']

   0   a
   1   a
   2   a
   3   a
   4   a
   5   a
   6   a
   7   a
   8   a
   9   a
----------
  10   a
  11   a
  12   a
  13   a
  14   a
  15   a
  16   a
  17   a
  18   a
  19   a
----------
  20   a
  21   a
  22   a
  23   a
  24   a
  25   a
  26   a
  27   a
  28   a
  29   a

Finished time 30

Per-CPU stats
  CPU 0  utilization 100.00 [ warm 0.00 ]

```

* 1个CPU核心
* 一个队列
	* 名字`a`
	* 100 执行时间
	* 200 工作size(how much cache
space it needs to run efficiently)

但是把工作size, 降低为100, 则可以利用`warm`机制，提升运行效率。默认`warm`时间为`10 time units`, 运行速度为`2`倍。
因为默认的`warm`的size为100, 所以上面队列的`200`size无法利用`warm`机制（结果可见 `warm 0.00`）

```
python multi.py -n 1 -L a:30:100 -c -t

Scheduler central queue: ['a']

   0   a
   1   a
   2   a
   3   a
   4   a
   5   a
   6   a
   7   a
   8   a
   9   a
----------
  10   a
  11   a
  12   a
  13   a
  14   a
  15   a
  16   a
  17   a
  18   a
  19   a

Finished time 20

Per-CPU stats
  CPU 0  utilization 100.00 [ warm 50.00 ]
```

`0-10`没有利用`warm`, `10-19`利用了`warm`。最终结果`warm 50.00`, 利用率为`50%`。

### 2

Now increase the cache size so as to make the job’s working set
(size=200) fit into the cache (which, by default, is size=100); for
example, run ./multi.py -n 1 -L a:30:200 -M 300. Can
you predict how fast the job will run once it fits in cache? (hint:
remember the key parameter of the warm rate, which is set by the
-r flag) Check your answer by running with the solve flag (-c) en-
abled.

使用`20`单位时间技能完成，因为`0-10`用于`warm`, 后续`11-20`是`2`倍效率运行。

```
python multi.py -n 1 -L a:30:200 -M 300 -t -c

Scheduler central queue: ['a']

   0   a
   1   a
   2   a
   3   a
   4   a
   5   a
   6   a
   7   a
   8   a
   9   a
----------
  10   a
  11   a
  12   a
  13   a
  14   a
  15   a
  16   a
  17   a
  18   a
  19   a

Finished time 20

Per-CPU stats
  CPU 0  utilization 100.00 [ warm 50.00 ]
```

### 3

One cool thing about multi.py is that you can see more detail
about what is going on with different tracing flags. Run the same
simulation as above, but this time with time left tracing enabled
(-T). This flag shows both the job that was scheduled on a CPU
at each time step, as well as how much run-time that job has left
after each tick has run. What do you notice about how that second
column decreases?

**利用warm**

```
python multi.py -n 1 -L a:30:200 -M 300 -t -c -T
...
Scheduler central queue: ['a']

   0   a [ 29]
   1   a [ 28]
   2   a [ 27]
   3   a [ 26]
   4   a [ 25]
   5   a [ 24]
   6   a [ 23]
   7   a [ 22]
   8   a [ 21]
   9   a [ 20]
----------------
  10   a [ 18]
  11   a [ 16]
  12   a [ 14]
  13   a [ 12]
  14   a [ 10]
  15   a [  8]
  16   a [  6]
  17   a [  4]
  18   a [  2]
  19   a [  0]

Finished time 20

Per-CPU stats
  CPU 0  utilization 100.00 [ warm 50.00 ]
```

**不利用warm**

```
python multi.py -n 1 -L a:30:200 -t -c -T
...
Scheduler central queue: ['a']

   0   a [ 29]
   1   a [ 28]
   2   a [ 27]
   3   a [ 26]
   4   a [ 25]
   5   a [ 24]
   6   a [ 23]
   7   a [ 22]
   8   a [ 21]
   9   a [ 20]
----------------
  10   a [ 19]
  11   a [ 18]
  12   a [ 17]
  13   a [ 16]
  14   a [ 15]
  15   a [ 14]
  16   a [ 13]
  17   a [ 12]
  18   a [ 11]
  19   a [ 10]
----------------
  20   a [  9]
  21   a [  8]
  22   a [  7]
  23   a [  6]
  24   a [  5]
  25   a [  4]
  26   a [  3]
  27   a [  2]
  28   a [  1]
  29   a [  0]

Finished time 30

Per-CPU stats
  CPU 0  utilization 100.00 [ warm 0.00 ]

```

### 4

Now add one more bit of tracing, to show the status of each CPU
cache for each job, with the -C flag. For each job, each cache will
either show a blank space (if the cache is cold for that job) or a ’w’(if the cache is warm for that job). At what point does the cache
become warm for job ’a’ in this simple example? What happens
as you change the warmup time parameter (-w) to lower or higher
values than the default?


`2`单位时间就可以`warm`

```
python multi.py -n 1 -L a:30:200 -M 300 -t -c -T -C -w 2
...
Scheduler central queue: ['a']

   0   a [ 29] cache[ ]
   1   a [ 28] cache[w]
   2   a [ 26] cache[w]
   3   a [ 24] cache[w]
   4   a [ 22] cache[w]
   5   a [ 20] cache[w]
   6   a [ 18] cache[w]
   7   a [ 16] cache[w]
   8   a [ 14] cache[w]
   9   a [ 12] cache[w]
-------------------------
  10   a [ 10] cache[w]
  11   a [  8] cache[w]
  12   a [  6] cache[w]
  13   a [  4] cache[w]
  14   a [  2] cache[w]
  15   a [  0] cache[w]

Finished time 16

Per-CPU stats
  CPU 0  utilization 100.00 [ warm 87.50 ]

```

### 5

At this point, you should have a good idea of how the simula-
tor works for a single job running on a single CPU. But hey, isn’t
this a multi-processor CPU scheduling chapter? Oh yeah! So let’s
start working with multiple jobs. Specifically, let’s run the follow-ing three jobs on a two-CPU system (i.e., type ./multi.py -n
2 -L a:100:100,b:100:50,c:100:50) Can you predict how
long this will take, given a round-robin centralized scheduler? Use
-c to see if you were right, and then dive down into details with -t to see a step-by-step and then -C to see whether caches got warmed effectively for these jobs. What do you notice?

`150`单位时间，没有利用`warm`

### 6

Now we’ll apply some explicit controls to study cache affinity, as
described in the chapter. To do this, you’ll need the -A flag. This
flag can be used to limit which CPUs the scheduler can place a par-
ticular job upon. In this case, let’s use it to place jobs ’b’ and ’c’ on
CPU 1, while restricting ’a’ to CPU 0. This magic is accomplished
by typing this ./multi.py -n 2 -L a:100:100,b:100:50,
c:100:50 -A a:0,b:1,c:1 ; don’t forget to turn on various trac-
ing options to see what is really happening! Can you predict how
fast this version will run? Why does it do better? Will other com-
binations of ’a’, ’b’, and ’c’ onto the two processors run faster or
slower?

```
python multi.py -n2 -L a:100:100,b:100:50,c:100:50 -A a:0,b:1,c:1 -t -c -T -C 
```

`110`时间

* `A`运行`54`时间
	* `10`(没有`warm`) + (90/2的`warm`)时间
* `B`和`C`轮换调度，因为各自`szie`都是`50`可以同时放入到`CPU`的`CACHE`默认`100`
	* 20次没有warm + 180/2的warm时间

总时间110, 取最长B和C运行的时间。

### 7

One interesting aspect of caching multiprocessors is the opportu-
nity for better-than-expected speed up of jobs when using multi-
ple CPUs (and their caches) as compared to running jobs on a sin-
gle processor. Specifically, when you run on N CPUs, sometimes
you can speed up by more than a factor of N, a situation entitled
super-linear speedup. To experiment with this, use the job descrip-
tion here (-L a:100:100,b:100:100,c:100:100) with a small
cache (-M 50) to create three jobs. Run this on systems with 1, 2,
and 3 CPUs (-n 1, -n 2, -n 3). Now, do the same, but with a
larger per-CPU cache of size 100. What do you notice about per-
formance as the number of CPUs scales? Use -c to confirm your
guesses, and other tracing flags to dive even deeper.

* `-n 1 -L a:100:100,b:100:100,c:100:100 -M 50 -t -c -T -C`
	* `300`时间单位, 且不能用缓存
* `-n 2 -L a:100:100,b:100:100,c:100:100 -M 50 -t -c -T -C`
	* `150`时间单位, 且不能用缓存
* `-n 3 -L a:100:100,b:100:100,c:100:100 -M 50 -t -c -T -C`
	* `100`时间单位, 且不能用缓存
* `-n 1 -L a:100:100,b:100:100,c:100:100 -M 100 -t -c -T -C`
	* `300`时间单位, 因为调度可以用缓存
* `-n 2 -L a:100:100,b:100:100,c:100:100 -M 100 -t -c -T -C`
	* `150`时间单位, 因为调度不可以用缓存
* `-n 3 -L a:100:100,b:100:100,c:100:100 -M 100 -t -c -T -C`
	* `55`时间单位, 可以用缓存(10 + 90/2)

### 8

One other aspect of the simulator worth studying is the per-CPU
scheduling option, the -p flag. Run with two CPUs again, and this
three job configuration (-L a:100:100,b:100:50,c:100:50).
How does this option do, as opposed to the hand-controlled affinity
limits you put in place above? How does performance change as
you alter the ’peek interval’ (-P) to lower or higher values? How
does this per-CPU approach work as the number of CPUs scales?

### 9 

Finally, feel free to just generate random workloads and see if you
can predict their performance on different numbers of processors,
cache sizes, and scheduling options. If you do this, you’ll soon be
a multi-processor scheduling master, which is a pretty awesome
thing to be. Good luck!