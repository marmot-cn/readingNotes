# rancher 的 server节点宕机后不能重启问题

---

`20171207`今天发现`master`节点的`server`宕机后, 不能重启. 发现如下错误:

```
...liquibase.exception.LockException: Could not acquire change log lock. Currently locked by <container_ID>
```

我用的是独立数据库, 做如下操作:

```
mysql> use cattle;

# Check that there is a lock in the table
mysql> select * from DATABASECHANGELOGLOCK;

# Update to remove the lock by the container
mysql> update DATABASECHANGELOGLOCK set LOCKED="", LOCKGRANTED=null, LOCKEDBY=null where ID=1;


# Check that the lock has been removed
mysql> select * from DATABASECHANGELOGLOCK;
+----+--------+-------------+----------+
| ID | LOCKED | LOCKGRANTED | LOCKEDBY |
+----+--------+-------------+----------+
|  1 |        | NULL        | NULL     |
+----+--------+-------------+----------+
1 row in set (0.00 sec)
```

然后重启`server`节点即可.