# ssh安全设置修改
---

## 修改端口号

打开`/etc/ssh/sshd_config`

```shell
sudo vi /etc/ssh/sshd_config
```
		
找见:

```shell
...
#Port 22
#ListenAddress 0.0.0.0
#ListenAddress ::
...
```
		
修改如下:

```shell
Port 22
Port 17456
#ListenAddress 0.0.0.0
#ListenAddress ::
```
		
**这里放开2个端口主要是用于测试,放置新的端口17456不行而不能ssh到服务器上**.
		
找其他机器尝试用 `17456`端口`ssh`,如果可以则屏蔽掉`22`端口好.最终配置文件如下:

		#Port 22
		Port 17456
		#ListenAddress 0.0.0.0
		#ListenAddress ::
		
## 禁止`root`用户登录		
		
修改 `PermitRootLogin`将`yes`改为`no`: 用户禁止root远程登录
		
## 只让可允许的用户登录

### 参数

* `AllowGroups`: 允许`ssh`登录的用户组. 用空格分开.
* `AllowUsers`: 允许`ssh`登录的用户. 用空格分开.
* `DenyGroups`: 拒绝`ssh`登录的用户组. 用空格分开.
* `DenyUsers`: 拒绝`ssh`登录的用户. 用空格分开.

### 修改示例

```shell
添加可以ssh登录的用户组 sshallow
groupadd sshallow

把给ansible用户添加到sshallow用户组, 这里必须要加上 -a, 否则会使用户离开自己的现有的组, 而只属于sshallow组
usermod -a -G sshallow ansible

id ansible
uid=1000(ansible) gid=1000(ansible) groups=1000(ansible),10(wheel),993(docker),1003(sshallow)

编辑 /etc/ssh/sshd_config, 在最后添加 
AllowGroups sshallow
```

这样我们只有在`sshallow`组的用户才能进行`ssh`.

## 重启`sshd`		
		
让`sshd`重新加载配置文件

		sudo /etc/init.d/sshd reload
		
或重启`sshd`服务

		sudo service sshd restart
		
