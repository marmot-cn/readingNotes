# README

## 作业

`paging-multilevel-translate.py`

## 问题

### 1

对于线性页表，你需要一个寄存器来定位页表，假设硬件在TLB未命中时进行查找。你需要多少个寄存器才能找到两级页表？三级页表

### 2

使用模拟器对随机种子0、1和2执行翻译，并使用`-c`标志检查你的答案。需要多少内存应用来执行每次查找

### 3

根据你对缓存内从的工作原理的理解，你认为对页表的内存引用如何在缓存中工作？它们是否会导致大量的缓存命中（并导致快速访问）或者很多未命中（并导致访问缓慢）？