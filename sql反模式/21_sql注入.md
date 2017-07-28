# 21 SQL注入

## 笔记

---

### 目标: 编写`SQL`动态查询

`SQL`动态查询, 是指将程序中的变量和基本`SQL`语句拼接成一个完整的查询语句.

### 反模式: 将未经验证的输入作为代码执行

`SQL`注入: 当往`SQL`查询的字符串中插入别的内容, 而这些被插入的内容以你不希望的方式修改了查询语句的语法时, `SQL`注入就成功了.

```php
$sql = "SELECT * FROM Bugs WHERE bug_id = $bug_id";

如果把$bug_id变为 1234; DELTE FROM Bugs
$sql = "SELECT * FROM Bugs WHERE bug_id = 1234; DELTE FROM Bugs";
```

#### 意外无处不在

```php
$sql = "SELECT * FROM Projects WHERE project_name = '$project_name'";
```

如果`$project_name`变为`O'Hare`, 上述查询语句就会返回一个语法错误.

#### 对`Web`安全的严重威胁

```php
$sql = "UPDATE Accounts SET password_hash = SHA2('$password') WHERE account_id = $userid";

如遭在链接里面修改userid=xx为userid=xxx OR TRUE

UPDATE Accounts SET password_hash = SHA2('$password')
WHERE account_id = xxx OR TRUE;
```

原来就是修改一个用户的密码, 现在变成了修改全部用户的密码.

#### 寻找治愈良方

##### 转义

`PHP`的`PDO`扩展中使用一个`quote()`函数.

`O'Hare` 转换为 `O\'Hare`.

`account_id = xxx OR TRUE` 转换为 `account_id = 'xxx OR TRUE'`.

##### 查询参数

查询参数的做法是在准备查询语句的时候, 在对应参数的地方使用**参数占位符**. 随后, 在执行这个预先准备好的查询时提供一个参数.

```php
<?php
$stmt = $pdo->prepare("SELECT * FROM Projects WHERE project_name = ?");
$params = array($_REQUEST["name"]);
$stmt->execute($params);
```

查询参数总被视为是一个字面值.

* 多个值得列表不可以当成单一参数
	
		WHERE bug_id IN ( ? )
		$stmt->execute(array("1234,3456,5678"));
		
		会变为
		SELECT * FROM Bugs WHERE bug_id IN ( '1234,3456,5678' );
		
* 表名无法作为参数

### 整理知识点

---