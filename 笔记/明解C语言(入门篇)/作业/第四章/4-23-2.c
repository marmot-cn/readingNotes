#include <stdio.h>

int main(void)
{
	int i, j, len;

	puts("生成直角在右上方的等腰直角三角形。");
	printf("短边：");
	scanf("%d", &len);

	for (i = 0; i < len; i++) {
		for (j = 0; j < i; j++)
			putchar(' ');
		for (j = 0; j < len-i; j++)
			putchar('*');
		putchar('\n');
	}
	return 0;
}