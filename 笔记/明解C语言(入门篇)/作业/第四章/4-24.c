#include <stdio.h>

int main(void)
{
	int len, count, maxcount;

	puts("让我们来画一个金字塔");
	printf("金字塔有几层：");
	scanf("%d", &len);

	for (int i = 1; i <= len; i++) {

		count = (i - 1) * 2 + 1;

		maxcount = (len - 1) * 2 + 1;

		for (int k = 1; k <= (maxcount - count)/2; k++)
			putchar(' ');

		for (int j = 1; j <= count; j++)
			putchar('*');

		// for (int k = 1; k <= (maxcount - count)/2; k++)
		// 	putchar(' ');

		putchar('\n');
	}
	return 0;
}