# 安装Hive

---

### 安装Mysql

我在其中的一台`slave`节点`120.25.161.1`安装了`mysql 5.6`.

### 安装Hive

#### 下载Hive包

下载链接`http://www-eu.apache.org/dist/hive/hive-2.3.0/apache-hive-2.3.0-bin.tar.gz`

下载后解压到`hadoop`目录下的`thirdparty`目录内.

```shell
master服务器
[root@iZ94xwu3is8Z ~]# wget http://www-eu.apache.org/dist/hive/hive-2.3.0/apache-hive-2.3.0-bin.tar.gz
...
[root@iZ94xwu3is8Z ~]# ls
apache-hive-2.3.0-bin.tar.gz  file1  file2  hadoop-2.7.3  hadoop-2.7.3.tar.gz  hadoop-mapreduce-examples-2.7.3.jar  jdk-8u141-linux-x64.rpm

[root@iZ94xwu3is8Z ~]# mkdir /home/hive
[root@iZ94xwu3is8Z ~]# tar zxvf apache-hive-2.3.0-bin.tar.gz -C /home/hive/
apache-hive-2.3.0-bin/lib/hive-common-2.3.0.jar
apache-hive-2.3.0-bin/lib/hive-shims-2.3.0.jar
apache-hive-2.3.0-bin/lib/hive-shims-common-2.3.0.jar
apache-hive-2.3.0-bin/lib/log4j-slf4j-impl-2.6.2.jar
apache-hive-2.3.0-bin/lib/log4j-api-2.6.2.jar
apache-hive-2.3.0-bin/lib/guava-14.0.1.jar
apache-hive-2.3.0-bin/lib/commons-lang-2.6.jar
apache-hive-2.3.0-bin/lib/libthrift-0.9.3.jar
apache-hive-2.3.0-bin/lib/httpclient-4.4.jar
apache-hive-2.3.0-bin/lib/httpcore-4.4.jar
apache-hive-2.3.0-bin/lib/commons-logging-1.2.jar
apache-hive-2.3.0-bin/lib/commons-codec-1.4.jar
...
```

#### 设置`Hive`环境变量

编辑`/etc/profile`文件

```shell
添加
export HIVE_HOME=/home/hive/apache-hive-2.3.0-bin/
export PATH=$PATH:$HIVE_HOME/bin

使环境变量生效
[root@iZ94xwu3is8Z ~]# source /etc/profile
[root@iZ94xwu3is8Z ~]# echo $HIVE_HOME
/home/hive/apache-hive-2.3.0-bin/
```

#### 配置`Hive`

```shell
[root@iZ94xwu3is8Z ~]# cd /home/hive/apache-hive-2.3.0-bin/conf/
[root@iZ94xwu3is8Z conf]# pwd
/home/hive/apache-hive-2.3.0-bin/conf
[root@iZ94xwu3is8Z conf]# cp hive-env.sh.template hive-env.sh
[root@iZ94xwu3is8Z conf]# cp hive-default.xml.template hive-site.xml
[root@iZ94xwu3is8Z conf]# cp hive-log4j2.properties.template hive-log4j2.properties
[root@iZ94xwu3is8Z conf]# cp hive-exec-log4j2.properties.template hive-exec-log4j2.properties
```

##### 修改`hive-env.sh`

因为 Hive 使用了 Hadoop, 需要在 hive-env.sh 文件中指定 Hadoop 安装路径：

```shell
编辑添加:
export JAVA_HOME=/usr/java/jdk1.8.0_141
export export HADOOP_HOME=/home/hadoop/hadoop-2.7.3/
export HIVE_HOME=/home/hive/apache-hive-2.3.0-bin/
export HIVE_CONF_DIR=$HIVE_HOME/conf
```

##### 修改`hive-site.xml`

```shell
...
    <name>hive.exec.scratchdir</name>
    <value>/tmp/hive-${user.name}</value>
...
    <name>hive.exec.local.scratchdir</name>
    <value>/tmp/${user.name}</value>
...
    <name>hive.downloaded.resources.dir</name>
    <value>/tmp/hive/resources</value>
...
    <name>hive.querylog.location</name>
    <value>/tmp/${user.name}</value>
...
    <name>hive.server2.logging.operation.log.location</name>
    <value>/tmp/${user.name}/operation_logs</value>
```

#### 配置`Hive Metastore`

修改`hive-site.xml`

配置使用`mysql`

```shell
    <name>javax.jdo.option.ConnectionURL</name>
    <value>jdbc:mysql://10.116.138.44:3306/hive?createDatabaseIfNotExist=true&amp;characterEncoding=UTF-8&amp;useSSL=false</value>
    ...
    <name>javax.jdo.option.ConnectionDriverName</name>
    <value>com.mysql.jdbc.Driver</value>
    ...
    <name>javax.jdo.option.ConnectionUserName</name>
    <value>hive</value>
    ...
    <name>javax.jdo.option.ConnectionPassword</name>
    <value>hive</value>
```

#### 为`Hive`创建`HDFS`目录

```shell
[root@iZ94xwu3is8Z hadoop-2.7.3]# ./bin/hdfs dfs -mkdir /tmp
[root@iZ94xwu3is8Z hadoop-2.7.3]# ./bin/hdfs dfs -mkdir -p /usr/hive/warehouse
[root@iZ94xwu3is8Z hadoop-2.7.3]# ./bin/hdfs dfs -chmod g+w /tmp
[root@iZ94xwu3is8Z hadoop-2.7.3]# ./bin/hdfs dfs -chmod g+w /usr/hive/warehouse
```

#### 给`mysql`创建用户`hive`密码`hive`

```sql
mysql> CREATE USER 'hive'@'%' IDENTIFIED BY "hive";
mysql> grant all privileges on *.* to 'hive'@'%' IDENTIFIED BY 'hive';
```

测试:

```shel
root@2bf2ab092e8d:/# mysql -uhive -h120.25.161.1
ERROR 1045 (28000): Access denied for user 'hive'@'113.137.209.168' (using password: NO)
root@2bf2ab092e8d:/# mysql -uhive -h120.25.161.1 -p
Enter password:
Welcome to the MySQL monitor.  Commands end with ; or \g.
Your MySQL connection id is 25
Server version: 5.6.37 MySQL Community Server (GPL)

Copyright (c) 2000, 2016, Oracle and/or its affiliates. All rights reserved.

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
+--------------------+
3 rows in set (0.06 sec)

mysql>
```

引入JDBC驱动包(该驱动包是我自己网上下载的):

```shel
[root@iZ94xwu3is8Z ~]# cp mysql-connector-java-5.1.18-bin.jar /home/hive/apache-hive-2.3.0-bin/lib/
```

#### 运行`hive`

从 Hive 2.1 版本开始, 我们需要先运行 schematool 命令来执行初始化操作.

```shell
[root@iZ94xwu3is8Z ~]# cd /home/hive/apache-hive-2.3.0-bin/
[root@iZ94xwu3is8Z apache-hive-2.3.0-bin]# bin/schematool -dbType mysql -initSchema
SLF4J: Class path contains multiple SLF4J bindings.
SLF4J: Found binding in [jar:file:/home/hive/apache-hive-2.3.0-bin/lib/log4j-slf4j-impl-2.6.2.jar!/org/slf4j/impl/StaticLoggerBinder.class]
SLF4J: Found binding in [jar:file:/home/hadoop/hadoop-2.7.3/share/hadoop/common/lib/slf4j-log4j12-1.7.10.jar!/org/slf4j/impl/StaticLoggerBinder.class]
SLF4J: See http://www.slf4j.org/codes.html#multiple_bindings for an explanation.
SLF4J: Actual binding is of type [org.apache.logging.slf4j.Log4jLoggerFactory]
Metastore connection URL:	 jdbc:mysql://10.116.138.44:3306/hive?createDatabaseIfNotExist=true&characterEncoding=UTF-8&useSSL=false
Metastore Connection Driver :	 com.mysql.jdbc.Driver
Metastore connection User:	 hive
Starting metastore schema initialization to 2.3.0
Initialization script hive-schema-2.3.0.mysql.sql
Initialization script completed
schemaTool completed
[root@iZ94xwu3is8Z apache-hive-2.3.0-bin]# bin/hive
which: no hbase in (/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/usr/java/jdk1.8.0_141/bin:/usr/java/jdk1.8.0_141/jre/bin:/home/hive/apache-hive-2.3.0-bin//bin:/root/bin)
SLF4J: Class path contains multiple SLF4J bindings.
SLF4J: Found binding in [jar:file:/home/hive/apache-hive-2.3.0-bin/lib/log4j-slf4j-impl-2.6.2.jar!/org/slf4j/impl/StaticLoggerBinder.class]
SLF4J: Found binding in [jar:file:/home/hadoop/hadoop-2.7.3/share/hadoop/common/lib/slf4j-log4j12-1.7.10.jar!/org/slf4j/impl/StaticLoggerBinder.class]
SLF4J: See http://www.slf4j.org/codes.html#multiple_bindings for an explanation.
SLF4J: Actual binding is of type [org.apache.logging.slf4j.Log4jLoggerFactory]

Logging initialized using configuration in file:/home/hive/apache-hive-2.3.0-bin/conf/hive-log4j2.properties Async: true
Hive-on-MR is deprecated in Hive 2 and may not be available in the future versions. Consider using a different execution engine (i.e. spark, tez) or using Hive 1.X releases.
hive> show tables;
hive> show tables;
OK
Time taken: 6.554 seconds
```