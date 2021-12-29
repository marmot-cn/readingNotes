#include <stdio.h>

#define swap(type, a, b) type temp = 0; temp = a; a = b; b = temp;

int main (void)
{
	int x = 5;
	int y = 10;

	swap(int, x, y);

	printf("x is %d, y is %d", x, y);
	return 0;
}