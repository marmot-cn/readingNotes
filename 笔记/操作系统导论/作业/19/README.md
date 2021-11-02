# README

## 问题

### 1

`gettimeofday()`

```
[ansible@k8s-agent-1 ~]$ cat gettimeofdat.c
#include<sys/time.h>
#include<unistd.h>
#include <stdio.h>

int main() {

  struct timeval tv;
  struct timezone tz;

  gettimeofday (&tv , &tz);

  printf("tv_sec; %d\n", tv.tv_sec) ;
  printf("tv_usec; %d\n", tv.tv_usec);
  printf("tz_minuteswest; %d\n", tz.tz_minuteswest);
  printf("tz_dsttime, %d\n", tz.tz_dsttime);

  return 0;
}

[ansible@k8s-agent-1 ~]$ gcc gettimeofdat.c
[ansible@k8s-agent-1 ~]$ ./a.out
tv_sec; 1634196227
tv_usec; 225954
tz_minuteswest; -480
tz_dsttime, 0
```

### 2

写一个程序，命名为`tlb.c`, 大体测算一下每个页的评价访问时间。程序的输入参数有：页的数目和尝试的次数。

页大小

```
[ansible@k8s-agent-1 ~]$ getconf PAGE_SIZE
4096
```

**19.2.c**

```
#include<stdio.h>
#include<sys/time.h>
#include<stdlib.h>
 
#define PAGESIZE 4096
 
int main(int argc, char *argv[])
{
	if(argc != 3) {
		fprintf(stderr, "error parameters!");
		exit(0);
	}
	struct timeval start, end;
	int pageNum = atoi(argv[1]);
	int i,j, num = atoi(argv[2]);
	char arr[PAGESIZE * pageNum];
	
	gettimeofday(&start, NULL);
	for(j=0; j<num; ++j) {
		for(i=0; i<pageNum; ++i) {
			arr[i*4096]=1;
		}
	}
	gettimeofday(&end, NULL);
	printf("%lf %d %d\n", (((double)end.tv_usec - start.tv_usec)/pageNum)/num, end.tv_usec, start.tv_usec);
	return 0;
}
```

### 5 

禁止优化`gcc -O0`

### 6 

代码从一个CPU移到了另一个，TLB未命中，会增加访问时间

**cpu.c**

通过`pthread_setaffinity_np`可以绑定`thread`到具体CPU, 同时通过`pthread_getaffinity_np`可以验证绑定到哪个CPU核心(把值写会到cpuset)。

```
#define _GNU_SOURCE
#include <pthread.h>
#include <stdio.h>
#include <stdlib.h>
#include <errno.h>

#define handle_error_en(en, msg) \
       do { errno = en; perror(msg); exit(EXIT_FAILURE); } while (0)

int main(int argc, char *argv[])
{
   //cpu core number
   printf("system cpu num is %d\n", get_nprocs_conf());
   printf("system enable num is %d\n", get_nprocs());

   int s, j;
   cpu_set_t cpuset;
   pthread_t thread;

   thread = pthread_self();

   /* Set affinity mask to include CPUs 0 to 7 */

   CPU_ZERO(&cpuset);
   //set run on cpu 0
   CPU_SET(0, &cpuset);
   CPU_SET(1, &cpuset);
   printf("cpu affinity before set on CPU 0 is %d\n", s);

   s = pthread_setaffinity_np(thread, sizeof(cpu_set_t), &cpuset);
   if (s != 0)
       handle_error_en(s, "pthread_setaffinity_np");

   /* Check the actual affinity mask assigned to the thread */
   
   CPU_ZERO(&cpuset);
   s = pthread_getaffinity_np(thread, sizeof(cpu_set_t), &cpuset);
   printf("cpu affinity after set on CPU 0 is %d\n", s);
   if (s != 0)
       handle_error_en(s, "pthread_getaffinity_np");

   printf("Set returned by pthread_getaffinity_np() contained:\n");
   for (j = 0; j < CPU_SETSIZE; j++)
       if (CPU_ISSET(j, &cpuset))
           printf("    CPU %d\n", j);
}
```

### 7

在开始计时前把整个数组初始化即可。

**19.7.c**

```
#include<stdio.h>
#include<sys/time.h>
#include<stdlib.h>
 
#define PAGESIZE 4096
 
int main(int argc, char *argv[])
{
	if(argc != 3) {
		fprintf(stderr, "error parameters!");
		exit(0);
	}
	struct timeval start, end;
	int pageNum = atoi(argv[1]);
	int i,j, num = atoi(argv[2]);
	char arr[PAGESIZE * pageNum];
	for(i=0; i<PAGESIZE*pageNum; ++i) {
		arr[i] = 0;
	}
	
	gettimeofday(&start, NULL);
	for(j=0; j<num; ++j) {
		for(i=0; i<pageNum; ++i) {
			arr[i*4096]=1;
		}
	}
	gettimeofday(&end, NULL);
	printf("%lf %d %d\n", (((double)end.tv_usec - start.tv_usec)/pageNum)/num, end.tv_usec, start.tv_usec);
	return 0;
}
```