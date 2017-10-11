# Linux 查看物理CPU个数,核数,逻辑CPU个数

---

### CPU逻辑核心数和物理核心数

---

物理CPU个数×每颗核数就应该等于逻辑CPU的个数

* 总核数 = 物理CPU个数 X 每颗物理CPU的核数 
* 总逻辑CPU数 = 物理CPU个数 X 每颗物理CPU的核数 X 超线程数

### 查看物理CPU个数

**`physical id`**

		cat /proc/cpuinfo| grep "physical id"| sort| uniq| wc -l
		
###  查看每个物理CPU中core的个数(即核数)

**`cpu cores`**

		cat /proc/cpuinfo| grep "cpu cores"| uniq
		
### 查看逻辑CPU的个数

**`processor`**

		cat /proc/cpuinfo| grep "processor"| wc -l