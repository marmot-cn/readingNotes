#include <stdio.h>

int main(void)
{
	int level;
	
	level = 0;

	printf("生成一个正方形\n正方形有几层：");
	scanf("%d", &level);

	for (int i=1; i <=level; i++) {
		for (int j=1; j<=level; j++) {
			putchar('*');
		}
		putchar('\n');
	}

	return 0;
}