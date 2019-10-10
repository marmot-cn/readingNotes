#修改DOCKER镜像的存储位置

Docker的镜像以及一些数据都是在`/var/lib/docker`目录下,它占用的是Linux的系统分区,我们可以把docker的数据挂载到数据盘.
		
**停止`docker`**

		sudo service docker stop
		
**用`rsync`同步`/var/lib/docker`到新位置**

		rsync -aXS /var/lib/docker/  /docker/
		
**修改`/etc/sysconfig/docker`**

添加

		other_args="-g /docker"
		
		centos7
		
		OPTIONS="-g /docker"
		
`注意`:

For `CentOS 7.x` and `RHEL 7.x`, the name of the variable is `OPTIONS` and for `CentOS 6.x` and `RHEL 6.x`, the name of the variable is `other_args`.

我本机:

		cat /etc/redhat-releas
		CentOS release 6.5 (Final)
		
所以使用`other_args`.
