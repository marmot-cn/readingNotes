# 29_04_samba之二

---

### 笔记

---

#### smbclient 访问共享

```shell
安装 smb 客户端
[root@iZ94ebqp9jtZ ~]# yum install -y samba-client
```

**参数**

* `-L NetBios_Name` 指定netbios主机名, 跟`ip`地址也可以
* `-U Username` 指定用户名
* `-P` 输入密码

**示例**

```shell
我测试的访问另外一台linux的samba共享目录, 没有链接windwos服务器

[root@iZ94ebqp9jtZ ~]# smbclient -L 120.25.87.35 -U chloroplast
Enter chloroplast's password:
Domain=[SAMBA] OS=[Windows 6.1] Server=[Samba 4.4.4]

	Sharename       Type      Comment
	---------       ----      -------
	print$          Disk      Printer Drivers
	tools           Disk      Share testing
	IPC$            IPC       IPC Service (Samba 4.4.4)
	chloroplast     Disk      Home Directories
Domain=[SAMBA] OS=[Windows 6.1] Server=[Samba 4.4.4]

	Server               Comment
	---------            -------

	Workgroup            Master
	---------            -------
	
直接访问 samba 服务器
[root@iZ94ebqp9jtZ ~]# smbclient //120.25.87.35 -U chloroplast

\\120.25.87.35: Not enough '\' characters in service
Usage: smbclient [-?EgBVNkPeC] [-?|--help] [--usage] [-R|--name-resolve NAME-RESOLVE-ORDER] [-M|--message HOST] [-I|--ip-address IP] [-E|--stderr] [-L|--list HOST] [-m|--max-protocol LEVEL] [-T|--tar <c|x>IXFqgbNan]
        [-D|--directory DIR] [-c|--command STRING] [-b|--send-buffer BYTES] [-t|--timeout SECONDS] [-p|--port PORT] [-g|--grepable] [-B|--browse] [-d|--debuglevel DEBUGLEVEL] [-s|--configfile CONFIGFILE]
        [-l|--log-basename LOGFILEBASE] [-V|--version] [--option=name=value] [-O|--socket-options SOCKETOPTIONS] [-n|--netbiosname NETBIOSNAME] [-W|--workgroup WORKGROUP] [-i|--scope SCOPE] [-U|--user USERNAME] [-N|--no-pass]
        [-k|--kerberos] [-A|--authentication-file FILE] [-S|--signing on|off|required] [-P|--machine-pass] [-e|--encrypt] [-C|--use-ccache] [--pw-nt-hash] service <password>
[root@iZ94ebqp9jtZ ~]# smbclient //120.25.87.35/tools -U chloroplast
Enter chloroplast's password:
Domain=[SAMBA] OS=[Windows 6.1] Server=[Samba 4.4.4]
smb: \> ls
  .                                   D        0  Sat Jul 22 15:41:03 2017
  ..                                  D        0  Sat Jul 22 14:12:24 2017
  222.txt                             N        3  Sat Jul 22 15:41:10 2017
  111                                 N        0  Sat Jul 22 15:40:53 2017

		20510332 blocks of size 1024. 17599396 blocks available
smb: \> get 222.txt
getting file \222.txt of size 3 as 222.txt (0.1 KiloBytes/sec) (average 0.1 KiloBytes/sec)
smb: \> exit
[root@iZ94ebqp9jtZ ~]# ls
222.txt
[root@iZ94ebqp9jtZ ~]# cat 222.txt
33
```

### 整理知识点

---