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
		
		$stmt  = $pdo->prepare("SELECT * FROM ? WHERE bug_id = 1234");
		$stmt->execute(array("Bugs"));
		将一个字符串插入表名所在的位置, 会得到一个语法错误的提示

* 列名无法作为参数
		
		$stmt  = $pdo->prepare("SELECT * FROM Bugs ORDER BY ?");
		$stmt->execute(array("date_reported"));
		排序是一个无效操作, 列名无法作为参数
		
* `SQL`关键字不能作为参数

		$stmt  = $pdo->prepare("SELECT * FROM Bugs ORDER BY date_reported ?");
		$stmt->execute(array("DESC"));
		
		会返回语法错误

##### 查询是如何完成的

* `RDBMS`服务器首先会解析你所准备好的`SQL`查询语句. 在完成这一步操作之后, 就没有任何方法能够改变那句`SQL`查询语句的结构.
* 在执行一个已经准备好的`SQL`查询时, 你需要提供对对应的参数, 每个你提供的参数都对应于预先准备好的查询中的一个占位符.

##### 存储过程

存储过程包含固定的`SQL`语句, 这些语句是在定义这个存储过程的时候被解析的.

#### 如何识别反模式

使用拼接字符串的形式或者将变量插入到字符串中的方法来构建`SQL`语句, 就会让应用程序暴露在`SQL`注入攻击的威胁下.

#### 合理使用反模式

没有任何合理的理由允许`SQL`注入让程序存在安全漏洞.

#### 解决方案: 不信任任何人

##### 过滤输入内容

将所有不合法的字符从用户输入中剔除掉. 如果你需要一个整数, 那就只是用输入中的整数部分.

##### 参数动态化内容

使用查询参数将其和`SQL`表达式分离.

如果你在`RDMBS`解析完`SQL`语句之后才插入这个参数值, 没有哪种`SQL`注入的攻击能够改变一个参数化了查询的语法结构.

即使攻击者尝试使用带有恶意的参数值, 注入`123 OR TRUE`, `RDBMS`会将这个字符串当成一个完成的值插入.

##### 给动态输入的值加引号

在有些很特殊的情况下, 参数的占位符会导致查询优化器无法正确选择使用哪个索引来进行优化.

比如`Accounts`表中有一个`is_active`列. 这一列中99%的记录都是真实值. 对`is_active=false`的查询会得益于这一列上的索引, 但对于`is_active=true`的查询却会在读取索引的过程浪费很多时间.

如果使用了一个参数`is_active=?`来构造这个表达式, 优化器不知道在预处理这条语句的时候你最终会传入哪个值, 因此很有可能就选择了错误的优化方案.

这个时候直接将变量内容插入到`SQL`语句中会是更好的办法.

##### 将用户与代码隔离

* 声明一个数组, 将用户的可选项作为索引值, 将`SQL`的列名作为对应的值.
* 当用户的选择不在这两个数组中时, 让变量等于默认值.
* 如果用户的选额在数组中, 就使用对应的值.
* 使用变量是安全的, 因为它们只能是在代码中预先定义的值.

优点:

* 从不讲用户的输入与`SQL`查询语句链接, 因此减少了`SQL`注入的风险.
* 可以让`SQL`语句中的任意部分变得动态化, 包括标识,`SQL`关键字, 甚至整句表达式.
* 是用这个办法验证用户的输入变得很简单且高效.
* 能将数据库查询的细节和用户界面解耦.

##### 找个可靠的人来帮你审查代码

* 找出所有使用了程序变量, 字符串练级额或者替换等方法组成的`SQL`语句.
* 跟踪在`SQL`语句中使用的倒台内容的来源. 找出所有的外部的输入, 比如用户输入, 文件, 系统环境, 网络服务, 第三方代码, 甚至于从数据库中获取的字符串.
* 假设任何外部内容都是潜在的威胁. 对于不受信任的内容都要进行过滤, 验证或者使用数组映射的方式来处理.
* 在将外部数据合并到`SQL`语句时, 使用查询参数, 或者用稳健的转义函数预先处理.
* 别忘了在存储过程的代码以及任何其他使用`SQL`动态查询语句的地方做同样的检查.

#### 总结

==让用户输入内容, 但永远别让用户输入代码.==

### 整理知识点

---