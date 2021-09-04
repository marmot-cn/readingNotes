#include <stdio.h>
#include <pthread.h>
#include <stdlib.h>


typedef struct { int a; int b; } myarg_t;
typedef struct { int x; int y; } myret_t;

void *mythread(void *arg) {

	myarg_t *m = (myarg_t *)arg;
	printf("%d %d\n", m->a, m->b);

	//局部变量的作用域在函数的内部，函数返回后，局部变量内存已经释放掉了。
	//因此，如果函数返回的是局部变量的值，不涉及地址，程序不会出错。
	//故函数是可以返回局部变量的。但是如果函数返回的是局部变量的地址(指针)的话，程序运行就会报错。
	//由于只是把指针（地址）复制后返回了，而指针指向的内存（存储的内容）已经被释放了，这样指针无访问该内存的权限，调用后就会报错。
	//准确的说，函数不能通过返回指向栈内存的指针(注意这里指的是栈，返回指向堆内存的指针是可以的

	// myret_t rvals;
	// rvals.x = 1;
	// rvals.y = 2;
	// return &rvals;

	//申请的内存
	myret_t *rvals = malloc(sizeof(myret_t));
	rvals->x = 1;
	rvals->y = 2;
	return (void *) rvals;
}

int main(int argc, char *argv[]) {
	
	pthread_t p;

	myret_t *rvals;
	myarg_t args = { 10, 20 };
	
	pthread_create(&p, NULL, mythread, &args);
	pthread_join(p, (void **) &rvals);
	
	printf("returned %d %d\n", rvals->x, rvals->y);
	free(rvals);
	
	return 0;
}