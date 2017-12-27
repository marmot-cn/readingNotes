# 生产环境docker部署mongo副本集

## 版本

部署了`mongo:3.6`

## 生产环境调整参数

简单部署完成后, 查看日志发现:

```
WARNING: Using the XFS filesystem is strongly recommended with the WiredTiger storage engine

WARNING: Access control is not enabled for the database.
Read and write access to data and configuration is unrestricted.

WARNING: /sys/kernel/mm/transparent_hugepage/enabled is 'always'.
We suggest setting it to 'never'

WARNING: /sys/kernel/mm/transparent_hugepage/defrag is 'always'.
We suggest setting it to 'never'
```

查看手册, 还有如下约定:

```
Disable the tuned tool if you are running RHEL 7 / CentOS 7 in a virtual environment.

When RHEL 7 / CentOS 7 run in a virtual environment, the tuned tool automatically invokes a performance profile derived from performance throughput, which automatically sets the readahead settings to 4MB. This can negatively impact performance.
```