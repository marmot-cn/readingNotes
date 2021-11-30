#include <stdio.h>

int main(void)
{	
	int i, j, temp;
	int matrix_a[3][4];
	int matrix_b[4][3];

	printf("请输入3行4列的元素\n");

	for (i = 0; i<3; i++) {
		for (j = 0; j<4; j++) {
			printf("元素 %d 行, %d 列：", i+1, j+1);
			scanf("%d", &matrix_a[i][j]);
		}
	}

	printf("计算3行4列的乘积\n");

	temp = 1;
	for (i = 0; i<3; i++) {
		for (j = 0; j<4; j++) {
			temp = temp * matrix_a[i][j];
		}
	}

	printf("计算3行4列的乘积为%d\n", temp);

	printf("请输入4行3列的元素\n");

	for (i = 0; i<4; i++) {
		for (j = 0; j<3; j++) {
			printf("元素 %d 行, %d 列：", i+1, j+1);
			scanf("%d", &matrix_b[i][j]);
		}
	}

	printf("计算4行3列的乘积\n");

	temp = 1;
	for (i = 0; i<4; i++) {
		for (j = 0; j<3; j++) {
			temp = temp * matrix_b[i][j];
		}
	}

	printf("计算4行3列的乘积为%d\n", temp);
}