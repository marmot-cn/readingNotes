#include <stdio.h>

#define COUNT 5

int min_of(const int v[], int n) {
	
	int min = v[0];

	for (int i = 1; i<n; i++) {
		if (min > v[i]) {
			min = v[i];
		}
	}

	return min;
}

int main(void)
{
	int v[COUNT] = {3,4,5,9,10};

	printf("minimum value is %d\n", min_of(v, COUNT));

	return 0;
}