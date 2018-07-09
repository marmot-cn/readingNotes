# innodb_autoinc_lock_mode

## 简介

`innodb_autoinc_lock_mode`这个参数控制着在向有`auto_increment`列的表插入数据时, 相关**锁**的行为.

### insert 的三种类型

* `simple insert`如insert into t(name) values('test')
* `bulk insert`如load data | insert into ... select .... from ....
* `mixed insert`如insert into t(id,name) values(1,'a'),(null,'b'),(5,'c');

### innodb_autoinc_lock_mode 说明

`innodb_auto_lockmode`有三个取值:

* **0**这个表示`tradition`传统
* **1**这个表示`consecutive`连续
* **2**这个表示`interleaved`交错

#### tradition(innodb_autoinc_lock_mode=0) 模式

* 提供了一个向后兼容的能力
* 在这一模式下, 所有的insert语句("insert like")都要在语句开始的时候得到一个
表级的auto_inc锁, 在语句结束的时候才释放这把锁, 注意呀, 这里说的是语句级而不是事务级的, **一个事务可能包涵有一个或多个语句**.
* 它能保证值分配的可预见性, 与连续性, 可重复性, 这个也就保证了insert语句在复制到slave的时候还能生成和master那边一样的值(它保证了基于语句复制的安全).
* 由于在这种模式下auto_inc锁一直要保持到语句的结束, 所以这个就影响到了并发的插入.

### consecutive(innodb_autoinc_lock_mode=1) 模式

* 这一模式下去simple insert 做了优化, 由于simple insert一次性插入值的个数可以立马得到确定, 所以mysql可以一次生成几个连续的值, 用于这个insert语句; 总的来说这个对复制也是安全的(它保证了基于语句复制的安全)
* 这一模式也**是mysql的默认模式**, 这个模式的好处是auto_inc锁不要一直保持到语句的结束, 只要语句得到了相应的值后就可以提前释放锁.

### interleaved(innodb_autoinc_lock_mode=2) 模式

由于这个模式下已经没有了auto_inc锁, 所以这个模式下的性能是最好的; 但是它也有一个问题, 就是对于同一个语句来说它所得到的auto_incremant值可能不是连续的.

## 示例

```
mysql> SHOW GLOBAL VARIABLES LIKE '%lock%';
+-----------------------------------------+----------------------+
| Variable_name                           | Value                |
+-----------------------------------------+----------------------+
| block_encryption_mode                   | aes-128-ecb          |
| innodb_api_disable_rowlock              | OFF                  |
| innodb_autoinc_lock_mode                | 1                    |
| innodb_deadlock_detect                  | ON
...
```
