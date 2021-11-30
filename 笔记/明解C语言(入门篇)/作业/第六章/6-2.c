#include <stdio.h>

int min3 (int a, int b, int c) {

	int min = a;

	if (b < min) min = b;
	if (c < min) min = c;
	return min;
}

int main(void)
{
	int a, b, c;

	puts("请输入三个数:");
	printf("整数a:"); scanf("%d", &a);
	printf("整数b:"); scanf("%d", &b);
	printf("整数v:"); scanf("%d", &c);

	printf("较小的数是 %d", min3(a, b, c));
	return 0;	
}