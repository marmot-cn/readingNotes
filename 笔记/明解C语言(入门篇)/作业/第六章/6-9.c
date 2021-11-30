#include <stdio.h>

#define COUNT 5

void rev_intary(int v[], int n) {
	
	for (int i = n-1; i>=0; i--) {
		printf("%d ", v[i]);
	}
}

int main(void)
{
	int v[COUNT] = {3,4,5,9,10};

	rev_intary(v, COUNT);

	return 0;
}