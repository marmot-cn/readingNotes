/*
 Author:chloroplast
 */
# MongoDB 副本集

---
MongoDB的主从集群分为两种:

1. Master-Slave 复制(主从复制)  
2. Replica Sets 复制(副本集)

Master-Slave模式的时候一旦Master停掉，客户端就会报异常，这个时候已经没有Master了，Slave不会自动接管Master.

Replica Sets 也是一种Master-Slave，但它更健壮，一旦Master停掉后，将会在Slave中选举一个作为Master，这种方式也是官方推荐的.

MongoDB在1.6版本开发了replica set,主要增加了故障自动切换和自动修复成员节点.各个DB之间数据完全一致,最为显著的区别在于,副本集没有固定的主节点,它是整个集群选举得出的一个主节点.当其不工作时变更其它节点.

这里只是描述基本的安装流程.

更多参考 [官方手册-复制][id]
[id]:http://docs.mongoing.com/manual-zh/core/replication-introduction.html


###部署环境
因为是在一台电脑上布置,用三个端口模拟3台机子
2个standard节点

	127.0.0.1:28010
	127.0.0.1:28011
	
1个Arbiter节点
	
	127.0.0.1:28012

####第一步:创建目录
		➜  data  mkdir ~/data/mongodb/replset/r0 //机子0号存放数据的目录
		➜  data  mkdir ~/data/mongodb/replset/r1 //机子1号存放数据的目录		
		➜  data  mkdir ~/data/mongodb/replset/r2 //机子2号存放数据的目录		
		➜  data  mkdir ~/data/mongodb/replset/key //存放密钥的目录
		➜  data  mkdir ~/data/mongodb/replset/log //存放日记的目录
		➜  data  chmod 600 ~/data/mongodb/replset/key/r* //600，防止其它程序改写此KEY

####第二步:创建Key
	
	➜  data  echo "replset1 key" > ~/data/mongodb/replset/key/r0
    ➜  data  echo "replset1 key" > ~/data/mongodb/replset/key/r1
    ➜  data  echo "replset1 key" > ~/data/mongodb/replset/key/r2

–keyFile implies –auth. –auth does not imply –keyFile. keyFile 其实包含了 auth 的作用. 如果使用 keyFile 即时没有设置 auth为true, 也需要用户密码(这里即为包含的作用).

如果副本及开启auth,需要key
[key官方资料][id]
[id]:http://docs.mongodb.org/manual/tutorial/deploy-replica-set-with-auth/
####第三步:已副本及的方式启动Mongodb
这里采用的是直接在命令行输出查看
		
		mongod --dbpath=~/data/mongodb/replset/r0 --replSet replset --port 28010 --directoryperdb
		mongod --dbpath=/Users/chloroplast1983/data/mongodb/replset/r1 --replSet replset --port 28011 --directoryperdb
		mongod --dbpath=/Users/chloroplast1983/data/mongodb/replset/r2 --replSet replset --port 28012 --directoryperdb

####第四步:初始化副本集

		mongo --port 28010
		config_replset =
		{
 		_id:"replset",
		members:
			[ 
				{_id:0,host:"127.0.0.1:28010",priority:4},
				{_id:1,host:"127.0.0.1:28011",priority:2},
				{_id:2,host:"127.0.0.1:28012",arbiterOnly : true}
			]
		}
		rs.initiate(config_replset);
![Smaller icon](/img/mongodb/initialReplset.png "初始化")
####第五步:数据同步测试
1. 向PRIMARY主节点写入一条数据
		use test
		db.say.insert({"text":"Hello World"})	
![Smaller icon](/img/mongodb/insertPrimary.png "主节点写入数据")
2. 进入SECONDARY(副节点)查看数据是否同步
![Smaller icon](/img/mongodb/findSecondary.png "副节点读取数据")		
*SECONDARY不能写,而设置slaveOk后,可能从SECONDARY读取数据
默认情况下SECONDARY不能读写，要设定db.getMongo().setSlaveOk();才可以从SECONDARY读取
replSet里只能有一个Primary库，只能从Primary写数据，不能向SECONDARY写数据*
3. ARBITER 读取写入都不能<br />
![Smaller icon](/img/mongodb/canNotRWarbiter.png "ARBITER读写都不能")

####第六步:故障切换测试:把主节点关掉，看副节点是否能接替主节点进行工作

1. 用ctrl+c把28010端口的mongodb服务停掉
![Smaller icon](/img/mongodb/shutDownPrimary.png "停掉PRIMARY")

2. 查看端口28011的情况,发现28011投给自己一票
![Smaller icon](/img/mongodb/secondaryVote.png "28011变为PRIMARY")

3. 查看端口28012的情况,PRIMARY节点28010 DOWN了之后, ARBITER就投票给SECONDARY 28011, SECONDARY 就成为新的PRIMARY节点
![Smaller icon](/img/mongodb/arbiterVote.png "28012投票")


###总结
1. 当副本集的总可投票数为偶数时,可能出现无法选举出主节点的情况
2. 2个Standard节点组成Replication Sets是不合理的,因为不具备故障切换能力
   a. 当SECONDARY Down掉,剩下一个PRIMARY,此时副本集运行不会出问题,因为不用选择PRIMARY节点
   b.当PRIMARY Down掉,此时副本集只剩下一个SECONDARY,它只有1票,不超过总节点数的半数,它不会选举自己为PRIMARY节点!
