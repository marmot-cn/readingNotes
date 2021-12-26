#include <stdio.h>
#include <math.h>

unsigned set_n(unsigned x, int pos, int n)
{	
	//需要设为1的个数 ^(0<<n)
	//需要从第几位开始 (^(0<<n)) << pos

	unsigned j = 0;
	for (int i = 0; i < n; i++) {//从第0位到第(n-1)位均为1
		j += pow(2, i);
	}

	return x | (j << pos);
}

unsigned reset_n(unsigned x, int pos, int n)
{
	unsigned j = 0;
	for (int i = 0; i < n; i++) {//从第0位到第(n-1)位均为1
		j += pow(2, i);
	}

	return x & ~(j << pos);
}

unsigned inverse_n(unsigned x, int pos, int n)
{
	unsigned j = 0;
	for (int i = 0; i < n; i++) {//从第0位到第(n-1)位均为1
		j += pow(2, i);
	}

	return x ^ (j<< pos);
}

int main (void)
{	
	//10001 = 17
	//11111 = 31
	//10101 = 21

	printf("x=17, pos=1, n=3 应为 111111 = 31 实际是: %d\n", set_n(17, 1, 3));

	printf("x=31, pos=1, n=3 应为 10001 = 17 实际是: %d\n", reset_n(31, 1, 3));

	printf("x=21, pos=1, n=3 应为 11011 = 27 实际是: %d\n", inverse_n(21, 1, 3));

	return 0;
}