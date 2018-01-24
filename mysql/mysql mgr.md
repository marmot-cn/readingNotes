# mysql mgr 

---

创建3个文件

```
mkdir -p mysql-1/config mysql-1/data
mkdir -p mysql-2/config mysql-2/data
mkdir -p mysql-3/config mysql-3/data
```

mysql-1, 配置文件:

```
cat my.cnf
[mysqld]

datadir=/var/lib/mysql
pid-file=/var/run/mysqld/mysqld.pid

port=24803  
socket=/var/run/mysqld/mysqld.sock

#
# Replication configuration parameters
#
server_id=1
gtid_mode=ON
enforce_gtid_consistency=ON
master_info_repository=TABLE
relay_log_info_repository=TABLE
binlog_checksum=NONE
log_slave_updates=ON
log_bin=binlog
binlog_format=ROW

#
# Group Replication configuration
#
transaction_write_set_extraction=XXHASH64
loose-group_replication_group_name="aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
loose-group_replication_start_on_boot=off
loose-group_replication_local_address= "10.44.88.189:24901"
loose-group_replication_group_seeds= "10.44.88.189:24901,10.44.88.189:24902,10.44.88.189:24903"
loose-group_replication_bootstrap_group= off

# multi-primary node
loose-group_replication_single_primary_mode=off
loose-group_replication_enforce_update_everywhere_checks=ON
```

每个配置文件的`port`都要修改, 使用 docker --net=host 模式.


启动

```
docker run --net=host -v /data/mysql-1/config/:/etc/mysql/conf.d -v /data/mysql-1/data/:/var/lib/mysql --name=mysql-1 -e MYSQL_ROOT_PASSWORD=123456 -d docker.io/mysql:5.7

docker run --net=host -v /data/mysql-2/config/:/etc/mysql/conf.d -v /data/mysql-2/data/:/var/lib/mysql --name=mysql-2 -e MYSQL_ROOT_PASSWORD=123456 -d docker.io/mysql:5.7

docker run --net=host -v /data/mysql-3/config/:/etc/mysql/conf.d -v /data/mysql-3/data/:/var/lib/mysql --name=mysql-3 -e MYSQL_ROOT_PASSWORD=123456 -d docker.io/mysql:5.7
```


mysql 内执行:

```
mysql> SET SQL_LOG_BIN=0;
mysql> CREATE USER rpl_user@'%' IDENTIFIED BY 'password';
mysql> GRANT REPLICATION SLAVE ON *.* TO rpl_user@'%';
mysql> FLUSH PRIVILEGES;
mysql> SET SQL_LOG_BIN=1;
mysql> CHANGE MASTER TO MASTER_USER='rpl_user', MASTER_PASSWORD='password' FOR CHANNEL 'group_replication_recovery';
mysql> INSTALL PLUGIN group_replication SONAME 'group_replication.so';
+----------------------------+----------+--------------------+----------------------+-------------+
| Name                       | Status   | Type               | Library              | License     |
+----------------------------+----------+--------------------+----------------------+-------------+
| binlog                     | ACTIVE   | STORAGE ENGINE     | NULL                 | PROPRIETARY |

(...)

| group_replication          | ACTIVE   | GROUP REPLICATION  | group_replication.so | PROPRIETARY |
+----------------------------+----------+--------------------+----------------------+-------------+

mysql> SET GLOBAL group_replication_bootstrap_group=ON;
mysql> START GROUP_REPLICATION;
mysql> SET GLOBAL group_replication_bootstrap_group=OFF;
mysql> SELECT * FROM performance_schema.replication_group_members;
```

创建测试数据:

```
mysql> CREATE DATABASE test;
mysql> USE test;
mysql> CREATE TABLE t1 (c1 INT PRIMARY KEY, c2 TEXT NOT NULL);
mysql> INSERT INTO t1 VALUES (1, 'Luis');
mysql> SELECT * FROM t1;
```

其他节点加入服务器:

```
mysql> SET SQL_LOG_BIN=0;
mysql> CREATE USER rpl_user@'%';
mysql> GRANT REPLICATION SLAVE ON *.* TO rpl_user@'%' IDENTIFIED BY 'password';
mysql> SET SQL_LOG_BIN=1;
mysql> CHANGE MASTER TO MASTER_USER='rpl_user', MASTER_PASSWORD='password' FOR CHANNEL 'group_replication_recovery';
mysql> INSTALL PLUGIN group_replication SONAME 'group_replication.so';
mysql> set global group_replication_allow_local_disjoint_gtids_join=ON;
mysql> START GROUP_REPLICATION;
```

一些问题:

### 1

一开始我在单台主机测试, 其他节点不能启动. 需要执行`set global group_replication_allow_local_disjoint_gtids_join=ON;`

### 2

发现其他节点同步数据的时候不能识别`host`, 我在容器内部给`hosts`文件写入主机名映射`ip`.

需要在主机名内映射主机名对应的`ip`地址, 在宿主机执行即可.

### 3

有子节点挂了, 我需要重新加入.

`STOP GROUP_REPLICATION;`

然后

`START GROUP_REPLICATION;`

### 安装 mysql-router

```
sudo rpm -ivh mysql-router-2.1.4-1.el7.x86_64.rpm
```

配置文件

```
[root@demo ansible]# cat /etc/mysqlrouter/mysqlrouter.conf
[DEFAULT]
logging_folder = /var/log/mysqlrouter/
plugin_folder = /usr/lib64/mysqlrouter
runtime_folder = /var/run/mysqlrouter
config_folder = /etc/mysqlrouter

[logger]
level = info

[routing:read_write]
bind_address = 120.24.3.210
bind_port = 7002
destinations = 10.44.88.189:24803,10.44.88.189:24804,10.44.88.189:24805
mode = read-write
```

启动

```
systemctl start mysqlrouter.service
```

问题:

默认用户和组是`mysqlrouter`, 我换成`root`才可以启动成功.

```
[root@demo ansible]# cat /usr/lib/systemd/system/mysqlrouter.service
[Unit]
Description=MySQL Router
After=syslog.target
After=network.target

[Service]
Type=simple
#User=mysqlrouter
#Group=mysqlrouter

PIDFile=/var/run/mysqlrouter/mysqlrouter.pid

ExecStart=/usr/bin/mysqlrouter -c /etc/mysqlrouter/mysqlrouter.conf

PrivateTmp=true

[Install]
WantedBy=multi-user.target
```