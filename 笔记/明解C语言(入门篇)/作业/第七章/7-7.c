#include <stdio.h>

int main (void)
{	

	float f;
	printf("请输入folat型: ");
	scanf("%f",&f);
	printf("输入的是: %f\n", f);

	double d;
	printf("请输入double型: ");
	scanf("%lf",&d);
	printf("输入的是: %lf\n", d);

	long double ld;
	printf("请输入long double型: ");
	scanf("%Lf",&ld);
	printf("输入的是: %Lf\n", ld);

	return 0;
}