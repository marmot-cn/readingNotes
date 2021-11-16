#include <stdio.h>

int main(void)
{
	int no, i;

	i = 1;

	printf("请输入一个整数：");
	scanf("%d", &no);

	while(i <= no) {

		if (i%2)
			putchar('+');
		else
			putchar('-');

		i++;
	}

	return 0;
}