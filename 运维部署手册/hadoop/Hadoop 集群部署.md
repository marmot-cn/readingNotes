# Hadoop 集群部署

---

三台机器:

* 120.24.3.210(公网IP) 10.44.88.189(内网IP): master
* 120.25.87.35(公网IP) 10.170.148.109(内网IP): master
* 120.25.161.1(公网IP) 10.116.138.44((内网IP): master

### readme

这里配置文件有的直接用了`ip`, 只是测试, 最好可以写在`/etc/hosts`里面, 这样不用直接写`ip`.

这里是用`root`用户安装的.

### SSH免密码登录 

这里用`root`用户测试.

#### 每台机器生成公钥

```shell
每台机器依次执行
[root@iZ94xwu3is8Z ~]# ssh-keygen -t rsa
Generating public/private rsa key pair.
Enter file in which to save the key (/root/.ssh/id_rsa):
Created directory '/root/.ssh'.
Enter passphrase (empty for no passphrase):
Enter same passphrase again:
Your identification has been saved in /root/.ssh/id_rsa.
Your public key has been saved in /root/.ssh/id_rsa.pub.
The key fingerprint is:
71:11:f6:58:ac:7e:7c:a4:e9:fa:b7:cc:99:6f:96:e0 root@iZ94xwu3is8Z
The key's randomart image is:
+--[ RSA 2048]----+
|          +o.    |
|         . =.    |
|        . o..    |
|         o.   .  |
|        S. . +   |
|          . =..  |
|           o... .|
|            .E.+o|
|          .o..*=.|
+-----------------+
```

#### 在`master`服务器合并

```shell
[root@iZ94xwu3is8Z .ssh]# cd ~/.ssh/
[root@iZ94xwu3is8Z .ssh]# ls
id_rsa  id_rsa.pub
[root@iZ94xwu3is8Z .ssh]# cat id_rsa.pub>> authorized_keys
[root@iZ94xwu3is8Z .ssh]# ls
authorized_keys  id_rsa  id_rsa.pub

合并slave机器的公钥
[root@iZ94xwu3is8Z .ssh]# ssh root@10.170.148.109  cat ~/.ssh/id_rsa.pub>> authorized_keys
The authenticity of host '10.170.148.109 (10.170.148.109)' can't be established.
ECDSA key fingerprint is eb:17:0d:f9:cc:2c:85:75:ff:d3:0e:4d:14:81:6f:7e.
Are you sure you want to continue connecting (yes/no)? yes
Warning: Permanently added '10.170.148.109' (ECDSA) to the list of known hosts.
root@10.170.148.109's password:

[root@iZ94xwu3is8Z .ssh]# ssh root@10.116.138.44  cat ~/.ssh/id_rsa.pub>> authorized_keys
The authenticity of host '10.170.148.109 (10.170.148.109)' can't be established.
ECDSA key fingerprint is eb:17:0d:f9:cc:2c:85:75:ff:d3:0e:4d:14:81:6f:7e.
Are you sure you want to continue connecting (yes/no)? yes
Warning: Permanently added '10.170.148.109' (ECDSA) to the list of known hosts.
root@10.170.148.109's password:


[root@iZ94xwu3is8Z .ssh]# cat authorized_keys
ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDRyzpy/1Vuf4oB80Y/3zpvOgtTbi0k9iq0IkqJU+WUHa5YC4iXii4CmIOU6KESzs3kE2kJDOh9AkiHkJhWr3jSYS36qXFhAjO/DXVqjn7rJUVc/L5NC+xs1gqeIv3nsw3GDxpfq7cnOwlkRW4dyFwpG/BRNOgKZfQSUHpudEgF3MFK9EX3eme+nVk03PJI5f0rSWefdYFt4FblL7tqYqSoqGbVZBLB94WL7mhT+gTvOYpCL20tPsMjm6nF40PaQtp/cx0KUjU6Y6DmNLm47Yc/0x/bMJ2odLV7uumVEChOapSfFAbFcex56B2ToH2PokIc/C+ivp04EQVa3DyGBeb9 root@iZ94ebqp9jtZ
ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC7HG7mZ6z5r28emLab1LEpQZj2YrzgV2tuy86gaVKSxFlMGQ8dwv9M95HpeV+YET1qM6afHcM7gQT5RK0Xn928lX43ACSdVqbyMRVzuPiQKfj7xn4Celk0QcuNjgUwxRfSQT6dv1ASoiEOtK4pCWi+Nfu6nh7qJhSiX7JleihmIPdj9bOLsVKw10UFXFEbmxT4nnthZMENX0u3ahDwgAkJ5cWoPLqGa/HSW6E6Z69VG3NTXuFwON537o8yQbCWgGVY2B5huBSRKAoIcmOvXbrJkjTOad0wfn8ynSsQed0292Fvet/zYSJP29ekJeOPUzPtDHWND/p3QiBRQ6RGcG1r root@iZ944l0t308Z
ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDRyzpy/1Vuf4oB80Y/3zpvOgtTbi0k9iq0IkqJU+WUHa5YC4iXii4CmIOU6KESzs3kE2kJDOh9AkiHkJhWr3jSYS36qXFhAjO/DXVqjn7rJUVc/L5NC+xs1gqeIv3nsw3GDxpfq7cnOwlkRW4dyFwpG/BRNOgKZfQSUHpudEgF3MFK9EX3eme+nVk03PJI5f0rSWefdYFt4FblL7tqYqSoqGbVZBLB94WL7mhT+gTvOYpCL20tPsMjm6nF40PaQtp/cx0KUjU6Y6DmNLm47Yc/0x/bMJ2odLV7uumVEChOapSfFAbFcex56B2ToH2PokIc/C+ivp04EQVa3DyGBeb9 root@iZ94ebqp9jtZ
```

#### 把`master`服务器的`authorized_keys`、`known_hosts`复制到`Slave`服务器的`/root/.ssh`目录

```shell
[root@iZ94xwu3is8Z .ssh]# scp authorized_keys root@10.170.148.109:/root/.ssh/
root@10.170.148.109's password:
authorized_keys                                                                                                                                                                                 100% 1197     1.2KB/s   00:00
[root@iZ94xwu3is8Z .ssh]# scp known_hosts root@10.170.148.109:/root/.ssh/
known_hosts

[root@iZ94xwu3is8Z .ssh]# scp authorized_keys root@10.116.138.44:/root/.ssh/
authorized_keys                                                                                                                                                                                 100% 1197     1.2KB/s   00:00
[root@iZ94xwu3is8Z .ssh]# scp known_hosts root@10.116.138.44:/root/.ssh/
known_hosts
```

#### 测试免密码登录

```shell
[root@iZ94xwu3is8Z .ssh]# ssh root@10.170.148.109
Last login: Wed Jul 26 20:59:32 2017 from 113.137.211.96

Welcome to aliyun Elastic Compute Service!

[root@iZ94xwu3is8Z ~]#

[root@iZ94ebqp9jtZ .ssh]# ssh root@10.116.138.44
Last login: Wed Jul 26 20:59:44 2017 from 113.137.211.96

Welcome to Alibaba Cloud Elastic Compute Service !

[root@iZ94ebqp9jtZ ~]#
```

### 安装JDK7

3台机器都需要安装JDK, 通过连接`http://www.oracle.com/technetwork/java/javase/downloads/jdk8-downloads-2133151.html`下载`jdk-8u141-linux-x64.rpm`, 我们直接使用`rpm`包安装省事.

#### `master`节点

```shell
[root@iZ94xwu3is8Z ~]# yum localinstall jdk-8u141-linux-x64.rpm -y
Failed to set locale, defaulting to C
Loaded plugins: langpacks
Examining jdk-8u141-linux-x64.rpm: 2000:jdk1.8.0_141-1.8.0_141-fcs.x86_64
Marking jdk-8u141-linux-x64.rpm to be installed
Resolving Dependencies
--> Running transaction check
---> Package jdk1.8.0_141.x86_64 2000:1.8.0_141-fcs will be installed
--> Finished Dependency Resolution

Dependencies Resolved

...

Installed:
  jdk1.8.0_141.x86_64 2000:1.8.0_141-fcs

Complete!
[root@iZ94xwu3is8Z ~]# java -version
java version "1.8.0_141"
Java(TM) SE Runtime Environment (build 1.8.0_141-b15)
Java HotSpot(TM) 64-Bit Server VM (build 25.141-b15, mixed mode)
[root@iZ94xwu3is8Z ~]# java -version
java version "1.8.0_141"
Java(TM) SE Runtime Environment (build 1.8.0_141-b15)
Java HotSpot(TM) 64-Bit Server VM (build 25.141-b15, mixed mode)
```

也可以使用`rpm -ivh jdk-8u141-linux-x64.rpm`来安装.

#### 配置环境变量

```shell
[root@iZ94xwu3is8Z ~]# ls -lrt /usr/bin/java
lrwxrwxrwx 1 root root 22 Jul 26 21:21 /usr/bin/java -> /etc/alternatives/java
[root@iZ94xwu3is8Z ~]# ls -lrt /etc/alternatives/java
lrwxrwxrwx 1 root root 35 Jul 26 21:21 /etc/alternatives/java -> /usr/java/jdk1.8.0_141/jre/bin/java

[root@iZ94xwu3is8Z ~]# vim /etc/profile
...
export JAVA_HOME=/usr/java/jdk1.8.0_141
export JRE_HOME=$JAVA_HOME/jre
export PATH=$PATH:$JAVA_HOME/bin:$JRE_HOME/bin
export CLASSPATH=:$JAVA_HOME/lib/dt.jar:$JAVA_HOME/lib/tools.jar:$JRE_HOME/lib/dt.jar
...

[root@iZ94xwu3is8Z ~]# source /etc/profile
[root@iZ94xwu3is8Z ~]# echo $JAVA_HOME
/usr/java/jdk1.8.0_141
```

#### `slave 节点`

`master`复制`rpm`包到`slave`服务器, 上述一样的命令安装.

### 安装 `Hadoop2.7`

只在`Master`服务器解压, 再复制到`Slave`服务器.

#### 下载`hadoop 2.7.3`**

```shell
[root@iZ94xwu3is8Z ~]# wget http://mirror.bit.edu.cn/apache/hadoop/common/hadoop-2.7.3/hadoop-2.7.3.tar.gz
```

#### 解压

```shell
[root@iZ94xwu3is8Z ~]# mkdir /home/hadoop
[root@iZ94xwu3is8Z ~]# cd /home/hadoop/
[root@iZ94xwu3is8Z hadoop]# mv ~/hadoop-2.7.3.tar.gz ./
[root@iZ94xwu3is8Z hadoop]# tar -xzvf hadoop-2.7.3.tar.gz
```

#### 创建数据存放文件夹

```shell
[root@iZ94xwu3is8Z hadoop]# pwd
/home/hadoop
[root@iZ94xwu3is8Z hadoop]# mkdir tmp
[root@iZ94xwu3is8Z hadoop]# mkdir hdfs
[root@iZ94xwu3is8Z hadoop]# mkdir hdfs/data
[root@iZ94xwu3is8Z hadoop]# mkdir hdfs/name
[root@iZ94xwu3is8Z hadoop]# mkdir var
[root@iZ94xwu3is8Z hadoop]# ls
hadoop-2.7.3  hdfs  tmp  var
```

#### 修改`etc/hadoop`中的一系列配置文件

```shell
[root@iZ94xwu3is8Z etc]# pwd
/home/hadoop/hadoop-2.7.3/etc/hadoop
```

#### `slaves`

```shell
10.170.148.109
10.116.138.44
```

#### `hadoop-env.sh`

不知道为什么这里不该`JAVA_HOME`就报错.

修改`JAVA_HOME`, 修改为本机路径.
```shell
export JAVA_HOME=/usr/java/jdk1.8.0_141
```

##### `core-site.xml`

```shell
<configuration>
	<property>
		<name>fs.defaultFS</name>
		<value>hdfs://10.44.88.189:9000</value>
	</property>
	<property>
		<name>hadoop.tmp.dir</name>
		<value>file:/home/hadoop/tmp</value>
		<description>Abase for other temporary directories.</description>
	</property>
</configuration>
```

##### `hdfs-site.xml`

```shell
<configuration>
         <property>
               <name>dfs.namenode.secondary.http-address</name>
               <value>10.44.88.189:9001</value>
         </property>
         <property>
                 <name>dfs.namenode.name.dir</name>
                 <value>file:/home/hadoop/hdfs/name</value>
         </property>
         <property>
                 <name>dfs.datanode.data.dir</name>
                 <value>file:/home/hadoop/hdfs/data</value>
         </property>
         <property>
                 <name>dfs.replication</name>
                 <value>3</value>
         </property>
         <property>
                 <name>dfs.webhdfs.enabled</name>
                 <value>true</value>
         </property>
</configuration>
```

##### `mapred-site.xml`

```shell
<configuration>
	<property>
   		<name>mapred.job.tracker</name>
   		<value>10.170.148.109:49001</value>
	</property>
	<property>
      		<name>mapred.local.dir</name>
       		<value>/home/hadoop/var</value>
	</property>
	<property>
       		<name>mapreduce.framework.name</name>
       		<value>yarn</value>
	</property>
</configuration>
```

##### `yarn-site.xml`

```shell
	<property>
        <name>yarn.resourcemanager.hostname</name>
        <value>10.44.88.189</value>
   </property>
   <property>
        <description>The address of the applications manager interface in the RM.</description>
        <name>yarn.resourcemanager.address</name>
        <value>${yarn.resourcemanager.hostname}:8032</value>
   </property>
   <property>
        <description>The address of the scheduler interface.</description>
        <name>yarn.resourcemanager.scheduler.address</name>
        <value>${yarn.resourcemanager.hostname}:8030</value>
   </property>
   <property>
        <description>The http address of the RM web application.</description>
        <name>yarn.resourcemanager.webapp.address</name>
        <value>${yarn.resourcemanager.hostname}:8088</value>
   </property>
   <property>
        <description>The https adddress of the RM web application.</description>
        <name>yarn.resourcemanager.webapp.https.address</name>
        <value>${yarn.resourcemanager.hostname}:8090</value>
   </property>
   <property>
        <name>yarn.resourcemanager.resource-tracker.address</name>
        <value>${yarn.resourcemanager.hostname}:8031</value>
   </property>
   <property>
        <description>The address of the RM admin interface.</description>
        <name>yarn.resourcemanager.admin.address</name>
        <value>${yarn.resourcemanager.hostname}:8033</value>
   </property>
   <property>
        <name>yarn.nodemanager.aux-services</name>
        <value>mapreduce_shuffle</value>
   </property>
   <property>
        <name>yarn.scheduler.maximum-allocation-mb</name>
        <value>2048</value>
        <discription>every node default memory8182MB</discription>
   </property>
   <property>
        <name>yarn.nodemanager.vmem-pmem-ratio</name>
        <value>2.1</value>
   </property>
   <property>
        <name>yarn.nodemanager.resource.memory-mb</name>
        <value>2048</value>
   </property>
   <property>
        <name>yarn.nodemanager.vmem-check-enabled</name>
        <value>false</value>
   </property
``` 

### 同步配置文件到`slave`节点

`master`服务器

```shell
scp -r /home/hadoop 10.170.148.109:/home/
scp -r /home/hadoop 10.116.138.44:/home/
```

### 启动`hadoop`

#### 格式化

在`master`启动`hadoop`, 从节点会自动启动, 进入`/home/hadoop/hadoop-2.7.3`目录.

初始化, 输入命令`bin/hdfs namenode -format`.

格式化成功后, 可以在看到在`/home/hadoop/hdfs/name/`目录多了一个`current`目录，而且该目录内有一系列文件

```shell
[root@iZ94xwu3is8Z hadoop-2.7.3]# ls /home/hadoop/hdfs/name/current/
VERSION  fsimage_0000000000000000000  fsimage_0000000000000000000.md5  seen_txid
```

#### 启动

全部启动`sbin/start-all.sh`，也可以分开`sbin/start-dfs.sh`、`sbin/start-yarn.sh`.

```shell
[root@iZ94xwu3is8Z hadoop-2.7.3]# sbin/start-all.sh
This script is Deprecated. Instead use start-dfs.sh and start-yarn.sh
Starting namenodes on [iZ94xwu3is8Z]
iZ94xwu3is8Z: starting namenode, logging to /home/hadoop/hadoop-2.7.3/logs/hadoop-root-namenode-iZ94xwu3is8Z.out
10.170.148.109: starting datanode, logging to /home/hadoop/hadoop-2.7.3/logs/hadoop-root-datanode-iZ944l0t308Z.out
10.116.138.44: starting datanode, logging to /home/hadoop/hadoop-2.7.3/logs/hadoop-root-datanode-iZ94ebqp9jtZ.out
Starting secondary namenodes [iZ94xwu3is8Z]
iZ94xwu3is8Z: starting secondarynamenode, logging to /home/hadoop/hadoop-2.7.3/logs/hadoop-root-secondarynamenode-iZ94xwu3is8Z.out
starting yarn daemons
resourcemanager running as process 28774. Stop it first.
10.170.148.109: starting nodemanager, logging to /home/hadoop/hadoop-2.7.3/logs/yarn-root-nodemanager-iZ944l0t308Z.out
10.116.138.44: starting nodemanager, logging to /home/hadoop/hadoop-2.7.3/logs/yarn-root-nodemanager-iZ94ebqp9jtZ.out
```

#### 检查爱是否启动成功

访问`http://120.24.3.210:50070/`

### 使用`nginx`反向代理外网访问页面

用外网的`18088`端口反向代理到内网的`8088`端口.

```shell
[root@iZ94xwu3is8Z conf.d]# cat hadoop8088.conf
server
{
    listen 18088;
    server_name hadoopproxy;

    location / {
        proxy_set_header Host $host:$server_port;
        proxy_redirect off;

        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_pass http://10.44.88.189:8088;
    }


    access_log /var/log/nginx/hadoopproxy.log;
}
```

### 测试`WordCount`

`jar`文件地址

```shell
[root@iZ94xwu3is8Z mapreduce]# pwd
/root/hadoop-2.7.3/share/hadoop/mapreduce
root@iZ94xwu3is8Z mapreduce]# ls hadoop-mapreduce-examples-2.7.3.jar
hadoop-mapreduce-examples-2.7.3.jar

复制到外面
[root@iZ94xwu3is8Z mapreduce]# cp hadoop-mapreduce-examples-2.7.3.jar ~

查看jar文件
[root@iZ94xwu3is8Z ~]# hadoop-2.7.3/bin/hadoop jar hadoop-mapreduce-examples-2.7.3.jar
An example program must be given as the first argument.
Valid program names are:
  aggregatewordcount: An Aggregate based map/reduce program that counts the words in the input files.
  aggregatewordhist: An Aggregate based map/reduce program that computes the histogram of the words in the input files.
  bbp: A map/reduce program that uses Bailey-Borwein-Plouffe to compute exact digits of Pi.
  dbcount: An example job that count the pageview counts from a database.
  distbbp: A map/reduce program that uses a BBP-type formula to compute exact bits of Pi.
  grep: A map/reduce program that counts the matches of a regex in the input.
  join: A job that effects a join over sorted, equally partitioned datasets
  multifilewc: A job that counts words from several files.
  pentomino: A map/reduce tile laying program to find solutions to pentomino problems.
  pi: A map/reduce program that estimates Pi using a quasi-Monte Carlo method.
  randomtextwriter: A map/reduce program that writes 10GB of random textual data per node.
  randomwriter: A map/reduce program that writes 10GB of random data per node.
  secondarysort: An example defining a secondary sort to the reduce.
  sort: A map/reduce program that sorts the data written by the random writer.
  sudoku: A sudoku solver.
  teragen: Generate data for the terasort
  terasort: Run the terasort
  teravalidate: Checking results of terasort
  wordcount: A map/reduce program that counts the words in the input files.
  wordmean: A map/reduce program that counts the average length of the words in the input files.
  wordmedian: A map/reduce program that counts the median length of the words in the input files.
  wordstandarddeviation: A map/reduce program that counts the standard deviation of the length of the words in the input files.
  
查看 wordcount 使用方法
[root@iZ94xwu3is8Z ~]# hadoop-2.7.3/bin/hadoop jar hadoop-mapreduce-examples-2.7.3.jar wordcount
Usage: wordcount <in> [<in>...] <out>
```

#### 准备2个文件

```shell
[root@iZ94xwu3is8Z ~]# cat file1
Hello world Hello me!
[root@iZ94xwu3is8Z ~]# cat file2
Hello Hadoop Hello you!
```

#### 在HDFS上创建输入文件夹

```shell
[root@iZ94xwu3is8Z ~]# hadoop-2.7.3/bin/hadoop fs -mkdir /input
[root@iZ94xwu3is8Z ~]# hadoop-2.7.3/bin/hadoop fs -put file* /input
[root@iZ94xwu3is8Z ~]# hadoop-2.7.3/bin/hadoop fs -ls /input
Found 2 items
-rw-r--r--   1 root root         22 2017-07-26 23:26 /input/file1
-rw-r--r--   1 root root         24 2017-07-26 23:26 /input/file2
```

#### 运行程序`word count`

```shell
[root@iZ94xwu3is8Z ~]# hadoop-2.7.3/bin/hadoop jar hadoop-mapreduce-examples-2.7.3.jar wordcount /input /output
...

[root@iZ94xwu3is8Z ~]# hadoop-2.7.3/bin/hadoop fs -ls /output/
Found 2 items
-rw-r--r--   1 root root          0 2017-07-26 23:27 /output/_SUCCESS
-rw-r--r--   1 root root         38 2017-07-26 23:27 /output/part-r-00000
[root@iZ94xwu3is8Z ~]# hadoop-2.7.3/bin/hadoop fs -cat /output/part-r-00000
Hadoop	1
Hello	4
me!	1
world	1
you!	1
```