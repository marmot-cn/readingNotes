#include <stdio.h>

int main(void)
{
	int n;

	printf("请输入一个整数：");
	scanf("%d", &n);

	for (int i =1; i <= n; i++) {
		printf("%d", i%10);
	}

	return 0;
}