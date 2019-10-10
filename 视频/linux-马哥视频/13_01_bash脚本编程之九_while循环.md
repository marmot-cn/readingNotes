#13_01_bash脚本编程之九 while循环

###笔记

---

**问题**

写一个脚本,完成以下功能:  
说明:此脚本能于同一个repo文件中创建多个Yum源的指向

1. 接收一个文件名作为参数,次文件存放至`/etc/yum.repos.d`目录中,且文件名以`.repo`为后缀;要求,此文件不能事先存,否则,报错.
2. 在脚本中,提醒用户输入`repo id`;如果为`quit`,则推出脚本;否则,继续完成下面的步骤;
3. `enabled`默认为1,而`gpgcheck`默认设定为0;
4. 此脚本会循环执行多次,除非用户为`repo id`指定为`quit`;

`代码`:

		#!/bin/bash
		#
		REPOFILE=/etc/yum.repos.d/$1
		
		if [ -e $REPOFILE ]; then
			echo "#1 exists."
			exit 3
		fi
		
		read -p "Repository ID:" REPOID
		until [ $REPOID == 'quit' ]; do
			echo "[$REPOID]" >> $REPOFILE
			read -p "Repository name: " REPONAME
			echo "name=$REPONAME" >> $REPOFILE
			read -p "Repository Baseurl: " REPOURL
			echo "baseurl= $REPOURL" >> $REPOFILE
			echo -e ' enabled=1\ngpgcheck=0' $REPOFILE
			read -p "Repository ID:" REPOID
		done
		
`break`: 提前退出循环.  
`continue`: 提前结束`本轮`循环,而进入下一轮循环.

**while的特殊用法一**

		while :;do
		
		done
		
**while的特殊用法二**
	
		每次读一行,然后放到LINE里面
		
		while read LINE; do
		
		done < /PATH/TO/SOMEFILE
		
**问题**

写一个脚本:

1. 判断一个指定的`bash`脚本是否有语法错误;如果有错误,则提醒用户键入Q或者q无视错误并退出,其它任何键可以通过`vim`打开这个指定脚本;
2. 如果用户通过`vim`打开编辑后保存退出时仍然有错误,则重复第1步中的内容;否则,就正常关闭退出.

`代码`:
		
		#!/bin/bash
		#
		//有错误送到/dev/null,循环
		until bash -n $1 &> /dev/null; do
			read -p "Syntax error, [Qq] to quit, others for editing:" CHOICE
			case $CHOICE in
			q|Q)
				echo "Something wrong, quiting..."
				exit 5
				;;
			*)
				vim + $1
				;;
			esac
		done

###整理知识点

---