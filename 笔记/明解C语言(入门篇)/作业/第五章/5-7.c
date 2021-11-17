#include <stdio.h>

#define MAX_COUNT 100

int main(void)
{
	int b[MAX_COUNT];

	int count, temp;

	printf("数据个数："); 
	scanf("%d", &count);

	for(int i=0; i<count; i++) {
		printf("%d号：", i+1);
		scanf("%d", &b[i]);
	}

	printf("{");
	for(int i=0; i<count; i++) {
		printf("%d", b[i]);

		//not last
		if (i != (count-1)) {
			printf(",");
		}
	}
	printf("}");
}