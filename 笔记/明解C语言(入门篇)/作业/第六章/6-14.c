#include <stdio.h>

#define COUNT 5

double a[COUNT] = {};

void func(void)
{
	for (int i = 0; i < COUNT; i++) {
		printf("%.1f ", a[i]);
	}
}

int main (void) {

	func();
	return 0;
}