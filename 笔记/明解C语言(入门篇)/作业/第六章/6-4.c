#include <stdio.h>

int pow4(int x) {

	return x * x * x * x;
}

int main(void)
{
	int x;

	printf("请输入整数:");
	scanf("%d", &x);

	printf("%d 的4次幂是: %d", x, pow4(x));
	return 0;	
}