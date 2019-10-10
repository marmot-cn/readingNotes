# grep

---

## 常用用法

```
grep [-acinv] [--color=auto] '搜寻字符串' filename
```

### 选项

* `-a`: 将`binary`文件以`text`文件的方式搜寻数据.
* `-c`：计算找到'搜寻字符串'的次数.
* `-i`：忽略大小写的不同, 所以**大小写视为相同**.
* `-n`：顺便输出行号.
* `-v`：反向选择, 亦即显示出没有'搜寻字符串'内容的那一行.
* `--color=auto`：可以将找到的关键词部分加上颜色的显示.

### 示例

将/etc/passwd，有出现 root 的行取出来,同时显示这些行在/etc/passwd的行号.

```
[ansible@localhost ~]$ grep -n root /etc/passwd
1:root:x:0:0:root:/root:/bin/bash
10:operator:x:11:0:operator:/root:/sbin/nologin
25:dockerroot:x:996:992:Docker User:/var/lib/docker:/sbin/nologin
```

### grep与正规表达式

#### 字符类

```
grep '正则' filename
```

匹配行首为`root`.

```
[ansible@localhost ~]$ grep -n '^root' /etc/passwd
1:root:x:0:0:root:/root:/bin/bash
```