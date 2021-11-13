#include <stdio.h>

int main(void)
{
	int a, b, c, min;

	//必须要使用if
	printf("请输入三个整数:");

	printf("整数A:");
	scanf("%d", &a);

	printf("整数B:");
	scanf("%d", &b);

	printf("整数C:");
	scanf("%d", &c);

	if (a < b)
		min = a;
	else 
		min =b;	

	if (min > c)
		min = c;

	printf("最小值为%d", min);
	return 0;
}