#12_03_Linux软件管理之三 yum

###笔记

---

**rpm包缺陷**

依赖关系

**yum**

C/S 架构

* Clinet,Server

`客户端`

* `配置文件`
	* yum仓库信息

`缓存目录`: 缓存当前指向yum仓库的信息(yum仓库的`元数据文件`)

**yum repository(yum仓库)**

repo文件是yum源(软件仓库)

* 文件服务
	* ftp(ftp://)
	* web(http://)
	* 本地文件路径(file:///,最后一个'/'代表根路径 )
	
**yum仓库中的元数据文件**

`primary.xm.gz`: 

* 包含了当前仓库内的所有软件包的列表  
* 包含了依赖关系
* 每一个RPM安装生成的文件列表

`filelist.xml.gz`:

* 当前仓库中所有RPM包的所有文件列表

`other.xml.gz`:

* 额外信息, RPM包的修改日志

`repomd.xml`:

* 记录的是上面三个文件的时间戳和校验和

`comps*.xml`:

* RPM包分组信息

**yum配置文件**

`/etc/yum.conf`:

		[chloroplast@dev-server ~]$ cat /etc/yum.conf
		[main]
		cachedir=/var/cache/yum/$basearch/$releasever
		keepcache=0
		debuglevel=2
		logfile=/var/log/yum.log
		exactarch=1
		obsoletes=1
		gpgcheck=1
		plugins=1
		installonly_limit=5
		bugtracker_url=http://bugs.centos.org/set_project.php?project_id=16&ref=http://bugs.centos.org/bug_report_page.php?category=yum
		distroverpkg=centos-release
		
		#  This is the default, if you make this bigger yum won't see if the metadata
		# is newer on the remote and so you'll "gain" the bandwidth of not having to
		# download the new metadata and "pay" for it by yum not having correct
		# information.
		#  It is esp. important, to have correct metadata, for distributions like
		# Fedora which don't keep old packages around. If you don't like this checking
		# interupting your command line usage, it's much better to have something
		# manually check the metadata once an hour (yum-updatesd will do this).
		# metadata_expire=90m
		
		# PUT YOUR REPOS HERE OR IN separate files named file.repo
		# in /etc/yum.repos.d
		timeout=3

* `[main]`: 片段的声明(名字叫为main的片段),低下的命令只对这个片段有效.
* `cachedir`: 缓存地址
* `keepcache`: 是否保存缓存文件
* `debuglevel`: 调试级别
* `logfile`: yum安装软件包的日志文件路径
* `tolerant`: 安装软件包装过了就不装,不报错
* `exactarch`: 安装rpm包在yum仓库中获得的版本必须与当前版本一致
* `obsoletes`: 
* `gpgcheck`: 使用gpg机制检查来源合法性和完整性
* `plugins`: 

**如果为yum定义repo文件**

		[Repo_ID]
		name=描述
		baseurl=repo仓库所在的具体访问路径(1.ftp://,2.http://,3.file:///)
		enabled={1|0} 1表示启用,0表示禁用
		gpgcheck={1|0} 使用gpg机制检查来源合法性和完整性
		gpgkey=指定gpg的文件路径
		
示例 `/etc/yum.repos.d/CentOS-Debuginfo.repo`:

		[root@dev-server chloroplast]# cat /etc/yum.repos.d/CentOS-Debuginfo.repo
		# CentOS-Debug.repo
		#
		# The mirror system uses the connecting IP address of the client and the
		# update status of each mirror to pick mirrors that are updated to and
		# geographically close to the client.  You should use this for CentOS updates
		# unless you are manually picking other mirrors.
		#
		
		# All debug packages from all the various CentOS-5 releases
		# are merged into a single repo, split by BaseArch
		#
		# Note: packages in the debuginfo repo are currently not signed
		#
		
		[debug]
		name=CentOS-6 - Debuginfo
		baseurl=http://debuginfo.centos.org/6/$basearch/
		gpgcheck=1
		gpgkey=file:///etc/pki/rpm-gpg/RPM-GPG-KEY-CentOS-Debug-6
		enabled=0
		

**yum使用格式**

`yum`

* `list`: 列表

		包名字	版本号	对应仓库(如果已经安装显示installed)

* `clean`: 清理缓存
* `repolist`: 可使用repo列表及其简要信息 
* `install`: 安装

		yum install package_name

* `update`: 升级	
* `update-to`: 指定升级到某个版本 	
* `remove|erase`: 卸载(如果其他包依赖这个包,则会一起卸载掉)
* `info`: 
* `provides|wahtprovides`: 查看指定的文件或特性是由哪个包安装生成的
* `groupinfo`: 显示组的信息
* `grouplist`: 显示组列表
* `groupremove`: 组卸载
* `groupinstall`: 组安装
* `search`: 查找软件包
	
`yum一些常用参数`:
	
* `-y`: 自动回答为yes
* `--nogpgcheck`

**如何创建yum仓库**

* `createrepo` (没有的话需要安装, `yum install createrepo`)


###整理知识点

---