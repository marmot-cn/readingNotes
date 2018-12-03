# 修改安装好后的`/home`空间大小

## 修改步骤

因为默认安装好以后`/home`空间过大.

修改步骤如下.

### `umount /home`

### 删除`/home`所在的lv逻辑卷

`lvremove /dev/mapper/centos-home`

### 扩展`/root`所在的lv 

增多多少根据自己的服务器空间设置

`lvextend -L +??G  /dev/mapper/centos-root`

### 扩展`/root`文件系统

* `xfs_info /dev/mapper/centos-root`
* `xfs_growfs /dev/mapper/centos-root`

### 重新创建`home lv`

`lvcreate -L ??G -n home centos`

### 重新创建文件系统

`mkfs.xfs  /dev/mapper/centos-home`

### 挂载

`mount  /dev/mapper/centos-home  /home`

### 验证