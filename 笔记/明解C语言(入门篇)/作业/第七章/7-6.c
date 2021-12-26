#include <stdio.h>

int main (void)
{	
	unsigned x = 65535;

	//均为 65536, 无符号整形能够表示的最大值 + 1
	printf("x + 1 = %u\n", x+1);
	printf("x + 2 = %u\n", x+1);
	printf("x + 100 = %u\n", x+1);

	return 0;
}