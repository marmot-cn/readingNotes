#include <stdio.h>

int main(void)
{
	int n1, n2, result;

	puts("请输入两个整数。");
	printf("整数1："); scanf("%d", &n1);
	printf("整数2："); scanf("%d", &n2);

	if (n1 > n2)
		result = n1-n2;
	else
		result = n2-n1;

	printf("它们的差是%d。\n", result);
}