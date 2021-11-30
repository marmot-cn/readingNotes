#include <stdio.h>

void put_count() {
	static int count = 0;

	count++;

	printf("put_count: 第 %d 次 \n", count);
}

int main(void) {
	put_count();
	put_count();
	put_count();

	return 0;
}