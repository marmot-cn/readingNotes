#include <stdio.h>

void mat_mul(const int a[4][3], const int b[3][4], int c[4][4]) {

	for (int i = 0; i < 4; i++) {
		for (int j = 0; j < 4; j++) {
			c[i][j] = a[i][0] * b[0][j] + a[i][1] * b[1][j] + a[i][2] * b[2][j];
		}
	}
}

int main(void)
{
	int c[4][4];

	int a[4][3] = {
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
		{10, 11, 12}
	};

	int b[3][4] = {
		{1, 2, 3, 4},
		{5, 6, 7, 8},
		{9, 10, 11, 12}
	};

	mat_mul(a, b, c);

	for (int i = 0; i < 4; i++) {
		for (int j = 0; j < 4; j++) {
			printf("%d ", c[i][j]);
		}
		putchar('\n');
	}

	return 0;
}