# README

## 问题

### 1

编写一个名为`null.c`的简单程序，它创建一个指向整数的指针，将其设置为NULL，然后尝试对其进行释放内存操作。把它编译成一个名为`null`的可执行文件。当你运行这个程序时会发生什么？

**null.c**

```c
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

int main(int argc, char *argv[])
{
  int *p = NULL;
  free(p);
  return 0;
}
```

无任何输出

### 2

编译该程序，其中包含符号信息（`-g`标志）。这样做可以将更多信息放入可执行文件中，使调试器可以访问有关变量名称等的更多有用信息。通过输入`gdb null`在调试器下运行该程序，然后，一旦gdb运行，输入run。gdb显示什么信息。

`-g`产生符号调试工具（GNU的gdb）所必要的符号信息，想要对源代码进行调试。如通过`gdb null`调试代码时候，通过`l`选项可以查看源代码。

* 创建符号表，符号表包含了程序中使用的变量名称的列表。
* 关闭所有的优化机制，以便程序执行过程中严格按照原来的C代码进行。

#### 带`-g`选项**

```
[ansible@k8s-agent-1 ~]$ gcc -g null.c -o null
[ansible@k8s-agent-1 ~]$ gdb null
GNU gdb (GDB) Red Hat Enterprise Linux 7.6.1-120.el7
Copyright (C) 2013 Free Software Foundation, Inc.
License GPLv3+: GNU GPL version 3 or later <http://gnu.org/licenses/gpl.html>
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.  Type "show copying"
and "show warranty" for details.
This GDB was configured as "x86_64-redhat-linux-gnu".
For bug reporting instructions, please see:
<http://www.gnu.org/software/gdb/bugs/>...
Reading symbols from /home/ansible/null...done.
(gdb) run
Starting program: /home/ansible/null
[Inferior 1 (process 11632) exited normally]
Missing separate debuginfos, use: debuginfo-install glibc-2.17-196.el7.x86_64
(gdb) l
1	#include <stdio.h>
2	#include <stdlib.h>
3	#include <unistd.h>
4
5	int main(int argc, char *argv[])
6	{
7	  int *p;
8	  free(p);
9	  return 0;
10	}
(gdb) quit
```

#### 不带`-g`选项

```
ansible@k8s-agent-1 ~]$ gcc null.c -o null
[ansible@k8s-agent-1 ~]$ gdb null
GNU gdb (GDB) Red Hat Enterprise Linux 7.6.1-120.el7
Copyright (C) 2013 Free Software Foundation, Inc.
License GPLv3+: GNU GPL version 3 or later <http://gnu.org/licenses/gpl.html>
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.  Type "show copying"
and "show warranty" for details.
This GDB was configured as "x86_64-redhat-linux-gnu".
For bug reporting instructions, please see:
<http://www.gnu.org/software/gdb/bugs/>...
Reading symbols from /home/ansible/null...(no debugging symbols found)...done.
(gdb) l
No symbol table is loaded.  Use the "file" command.
(gdb) quit
```

### 3

对这个程序使用`valgrind`工具。我们将使用属于`valgrind`的`memcheck`工具来分析发送的情况。输入以下命令来运行程序: `valgrind --leak-check=yes null`。当你运行它时会发生什么？你能解释工具的输出吗?

```
[ansible@k8s-agent-1 ~]$ valgrind --leak-check=yes ./null
==21258== Memcheck, a memory error detector
==21258== Copyright (C) 2002-2017, and GNU GPL'd, by Julian Seward et al.
==21258== Using Valgrind-3.15.0 and LibVEX; rerun with -h for copyright info
==21258== Command: ./null
==21258==
==21258==
==21258== HEAP SUMMARY:
==21258==     in use at exit: 0 bytes in 0 blocks
==21258==   total heap usage: 0 allocs, 0 frees, 0 bytes allocated
==21258==
==21258== All heap blocks were freed -- no leaks are possible
==21258==
==21258== For lists of detected and suppressed errors, rerun with: -s
==21258== ERROR SUMMARY: 0 errors from 0 contexts (suppressed: 0 from 0)
```

那么有发送内存泄露。

#### 没有初始化变量的代码

在`null.c`中不初始化指针为`NULL`, 则会发现未初始化错误。

```
[ansible@k8s-agent-1 ~]$ valgrind --leak-check=yes ./null
==19661== Memcheck, a memory error detector
==19661== Copyright (C) 2002-2017, and GNU GPL'd, by Julian Seward et al.
==19661== Using Valgrind-3.15.0 and LibVEX; rerun with -h for copyright info
==19661== Command: ./null
==19661==
==19661== Conditional jump or move depends on uninitialised value(s)
==19661==    at 0x4C2B020: free (vg_replace_malloc.c:540)
==19661==    by 0x400547: main (null.c:8)
==19661==
==19661==
==19661== HEAP SUMMARY:
==19661==     in use at exit: 0 bytes in 0 blocks
==19661==   total heap usage: 0 allocs, 0 frees, 0 bytes allocated
==19661==
==19661== All heap blocks were freed -- no leaks are possible
==19661==
==19661== Use --track-origins=yes to see where uninitialised values come from
==19661== For lists of detected and suppressed errors, rerun with: -s
==19661== ERROR SUMMARY: 1 errors from 1 contexts (suppressed: 0 from 0)
```

### 4

编写一个使用`malloc()`来分配内存的简单程序，但在退出之前忘记释放它。这个程序运行时会发生什么？你可以用`gdb`来查找它的任何问题吗？用`valgrind`(`--leak-check=yes`标志)

**malloc.c**

```c
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

int main(int argc, char *argv[])
{
  char *str;
  str = (char *)malloc(15);
  return 0;
}
```

**gdb**

```
[ansible@k8s-agent-1 ~]$ gdb malloc
GNU gdb (GDB) Red Hat Enterprise Linux 7.6.1-120.el7
Copyright (C) 2013 Free Software Foundation, Inc.
License GPLv3+: GNU GPL version 3 or later <http://gnu.org/licenses/gpl.html>
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.  Type "show copying"
and "show warranty" for details.
This GDB was configured as "x86_64-redhat-linux-gnu".
For bug reporting instructions, please see:
<http://www.gnu.org/software/gdb/bugs/>...
Reading symbols from /home/ansible/malloc...done.
(gdb) run
Starting program: /home/ansible/malloc
[Inferior 1 (process 26606) exited normally]
```

**valgrind**

可见内存泄露

```
[ansible@k8s-agent-1 ~]$  valgrind --leak-check=yes ./malloc
==26933== Memcheck, a memory error detector
==26933== Copyright (C) 2002-2017, and GNU GPL'd, by Julian Seward et al.
==26933== Using Valgrind-3.15.0 and LibVEX; rerun with -h for copyright info
==26933== Command: ./malloc
==26933==
==26933==
==26933== HEAP SUMMARY:
==26933==     in use at exit: 15 bytes in 1 blocks
==26933==   total heap usage: 1 allocs, 0 frees, 15 bytes allocated
==26933==
==26933== 15 bytes in 1 blocks are definitely lost in loss record 1 of 1
==26933==    at 0x4C29F73: malloc (vg_replace_malloc.c:309)
==26933==    by 0x400545: main (malloc.c:7)
==26933==
==26933== LEAK SUMMARY:
==26933==    definitely lost: 15 bytes in 1 blocks
==26933==    indirectly lost: 0 bytes in 0 blocks
==26933==      possibly lost: 0 bytes in 0 blocks
==26933==    still reachable: 0 bytes in 0 blocks
==26933==         suppressed: 0 bytes in 0 blocks
==26933==
==26933== For lists of detected and suppressed errors, rerun with: -s
==26933== ERROR SUMMARY: 1 errors from 1 contexts (suppressed: 0 from 0)
```

### 5

编写一个程序，使用`malloc()`创建一个名为`data`、大小为`100`的整数数组。然后，然后将`data[100]`设置为`0`。当你运行这个程序时会发生什么？当你使用`valgrind`运行这个程序时会发生什么？程序是否正确？

**data.c**

```
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

int main(int argc, char *argv[])
{
  int *data = (int *)malloc(100 * sizeof(int));
  data[0] = 100;
  printf("data 0 is %d", data[0]);
  return 0;
}
```

`gdb`没有发现问题。

`valgrind`检测出错误。

```
[ansible@k8s-agent-1 ~]$  valgrind --leak-check=yes -s ./data
==7815== Memcheck, a memory error detector
==7815== Copyright (C) 2002-2017, and GNU GPL'd, by Julian Seward et al.
==7815== Using Valgrind-3.15.0 and LibVEX; rerun with -h for copyright info
==7815== Command: ./data
==7815==
data 0 is 100==7815==
==7815== HEAP SUMMARY:
==7815==     in use at exit: 400 bytes in 1 blocks
==7815==   total heap usage: 1 allocs, 0 frees, 400 bytes allocated
==7815==
==7815== 400 bytes in 1 blocks are definitely lost in loss record 1 of 1
==7815==    at 0x4C29F73: malloc (vg_replace_malloc.c:309)
==7815==    by 0x400595: main (data.c:7)
==7815==
==7815== LEAK SUMMARY:
==7815==    definitely lost: 400 bytes in 1 blocks
==7815==    indirectly lost: 0 bytes in 0 blocks
==7815==      possibly lost: 0 bytes in 0 blocks
==7815==    still reachable: 0 bytes in 0 blocks
==7815==         suppressed: 0 bytes in 0 blocks
==7815==
==7815== ERROR SUMMARY: 1 errors from 1 contexts (suppressed: 0 from 0)
```

### 6

创建一个分配整数数组的程序（如上所述），释放它们，然后尝试打印数组中某个元素的值。程序会运行吗？当你使用`valgrind`时会发生什么？

**data.c**

```
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

int main(int argc, char *argv[])
{
  int *data = (int *)malloc(100 * sizeof(int));
  data[0] = 100;
  free(data);
  printf("data 0 is %d", data[0]);
  return 0;
}
```

结果正常输出值，因为`free`释放该内存后, 只是标记为脏数据。

* free本身不会改变传入的指针指向
* 指针指向的那块内存的内容free不会进行更改
* 如果继续访问这块内存可能还是原来的内容、隔一段时间被其他程序修改后可能会变成其他内容。
* 访问这块内存空间暂时不非法
* free后可以赋值null来提高程序安全性

```
[ansible@k8s-agent-1 ~]$ ./data
data 0 is 100
```

`valgrind`提示非法读取。

```
[ansible@k8s-agent-1 ~]$  valgrind --leak-check=yes -s ./data
==11653== Memcheck, a memory error detector
==11653== Copyright (C) 2002-2017, and GNU GPL'd, by Julian Seward et al.
==11653== Using Valgrind-3.15.0 and LibVEX; rerun with -h for copyright info
==11653== Command: ./data
==11653==
==11653== Invalid read of size 4
==11653==    at 0x4005F4: main (data.c:10)
==11653==  Address 0x51fa040 is 0 bytes inside a block of size 400 free'd
==11653==    at 0x4C2B06D: free (vg_replace_malloc.c:540)
==11653==    by 0x4005EF: main (data.c:9)
==11653==  Block was alloc'd at
==11653==    at 0x4C29F73: malloc (vg_replace_malloc.c:309)
==11653==    by 0x4005D5: main (data.c:7)
==11653==
data 0 is 100==11653==
==11653== HEAP SUMMARY:
==11653==     in use at exit: 0 bytes in 0 blocks
==11653==   total heap usage: 1 allocs, 1 frees, 400 bytes allocated
==11653==
==11653== All heap blocks were freed -- no leaks are possible
==11653==
==11653== ERROR SUMMARY: 1 errors from 1 contexts (suppressed: 0 from 0)
==11653==
==11653== 1 errors in context 1 of 1:
==11653== Invalid read of size 4
==11653==    at 0x4005F4: main (data.c:10)
==11653==  Address 0x51fa040 is 0 bytes inside a block of size 400 free'd
==11653==    at 0x4C2B06D: free (vg_replace_malloc.c:540)
==11653==    by 0x4005EF: main (data.c:9)
==11653==  Block was alloc'd at
==11653==    at 0x4C29F73: malloc (vg_replace_malloc.c:309)
==11653==    by 0x4005D5: main (data.c:7)
==11653==
==11653== ERROR SUMMARY: 1 errors from 1 contexts (suppressed: 0 from 0)
```

### 7

现在传递一个有趣的值来释放（如，在上面分配的数组中间的一个指针）。会发生什么？你是否需要工具来找到这种类型的问题？


**free.c**

```
#include<stdlib.h>

int main()
{
  int *p = (int *)malloc(100*sizeof(int));
  free(p+50);
  return 0;
}
```

运行报错, 可通过`gdb`设置断点调试。

```
[ansible@k8s-agent-1 ~]$ ./free
*** Error in `./free': free(): invalid pointer: 0x00000000017890d8 ***
======= Backtrace: =========
/lib64/libc.so.6(+0x7c619)[0x7f6f2336c619]
./free[0x4005a5]
/lib64/libc.so.6(__libc_start_main+0xf5)[0x7f6f23311c05]
./free[0x4004b9]
======= Memory map: ========
00400000-00401000 r-xp 00000000 fd:01 655679                             /home/ansible/free
00600000-00601000 r--p 00000000 fd:01 655679                             /home/ansible/free
00601000-00602000 rw-p 00001000 fd:01 655679                             /home/ansible/free
01789000-017aa000 rw-p 00000000 00:00 0                                  [heap]
7f6f1c000000-7f6f1c021000 rw-p 00000000 00:00 0
7f6f1c021000-7f6f20000000 ---p 00000000 00:00 0
7f6f230da000-7f6f230ef000 r-xp 00000000 fd:01 1050972                    /usr/lib64/libgcc_s-4.8.5-20150702.so.1
7f6f230ef000-7f6f232ee000 ---p 00015000 fd:01 1050972                    /usr/lib64/libgcc_s-4.8.5-20150702.so.1
7f6f232ee000-7f6f232ef000 r--p 00014000 fd:01 1050972                    /usr/lib64/libgcc_s-4.8.5-20150702.so.1
7f6f232ef000-7f6f232f0000 rw-p 00015000 fd:01 1050972                    /usr/lib64/libgcc_s-4.8.5-20150702.so.1
7f6f232f0000-7f6f234a8000 r-xp 00000000 fd:01 1049989                    /usr/lib64/libc-2.17.so
7f6f234a8000-7f6f236a8000 ---p 001b8000 fd:01 1049989                    /usr/lib64/libc-2.17.so
7f6f236a8000-7f6f236ac000 r--p 001b8000 fd:01 1049989                    /usr/lib64/libc-2.17.so
7f6f236ac000-7f6f236ae000 rw-p 001bc000 fd:01 1049989                    /usr/lib64/libc-2.17.so
7f6f236ae000-7f6f236b3000 rw-p 00000000 00:00 0
7f6f236b3000-7f6f236d4000 r-xp 00000000 fd:01 1049982                    /usr/lib64/ld-2.17.so
7f6f238c8000-7f6f238cb000 rw-p 00000000 00:00 0
7f6f238d2000-7f6f238d4000 rw-p 00000000 00:00 0
7f6f238d4000-7f6f238d5000 r--p 00021000 fd:01 1049982                    /usr/lib64/ld-2.17.so
7f6f238d5000-7f6f238d6000 rw-p 00022000 fd:01 1049982                    /usr/lib64/ld-2.17.so
7f6f238d6000-7f6f238d7000 rw-p 00000000 00:00 0
7ffd83d3a000-7ffd83d5b000 rw-p 00000000 00:00 0                          [stack]
7ffd83db1000-7ffd83db3000 r-xp 00000000 00:00 0                          [vdso]
ffffffffff600000-ffffffffff601000 r-xp 00000000 00:00 0                  [vsyscall]
Aborted
```

**valgrind**

```
[ansible@k8s-agent-1 ~]$  valgrind --leak-check=yes -s ./free
==19368== Memcheck, a memory error detector
==19368== Copyright (C) 2002-2017, and GNU GPL'd, by Julian Seward et al.
==19368== Using Valgrind-3.15.0 and LibVEX; rerun with -h for copyright info
==19368== Command: ./free
==19368==
==19368== Invalid free() / delete / delete[] / realloc()
==19368==    at 0x4C2B06D: free (vg_replace_malloc.c:540)
==19368==    by 0x4005A4: main (free.c:6)
==19368==  Address 0x51fa108 is 200 bytes inside a block of size 400 alloc'd
==19368==    at 0x4C29F73: malloc (vg_replace_malloc.c:309)
==19368==    by 0x40058E: main (free.c:5)
==19368==
==19368==
==19368== HEAP SUMMARY:
==19368==     in use at exit: 400 bytes in 1 blocks
==19368==   total heap usage: 1 allocs, 1 frees, 400 bytes allocated
==19368==
==19368== 400 bytes in 1 blocks are definitely lost in loss record 1 of 1
==19368==    at 0x4C29F73: malloc (vg_replace_malloc.c:309)
==19368==    by 0x40058E: main (free.c:5)
==19368==
==19368== LEAK SUMMARY:
==19368==    definitely lost: 400 bytes in 1 blocks
==19368==    indirectly lost: 0 bytes in 0 blocks
==19368==      possibly lost: 0 bytes in 0 blocks
==19368==    still reachable: 0 bytes in 0 blocks
==19368==         suppressed: 0 bytes in 0 blocks
==19368==
==19368== ERROR SUMMARY: 2 errors from 2 contexts (suppressed: 0 from 0)
==19368==
==19368== 1 errors in context 1 of 2:
==19368== Invalid free() / delete / delete[] / realloc()
==19368==    at 0x4C2B06D: free (vg_replace_malloc.c:540)
==19368==    by 0x4005A4: main (free.c:6)
==19368==  Address 0x51fa108 is 200 bytes inside a block of size 400 alloc'd
==19368==    at 0x4C29F73: malloc (vg_replace_malloc.c:309)
==19368==    by 0x40058E: main (free.c:5)
==19368==
==19368== ERROR SUMMARY: 2 errors from 2 contexts (suppressed: 0 from 0)
```

### 8

尝试一些其他接口来分配内存。例如，创建一个简单的向量似的数据结构，以及使用`realloc()`来管理向量的相关函数。使用数组来存储向量元素。当用户在向量中添加条目时，使用`realloc()`为其分配更多空间。这样的向量表现如何？它与链表相比如何？使用`valgrind`来帮助你发现错误。

`void *realloc(void *ptr, size_t size)` 尝试重新调整之前调用 `malloc` 或 `calloc` 所分配的 `ptr` 所指向的内存块的大小。

```
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
 
int main()
{
   char *str;
 
   /* 最初的内存分配 */
   str = (char *) malloc(15);
   strcpy(str, "runoob");
   printf("String = %s,  Address = %p\n", str, str);
 
   /* 重新分配内存 */
   str = (char *) realloc(str, 25);
   strcat(str, ".com");
   printf("String = %s,  Address = %p\n", str, str);
 
   free(str);
   
   return(0);
}
```

### 9

花更多时间阅读有关使用`gdb`和`valgrind`的信息。了解你的工具至关重要，花时间学习如何成为`UNIX`和`C`环境中的调试器专家。

`valgrind`内存泄露检查工具。

## 其他

因为做 gdb 调试的时候，需要一些依赖源安装。

```
[root@iZbp1b0n9ivu1hyf17tgsfZ ~]# cat /etc/yum.repos.d/CentOS-Debug.repo

#Debug Info
[debug]
name=CentOS-$releasever - DebugInfo
baseurl=http://debuginfo.centos.org/$releasever/$basearch/
gpgcheck=0
enabled=1
protect=1
priority=1
```
