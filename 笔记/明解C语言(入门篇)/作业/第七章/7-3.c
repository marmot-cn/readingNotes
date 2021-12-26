#include <stdio.h>

unsigned rrotate(unsigned x, int n) {

	return x >> n;
}

unsigned lrotate(unsigned x, int n) {

	return x << n;
}	

int main (void)
{
	int n = 4;

	printf("4 左移后是 %d\n", lrotate(n, 1));
	printf("4 右移后是 %d\n", rrotate(n, 1));

	return 0;
}