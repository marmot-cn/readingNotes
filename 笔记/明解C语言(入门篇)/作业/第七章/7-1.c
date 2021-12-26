#include <stdio.h>

int main (void)
{
	int n = 0;

	//4
	printf("sizeof 1 is %lu\n", sizeof 1);
	//4
	printf("sizeof +1 is %lu\n", sizeof +1);
	//4
	printf("sizeof -1 is %lu\n", sizeof -1);

	//3 (4-1)
	printf("sizeof -1 is %lu\n", sizeof(unsigned)-1);

	//7 (8-1)
	printf("sizeof -1 is %lu\n", sizeof(double)-1);

	//8 (double)
	printf("sizeof -1 is %lu\n", sizeof((double)-1));

	//6 (4+2)
	printf("sizeof -1 is %lu\n", sizeof n+2);
	//4
	printf("sizeof -1 is %lu\n", sizeof(n+2));
	//8 (double)
	printf("sizeof -1 is %lu\n", sizeof(n+2.0));

	return 0;
}