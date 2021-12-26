#include <stdio.h>

int main (void)
{
	int n = 4;

	printf("4 左移后是 %d\n", n<<1);
	printf("4 右移后是 %d\n", n>>1);

	return 0;
}