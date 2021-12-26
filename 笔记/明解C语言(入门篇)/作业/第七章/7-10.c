#include <stdio.h>

int main (void)
{	
	float x = -0.01;

	for (int i = 0; i <= 100; i++)
		printf("x = %f\tx = %f\n", x+=0.01, i/100.0);

	return 0;
}