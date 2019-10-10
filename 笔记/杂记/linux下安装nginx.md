## linux下安装nginx

## 安装

去官方网站找见`nginx`最新稳定版(https://nginx.org/en/download.html)

也可以直接使用`yum`安装
`sudo rpm -Uvh http://nginx.org/packages/centos/7/noarch/RPMS/nginx-release-centos-7-0.el7.ngx.noarch.rpm`

```shell
[root@localhost ~]# wget -c https://nginx.org/download/nginx-1.12.1.tar.gz
...
[root@localhost ~]# tar -zxvf nginx-1.12.1.tar.gz
...
[root@localhost ~]# cd nginx-1.12.1
[root@localhost nginx-1.12.1]# ./configure
checking for OS
 + Linux 3.10.0-514.el7.x86_64 x86_64
checking for C compiler ... not found

./configure: error: C compiler cc is not found

安装 gcc 
[root@localhost nginx-1.12.1]# yum install gcc -y

安装 pcre-devel, 否则会提示 the HTTP rewrite module requires the PCRE library
[root@localhost nginx-1.12.1]# yum install pcre-devel -y

安装 zlib-devel, 否则会提示 requires the zlib library
[root@localhost nginx-1.12.1]# yum install zlib-devel -y

[root@localhost nginx-1.12.1]# ./configure
...

[root@localhost nginx-1.12.1]# make && make install
...
[root@localhost nginx-1.12.1]# whereis nginx
nginx: /usr/local/nginx
...
nginx命令在`sbin`下
```