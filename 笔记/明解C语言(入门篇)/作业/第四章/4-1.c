#include <stdio.h>

int main(void)
{
	int cnt = 0;

	puts("请输入循环次数");

	printf("循环次数:");
	scanf("%d", &cnt);

	do {

		int no;

		printf("请输入一个整数：");
		scanf("%d", &no);

		if (no == 0)
			puts("该整数为0。");
		else if(no > 0)
			puts("该整数为正数。");
		else
			puts("该整数为负数。");

		cnt--;
	} while(cnt > 0);

	return 0;
}