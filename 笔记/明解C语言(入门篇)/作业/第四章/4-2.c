#include <stdio.h>

int main(void)
{
	int a, b, min, max;

	int cnt = 0;
	int sum = 0;

	puts("请输入两个整数");

	printf("整数a:");
	scanf("%d", &a);

	printf("整数b:");
	scanf("%d", &b);

	a < b ? (min = a, max = b) : (min = b, max = a);

	do {
		sum = sum + min;
	} while(++min <= max);

	printf("大于等于%d小于等于%d的所有整数的和是%d", min, max, sum);
	return 0;
}