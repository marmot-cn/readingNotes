# 服务器添加ftp功能

---

因为需求不能`ssh`登录, 需要`ftp`登录到服务器上.

我们需要给服务器添加`ftp`服务, 并添加对应账户.

**下载配置`ftp`服务器**

```shell
下载 vsftpd 软件
sudo yum install vsftpd
...
Y
...

sudo vi /etc/vsftpd/vsftpd.conf

设置指定的用户执行chroot
chroot_local_user=YES (默认所有用户不能chroot)
#chroot_list_enable=YES
#chroot_list_file=/etc/vsftpd/chroot_list (只有该列表的用户可以chroot)

添加读取用户配置目录
原来没有,需要自己添加该配置项
user_config_dir=/etc/vsftpd/userconf

allow_writeable_chroot=YES


如果上面没有配置启用该项, 则可以不建立该文件
sudo vi  /etc/vsftpd/chroot_list
...
xxxx(用户名, 在该文件内的用户可以 chroot)
...

mkdir -p /etc/vsftpd/userconf

sudo vi /etc/vsftpd/userconf/discuz
...
local_root=/data/html/DiscuzX/upload/
...

service vsftpd restart
```

**添加`ftp`使用账户**

```shell
添加用户 xxx 不能让其登录

或可以用此命令设置用户不能登录
usermod -s /usr/sbin/nologin xxx
```

### vsftpd 配置参数

#### 常见报错

`500 OOPS: vsftpd: refusing to run with writable root inside chroot ()`

从2.3.5之后, vsftpd增强了安全检查,如果用户被限定在了其主目录下,则该用户的主目录不能再具有写权限了！如果检查发现还有写权限,就会报该错误.

解决方案:

* 去掉用户主目录的写权限
* 在配置文件增加`allow_writeable_chroot=YES`

#### 1. 匿名用户设置

* `anonymous_enable=YES/NO(YES)`: 控制是否允许匿名用户登入, YES为允许匿名登入, NO为不允许. 默认值为YES.
* `no_anon_password=YES/NO(NO)`: 若是启动这项功能,则使用匿名登入时,不会询问密码.默认值为NO.
* `ftp_username=ftp`: 定义匿名登入的使用者名称. 默认值为`ftp`.
* `anon_root=/var/ftp`: 使用匿名登入时,所登入的目录.默认值为`/var/ftp`.注意ftp目录不能是777的权限属性,即匿名用户的家目录不能有777的权限.
* `anon_upload_enable=YES/NO(NO)`: 如果设为YES,则允许匿名登入者有上传文件(非目录)的权限,只有在`write_enable=YES`时,此项才有效. 当然, 匿名用户必须要有对上层目录的写入权.默认值为NO.
* `anon_world_readable_only=YES/NO(YES)`: 如果设为YES,则允许匿名登入者下载可阅读的档案(可以下载到本机阅读,不能直接在FTP服务器中打开阅读). 默认值为YES.
* `anon_mkdir_write_enable=YES/NO(NO)`: 如果设为YES, 则允许匿名登入者有新增目录的权限, 只有在write_enable=YES时, 此项才有效. 当然, 匿名用户必须要有对上层目录的写入权. 默认值为NO.
* `anon_other_write_enable=YES/NO(NO)`: 如果设为YES,则允许匿名登入者更多于上传或者建立目录之外的权限,譬如删除或者重命名.(如果anon_upload_enable=NO,则匿名用户不能上传文件,但可以删除或者重命名已经存在的文件；如果anon_mkdir_write_enable=NO, 则匿名用户不能上传或者新建文件夹,但可以删除或者重命名已经存在的文件夹.)默认值为NO.
* `chown_uploads=YES/NO(NO)`: 设置是否改变匿名用户上传文件(非目录)的属主.默认值为NO.
* `chown_username=username`: 设置匿名用户上传文件(非目录)的属主名.建议不要设置为root.
* `anon_umask=077`: 设置匿名登入者新增或上传档案时的umask 值.默认值为077,则新建档案的对应权限为700.
* `deny_email_enable=YES/NO(NO)`: 若是启动这项功能,则必须提供一个档案/etc/vsftpd/banner_emails,内容为email address.若是使用匿名登入,则会要求输入email address,若输入的email address 在此档案内,则不允许进入.默认值为NO.
* `banned_email_file=/etc/vsftpd/banner_emails`: 此文件用来输入email address,只有在deny_email_enable=YES时,才会使用到此档案.若是使用匿名登入,则会要求输入email address,若输入的email address 在此档案内,则不允许进入.

#### 2. 本地用户设置

* `local_enable=YES/NO(YES)`: 控制是否允许本地用户登入,YES 为允许本地用户登入,NO为不允许.默认值为YES.
* `local_root=/home/username`: 当本地用户登入时,将被更换到定义的目录下.默认值为各用户的家目录.
* `write_enable=YES/NO`: 是否允许登陆用户有写权限. 属于全局设置, 默认值为YES.
* `local_umask=022`: 本地用户新增档案时的umask 值.默认值为077.
* `file_open_mode=0755`: 本地用户上传档案后的档案权限,与chmod 所使用的数值相同.默认值为0666.

#### 3. 欢迎语设置

* `dirmessage_enable=YES/NO(YES)`: 如果启动这个选项,那么使用者第一次进入一个目录时,会检查该目录下是否有.message这个档案,如果有,则会出现此档案的内容,通常这个档案会放置欢迎话语,或是对该目录的说明.默认值为开启.
* `message_file=.message`: 设置目录消息文件,可将要显示的信息写入该文件.默认值为.message.
* `banner_file=/etc/vsftpd/banner`: 当使用者登入时,会显示此设定所在的档案内容,通常为欢迎话语或是说明.默认值为无.如果欢迎信息较多,则使用该配置项.
* `ftpd_banner=Welcome to BOB's FTP server`: 这里用来定义欢迎话语的字符串,banner_file是档案的形式,而ftpd_banner 则是字符串的形式.预设为无.

#### 4. 数据传输模式设置

* `ascii_upload_enable=YES/NO(NO)`: 设置是否启用ASCII 模式上传数据.默认值为NO.
* `ascii_download_enable=YES/NO(NO)`: 设置是否启用ASCII 模式下载数据.默认值为NO.

#### 5. 访问控制设置

##### 控制主机访问

`tcp_wrappers=YES/NO(YES)`: 设置vsftpd是否与tcp wrapper相结合来进行主机的访问控制.默认值为YES.如果启用,则vsftpd服务器会检查/etc/hosts.allow 和/etc/hosts.deny 中的设置,来决定请求连接的主机,是否允许访问该FTP服务器.这两个文件可以起到简易的防火墙功能.
比如：若要仅允许192.168.0.1—192.168.0.254的用户可以连接FTP服务器,则在/etc/hosts.allow文件中添加以下内容：

```
vsftpd:192.168.0. :allow
all:all :deny
```

##### 控制用户访问

对于用户的访问控制可以通过/etc目录下的vsftpd.user_list和ftpusers文件来实现.

* `userlist_file=/etc/vsftpd.user_list`: 控制用户访问FTP的文件,里面写着用户名称.一个用户名称一行.
* `userlist_enable=YES/NO(NO)`: 是否启用vsftpd.user_list文件.
* `userlist_deny=YES/NO(YES)`: 决定vsftpd.user_list文件中的用户是否能够访问FTP服务器.若设置为YES,则vsftpd.user_list文件中的用户不允许访问FTP,若设置为NO,则只有vsftpd.user_list文件中的用户才能访问FTP.
* `/etc/vsftpd/ftpusers`: 文件专门用于定义不允许访问FTP服务器的用户列表(注意:如果userlist_enable=YES,userlist_deny=NO,此时如果在vsftpd.user_list和ftpusers中都有某个用户时,那么这个用户是不能够访问FTP的,即ftpusers的优先级要高).默认情况下vsftpd.user_list和ftpusers,这两个文件已经预设置了一些不允许访问FTP服务器的系统内部账户.如果系统没有这两个文件,那么新建这两个文件,将用户添加进去即可.

#### 6. 访问速率设置

* `anon_max_rate=0`: 设置匿名登入者使用的最大传输速度,单位为B/s,0 表示不限制速度.默认值为0.
* `local_max_rate=0`: 本地用户使用的最大传输速度,单位为B/s,0 表示不限制速度.预设值为0.

#### 7. 超时时间设置

* `accept_timeout=60`: 设置建立FTP连接的超时时间,单位为秒.默认值为60.
* `connect_timeout=60`: PORT 方式下建立数据连接的超时时间,单位为秒.默认值为60.
* `data_connection_timeout=120`: 设置建立FTP数据连接的超时时间,单位为秒.默认值为120.
* `idle_session_timeout=300`: 设置多长时间不对FTP服务器进行任何操作,则断开该FTP连接,单位为秒.默认值为300 .

#### 8. 日志文件设置

* `xferlog_enable= YES/NO(YES)`: 是否启用上传/下载日志记录.如果启用,则上传与下载的信息将被完整纪录在xferlog_file 所定义的档案中.预设为开启.
* `xferlog_file=/var/log/vsftpd.log`: 设置日志文件名和路径,默认值为/var/log/vsftpd.log.
* `xferlog_std_format=YES/NO(NO)`: 如果启用,则日志文件将会写成xferlog的标准格式, 默认为关闭.
* `log_ftp_protocol=YES|NO(NO)`: 如果启用此选项,所有的FTP请求和响应都会被记录到日志中,默认日志文件在/var/log/vsftpd.log.启用此选项时,xferlog_std_format不能被激活.这个选项有助于调试.默认值为NO.

#### 9. 定义用户配置文件 

可以通过定义用户配置文件来实现不同的用户使用不同的配置.

`user_config_dir=/etc/vsftpd/userconf`.

设置用户配置文件所在的目录.当设置了该配置项后,用户登陆服务器后,系统就会到/etc/vsftpd/userconf目录下,读取与当前用户名相同的文件,并根据文件中的配置命令,对当前用户进行更进一步的配置.
例如：定义user_config_dir=/etc/vsftpd/userconf,且主机上有使用者 test1,test2,那么我们就在user_config_dir 的目录新增文件名为test1和test2两个文件.若是test1 登入,则会读取user_config_dir 下的test1 这个档案内的设定.默认值为无.利用用户配置文件,可以实现对不同用户进行访问速度的控制,在各用户配置文件中定义local_max_rate=XX,即可.

#### 10. 控制用户是否允许切换到上级目录

* `chroot_local_user`: 是否将所有用户限制在主目录,`YES`为启用,`NO`禁用. 默认为`NO`, 即用户可以向上切换到目录之外.
* `chroot_list_enable`: 是否启动限制用户的名单,`YES`为启用,`NO`禁用(包括注释掉也为禁用).
* `chroot_list_file=/etc/vsftpd/chroot_list`: 是否限制在主目录下的用户名单, 至于是限制名单还是排除名单, 这取决于`chroot_local_user`的值.

##### `chroot_local_user=YES`

* `chroot_list_enable=YES`
	* 所有用户都被限制在其主目录下.
	* 使用`chroot_list_file`指定的用户列表,这些用户作为“例外”,不受限制.
* `chroot_list_enable=NO`
	* 所有用户都被限制在其主目录下.
	* 不使用`chroot_list_file`指定的用户列表, 没有任何"例外"用户.

##### `chroot_local_user=NO`

* `chroot_list_enable=YES`
	* 所有用户都不被限制其主目录下.
	* 使用`chroot_list_file`指定的用户列表,这些用户作为"例外", 受到限制.
* `chroot_list_enable=NO`
	* 所有用户都不被限制其主目录下.
	* 不使用`chroot_list_file`指定的用户列表,没有任何“例外”用户.

#### 11. FTP的工作方式与端口设置

FTP有两种工作方式：PORT FTP(主动模式)和PASV FTP(被动模式)

* `listen_port=21`: 设置FTP服务器建立连接所监听的端口,默认值为21.
* `connect_from_port_20=YES/NO`: 指定FTP使用20端口进行数据传输,默认值为YES.
* `ftp_data_port=20`: 设置在PORT方式下,FTP数据连接使用的端口,默认值为20.
* `pasv_enable=YES/NO(YES)`: 若设置为YES,则使用PASV工作模式；若设置为NO,则使用PORT模式.默认值为YES,即使用PASV工作模式.
* `pasv_max_port=0`: 在PASV工作模式下,数据连接可以使用的端口范围的最大端口,0 表示任意端口.默认值为0.
* `pasv_min_port=0`: 在PASV工作模式下,数据连接可以使用的端口范围的最小端口,0 表示任意端口.默认值为0.

#### 12. 与连接相关的设置

* `listen=YES/NO(YES)`: 设置vsftpd服务器是否以standalone模式运行.以standalone模式运行是一种较好的方式,此时listen必须设置为YES,此为默认值.建议不要更改,有很多与服务器运行相关的配置命令,需要在此模式下才有效.若设置为NO,则vsftpd不是以独立的服务运行,要受到xinetd服务的管控,功能上会受到限制.
* `max_clients=0`: 设置vsftpd允许的最大连接数,默认值为0,表示不受限制.若设置为100时,则同时允许有100个连接,超出的将被拒绝.只有在standalone模式运行才有效.
* `max_per_ip=0`: 设置每个IP允许与FTP服务器同时建立连接的数目.默认值为0,表示不受限制.只有在standalone模式运行才有效.
* `listen_address=IP地址`: 设置FTP服务器在指定的IP地址上侦听用户的FTP请求.若不设置,则对服务器绑定的所有IP地址进行侦听.只有在standalone模式运行才有效.
* `setproctitle_enable=YES/NO(NO)`: 设置每个与FTP服务器的连接,是否以不同的进程表现出来.默认值为NO,此时使用ps aux |grep ftp只会有一个vsftpd的进程.若设置为YES,则每个连接都会有一个vsftpd的进程.

#### 13. 虚拟用户设置

* `pam_service_name=vsftpd`: 设置PAM使用的名称,默认值为/etc/pam.d/vsftpd.
* `guest_enable= YES/NO(NO)`: 启用虚拟用户.默认值为NO.
* `guest_username=ftp`: 这里用来映射虚拟用户.默认值为ftp.
* `virtual_use_local_privs=YES/NO(NO)`: 当该参数激活(YES)时,虚拟用户使用与本地用户相同的权限.当此参数关闭(NO)时,虚拟用户使用与匿名用户相同的权限.默认情况下此参数是关闭的(NO).

#### 14. 其他设置

* `text_userdb_names= YES/NO(NO)`: 设置在执行ls –la之类的命令时,是显示UID、GID还是显示出具体的用户名和组名.默认值为NO,即以UID和GID方式显示.若希望显示用户名和组名,则设置为YES.
* `ls_recurse_enable=YES/NO(NO)`: 若是启用此功能,则允许登入者使用ls –R(可以查看当前目录下子目录中的文件)这个指令.默认值为NO.
* `hide_ids=YES/NO(NO)`: 如果启用此功能,所有档案的拥有者与群组都为ftp,也就是使用者登入使用ls -al之类的指令,所看到的档案拥有者跟群组均为ftp.默认值为关闭.
* `download_enable=YES/NO(YES)`: 如果设置为NO,所有的文件都不能下载到本地,文件夹不受影响.默认值为YES.