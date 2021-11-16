#include <stdio.h>

int main(void)
{
	int start, end, interval;

	printf("开始数值（cm）：");
	scanf("%d", &start);

	printf("结束数值（cm）：");
	scanf("%d", &end);

	printf("间隔数值（cm）：");
	scanf("%d", &interval);

	for (int i = start; i <= end; i+=interval) {
		printf("%dcm    %.2fkg\n", i, (i-70)*0.6);
	}

	return 0;
}