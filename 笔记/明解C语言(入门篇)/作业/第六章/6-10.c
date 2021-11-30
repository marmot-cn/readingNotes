#include <stdio.h>

#define COUNT 5

void intary_rcpy(int v1[], const int v2[], int n) {
	
	int temp = 0;
	for (int i = n-1; i>=0; i--) {
		v1[temp++] = v2[i];
	}
}

int main(void)
{
	int v1[COUNT];

	int v2[COUNT] = {3,4,5,9,10};

	intary_rcpy(v1, v2, COUNT);

	for (int i = 0; i < COUNT; i++) {
		printf("%d ", v1[i]);
	}

	return 0;
}