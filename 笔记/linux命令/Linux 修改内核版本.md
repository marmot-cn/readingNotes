# Linux 修改内核版本

---

因为这几天测试, 像更新`lvm2`包, 因为旧的包中的`lvchange`命令不支持`--metadataprofile <ProfileName>`选项.

想着直接升级了省事, 就运行了`yun update -y`, 结果连同内核也一起升级了. 随后重启了系统一次, 更新为新的内核了.

结果`ranchet`的网络节点提示连通不了.

后来降级内核解决了问题.

## 查看系统中有全部的内核

```shell
rpm -q kernel 
```

## 验证默认启动项

```shell
[root@localhost ansible]# grub2-editenv list
saved_entry=CentOS Linux (3.10.0-514.el7.x86_64) 7 (Core)
```

## 修改默认内核启动

```shell
cat /etc/grub2.cfg
```

比如我在我的文件找到`CentOS Linux (3.10.0-514.el7.x86_64) 7 (Core)`这项内容.

```shell
grub2-set-default "CentOS Linux (3.10.0-514.el7.x86_64) 7 (Core)"
```

再次使用`grub2-editenv list`验证启动项.

## 删除多余的内核

```shell
使用命令查看出多余的内核
rpm -q kernel 

使用yum remove 删除无用的内核
yum remove xxx 
```