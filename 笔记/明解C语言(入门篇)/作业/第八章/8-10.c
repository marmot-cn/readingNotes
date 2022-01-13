#include <stdio.h>

int main (void)
{
	int i, ch;
	int cnt[10] = {0};

	while((ch = getchar()) != EOF) {
		if (ch >= '0' && ch <= '9')
			cnt[ch - '0']++;
	}

	puts("数字字符的出现次数");
	for (i = 0; i < 10; i++) {
		printf("%d : ", i);
		for (int j = 0; j < cnt[i]; j++) {
			printf("%c", '*');
		}
		printf("\n");
	}

	return 0;
}