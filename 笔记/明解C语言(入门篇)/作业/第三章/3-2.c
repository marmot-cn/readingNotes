#include <stdio.h>

int main(void)
{
	int no;

	printf("请输入一个整数:");
	scanf("%d", &no);

	if (no==0)
		puts("该整数为0。");
	else if(no>0)
		puts("该整数位正数。");
	// else 
	// 	puts("该整数为负数。");
	else if(no<0)
		puts("该整数为负数。");

	//如果最后一个else改为else if(no<0), 结果不会有影响
	return 0;
}