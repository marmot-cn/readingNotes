#include <stdio.h>

#define max(x,y) ((x) > (y)) ? (x) : (y)

int main (void)
{
	
	int a = 1;
	int b = 2;
	int c = 3;
	int d = 4;

	//两两比较
	printf("max is %d \n", max(max(a,b), max(c,d)));

	//逐个比较
	printf("max is %d \n", max(max(max(a,b),c),d));

	return 0;
}