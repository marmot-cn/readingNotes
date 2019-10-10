# 21 动态内存,typedefine,Makefile

[动态内存分配-malloc和realloc的使用](http://www.wyzc.com/play/8704/2992/#12652 "动态内存分配-malloc和realloc的使用")

[动态内存分配-free的使用及微型学生管理系统的代码重构](http://www.wyzc.com/play/8704/2992/#12653 "动态内存分配-free的使用及微型学生管理系统的代码重构")

[如何使用重定义typedefine](http://www.wyzc.com/play/8704/2992/#12654 "如何使用重定义typedefine")

[Makefile工程文件的编写规则](http://www.wyzc.com/play/8704/2983/#12552 "Makefile工程文件的编写规则")

###笔记

---

####动态内存

**malloc ralloc realloc free**

`void *malloc(size_t size)`;

`size_t` : `typedefine` 的 `整型`.

在堆上申请size个连续的内存空间. 把起始地址返回(`void *`赋值给任何类型的指针,因为只是申请连续的内存空间,函数不知道存放的`类型`,所以使用`void *`).

		#include <stdlib.h>//必须包含头文件
		
		int main(){
		
			int *p = NULL;
			p = malloc(sizeof(int));//没有必要强制转换(int *)malloc(sizeof(int));
			if(p == NULL){
				printf("malloc() error!\n");
				exit(1);
			}
		
			*p = 10;
			free(p);//free后空间已经分配给其他人使用了
			p = NULL;//free后把指针指为空
			
			exit(0);
		}
		
`示例:动态实现数组`:

		int main(){
			
			int *p;
			int num = 5;
			p = malloc(sizeof(int)*num);
			if(p == NULL){
				printf("malloc() error!\n");
				exit(1);
			}
			
			free(p);						
			exit(0);
		}	
		
**原则**

谁申请,谁释放. 防止产生内存泄露.

####typedef

为已有的数据类型改名.

`typedef` `已有的数据类型` `新名字` `;`

**typedef 函数**

1. `typedef` `int` `(init_fnc_t)` `(void)`;表示定义init_fnc_t为`函数类型`,该函数返回int型,无参数

		init_fnc_t *init_sequence[]={ cpu_init,  board_init };
		
		表示用init_fnc_t(函数类型)去定义一个一维指针数组,数组中的元素都是指针变量,而且都是指向函数的指针,这些函数返回值都是int型,无参数的,数组中的每个元素是用来存放函数入口首地址的.

2. `typedef` `int` `(*init_fnc_t)` `(void)`;表示定义一个`函数指针类型`
		
		init_fnc_t init_sequence[]={cpu_init,  board_init }
		
		表示用init_fnc_t(函数指针类型)去定义一个数组,数组里面的元素都是一个函数指针,cpu_init,  board_init都是函数指针,存放的都是函数的首地址.


`1`中的数组元素是`指向函数`的`指针变量`,`2`中的数组的元素是`函数指针`.(`2者定义的意思一样,写法不一样`)得到的结果是一样的.

**示例**		
		
		typedef int INT;
		
		#define INT int;
		
		int main(){
		
			INT i;//等同于 int i
			
			exit(0);
		}


* `typedef int INT` 和  `#define INT int`

		INT i; --> int i;一样
		
* `typedef int *IP` 和 `#define IP int *`
		
		//type 定义2个指针
		IP p,q; -> int *p,*q;
		//#define 定义一个整型变量和一个整型指针
		IP p,q; -> int *p,q;
		
		2者不同
		
* `typedef int ARR[6]` : `int[6] -> ARR`

		ARR a; -> int a[6];

* `struct`

		struct node_st{
			int i;
			float f;
		};
		typedef struct node_st NODE;
		
		NODE a; --> struct node_st a;
		
		typedef struct node_st *NODEP;
		
		NODEP p; -> struct node_st *p;
		
		//也可以这样定义
		typedef struct{
			int i;
			float f;
		}NODE, *NODEP;
		
* `函数`

		typedef int FUNC(int); -> int(int) FUNC;
		FUNC f; --> int f(int);
		
		//指针函数 
		typedef int *FUNCP(int);
		
		FUNCP p; -> int *p(int);
		
		//函数指针
		typedef int *(*FUNCP)(int);
		FUNCP p; --> int *(*p)(int);


####makefile工程文件

**make**

一个工程管理器.执行的脚本是`makefile`.

**makefile**

###整理知识点

---