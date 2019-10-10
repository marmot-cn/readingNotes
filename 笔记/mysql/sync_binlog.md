# sync_binlog

---

## 简介

是 MySQL 的二进制日志（binary log）同步到磁盘的频率. MySQL server 在`binary log`每写入`sync_binlog`次后，刷写到磁盘,


如果`autocommit`开启, 每个语句都写一次`binary log`, 否则每次事务写一次. 默认值是`0`, 不主动同步, 而依赖操作系统本身不定期把文件内容`flush`到磁盘. 设为`1`最安全, 在每个语句或事务后同步一次`binary log`, 即使在崩溃时也最多丢失一个语句或事务的日志, 但因此也最慢.