# 29_01_iptables系列之layer7

---

### 笔记

---

#### 地址转发

`IP`属于主机不属于网卡, 主机只要开启转发功能, 一个网卡也可以转发.

#### DMZ

非军事化区域, 外网服务器.

内网主机是军事化区域.

服务器有3个网卡, 一个对外, 一个对自己的对外服务器, 一个对自己的对内服务器.

#### DNS 视图 

内网找DNS服务器解析域名需要解析成内网地址.

外网找DNS服务器解析域名需要解析成外网地址.

可以用DSN的**`view`视图**来处理.

#### 内核编译

* 单内核: 模块化(文件系统, 驱动, 安全)
	* 配置, 保存为`.config`文件
		* make menuconfig
		* make gconfig
		* make kconfig
		* make config
		* make oldconfig
		* ... 
	* 编译
		* make SUBDIR=arch/ 编译根CPU内核核心相关的代码
		* make dir/ 只编译该目录下的内核源码
			
				make drivers/net/pcnet32.ko 只编译网卡
		* make O=/path/to/somewhere 编译转存结果
	* 安装内核模块
		* make modules_install
	* 安装内核
		* make install

**示例**

重新编译

```shell
清楚编译好的结果
make clean

删除所有的编译生成文件、内核配置文件(.config文件)和各种备份文件. make mrproper删除的范围比make clean大，实际上, make mrproper在具体执行时第一步就是调用make clean,
make mrproper
```

#### iptables layer7 

`iptables` 是一个`OSI`的第二(`mac`扩展),第三和第四层的工具.

检查的层越高, 消耗的系统资源越多, 效率越低.

`netfilter`: 基于扩展可以对 `http`, `smtp` 根据协议实现屏蔽. 需要对`netfilter`提供额外的过滤框架(只能针对内核打补丁).

`iptables` 语法, 是根据`netfilter`实现编译好的. 如果对`netfilter`提供额外的框架, 还需要给`iptables`打补丁让`iptables`支持相关语法.

只能用于`linux`内核`2.4~2.6`.


### 整理知识点

---