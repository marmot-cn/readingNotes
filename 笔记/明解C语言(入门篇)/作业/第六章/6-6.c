#include <stdio.h>

void alert(int n) {
	for(int i=0; i<n; i++) {
		putchar('\a');
	}
}

int main(void)
{
	int n;

	printf("请输入响铃次数:");
	scanf("%d", &n);

	alert(n);
	return 0;
}