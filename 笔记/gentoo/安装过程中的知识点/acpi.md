# acpi

`ACPI`表示高级配置和电源管理接口(Advanced Configuration and Power Management Interface).

早期Advanced Power Management模型(APM)将电源管理几乎完全分配给`BIOS`控制, 这大大的限制了操作系统在控制电能消耗方面的功能.

ACPI可以实现的功能包括:

* 系统电源管理(System power management)
* 设备电源管理(Device power management)
* 处理器电源管理(Processor power management)
* 设备和处理器性能管理(Device and processor performance management)
* 配置/即插即用(Configuration/Plug and Play)
* 系统事件(System Event)
* 电池管理(Battery management)
* 温度管理(Thermal management)
* 嵌入式控制器(Embedded Controller)
* SMBus控制器(SMBus Controller)

### 节电方式

1. (`suspend`即挂起), 显示屏自动断电. 只是主机通电. 这时敲任意键即可恢复原来状态. **屏幕断电**
2. (`save to ram`或`suspend to ram`即挂起到内存), 系统把当前信息储存在内存中, 只有内存等几个关键部件通电, 这时计算机处在高度节电状态, 按任意键后, 计算机从内存中读取信息很快恢复到原来状态. **存储到内存**
3. (`save to disk`或`suspend to disk`即挂起到硬盘)计算机自动关机, 关机前将当前数据存储在硬盘上, 用户下次按开关键开机时计算机将无须启动系统, 直接从硬盘读取数据, 恢复原来状态. **存储到硬盘**

### 实现功能

1. 用户可以使外设在指定时间开关.
2. 使用笔记本电脑的用户可以指定计算机在低电压的情况下进入低功耗状态, 以保证重要的应用程序运行.
3. 操作系统可以在应用程序对时间要求不高的情况下降低时钟频率.
4. 操作系统可以根据外设和主板的具体需求为它分配能源.
5. 在无人使用计算机时可以使计算机进入休眠状态, 但保证一些通信设备打开.
6. 即插即用设备在插入时能够由ACPI来控制.

### 六种状态

* `S0`: 平常的工作状态, 所有设备全开, 功耗一般会超过80W.
* `S1`: 也称为`POS`(Power On Suspend), 这时除了通过CPU时钟控制器将CPU关闭之外, 其他的部件仍然正常工作, 这时的功耗一般在30W以下;(其实有些CPU降温软件就是利用这种工作原理).
* `S2`: 这时CPU处于停止运作状态, 总线时钟也被关闭, 但其余的设备仍然运转.
* `S3`: 这就是我们熟悉的STR(Suspend to RAM), 这时的功耗不超过10W.
* `S4`: 也称为STD(Suspend to Disk), 这时系统主电源关闭, 硬盘存储S4前数据信息, 所以S4是比S3更省电状态.
* `S5`: 这种状态是最干脆的,就是连电源在内的所有设备全部关闭,即关机(shutdown), 功耗为0.