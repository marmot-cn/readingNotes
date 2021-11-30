#include <stdio.h>

int main(void)
{
	int i, j, sum;
	int students[6][2];

	printf("输入学生的分数\n");

	for (i = 0; i<6; i++) {
		printf("学生 %d, 的语文和数学分数是(用空格分隔): ", i+1);
		scanf("%d%d", &students[i][0], &students[i][1]);
	}	

	printf("计算总数与平均数\n");
	sum = 0;
	for (i = 0; i<6; i++) {
		sum = sum + students[i][0] + students[i][1];
	}

	printf("总数为是 %d, 平均数是 %d\n", sum, sum/6);

	printf("计算每个学生的总数与平均数\n");
	for (i = 0; i<6; i++) {
		printf("学生 %d, 的总数是 %d 平均数是%d\n: ", i+1, students[i][0] + students[i][1], (students[i][0] + students[i][1])/2);
	}
}