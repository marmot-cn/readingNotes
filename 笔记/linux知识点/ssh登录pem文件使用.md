# ssh登录pem文件使用

---

## 命令

```
ssh -i key.pem root@IP
```

如果出现报错说明这个问题是文件的权限太大了, 需要给小点:

```
sudo chmod 600 key.pem
```

## 添加`key`文件

```
ssh-add -k key.pem
```

### 生成`pem`文件测试

#### 生成公钥和私钥

```
ssh root@120.25.161.1
root@120.25.161.1's password:
Last failed login: Mon Apr  9 02:52:42 CST 2018 from 179.132.148.187 on ssh:notty
There were 9 failed login attempts since the last successful login.
Last login: Tue Dec 26 16:05:43 2017 from 222.90.89.69

Welcome to Alibaba Cloud Elastic Compute Service !

[root@iZ94ebqp9jtZ ~]# useradd tom
[root@iZ94ebqp9jtZ ~]# su tom
[tom@iZ94ebqp9jtZ root]$ ssh-keygen -t rsa -C "41893204@qq.com"
Generating public/private rsa key pair.
Enter file in which to save the key (/home/tom/.ssh/id_rsa):
Created directory '/home/tom/.ssh'.
Enter passphrase (empty for no passphrase):
Enter same passphrase again:
Your identification has been saved in /home/tom/.ssh/id_rsa.
Your public key has been saved in /home/tom/.ssh/id_rsa.pub.
The key fingerprint is:
36:30:ad:8c:69:d6:87:dd:41:d1:7a:de:a0:b0:b7:4f 41893204@qq.com
The key's randomart image is:
+--[ RSA 2048]----+
|          oo     |
|       . .  .    |
|      o . ..     |
|     = *....o    |
|    = = So.+ o   |
|   o   o..o . .  |
|         . .E    |
|          ..     |
|           ..    |
+-----------------+
[tom@iZ94ebqp9jtZ root]$ cat /home/tom/.ssh/id_rsa.pub
ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDoLlYDPXYMNkgHKdxUPSBDwgQn6tS8PGeDSScZUlfaRkNDlrfmFHCEL9fD2vuwciJp1pfoMj7mybiEVuIb0pfQ23WTCtPHTgKEeyOcT7tbPVDZ9DqsYPg4RMfo2USpf8MA3oN2B1INnlHbYr0sD+RoLX1fLQuPLQrSL7XFeEZjU8A5C0YvBwBsHAWDp+duuK4LiLh9L/1gvjlskDVU9PAH2ERSwvtblIzqhM77ITrc/jRXJT3O0Fkwzcab5DrBCePfyxgE2YxkbI0B8z+zHybYIcT25rE2/Rf2XWNiUAxb6IRVB/6eDqm5uKxD2yzPxGheC2rXZ/aHrTFKZCQM8IZ3 41893204@qq.com
```

#### 公钥复制到`authorized_keys`文件内,`authorized_keys`为`0600`权限

```
[tom@iZ94ebqp9jtZ .ssh]$ cat authorized_keys
ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCgP3XI+g6um+beilzlWefT7AslrwcwV+MvVsq3mIcqlShepH5z3GT3i0bk7fOAUWYjcNjAm6QaYDTwrC3bN6BxgCC8ypaJesBx0lttQOVp1PYedSQjRsLXvOYulPdEziQIe2VW3CyiHzLAt2jugVE92vDKGwcs16vMn0S7a3ruJM9kJ1RVfY7HFWi4IIIlTchcAJd+SXrES8SgHcDrqEidLxUSwVvR1EWCq+1NFUVl2OogCaKULdTT1I4Blr5jLv6j0Qf+nhWfX8/wwWJnGe55rtETB3sYkqvsU8m3+UQdlBL95r6jMM6eiJNfWoJ/rX6eiWj2eNrH7WmDwZQnpKst 41893204@qq.com
[tom@iZ94ebqp9jtZ .ssh]$ pwd
/home/tom/.ssh
[tom@iZ94ebqp9jtZ .ssh]$ ll -s
total 12
4 -rw------- 1 tom tom  397 Apr  9 17:26 authorized_keys
4 -rw------- 1 tom tom 1679 Apr  9 17:20 id_rsa
4 -rw-r--r-- 1 tom tom  397 Apr  9 17:20 id_rsa.pub
```

#### 把私钥`id_rsa`复制到各台单独服务器就可以实现免密钥登录

在另外一台机器使用刚才服务器的私钥进行登录.

```
ssh -i ./id_rsa tom@120.25.161.1
Last login: Mon Apr  9 17:26:17 2018

Welcome to Alibaba Cloud Elastic Compute Service !

[tom@iZ94ebqp9jtZ ~]$ exit
```

#### 转换成`pem`文件

```
openssl rsa -in id_rsa -out tom.pem
writing RSA key
```

则用户可以使用`tom.pem`文件来登录.

## 原理

其实就和我正常使用的免密钥登录原理一致.

正常情况下我在我本机创建私钥公钥, 然后把公钥复制到远程服务器某个账户下的`authorized_keys`中. 然后因为本机默认私钥位置没有改变所以不用`-i`指定私钥位置.

这个原理是在服务器创建账户, 把自己的公钥写入`authorized_keys`中, 然后把私钥传给需要登录的机器. 则用户通过指定`-i`即可免密钥登录, 或者传到默认文件位置亦可(不使用`-i`).

