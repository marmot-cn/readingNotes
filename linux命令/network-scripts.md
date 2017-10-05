# network-scripts

---

网络接口配置文件.

## 配置参数说明

* `TYPE`: 配置文件接口类型.
	* `Ethernet`
	* `IPsec`
	* `...`
* `DEVICE`: 网络接口名称
* `BOOTPROTO`: 系统启动地址协议
	* `none`: 不使用启动地址协议
	* `bootp`: `BOOTP`协议
	* `dhcp`: `DHCP`动态地址协议
	* `static`: 静态地址协议
* `ONBOOT`: 系统启动时是否激活
	* `yes`: 系统启动时激活该网络接口
	* `no`: 系统启动时不激活该网络接口
* `IPADDR`: IP地址
* `NETMASK`: 子网掩码
* `GATEWAY`: 网关地址
* `BROADCAST`: 广播地址
* `HWADDR/MACADDR`: MAC地址, 只需设置一个.
* `PEERDNS`: 是否指定DNS. 如果使用DHCP协议, 默认为`yes`.
	* `yes`: 如果DNS设置, 修改`/etc/resolv.conf`中的`DNS`.
	* `no`: 不修改`/etc/resolv.conf`中的`DNS`.
* `NDS{1,2}`: `DNS`地址. 当`PEERDNS`为yes时会被写入`/etc/resolv.conf`中.
* `NM_CONTROLLED`: 是否由Network Manager控制该网络接口. **一般设置为`no`**
	* `yes`：由Network Manager控制.
	* `no`：不由Network Manager控制
* `USERCTL`: 用户权限控制
	* `yes`：非root用户允许控制该网络接口
	* `no`：非root用户不允许控制该网络接口
* `IPV6INIT`：是否执行IPv6
	* `yes`：支持IPv6
	* `no`：不支持IPv6
* `IPV6ADDR`：IPv6地址/前缀长度
* `UUID`, 网卡的`uuid`. 可以通过`uuidgen 网卡名称`来生成.