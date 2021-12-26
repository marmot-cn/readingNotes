#include <stdio.h>


//第pos为设为1
unsigned set(unsigned x, int pos) {

	return x | (1 << pos);
}

//第pos位设为0
unsigned reset(unsigned x, int pos) {

	//1 左移后取反，如 1 左移2位后是  100, 取反为 011
	//任何数 & 0 为0, 任何数 & 1 为其本身
	return x & ~(1 << pos);
}	

//第pos为取反
unsigned inverse(unsigned x, int pos) {
	return x ^ (1 << pos);
}

int main (void)
{	
	//101
	int n = 5;

	printf("5 第1位设为1应为7, = %d\n", set(n, 1));
	printf("5 第2位设为0应为1, =  %d\n", reset(n, 2));
	printf("5 第1位取反应为7, =  %d\n", inverse(n, 1));

	return 0;
}