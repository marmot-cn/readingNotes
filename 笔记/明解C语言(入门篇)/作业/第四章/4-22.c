#include <stdio.h>

int main(void)
{
	int i, j;
	int sidea, sideb, maxside, minside;

	puts("让我们来画一个长方形。");
	printf("一边："); scanf("%d", &sidea);
	printf("另一边："); scanf("%d", &sideb);

	sidea > sideb ? (maxside = sidea, minside = sideb) : (maxside = sideb, minside = sidea);

	for (i=1; i <= minside; i++) {
		for (j=1; j <= maxside; j++)
			putchar('*');
		putchar('\n');
	}

	return 0;
}