#include <stdio.h>

int main(void)
{
	double a;
	int b;

	a = b = 1.5;

	//1.000000, 小数点后默认保留6位
	printf("a: %f\n", a);

	//1
	printf("b: %d\n", b);

	return 0;
}