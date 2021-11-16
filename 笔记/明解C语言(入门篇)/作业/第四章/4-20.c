#include <stdio.h>

int main(void)
{
	int x, y;

	putchar('|');

	for (x = 1; x <= 9; x++) {
		printf("%d ", x);
	}
	putchar('\n');
	puts("--------------------------");

	for (x = 1; x <= 9; x++) {
		printf("%d | ", x);
		for (y = 1; y <= 9; y++) {
			printf("%d ", x*y);
		}
		putchar('\n');
	}

	return 0;
}