#include <stdio.h>

#define NUMBER 80

int main(void)
{
	int i, j, m, max, space;
	int num;
	int tensu[NUMBER];
	int bunpu[11] = {0};

	printf("请输入学生人数：");

	do {
		scanf("%d", &num);
		if (num <1 || num > NUMBER)
			printf("\a 请输入 1 ~ %d的数：", NUMBER);
	} while(num < 1 || num > NUMBER);

	printf("请输入%d人的分数。\n", num);

	for (i = 0; i < num; i++) {
		printf("%2d号：", i + 1);
		do {
			scanf("%d", &tensu[i]);
			if (tensu[i] < 0 || tensu[i] > 100)
				printf("\a 请输入 1~100的数：");
		} while (tensu[i] < 0 || tensu[i] > 100);

		bunpu[tensu[i] / 10]++;
	}


	//找见最大的分数区间个数
	for (i = 0; i<=10; i++) {
		if (max <= bunpu[i]) {
			max = bunpu[i];
			// printf("max: %d, i is %d \n", bunpu[i], i);
		}
	}

	for (i = max; i >0; i--) {
		space = 0;
		for (j = 0; j <=10; j++)  {
			//如果数量匹配则显示
			if (bunpu[j] >= i) {
				printf(" *  ");
			} else {
				printf("    ");
			}
			
		}
		printf("\n");
	}

	putchar('\n');
    printf("---------------------------------------------\n");
    for (j = 0; j <= 100; j += 10)
    {
        printf(" %d ", j);
    }

	return 0;
}