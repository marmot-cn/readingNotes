# chage

---

## 简介

是用来修改**帐号**和**密码**的**有效期限**.

## 语法

`chage [选项] 用户名`

## 选项

* `-m`: 密码可更改的最小天数. 为零时代表任何时候都可以更改密码.
* `-M`: 密码保持有效的最小天数. 超过该天数就需要更改密码.
* `-W`: 用户密码到期前, 提前收到警告信息的天数.
* `-E`: 账户到期的日期. 过了这天, 此账号将不可用.
* `-d`: 上一次更改的日期.
* `-I`: 停滞时期. 如果一个密码已过期这些天, 那么此账号将不可用.
* `-l`: 列出当前的设置. 由非特权用户来确定他们的密码或帐号何时过期.


```shell
[ansible@iZ94ebqp9jtZ ~]$ chage -l ansible
Last password change					: Aug 24, 2017
Password expires					: Nov 22, 2017
Password inactive					: never
Account expires						: never
Minimum number of days between password change		: 0
Maximum number of days between password change		: 90
Number of days of warning before password expires	: 7
```

* `Last password change`: 最近一次修改密码的时间.
* `Password expires`: 密码过期的时间.
* `Password inactive`: 密码失效的时间.
* `Account expires`: 账户过期时间.
* `Minimum number of days between password change`: 两次改变密码之间相距最小天数.
* `Maximum number of days between password change`: 两次密码改变密码相距最大天数.
* `Number of days of warning before password expires`: 密码过期前开始警告的天数.

`Last passwrd change + M天数 = Password expires`

`Maximum number of days between password chang=M天数`