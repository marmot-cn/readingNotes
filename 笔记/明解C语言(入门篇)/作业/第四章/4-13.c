#include <stdio.h>

int main(void)
{
	int n, sum;

	sum = 0;

	printf("n的值为：");
	scanf("%d", &n);

	printf("1到%d的和为", n);
	for (int i =1; i <= n; i++) {
		sum+=i;
	}
	
	printf("%d。", sum);

	return 0;
}