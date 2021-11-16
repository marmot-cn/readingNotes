#include <stdio.h>

int main(void)
{
	int len, count, maxcount;

	puts("让我们来画一个向下的金字塔");
	printf("金字塔有几层：");
	scanf("%d", &len);

	for (int i = len; i >0; i--) {

		count = (i - 1) * 2 + 1;

		maxcount = (len - 1) * 2 + 1;

		for (int k = 1; k <= (maxcount - count)/2; k++)
			putchar(' ');

		for (int j = 1; j <= count; j++)
			printf("%d", (len-i+1)%10);

		// for (int k = 1; k <= (maxcount - count)/2; k++)
		// 	putchar(' ');

		putchar('\n');
	}
	return 0;
}