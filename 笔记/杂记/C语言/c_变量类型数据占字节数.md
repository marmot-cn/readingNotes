#C 变量类型数据占字节数

---

* `bit`(比特)描述电脑数据的最小单位. `二进制`系统中, 每个`0`或`1`就是一个位(`bit`).
* `Byte` 就是字节, `1`个由`8`个二进制位组成. `1Byte = 8bit`(1个字节有8位(8个比特))

int，long int，short int的宽度都可能`随编译器而异`,有几条铁定的原则(`ANSI/ISO`制订的):

1. sizeof(short int)<=sizeof(int) 
2. sizeof(int)<=sizeof(long int) 
3. short int至少应为16位(2字节) 
4. long int至少应为32位

**16位编辑器**

* `char`: 1个字节
* `char*(即指针变量)`: 2个字节
* `short int`: 2个字节
* `int`: 2个字节
* `unsigned int`: 2个字节
* `float`: 4个字节
* `double`: 8个字节
* `long`: 4个字节
* `long long`: 8个字节
* `unsigned long`: 4个字节

**32位编辑器**

* `char`: 1个字节
* **`char*(即指针变量)`**: 4个字节(32位的寻址空间是2^32, 即32个bit, 也就是4个字节.同理64位编译器)(16位机,32位机,64位机各不相同)
* `short int`: 2个字节
* **`int`**: 4个字节（16位机是2B，32位&64位是4B）
* **`unsigned int`**: 4个字节(16位机是2B,32位&64位是4B)
* `float`: 4个字节
* `double`: 8个字节
* **`long`**: 4个字节(16位&32位是4B，64位是8B)
* `long long`: 8个字节
* **`unsigned long`**: 4个字节(16位&32位是4B，64位是8B)

**64位**

* `char`: 1个字节
* `char*(即指针变量)`: 8个字节
* `short int`: 2个字节
* `int`: 4个字节
* `unsigned int`: 4个字节
* `double`: 8个字节
* `long`: 4个字节
* `long long`: 8个字节
* `unsigned long`: 8个字节