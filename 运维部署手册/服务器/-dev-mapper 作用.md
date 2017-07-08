# /dev/mapper 作用

### Device mapper 机制

`Device mapper` 是 Linux2.6内核中提供的一种从逻辑设备到物理设备的映射机制.`LVM2`必须使用这个模块.

### /dev/mapper 目录

名称的格式

`/dev/{vg_name}/{lv_name} -> /dev/mapper/{vg_name}-{lv_name}`