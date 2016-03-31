#sql反模式

###目录

---

1. 引言

####逻辑性数据库设计模式

2. 乱穿马路
		
		解决多对多关系
		
		一个字段存储用逗号分隔的数据(违反第一范式)

 		uid     pid
 		1       1,2,3
		
		交叉表(效率高)

 		uid     pid
 		1       1
 		1       2
 		1       3
 		2		2
	
3. 单纯的树

4. 需要ID

		1.理解主键和伪主键
		2.用自然键作主键
		3.用复合键作主键
		
		不要盲目的使用id,分析清楚"主键"和"伪主键",不要为了主键而主键.
		但是具体业务场景要具体分析,我再笔记内也记录了在使用mysql.innodb引擎的一些关于主键的考量.

5. 不用钥匙的入口

		外键的重要性:
		
		1.使用外键保证数据的一致性.
		2.使用外键减少代码的量(用额外的代码来维护数据的完整性),越少的代码意味着越少的错误

6. 实体 - 属性 - 值

7. 多态关联

8. 多列属性

		解决一对多关系 (和第2章的区别是 一对多 和 多对多 关系)

		CREATE TABLE Bugs (
			bug_id		SERIAL	PRIMARY KEY,
			description	CARCHAR(1000),
			tag1		VARCHAR(20),
			tag2		VARCHAR(20),
			tag3		VARCHAR(20),
		);
		
		bug_id	description		tag1	tag2	tag3
		1		xxxx			xx		xx		xx
		
		修改为添加一个'从属表':
		
		CREATE TABLE Tags (
			bug_id		BIGINT	UNSIGNED NOT NULL,
			tag			VARCHAR(20),
			PRIMARY KEY (bug_id, tag),
			FOREIGN KET (bug_id) REFERENCES Bugs(bug_id)
		);	
	
9. 元数据分裂

####物理数据库设计模式

10. 取值错误
	
		FLOAT 遵循IEEE754标准,表示十进制数有误差, 是非精确值.
		
		使用 NUMERIC 或 DECIMAL 来代替 FLOAT 使用

11. 每日新花样

12. 幽灵文件

13. 乱用索引


####查询反模式

14. 对未知的恐惧 

		NULL + 12345 结果为 NULL
		NULL || "string" 结果为 NULL
		
		NULL = 0 结果为 NULL
		NULL <> 12345 结果为 NULL
		NULL = NULL 结果为 NULL
		NULL <> NULL 结果为 NULL
		
		NULL AND TRUE 结果为 NULL
		NULL AND FALSE 结果为 FALSE
		NULL OR FALSE 结果为 FALSE
		NULL OR TRUE 结果为 TRUE
		NOT(NULL) 结果为 NULL
	
15. 模棱两可的分组

16. 随机选择

17. 可怜人的索引

18. 意大利面条式查询

19. 隐式的列


#####应用程序开发反模式

20. 明文密码

21. SQL注入

22. 伪键洁癖

23. 非礼勿视

24. 外交豁免权

25. 魔豆

