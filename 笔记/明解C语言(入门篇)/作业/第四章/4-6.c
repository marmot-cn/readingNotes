#include <stdio.h>

int main(void)
{
	int i, no;

	printf("请输入一个整数：");
	scanf("%d", &no);

	i = 1;
	while( i <= no) {
		if (!(i%2))
			printf("%d ", i);
		i++;
	}

	return 0;
}