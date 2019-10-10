# 42_04_MySQL主从复制——mysql-proxy.0.8.3实现MySQL-5.6读写分离

---

## 笔记

### MySQL Proxy

插件化, 内置了`lua`的引擎. 

`MySQL Proxy`提供了在某种场景下, 让你调用某种特定脚本, 实现你的特定功能的**框架**.  严重依赖脚本. 

`lua`嵌入式脚本.

### 配置文件

配置文件可以放在`mysql.cnf`中的`[mysql-proxy]`块内. `mysql proxy`只读该区域的配置信息.

## 整理知识点

---