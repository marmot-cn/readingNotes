#12_02_Linux软件管理之二 rpm

###笔记

---

**rpm命令**

`rpm`: 管理软件包.

`rpm`数据库在`/var/lib/rpm`用于追踪安装文件,方便以后卸载.查询文件.

`rpmbuild`: 创建软件包.

rpm: 安装,查询,卸载,升级,校验,数据库的重建,验证数据包等工作.

**服务器安装软件**

* 验证软件包来源合法性
* 验证软件包是否完整

**rpm命名**

包: 组成部分(软件`bind`),一个软件包可能有1个主包和多个子包组成.

* 主包:
	* bind-9.7.1-1.el5.i586.rpm
* 子包:
    * bind-libs-9.7.1-1.el5.i586.rpm

包名格式:

name-version(版本号:major.minor.release)-release(rpm包制作者的修正号).arch(平台).rpm

`version`:

* 主版本号:重大改进.
* 次版本号:某个子功能发生重大变化.
* 发行号(release): 修正了bug,调整了一点功能.rmp包发行者的修正号.

el5: 红帽5

`rpm包`:

* 二进制格式: 装了以后直接可以运行
* 源码格式: 装了以后编译

**rpm安装**

`rpm` `-i /PATH/TP/PACKAGE_FILE`

* `-h`: 以#显示进度,每个#表示2%
* `-v`: 显示详细过程
* `-vv`: 更详细的过程
* `--nodeps`: 忽略依赖关系
* `--replacepkgs`: 重新安装,替换原有的安装
* `--oldpackage`: 降级安装
* `--force`: 强行安装,可以实现重装或降级
* `--test`

**查询**

`rpm` `-q PACKAGE_NAME`: 查询指定的包是否已经安装

`rpm` `-qa`: 查询已经安装的所有包

`rpm -qi PACKAGE_NAME`: 查询指定包的说明信息

`rpm -ql PACKAGE_NAME`: 查询指定包安装后生成的文件列表

`rpm -qf /path/to/somefile`: 查询指定的文件是由哪个rpm包安装生成的

`rpm -qc PACKAGE_NAME`: 查询指定软件包安装的配置文件

`rpm -qd PACKAGE_NAME`: 查询指定包安装的帮助文件

`rpm -q --scripts PACKAGE_NAME`: 查询指定包中包含的脚本(四类脚本)

* 安装前(`preinstall`)
* 安装后(`postinstall`)
* 卸载前(`preuninstall`)
* 卸载后(`postuninstall`)

如果某rpm包尚未安装,我们需查询其说明信息,安装以后会生成的文件:

`rpm -qpi /PATH/TO/PACKAGE_FILE`

**升级**

`rpm -Uvh /PATH/TO/NEW_APCKAGE_FILE`: 如果装有老版本,则升级;否则,则安装

`RPM -Fvh /PATH/TO/NEW_APCKAGE_FILE`: 如果装有老板,则升级;否则,退出

`--oldpackage`: 降级

**卸载**

`rpm -e PACKAGE_NAME`

* `--nodeps` 

**校验**

`rpm -V PAKCAGE_NAME`

**检验来源合法性,及软件包完整性**

加密类型:

* 对称: 加密解密使用同一个密钥
* 公钥: 一对密钥,公钥,私钥.公钥隐含与私钥中,可以提取出来,并公开出去 
* 单向: 提取md5码,用私钥加密md5

`/etc/pki/rpm-gpg/`

`rpm --import /etc/pki/rpm-gpg/密钥文件  -K /PATH/TO//PACKAGE_FILE`

* `dsa,gpg`: 验证来源合法性,即验证签名,也可以使用 `--nosignature`,略过此项
* `sha1,md5`: 验证软件包完整性,也可以使用 `--nodigest`,略过此项

导入密钥文件,进行验证.

**重建数据库**

`rpm`

* `--rebuilddb`: 重建数据库,使用所有以安装软件包的首部信息.有也重建.
* `--initdb`: 初始化数据库,没有才建立,有不会建立.

###整理知识点

---
