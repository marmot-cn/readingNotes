#include <stdio.h>

int main(void)
{

	//计算机不能保证其内部转换为二进制的浮点数的每一位都不发生数据丢失

	float x, y, z = 0.0;

	printf("7-11 \n");
	for (x = 0.0; x <= 1.0; x += 0.01) {
		printf("x=%f  ", x);
		y += x;
		printf("x的累计和=%f\n",y);
	}

	printf("7-12 \n");
	for (int i = 0; i <= 100; i++) {
		printf("x=%f   ", i / 100.0);
		printf("x的累计和=%f\n", z += i/100.0);
	}
	return 0;
}