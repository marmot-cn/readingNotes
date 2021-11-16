#include <stdio.h>

int main(void)
{
	int n;

	printf("整数值：");
	scanf("%d", &n);

	for (int i =1; i <= n; i++) {
		if (i%2) {
			printf("%d ", i);
		}
	}

	return 0;
}