# 物理CPU, 核数以及逻辑CPU个数

---

### 概述

```
CPU总核数 = 物理CPU个数 * 每颗物理CPU的核数
总逻辑CPU数 = 物理CPU个数 * 每颗物理CPU的核数 * 超线程数
```

### 查看CPU信息(型号)

```
[root@iZ94ebqp9jtZ ~]# cat /proc/cpuinfo | grep name | cut -f2 -d:
 Intel(R) Xeon(R) CPU E5-2650 v2 @ 2.60GHz
```

### 查看物理CPU个数

```
[root@iZ94ebqp9jtZ ~]# cat /proc/cpuinfo| grep "physical id"| sort| uniq| wc -l
1
```

### 查看每个物理CPU中core的个数(即核数)

```
[root@iZ94ebqp9jtZ ~]# cat /proc/cpuinfo| grep "cpu cores"| uniq
cpu cores	: 1
```

### 查看逻辑CPU的个数

```
[root@iZ94ebqp9jtZ ~]# cat /proc/cpuinfo| grep "processor"| wc -l
1
```

