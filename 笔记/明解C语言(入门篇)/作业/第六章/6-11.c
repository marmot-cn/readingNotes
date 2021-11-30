#include <stdio.h>

#define COUNT 7

int search_idx(const int v[], int idx[], int key, int n) {
	
	int count = 0;
	for (int i = 0; i < n; i++) {
		if (v[i] == key) {
			idx[count++] = i;
		}
	}

	return count;
}

int main(void)
{
	int duplicateCount;
	int idx[COUNT] = {-1, -1, -1, -1, -1, -1, -1};
	int v[COUNT] = {1, 7, 5, 7, 2, 4, 7};

	int key = 7;

	duplicateCount = search_idx(v, idx, key, COUNT);

	printf("duplicate item count is %d \n", duplicateCount);
	for (int i = 0; i < COUNT; i++) {

		if (idx[i] > 0) printf("%d ", idx[i]);
	}

	return 0;
}