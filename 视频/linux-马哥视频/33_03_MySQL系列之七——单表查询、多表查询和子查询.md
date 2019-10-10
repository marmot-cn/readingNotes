# 33_03_MySQL系列之七——单表查询、多表查询和子查询

---

## 笔记

---

### DML

* SELECT
* INSERT INTO
* DELETE
* UPDATE

#### SELECT

`SELECT select-list FROM tb WHERE qualification ORDER BY field_name {ASC|DESC} LIMIT [offset,]count;`

查询语句类型:

* 单表查询
* 多表查询
* 子(嵌套)查询

`SELECT * FROM tb_name;` 显示表中的所有数据.

`SELECT field1,field2 FROM tb_name;` 显示表中所有行的指定字段相关信息. **投影**

`SELECT * FROM tb_name WHERE qualification;` 显示符合搜索指定条件行. **选择**

`FROM`子句, 要查询的关系:

* 表
* 多个表
* 其他`SELECT`语句

`WHERE`子句: 布尔关系表达式. 比较数值不加引号, 比较字符需要加引号.

* `<=>` 和空值(`NULL`)比较, 就算有空值也可以正确比较.
* 在搜索码中进行表达值计算, 无法使用索引.(`WHERE Age+1>20`)
* 逻辑关系:
	* AND or &&
	* OR or ||
	* NOT or !
	* XOR 异或
* LIKE
	* `%`: 任意长度任意字符. (`Y%`, Y开头)
	* `_`: 任意单个字符. (`Y____`, Y开头后面跟4个字符)
* 正则表达式, 索引无效.
	* `REGEXP` or `RLIKE`. (以`M,N,Y`开头 RLIKE '^[MNY].*$`)

##### 排序`ORDER BY`

`ORDER BY field_name {ASC|DESC}`

##### 偏移`LIMIT [offset,]count`

* `offset`: 偏移多少个
* `count`: 取多少个

##### 聚合计算

* `AVG`平均值
* `MAX`最大值
* `MIN`最小值
* `SUM`求和
* `COUNT`总数,个数之和

##### 分组

`GROUP BY field_name HAVING qualification;`

`HAVING`只能和`GROUP BY`一起用, 用来把`GROUP BY`的结果再做一次过滤.

### 执行次序

1. FROM 从哪张表过滤
2. ON 筛选, 只有那些符合<join-condition>的行才会被记录
3. JOIN
4. WHERE 指定条件
5. GROUP BY 分组
6. HAVING 对分组结果过滤
7. SELECT
8. DISTINCT
9. ORDER BY 排序
10. LIMIT

### 多表查询

#### 笛卡尔乘积, 交叉连接

`SELECT * FROM tb1,tb2;`. 两张表交叉出来结果.

#### 自然连接
 
两张表相同字段的值逐一比较, 只将等值关系保留起来.

`SELECT * FROM tb1,tb2 WHERE tb1.field = tb2.field;`
 
#### 外连接

* 左外连接`LEFT JOIN...ON...`
* 右外连接`RIGHT JOIN...ON...`

#### 自连接

自己和自己连接.

### 子查询

比较操作中使用子查询, 子查询只能返回单个值:

`SELECT Name FROM students WHERE Age > (SELECT AVG(age) FROM students);`

在`IN()`中使用子查询, 子查询可以不限定返回单值:

`SELECT Name FROM students WHERE Age IN (SELECT Age FROM tutors);`

在`FROM`中使用子查询:

`SELECT Name, Age FROM (SELECT Name, Age FROM students) AS t WHERE t.Age >=20;`

### 联合查询

`UNION`

## 整理知识点

---

### `WHERE`和`HAVING`的区别

* “Where” 是一个约束声明,使用Where来约束来之数据库的数据,Where是在结果返回之前起作用的,且Where中不能使用聚合函数.
* “Having”是一个过滤声明,是在查询返回结果集以后对查询结果进行的过滤操作,在Having中可以使用聚合函数.