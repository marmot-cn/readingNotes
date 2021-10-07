#include <stdio.h>

int main(void)
{
	int i;
	
	printf("请输入一个整数:");
	scanf("%d", &i);

	printf("该整数减去 6 之后的结果是: %d\n", i-6);
}