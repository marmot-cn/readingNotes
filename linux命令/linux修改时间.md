# linux修改时间

---

## 设置时区

### 查看当前时区

```shell
[ansible@iZ944l0t308Z ~]$ timedatectl
      Local time: Wed 2017-08-16 20:39:44 CST
  Universal time: Wed 2017-08-16 12:39:44 UTC
        RTC time: Wed 2017-08-16 20:39:42
       Time zone: Asia/Shanghai (CST, +0800)
     NTP enabled: no
NTP synchronized: yes
 RTC in local TZ: yes
      DST active: n/a
      
[ansible@iZ944l0t308Z ~]$ cat /etc/timezone
Asia/Shanghai
```

### 查看所有时区文件

```shell
ls -F /usr/share/zoneinfo/
Africa/      Arctic/    Australia/  CST6CDT  Cuba  EST5EDT
```

### 如果想看对于每个`time zone`当前的时间我们可以用`zdump`命令

```shell
[ansible@iZ944l0t308Z ~]$ zdump Africa
Africa  Wed Aug 16 12:36:53 2017 Africa
```

### 修改时区

#### 使用`tzselect`命令

一步一步选择即可,

```shell
[ansible@iZ944l0t308Z ~]$ tzselect
Please identify a location so that time zone rules can be set correctly.
Please select a continent or ocean.
 1) Africa
 2) Americas
 3) Antarctica
 4) Arctic Ocean
 5) Asia
 6) Atlantic Ocean
 7) Australia
 8) Europe
 9) Indian Ocean
10) Pacific Ocean
11) none - I want to specify the time zone using the Posix TZ format.
...
```

### 修改`/etc/localtime`

复制相应的时区文件，替换系统时区文件；或者创建链接文件.

`cp /usr/share/zoneinfo/主时区/次时区 /etc/localtime`

```shell
在设置中国时区使用亚洲/上海

cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

或 
rm /etc/localtime
ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
```

## 硬件时间时钟,软件时间时钟

* `RTC`, 硬件时间时钟
* `System Clock`, 软件时间时钟

硬件时钟是指嵌在主板上的特殊的电路, 它的存在就是平时我们关机之后还可以计算时间的原因.

软件时钟就是操作系统的kernel所用来计算时间的时钟. 它从1970年1月1日00:00:00 UTC时间到目前为止秒数总和的值 在Linux下系统时间在开机的时候会和硬件时间同步(synchronization),之后也就各自独立运行了.

### 查看软件时间, 硬件时间

```shell
软件时间
[ansible@iZ944l0t308Z ~]$ date
Wed Aug 16 20:52:05 CST 2017

硬件时间
[ansible@iZ944l0t308Z ~]$ sudo hwclock --show
Wed Aug 16 20:52:38 2017  -0.271194 seconds
```

### 把硬件时间设置成软件时间

`hwclock --hctosys`

`-s, --hctosys        set the system time from the hardware clock`

### 把软件时间设置成硬件时间

`hwclock --systohc`

`-w, --systohc        set the hardware clock from the current system time`

### 修改硬件时间

`hwclock --set --date="mm/dd/yy hh:mm:ss"`

### 修改软件时间 

` date -s "dd/mm/yyyy hh:mm:ss"`

## ntp

指软件时间和网络服务器之间的同步.

* `ntpd`在实际同步时间是一点点的校准过来时间的, 最终把时间慢慢的**校正**对.
* `ntpdate`不会考虑其他程序是否会阵痛, 直接**调整**时间.

`ntpd`在和时间服务器的同步过程中,会把`BIOS`计时器的振荡频率偏差——或者说`Local Clock`的自然漂移(drift)——记录下来. 这样即使网络有问题,本机仍然能维持一个相当精确的走时. 如果`client`与`server`时差异常大或过小,`ntpd`将会拒绝`server`参考时间.

### ntpdate 使用

通常采用`ntpdate`同步时间时是设置一个`crontab`任务, 一个周期内重复执行同步命令.

### ntpd 使用

`ntpd`是**守护进程**,运行状态可以随时监控.

### 配置文件`/etc/ntp.conf`

#### server

`server`用于设定ntp同步时间的外网时间服务器.

`server + ip`或者`server + hostname`

参数:

* `prefer`: 优先级.
* `burst`: 当一个运程NTP服务器**可用**时, 向它发送一系列的并发包进行检测.
* `iburst`: 当一个运程NTP服务器不可用时, 向它发送一系列的并发包进行检测.
* `minpoll/maxpoll`: 规定查询的间隔,以2的幕的形式,取值范围在4-17.minpoll 3表示2的3次方,也就是最短8秒钟后主动与上层NTP服务器同步一次,maxpoll 4表示2的4次方,也就是最长16秒钟后主动与上层NTP服务器同步一次

如果连接NTP上层服务器失败,`ntp`服务会跳过失败的NTP服务器,而以配置项的顺序依次与`ntp`服务器同步.

#### restrict

`restrict`用于对访问`ntp`的客户端的限制.

参数:

* `kod`: 使用`kod`技术防范"kiss of death"攻击.
* `ignore`: 拒绝任何`ntp`链接.
* `nomodify`: 客户端不能使用ntpc,ntpq修改时间服务器参数,可以进行网络校时.
* `noquery`: 客户端不能使用ntpq,ntpc等指令来查询服务器时间,等于不提供ntp的网络校时.
* `notrap`: 不提供远程日志功能.
* `notrust`: 拒绝没有认证的客户端,
* `nopeer`：不与其他同一层的ntp服务器进行时间同步.

`restrict ip`或者`restrict IP地址 + mask + 子网掩码 + 参数`

```shell
restrict default nomodify notrap nopeer noquery   #默认拒绝所有访问 只可以同步时间, ipv4
restrict -6 default kod nomodify notrap nopeer noquery #ipv6
restrict 211.71.14.254 mask 255.255.255.0 #添加允许211.71.14.254/24网段访问
restrict 10.111.1.1 mask 255.0.0.0 nomodify #添加10.0.0.0/8网段访问，不可以修改服务器时间参数
```

### logfile

`logfile  /var/log/ntpd.log` 日志文件保存地址.

### key认证文件

`driftfile /var/lib/ntp/drift` 本地与上层服务器BIOS晶片振荡频率差值保存目录,不需要修改.

`keys /etc/ntp/keys` 可以借此来给客户端设置认证信息.

### 国内`NTP`服务器列表

国内(www.pool.ntp.org/zone/cn):

* server 0.cn.pool.ntp.org
* server 1.cn.pool.ntp.org
* server 2.cn.pool.ntp.org
* server 3.cn.pool.ntp.org

阿里云:

* ntp1.aliyun.com
* ntp2.aliyun.com
* ntp3.aliyun.com
* ntp4.aliyun.com
* ntp5.aliyun.com
* ntp6.aliyun.com
* ntp7.aliyun.com

### 本机修改服务文件

```shell

server ntp1.aliyun.com iburst
server ntp2.aliyun.com iburst
server ntp3.aliyun.com iburst
server ntp4.aliyun.com iburst
server ntp5.aliyun.com iburst
server ntp6.aliyun.com iburst
server ntp7.aliyun.com iburst
server 0.cn.pool.ntp.org iburst
server 1.cn.pool.ntp.org iburst
server 2.cn.pool.ntp.org iburst
server 3.cn.pool.ntp.org iburst

```



