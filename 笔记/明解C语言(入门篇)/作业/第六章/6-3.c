#include <stdio.h>

int cube(int x) {

	return x * x * x;
}

int main(void)
{
	int x;

	printf("请输入整数:");
	scanf("%d", &x);

	printf("%d 的立方是: %d", x, cube(x));
	return 0;	
}