#include <stdio.h>

int main(void)
{
	int n1, n2;

	puts("请输入两个整数。");
	printf("整数1："); scanf("%d", &n1);
	printf("整数2："); scanf("%d", &n2);

	if ((n1 - n2) < -10 || (n1 - n2) > 10) 
		puts("它们的差大于11");
	else
		puts("它们的差小于等于10");
}