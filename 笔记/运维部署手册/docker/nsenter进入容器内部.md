#nsenter进入容器内部

---

####下载并安装nsenter

		cd /tmp
		curl https://www.kernel.org/pub/linux/utils/util-linux/v2.24/util-linux-2.24.tar.gz \
		     | tar -zxf-
		cd util-linux-2.24
		./configure --without-ncurses
		make nsenter
		cp nsenter /usr/local/bin
		
####docker-enter

`shell脚本`:

	#!/bin/sh

    if [ -e $(dirname "$0")/nsenter ]; then
        # with boot2docker, nsenter is not in the PATH but it is in the same folder
        NSENTER=$(dirname "$0")/nsenter
    else
        NSENTER=nsenter
    fi

    if [ -z "$1" ]; then
        echo "Usage: `basename "$0"` CONTAINER [COMMAND [ARG]...]"
        echo ""
        echo "Enters the Docker CONTAINER and executes the specified COMMAND."
        echo "If COMMAND is not specified, runs an interactive shell in CONTAINER."
    else
        PID=$(docker inspect --format "{{.State.Pid}}" "$1")
        if [ -z "$PID" ]; then
            exit 1
        fi
        shift

        OPTS="--target $PID --mount --uts --ipc --net --pid --"

        if [ -z "$1" ]; then
            # No command given.
            # Use su to clear all host environment variables except for TERM,
            # initialize the environment variables HOME, SHELL, USER, LOGNAME, PATH,
            # and start a login shell.
            "$NSENTER" $OPTS su - root
        else
            # Use env to clear all host environment variables.
            "$NSENTER" $OPTS env --ignore-environment -- "$@"
        fi
    fi
 
放到`$PATH`内.  
		
		echo $PATH
		/usr/local/bin:/bin:/usr/bin:/usr/local/sbin:/usr/sbin:/sbin:/home/chloroplast/bin
		
		ls /home/chloroplast/bin/
		docker-enter 
    
####一些问题

**nsenter: cannot open /proc/16718/ns/ipc: 权限不够**


		docker-enter daemon_mysql_server
		nsenter: cannot open /proc/16718/ns/ipc: 权限不够

需要`sudo`运行

**docker-enter：找不到命令**
		
		udo docker-enter daemon_mysql_server
		sudo：docker-enter：找不到命令

`sudo`运行后又提示找不到命令.

因为当 sudo以管理权限执行命令的时候,linux将PATH环境变量进行了重置,当然这主要是因为系统安全的考虑,但却使得sudo搜索的路径不是我们想要的PATH变量的路径,当然就找不到我们想要的命令了
		
编辑.bashrc,最后添加alias sudo='sudo env PATH=$PATH'

		cat /home/chloroplast/.bashrc
		# .bashrc

		# Source global definitions
		if [ -f /etc/bashrc ]; then
			. /etc/bashrc
		fi
		
		# User specific aliases and functions
		alias sudo='sudo env PATH=$PATH'
		
####进入容器

进入redmine数据库容器,可以看见redmine数据库

		sudo docker-enter mysql-redmine
		root@56ba37f2af34:~#mysql -uroot -p123456
		
		Warning: Using a password on the command line interface can be insecure.
		Welcome to the MySQL monitor.  Commands end with ; or \g.
		Your MySQL connection id is 7
		Server version: 5.6.28 MySQL Community Server (GPL)
		
		Copyright (c) 2000, 2015, Oracle and/or its affiliates. All rights reserved.
		
		Oracle is a registered trademark of Oracle Corporation and/or its
		affiliates. Other names may be trademarks of their respective
		owners.
		
		Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.
	
		mysql> show databases;
		+--------------------+
		| Database           |
		+--------------------+
		| information_schema |
		| mysql              |
		| performance_schema |
		| redmine            |
		+--------------------+
		4 rows in set (0.01 sec)
		mysql> exit
		Bye
		root@56ba37f2af34:~# exit
		logout
		[chloroplast@iZ94ebqp9jtZ /]$

