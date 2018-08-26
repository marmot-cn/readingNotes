# 41_04_Linux集群系列之二十二——keepalived详解

---

## 笔记

### keepalived

提供`HA`的一个底层工具.

最早设计为`ipvs`提供`HA`功能. `ipvs`是在内核中提供, `keepalived`添加了能够将`vip`在节点之间流转的功能. `vip`流转基于`vrrp`协议.

#### vrrp

`vrrp`协议, 将多个物理设备虚拟为一个物理设备.

一主多从, 从始终处于空闲状态.

多主模式:

* `eth0:0` 对应虚拟路由器1
* `eth1:1` 对应虚拟路由器2

### keepalived 和 corosync

* keepalived 轻量级

### ipvs HA

* ipvs 高可用
* health check
	* fall_back: server

### keepalived 模拟宕机

```
vrrp_script chk_schedown {
    # 检测有没有这个文件, 有就返回1, 认为宕机了. 否则返回0, 表示成功
    script "[ -e /etc/keepalived/down ] && exit 1 || exit 0"
}

# 在 vrrp_instance 内
# 定义什么时候执行脚本
track_script {
    chk_schedown
}
```

通知脚本

```
#!/bin/bash
#

contact='root@localhost'

[ $# -lt 2 ] && Usage && exit

Usage() {
  echo "Usage: `basename $0` {master|backup|falut} VIP"
}

Notify() {
  subject="`hostname`'s state changed to $1"
  mailbody="`date "+%F %T"`: `hostname`'s state change to $1, vip floating."
  
  #$subject用引号括起来, 是因为subject中间有空格.
  echo $mailbody | mail -s "$subject" $contact
}

VIP=$2

case $1 in
  master)
    Notify master
    ;;
  backup)
    Notify backup
    ;;
  fault)
    Notify fault
    ;;
  *)
    Usage
    exit 1
    ;;
esac   
```


触发:

```
# 在 vrrp_instance 内

# 我们上文编写的脚本名称
notify_master "/etc/keepalived/notify.sh master xxx.xxx.xxx.xxx(ip)"
notify_backup "/etc/keepalived/notify.sh backup xxx.xxx.xxx.xxx(ip)"
notify_fault "/etc/keepalived/notify.sh fault xxx.xxx.xxx.xxx(ip)"
```

说明

```
basename $0
```

`$0`是当前脚本名称.

`basename`则打印出基本名称.

## 整理知识点

---