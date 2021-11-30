#include <stdio.h>

int min2 (int a, int b) {

	return a < b ? a : b;
}

int main(void)
{
	int a, b;

	printf("请输入两个数:");
	scanf("%d%d", &a, &b);

	printf("较小的数是 %d", min2(a,b));
	return 0;	
}