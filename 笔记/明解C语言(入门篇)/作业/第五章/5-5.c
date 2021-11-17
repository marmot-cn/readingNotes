#include <stdio.h>

#define COUNT 5

int main(void)
{
	int i;
	int x[COUNT];

	for (i = 0; i < COUNT; i++) {
		printf("x[%d]: ", i);
		scanf("%d", &x[i]);
	}

	for (i = 0; i < COUNT/2; i++) {
		int temp = x[i];
		x[i] = x[COUNT-1- i];
		x[COUNT-1-i] = temp;
	}

	puts("倒序排列了。");
	for (i = 0; i < COUNT; i++)
		printf("x[%d] = %d\n", i, x[i]);

	return 0;
}