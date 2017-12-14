# mysql explain

## 执行结果

## id

`id`列数字越大越先执行, 如果说数字一样大, 那么就**从上往下**依次执行, `id`列为`null`的就表是这是一个结果集, 不需要使用它来进行查询.

## select_type

#### simple

表示不需要`union`操作或者不包含子查询的简单`select`查询. 有连接查询时, 外层的查询为simple, 且只有一个.

#### primary

一个需要`union`操作或者含有子查询的`select`，位于最外层的单位查询的select_type即为primary。且只有一个.

#### union

`union`连接的两个`select`查询, 第一个查询是`dervied`派生表, 除了第一个表外, 第二个以后的表`select_type`都是`union`.

#### dependent union

`dependent union`：与union一样，出现在union 或union all语句中，但是这个查询要受到外部查询的影响.

#### union result

包含union的结果集, 在union和union all语句中, 因为它不需要参与查询，所以id字段为`null`.

#### subquery

除了`from`字句中包含的子查询外, 其他地方出现的子查询都可能是subquery.

#### dependent subquery

dependent subquery：与dependent union类似，表示这个subquery的查询要受到外部表查询的影响.

#### derived

derived：`from`字句中出现的子查询，也叫做派生表，其他数据库中可能叫做内联视图或嵌套`select`.

## table

显示的查询表名, 如果查询使用了别名, 那么这里显示的是别名, 如果不涉及对数据表的操作, 那么这显示为`null`, 如果显示为尖括号括起来的`<derived N>`就表示这个是临时表, 后边的`N`就是执行计划中的`id`, 表示结果来自于这个查询产生. 如果是尖括号括起来的`<union M,N>`, 与`<derived N>`类似, 也是一个临时表, 表示这个结果来自于`union`查询的`id`为`M,N`的结果集.

## type

依次从好到差: system，const，eq_ref，ref，fulltext，ref_or_null，unique_subquery，index_subquery，range，index_merge，index，ALL，除了all之外，其他的type都可以使用到索引，除了index_merge之外，其他的type只可以用到一个索引.

### system

表中只有一行数据或者是空表, 且只能用于myisam和memory表. 如果是Innodb引擎表, type列在这个情况通常都是`all`或者`index`.

### const

使用唯一索引或者主键, 返回记录一定是1行记录的等值where条件时, 通常type是const. 其他数据库也叫做唯一索引扫描.

### eq_ref

出现在要连接过个表的查询计划中, 驱动表只返回一行数据, 这行数据是第二个表的主键或者唯一索引, 且必须为not null, 唯一索引和主键是多列时, 只有所有的列都用作比较时才会出现eq_ref.

### ref

不像`eq_ref`那样要求连接顺序, 也没有主键和唯一索引的要求, 只要使用相等条件检索时就可能出现, 常见与辅助索引的等值查找. 或者多列主键, 唯一索引中, 使用第一个列之外的列作为等值查找也会出现, 总之, 返回数据不唯一的等值查找就可能出现.

### ulltext

全文索引检索, 要注意, 全文索引的优先级很高, 若全文索引和普通索引同时存在时, `mysql`不管代价, 优先选择使用全文索引.

### ref_or_null

与ref方法类似, 只是增加了`null`值的比较. 实际用的不多.

### unique_subquery

用于where中的in形式子查询, 子查询返回不重复值唯一值.

### index_subquery

用于in形式子查询使用到了辅助索引或者in常数列表, 子查询可能返回重复值, 可以使用索引将子查询去重.

### range

索引范围扫描, 常见于使用>,<,is null,between ,in ,like等运算符的查询中.

### index_merge

表示查询使用了两个以上的索引, 最后取交集或者并集, 常见and,or的条件使用了不同的索引, 官方排序这个在ref_or_null之后, 但是实际上由于要读取所有索引, 性能可能大部分时间都不如range.

### index

索引全表扫描, 把索引从头到尾扫一遍, 常见于使用索引列就可以处理不需要读取数据文件的查询,可以使用索引排序或者分组的查询,

### all

这个就是全表扫描数据文件，然后再在server层进行过滤返回符合要求的记录.

## possible_keys

查询可能使用到的索引都会在这里列出来

## key

查询真正使用到的索引, `select_type`为`index_merge`时, 这里可能出现两个以上的索引, 其他的select_type这里只会出现一个.

## key_len

用于处理查询的索引长度, 如果是单列索引, 那就整个索引长度算进去, 如果是多列索引, 那么查询不一定都能使用到所有的列, 具体使用到了多少个列的索引, 这里就会计算进去, 没有使用到的列, 这里不会计算进去. `key_len`只计算where条件用到的索引长度, 而排序和分组就算用到了索引, 也不会计算到key_len中.

## ref

如果是使用的常数等值查询, 这里会显示`const`, 如果是连接查询, 被驱动表的执行计划这里会显示驱动表的关联字段, 如果是条件使用了表达式或者函数, 或者条件列发生了内部隐式转换, 这里可能显示为func.

## rows

这里是执行计划中估算的扫描行数, 不是精确值.

## extra

### distinct

在`select`部分使用了`distinc`关键字.

### no tables used

查询没有from子句, 或者有一个from dual(dual：虚拟表，是为了满足select...from...习惯)子句.

```
EXPLAIN SELECT VERSION();
```

### not exists

在左连接中, 优化器可以通过改变原有的查询组合而使用的优化方法. 当发现一个匹配的行之后, 不再为前面的行继续检索, 可以部分减少数据访问的次数. 例如, 表t1、t2，其中t2.id为not null, 对于SELECT * FROM t1 LEFT JOIN t2 ON t1.id=t2.id WHERE t2.id IS NULL;由于 t2.id非空，所以只可能是t1中有，而t2中没有，所以其结果相当于求差。left join原本是要两边join, 现在Mysql优化只需要依照 t1.id在t2中找到一次t2.id即可跳出.

### const row not found

涉及到的表为空表, 里面没有数据.

### Full scan on NULL key

是优化器对子查询的一种优化方式, 无法通过索引访问NULL值的时候会做此优化.

### Umpossible having

`having`子句总是`false`而不能选择任何列. 

例如

```
having 1=0
```

### Impossible WHERE

`where`子句总是`false`而不能选择任何列, 例如where 1=0

### Impossible WHERE noticed after reading const tables

mysql通过读取"`const/system tables`", 发现Where子句为false. 也就是说: 在where子句中false条件对应的表应该是const/system tables. 这个并不是mysql通过统计信息做出的, 而是真的去实际访问一遍数据后才得出的结论. 当对某个表指定了主键或者非空唯一索引上的等值条件, 一个query最多只可能命中一个结果, mysql在explain之前会优先根据这一条件查找对应记录, 并用记录的实际值替换query中所有用到来自该表属性的地方.

例如: select * from a,b where a.id = 1 and b.name = a.name

执行过程如下: 先根据a.id = 1找到一条记录(1, 'name1'),然后将b.name换成'name1',然后通过a.name = 'name1'查找,发现没有命中记录,最终返回“Impossible WHERE noticed after reading const tables”.

### No matching min/max row

没有行满足如下的查询条件.

例如

```
EXPLAIN SELECT MIN(actor_id) FROM actor WHERE actor_id > 3(只有两条记录)
```

actor_id为唯一性索引时, 会显示"No matching min/max row", 否则会显示"using where".

### no matching row in const table

对一个有`join`的查询, 包含一个空表或者没有数据满足一个唯一索引条件.

### Range checked for each record (index map: N)

### Select tables optimized away

### unique row not found

### Using filesort

### Using index

### Using index for group-by

### Using temporary

### Using where

### Using join buffer

### Scanned N databases

### Using sort_union(…), Using union(…), Using intersect(…)

### Using where with pushed conditio

## filtered

使用`explain extended`时会出现这个列. 这个字段表示存储引擎返回的数据在`server`层过滤后, **剩下多少满足查询的记录数量的比例**, 是百分比, 不是具体记录数.


