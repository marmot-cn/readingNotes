#include <stdio.h>

int main(void)
{
	int i, j, m;
	int tensu[2][4][3] = {
		{{91, 63, 78}, {67, 72, 46}, {89, 34, 53}, {32, 54, 34}},
		{{97, 67, 82}, {73, 43, 46}, {97, 56, 21}, {85, 46, 35}}
	};
	int sum[4][3];

	for (i = 0; i < 4; i++) {
		for (j = 0; j < 3; j++) {
			sum[i][j] = tensu[0][i][j] + tensu[1][i][j];
		}
	}

	for (m = 0; m <= 1; m++) {
		printf("第 %d 次考试的分数\n", m+1);
		for (i = 0; i < 4; i++) {
			for (j = 0; j < 3; j++) {
				printf("%4d", tensu[m][i][j]);
			}
			putchar('\n');
		}
	}

	puts("总分");
	for (i = 0; i < 4; i++) {
		for (j = 0; j < 3; j++) {
			printf("%4d", sum[i][j]);
		}
		putchar('\n');
	}

	return 0;
}