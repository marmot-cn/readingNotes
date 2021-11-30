#include <stdio.h>

int sumup(int n) {

	int temp = 0;

	for (int i=1; i<=n; i++) {
		temp += i;
	}

	return temp;
}

int main(void)
{
	int n;

	printf("请输入整数:");
	scanf("%d", &n);

	printf("1 到 %d 的和是 %d", n, sumup(n));
	return 0;	
}