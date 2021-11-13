#include <stdio.h>

int main(void)
{
	int a, b, c, d, max;

	//必须要使用if
	puts("请输入四个整数:");

	printf("整数A:");
	scanf("%d", &a);

	printf("整数B:");
	scanf("%d", &b);

	printf("整数C:");
	scanf("%d", &c);

	printf("整数D:");
	scanf("%d", &d);

	if (a > b)
		max = a;
	else 
		max =b;	

	if (max < c)
		max = c;

	if (max < d)
		max = d;

	printf("最大值为%d", max);
	return 0;
}