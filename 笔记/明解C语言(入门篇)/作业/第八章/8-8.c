#include <stdio.h>

int gcd(int x, int y)
{
	int remainder;
	remainder = x%y;

	if (remainder == 0) {
		return y;
	}

	return (remainder > y) ? gcd(remainder, y) : gcd(y, remainder);
}

int main (void)
{
	int x, y, result;

	printf("请输入第一个数：");
	scanf("%d", &x);

	printf("请输入第二个数：");
	scanf("%d", &y);

	if (x > y) {
		result = gcd(x, y);
	} else {
		result = gcd(y, x);
	}

	printf("公约数数为 %d", result);
	return 0;
}