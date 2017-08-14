# 33_04_MySQL系列之八——多表查询、子查询及视图

---

## 笔记

---

### 练习

1. 挑选出`courses`表中没有被`students`中的`CID2`学习的课程的课程名称
2. 显示每一位老师及其所教授的课程, 没有教授的课程保持为`NULL`
3. 显示每一个课程及其相关老师, 没有老师教授的课程将其老师显为空.
4. 显示每位同学`CID1`课程的课程名称及其讲授了相关课程的老师名称.
5. 挑选出没有教授任何课程的老师
6. 挑出`students`表中`CID1`有两个或两个以上同学学习了的同一门课程的课程名称.

#### `1`

```sql
SELECT Cname FROM courses WHERE CID NOT IN (SELECT DISTINCT CID2 FROM students WHERE CID2 is NOT NULL
```

#### `2`

```sql
SELECT t.Tname, C.Cname FROM tutors AS t LEFT JOIN courses AS c ON t.TID=c.TID;
```

#### `3`

```sql
SELECT t.Tname, C.Cname FROM tutors AS t RIGHT JOIN courses AS c ON t.TID=c.TID;
```

#### `4`

```sql
SELECT Name,Cname,Tname FROM students,courses,tutors WHERE students.CID1=courses.CID AND courses.TID=tutors.TID
```

#### `5`

`SELECT Tname FROM tutors WHERE TID NOT IN (SELECT DISTINCT TID FROM courses)`

#### `6`

```sql
SELECT Cname FROM courses WHERE CID IN (SELECT CID1 FROM students GROUP BY CID1 HAVING COUNT(CID1) >=2
```

### 视图

存储下来的`SELECT`语句, 不存储任何内容, 把这个`SELECT`语句当做表来使用. 基于基表的查询结果.

视图`VIEW`.

创建视图`CREATE VIEW`.

```sql
CREATE VIEW xxx AS sql语句

使用 show tables 视图会被当做一张表.
```

不允许往视图插入数据(理论上可以插入).

视图也被称为虚表, 背后的表称为基表.

`SHOW CRATE VIEW view_name`可以看视图怎么创建的.

#### 删除视图

`DROP VIEW view_name;`

#### 物化视图

可以把查询结果保存起来.

### `mysql -e SQL语句`

可以在`shell`脚本内使用, 获取返回结果.

```shell
root@90af5cf2869e:/# mysql -uroot -p123456 -e "show databases;"
Warning: Using a password on the command line interface can be insecure.
+------------------------+
| Database               |
+------------------------+
| information_schema     |
| hellodb                |
| marmot                 |
| marmot_test            |
| mysql                  |
| performance_schema     |
| qixinyun_purview       |
| qixinyun_purview_test  |
| qixinyun_workflow      |
| qixinyun_workflow_test |
| saas_member            |
| saas_member_test       |
| saas_order             |
| saas_order_test        |
| saas_product           |
| saas_product_test      |
| test                   |
+------------------------+
```

## 整理知识点

---