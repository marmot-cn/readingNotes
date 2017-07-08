# sql反模式

### 目录

---

**1** 引言

#### 逻辑性数据库设计模式

**2** 乱穿马路
		
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
	
**3** 单纯的树

**4** 需要ID

		1.理解主键和伪主键
		2.用自然键作主键
		3.用复合键作主键
		
		不要盲目的使用id,分析清楚"主键"和"伪主键",不要为了主键而主键.
		但是具体业务场景要具体分析,我再笔记内也记录了在使用mysql.innodb引擎的一些关于主键的考量.

**5** 不用钥匙的入口

		外键的重要性:
		
		1.使用外键保证数据的一致性.
		2.使用外键减少代码的量(用额外的代码来维护数据的完整性),越少的代码意味着越少的错误

**6** 实体 - 属性 - 值

		EAV 设计:
		每个属性都是一行数据,即:
		
		id,
		attr_name,
		attr_value
		
		不要使用EAV设计,可以使用如下方法:
		
		1.单表继承
		所有的字段都放在一张表中,用单独一个type字段来区分每行数据的归属.会有冗余数据.
		
		一张表
		
		order表:
		//order公共属性
		type 字段,区分属于 order_car 还是 order_computer
		//order_car 汽车订单的字段
		//order_computer 电脑订单的字段
		
		2.实体表继承
		每个实体都存放在单独的表中,相同的属性在每张表中都会重复出现.
		
		二张表
		
		order_car表:
		//order 属性
		//order_car 汽车订单的字段
		
		order_computer表:
		//order 属性
		//order_computer 电脑订单的字段

		
		3.类表继承
		模拟面向对象的继承方法,把公共的属性抽离成一张表.其他自己私有的属性一张表.列入:
		order,//包含订单的通用属性,状态,支付时间等等...
		order_car//汽车订单中独有的字段,表内的 order_id 和 order 表的 id 为对应关系
		order_computer//电脑订单中独有的字段,表内的 order_id 和 order 表的 id 为对应关系
		
		三张表
		
		order表:
		//order公共属性
		
		order_car表:
		//order_car 汽车订单的字段
		
		order_computer表:
		//order_computer 电脑订单的字段
		
		4.半结构化数据
		类似单表继承,但是把每个种类的私有的字段统一以序列化的形式存放在一个text字段中.用的时候通过程序反序列化取出.
		
		order表:
		//order公共属性
		type 字段,区分属于 order_car 还是 order_computer
		//attributes,序列化存储.如果 type == order_car 则该字段为序列化的 order_car 汽车订单私有字段.如果 type == order_computer 则该字段为序列化的 order_computer 电脑订单私有字段,


**7** 多态关联

**8** 多列属性

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
	
**9** 元数据分裂

####物理数据库设计模式

**10** 取值错误
	
		FLOAT 遵循IEEE754标准,表示十进制数有误差, 是非精确值.
		
		使用 NUMERIC 或 DECIMAL 来代替 FLOAT 使用

**11** 每日新花样

**12** 幽灵文件

**13** 乱用索引

		正确的使用索引: (无索引, 过多索引)
		
		复合索引的规则: 最左前缀
		
		不能使用索引的情况(具体看连接的笔记)
		
		索引分离率: 所有不重复的值得数量和总记录条数之比. 分离率越低,索引的效率越低
		
		使用Mysql的 EXPLAIN 分析建立的索引
		
		定期重建索引 ALALYZE TABLE, OPTIMIZE TABLE

#### 查询反模式

**14** 对未知的恐惧 

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
	
**15** 模棱两可的分组

**16** 随机选择

**17** 可怜人的索引

**18** 意大利面条式查询

**19** 隐式的列


##### 应用程序开发反模式

**20** 明文密码

**21** SQL注入

**22** 伪键洁癖

**23** 非礼勿视

**24** 外交豁免权

**25** 魔豆

