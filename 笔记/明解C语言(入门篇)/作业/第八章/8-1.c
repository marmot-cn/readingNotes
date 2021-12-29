#include <stdio.h>

#define diff(x, y) ((x) - (y))

int main (void)
{
	int x = 5;
	int y = 4;

	printf("diff is %d", diff(x,y));
	return 0;
}