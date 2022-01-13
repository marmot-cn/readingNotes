#include <stdio.h>

int factorial(int n)
{
    if (n > 0)
        return n * factorial(n - 1);
    else
        return 1;
}

int combination(int n, int r) {

    if (n > 0)
        return  factorial(n) / (factorial(r) * factorial(n - r));
    else
        return 1;
}
int main() {
    int n;
    int r;
    printf("n=");
    scanf("%d", &n);
    printf("r=");
    scanf("%d", &r);
    printf("%d", combination(n, r));
}