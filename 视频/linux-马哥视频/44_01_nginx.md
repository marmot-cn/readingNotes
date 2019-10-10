# 44_01_nginx

---

## 笔记

### Nginx

* mmap
* event-driven
	* 一个进程响应多个请求, 单线程进程. 对内存占用量低.
* aio

### LNMP

Nginx + (FastCGI) + php-fpm

`Apache`是基于模块与`php`通信.

`php 5.6`的`opcode`只能在单个进程能共享, 可以使用`xcache`缓存共享`opcode`.

### memcached

缓存可序列化数据.

## 整理知识点