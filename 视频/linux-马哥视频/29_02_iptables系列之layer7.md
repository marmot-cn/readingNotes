# 29_02_iptables系列之layer7

---

### 笔记

---

#### layer7 模块

需要给内核打完补丁, 需要重新编译内核.

复制`netfilter`新的扩展, 重新编译`iptables`

复制`iptables layer7`的协议特征包.

* `-m layer7`
	* `--l7proto [protocol name] -j [action]`

**示例: 禁止QQ上网**

* 服务器A, 2个网卡
	* 172.16.100.7 外网地址, 可以访问外网
	* 192.168.10.6 内网地址, 可以访问内网
* 服务器B要通过服务器A上外网.

需要添加`SNAT`规则.

```shell
服务器网关指向 192.168.10.6

服务器B可以上互联网
iptables -t nat -A POSTROUTING -s 192.168.10.0/24 -j SNAT --to-source 172.16.100.7

禁止上QQ
iptables -A FORWARD -s 192.168.10.0/24 --l7proto qq -j DROP
```

#### 根据时间控制

* `-m time`
	* `--datestart`
	* `--datestop`
	* `--timestart`
	* `--timestop`
		* `--monthdays`
		* `--weekdays`
	
**示例: 拒绝内网用户上午8点10到中午12点上互联网**	
	
```shell
iptables -A FORWARD -s 192.168.10.0/24 -m time --timestart 08:10:00 --timestop 12:00:00 -j DROP
```	

#### iptables 命令

* `iptables-save > xxxx` 保存 
* `iptables-restore < xxx` 恢复脚本

**iptables脚本**

```shell
#!/bin/bash
#
ipt=/usr/bin/iptables
einterface=eth1
iinterface=eth0

eip=172.16.100.7
iip=192.168.10.6

#清理所有链
$ipt -t nat -f
$ipt -t filter -f
$ipt -t mangle -f

#添加一些链
$ipt -N clean_up
$ipt -A clean_up -d 255.255.255.255 -p -icmp -j DROP
$ipt -A clean_up -j RETURN
```

开机需要自动执行.

`POST->MBR(bootloader)->Kernel(initrd)->init(/etc/inittab)`

`(/etc/inittab`:

1. 设定默认级别
2. 系统初始化脚本
3. 运行指定级别的服务 
	
```shell
写到 /etc/rc.d/rc.local

vim rc.local
...
要运行的shell脚本
...


```	

#### IDS 入侵检测系统

* nids 网络入侵检测系统
	* snort(开源实现版本) + iptables: 根据实现定义好的规则发现攻击, 自动生成攻击放到`iptables` = NIPS
* hids 主机入侵检测系统


	
### 整理知识点

---