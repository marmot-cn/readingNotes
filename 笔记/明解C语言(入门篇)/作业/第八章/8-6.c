#include <stdio.h>

int factorial(int n)
{
	int i, result = 1;
	for (i = 1; i <= n; i++)
	{
		result = result * i;
	}

	return result;
}

int main (void)
{
	int num;

	printf("请输入一个整数：");
	scanf("%d", &num);
	printf("%d 的阶乘为 %d", num, factorial(num));
	return 0;
}