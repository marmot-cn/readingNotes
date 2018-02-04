# mongo export 和 import

---

### 导出备份整个库

```
mongodump -d 数据库 -o ./
```


### 导入整个库

```
mongorestore -h 139.224.65.104:27018 -d 数据库 数据库文件/
```