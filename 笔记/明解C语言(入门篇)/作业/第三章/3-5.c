#include <stdio.h>

int main(void)
{
	int a, b;

	printf("请输入两个整数:");

	printf("整数A:");
	scanf("%d", &a);

	printf("整数B:");
	scanf("%d", &b);

	//相等 a==b 返回 1, 否则返回 0
	printf("A和B的比较结果是%d", a==b);

	return 0;
}