# ElastticSearch 性能测试

## 概述

* 导入数据测试
* 搜索数据测试
* 备份数据测试
* 还原数据测试

## 测试

### 导入数据样本

### 导入数据方式

我们以文件形式导入

```
curl -XPOST xxx.xxx.xxx.xxx:9200/resource-catelog/_bulk --data-binary @xxxx.json
```