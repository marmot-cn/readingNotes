# README 

In this homework, we’ll just learn about a few useful tools to examine virtual memory usage on Linux-based systems. This will only be a brief hint at what is possible; you’ll have to dive deeper on your own to truly become an expert (as always!).

### 1

The first Linux tool you should check out is the very simple tool
free. First, type man free and read its entire manual page; it’s
short, don’t worry!

```
Usage:
 free [options]

Options:
 -b, --bytes         show output in bytes
 -k, --kilo          show output in kilobytes
 -m, --mega          show output in megabytes
 -g, --giga          show output in gigabytes
     --tera          show output in terabytes
 -h, --human         show human-readable output
     --si            use powers of 1000 not 1024
 -l, --lohi          show detailed low and high memory statistics
 -t, --total         show total for RAM + swap
 -s N, --seconds N   repeat printing every N seconds
 -c N, --count N     repeat printing N times, then exit
 -w, --wide          wide output

     --help     display this help and exit
 -V, --version  output version information and exit
```

### 2

Now, run free, perhaps using some of the arguments that might
be useful (e.g., -m, to display memory totals in megabytes). How
much memory is in your system? How much is free? Do these
numbers match your intuition?

```
[ansible@k8s-agent-1 ~]$ free -m
              total        used        free      shared  buff/cache   available
Mem:          32013       14464        1142        1040       16406       14294
Swap:             0           0           0
```

### 3

Next, create a little program that uses a certain amount of memory,
called memory-user.c. This program should take one command-
line argument: the number of megabytes of memory it will use.
When run, it should allocate an array, and constantly stream through
the array, touching each entry. The program should do this indefi-
nitely, or, perhaps, for a certain amount of time also specified at the command line.

```C
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <string.h>

int main(int argc, char *argv[])
{
  int mem = atoi(argv[1]);
  const int size = mem*1024*1024/sizeof(int);
  int array[size];
  int b = size*sizeof(int);
  int i = 0;
  printf("arg is %d\n", size);
  while(1) {
        int i;
	for(i=0; i<size; i=i+1) {
		array[i] = 0;
		printf("index %d is 0\n",i);
        }
  }
  return 0;
}
```

### 4

Now, while running your memory-user program, also (in a dif-
ferent terminal window, but on the same machine) run the free
tool. How do the memory usage totals change when your program
is running? How about when you kill the memory-user program?
Do the numbers match your expectations? Try this for different
amounts of memory usage. What happens when you use really
large amounts of memory?

### 5

Let’s try one more tool, known as pmap. Spend some time, and read
the pmap manual page in detail.

```
[ansible@k8s-agent-1 ~]$ pmap -h

Usage:
 pmap [options] PID [PID ...]

Options:
 -x, --extended              show details
 -X                          show even more details
            WARNING: format changes according to /proc/PID/smaps
 -XX                         show everything the kernel provides
 -c, --read-rc               read the default rc
 -C, --read-rc-from=<file>   read the rc from file
 -n, --create-rc             create new default rc
 -N, --create-rc-to=<file>   create new rc to file
            NOTE: pid arguments are not allowed with -n, -N
 -d, --device                show the device format
 -q, --quiet                 do not display header and footer
 -p, --show-path             show path in the mapping
 -A, --range=<low>[,<high>]  limit results to the given range

 -h, --help     display this help and exit
 -V, --version  output version information and exit

For more details see pmap(1).
```

**扩展模式**

* Address: 内存开始地址
* Kbytes: 占用内存的字节数（KB）
* RSS: 保留内存的字节数（KB）
* Dirty: 脏页的字节数（包括共享和私有的）（KB）
* Mode: 内存的权限：read、write、execute、shared、private (写时复制)
* Mapping: 占用内存的文件、或[anon]（分配的内存）、或[stack]（堆栈）
* Offset: 文件偏移
* Device: 设备名 (major:minor)

### 6

To use pmap, you have to know the process ID of the process you’re
interested in. Thus, first run ps auxw to see a list of all processes; then, pick an interesting one, such as a browser. You can also use

your memory-user program in this case (indeed, you can even
have that program call getpid() and print out its PID for your
convenience).

```
# inux 系统中对每个线程都有自己的栈，linux系统中默认大小为8M，因为栈中还有其他的，所为会比8M小一点才可以。这里我们增大了限制
ulimit -s 10240000

[root@k8s-agent-1 memory-user]# ./memory-user.out 10
pid is 28722
arg is 2621440
```


```
[ansible@k8s-agent-1 ~]$ sudo pmap 26720
26720:   ./memory-user.out 70
0000000000400000      4K r-x-- memory-user.out
0000000000600000      4K r---- memory-user.out
0000000000601000      4K rw--- memory-user.out
00007fb7b644e000   1760K r-x-- libc-2.17.so
00007fb7b6606000   2048K ----- libc-2.17.so
00007fb7b6806000     16K r---- libc-2.17.so
00007fb7b680a000      8K rw--- libc-2.17.so
00007fb7b680c000     20K rw---   [ anon ]
00007fb7b6811000    132K r-x-- ld-2.17.so
00007fb7b6a26000     12K rw---   [ anon ]
00007fb7b6a30000      8K rw---   [ anon ]
00007fb7b6a32000      4K r---- ld-2.17.so
00007fb7b6a33000      4K rw--- ld-2.17.so
00007fb7b6a34000      4K rw---   [ anon ]
00007ffff3bd1000  71692K rw---   [ stack ]
00007ffff81e3000      8K r-x--   [ anon ]
ffffffffff600000      4K r-x--   [ anon ]
 total            75732K
```

### 7

Now run pmap on some of these processes, using various flags (like
-X) to reveal many details about the process. What do you see?
How many different entities make up a modern address space, as
opposed to our simple conception of code/stack/heap?

```
[ansible@k8s-agent-1 ~]$ sudo pmap 28722 -X
28722:   ./memory-user.out 10
         Address Perm   Offset Device   Inode  Size   Rss   Pss Referenced Anonymous Swap Locked Mapping
        00400000 r-xp 00000000  fd:01  660707     4     4     4          4         0    0      0 memory-user.out
        00600000 r--p 00000000  fd:01  660707     4     4     4          4         4    0      0 memory-user.out
        00601000 rw-p 00001000  fd:01  660707     4     4     4          4         4    0      0 memory-user.out
    7fb104ad9000 r-xp 00000000  fd:01 1049989  1760   272     6        272         0    0      0 libc-2.17.so
    7fb104c91000 ---p 001b8000  fd:01 1049989  2048     0     0          0         0    0      0 libc-2.17.so
    7fb104e91000 r--p 001b8000  fd:01 1049989    16    16    16         16        16    0      0 libc-2.17.so
    7fb104e95000 rw-p 001bc000  fd:01 1049989     8     8     8          8         8    0      0 libc-2.17.so
    7fb104e97000 rw-p 00000000  00:00       0    20    12    12         12        12    0      0
    7fb104e9c000 r-xp 00000000  fd:01 1049982   132   108     2        108         0    0      0 ld-2.17.so
    7fb1050b1000 rw-p 00000000  00:00       0    12    12    12         12        12    0      0
    7fb1050bb000 rw-p 00000000  00:00       0     8     8     8          8         8    0      0
    7fb1050bd000 r--p 00021000  fd:01 1049982     4     4     4          4         4    0      0 ld-2.17.so
    7fb1050be000 rw-p 00022000  fd:01 1049982     4     4     4          4         4    0      0 ld-2.17.so
    7fb1050bf000 rw-p 00000000  00:00       0     4     4     4          4         4    0      0
    7fff60a88000 rw-p 00000000  00:00       0 10252 10252 10252      10252     10252    0      0 [stack]
    7fff615cc000 r-xp 00000000  00:00       0     8     4     0          4         0    0      0 [vdso]
ffffffffff600000 r-xp 00000000  00:00       0     4     0     0          0         0    0      0 [vsyscall]
                                              ===== ===== ===== ========== ========= ==== ======
                                              14292 10716 10340      10716     10328    0      0 KB
```                                             

### 8

Finally, let’s run pmap on your memory-user program, with dif-
ferent amounts of used memory. What do you see here? Does the
output from pmap match your expectations?