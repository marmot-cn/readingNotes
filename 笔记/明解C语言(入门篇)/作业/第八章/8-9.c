#include <stdio.h>

int main (void)
{
	int ch, i;

	i = 0;

	while ((ch = getchar()) != EOF) {
		if (ch == '\n') {
			i++;
		}
	}

	printf("换行出现 %d 次", i);

	return 0;
}