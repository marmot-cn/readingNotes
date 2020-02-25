# 22 | 撬动离线业务：Job与CronJob

## 笔记

"离线业务", 或者叫作`Batch Job`(计算业务). 这种业务在**计算完成后直接退出了**.

作业分类处理：

* `LRS`(Long Running Service)
* `Batcg Jobs`

### Job API 对象的定义

```
apiVersion: batch/v1
kind: Job
metadata:
  name: pi
spec:
  template:
    spec:
      containers:
      - name: pi
        image: resouer/ubuntu-bc 
        command: ["sh", "-c", "echo 'scale=10000; 4*a(1)' | bc -l "]
      restartPolicy: Never
  backoffLimit: 4
```

"Pod"模板, `spec.template`字段.

Job 对象并不邀请定义一个`spec.selector`来描述要控制哪些`Pod`.

创建一个`Job`后.

```
$ kubectl describe jobs/pi
Name:             pi
Namespace:        default
Selector:         controller-uid=c2db599a-2c9d-11e6-b324-0209dc45a495
Labels:           controller-uid=c2db599a-2c9d-11e6-b324-0209dc45a495
                  job-name=pi
Annotations:      <none>
Parallelism:      1
Completions:      1
..
Pods Statuses:    0 Running / 1 Succeeded / 0 Failed
Pod Template:
  Labels:       controller-uid=c2db599a-2c9d-11e6-b324-0209dc45a495
                job-name=pi
  Containers:
   ...
  Volumes:              <none>
Events:
  FirstSeen    LastSeen    Count    From            SubobjectPath    Type        Reason            Message
  ---------    --------    -----    ----            -------------    --------    ------            -------
  1m           1m          1        {job-controller }                Normal      SuccessfulCreate  Created pod: pi-rq5rl
```

* `Pod`模板被自动加上了一个"controller-uid=< 一个随机字符串 >"这样的`Label`.
* `Job`独享本身, 则被**自动**加上了这个`Label`对应的`Selector`. **保证了Jon与它管理的Pod之间的匹配关系**.

**这种自动生成的`Label`对用户来说并不友好, 所以不太适合推广到`Deployment`等长作业编排对象上**.

`Job`创建的`Pod`进入了`Running`状态, 以为着它正在运行. 运行结束后会进入`Completed`状态. 这就是为什么在`Pod`模板中定义`restartPolicy=Never`的原因. **离线计算的`Pod`永远都不应该被重启, 否则它们会被再重新计算一遍**.

`restartPolicy`:

* `Job`: 只允许被设置为`Never`和`OnFailure`.
* `Deployment`: 只允许被设置为`Always`.

### 离线业务失败了怎么办

`restartPolicy=Never`, 离线作业失败后`Job Controller`会不断地尝试创建一个新`Pod`. `Job`对象的`spec.backoffLimit`字段里定义了重试次数, 这个字段的默认值是`6`. `Job Controller`重新创建`Pod`的间隔是呈指数增加的, 即10s, 20s, 40s...

`restartPolicy=OnFailure`, 离线作业失败后, `Job Controller`就不会去尝试创建新的`Pod`. 但是, 它会不断尝试重启`Pod`里的容器.

### 离线业务最长运行时间

`spec.activeDeadlineSeconds`可以设置离线业务最长运行时间, 超过这个时间, 这个`Job`的所有`Pod`都会被终止. `Pod`的状态里看到的终止原因是`reason: DeadlineExceeded`.

### 并行作业的控制方法

负责并行控制的参数:

1. `spec.parallelism`: 一个`Job`在任意时间最多可以启动多少个`Pod`同时运行.
2. `spec.completions`: `Job`至少要完成的`Pod`数目, 即`Job`的最小完成数.

示例:

```
apiVersion: batch/v1
kind: Job
metadata:
  name: pi
spec:
  parallelism: 2
  completions: 4
  template:
    spec:
      containers:
      - name: pi
        image: resouer/ubuntu-bc
        command: ["sh", "-c", "echo 'scale=5000; 4*a(1)' | bc -l "]
      restartPolicy: Never
  backoffLimit: 4
```
该`Job`最大的并行数是`2`, 最小的完成数是`4`.

`Job`维护了两个状态字段, 即`DESIRED`和`SUCCESSFUL`.

```

$ kubectl get job
NAME      DESIRED   SUCCESSFUL   AGE
pi        4         0            3s
```

**DESIRED, 正是`completions`定义的最小完成数**.

运行过程:

1. 这个`Job`首先创建两个并行运行的`Pod`来计算.
2. 每当有一个`Pod`完成计算进入`Completed`状态时, 就会有一个新的`Pod`被自动创建出来, 并且快速地从`Pending`状态进入到`ContainerCreating`状态.
3. `Job Controller`第二次创建出来的两个并行的 Pod 也进入了 Running 状态.
4. 后面创建的这两个`Pod`也完成了计算，进入了`Completed`状态.

```
# 步骤2

$ kubectl get pods
NAME       READY     STATUS    RESTARTS   AGE
pi-gmcq5   0/1       Completed   0         40s
pi-84ww8   0/1       Pending   0         0s
pi-5mt88   0/1       Completed   0         41s
pi-62rbt   0/1       Pending   0         0s

$ kubectl get pods
NAME       READY     STATUS    RESTARTS   AGE
pi-gmcq5   0/1       Completed   0         40s
pi-84ww8   0/1       ContainerCreating   0         0s
pi-5mt88   0/1       Completed   0         41s
pi-62rbt   0/1       ContainerCreating   0         0s
```

```
步骤 3

$ kubectl get pods 
NAME       READY     STATUS      RESTARTS   AGE
pi-5mt88   0/1       Completed   0          54s
pi-62rbt   1/1       Running     0          13s
pi-84ww8   1/1       Running     0          14s
pi-gmcq5   0/1       Completed   0          54s
```

最终所有`Pod`均已经成功退出, 这个`Job`也就执行完了. 会看见`SUCCESSFUL`字段的值变成了`4`

```
$ kubectl get pods 
NAME       READY     STATUS      RESTARTS   AGE
pi-5mt88   0/1       Completed   0          5m
pi-62rbt   0/1       Completed   0          4m
pi-84ww8   0/1       Completed   0          4m
pi-gmcq5   0/1       Completed   0          5m

$ kubectl get job
NAME      DESIRED   SUCCESSFUL   AGE
pi        4         4            5m
```

### Job Controller 的工作原理

* `Job Controller`控制的对象, 直接就是`Pod`.
* `Job Controller`在控制循环中进行调谐(`Reconcile`)操作, 根据下列值共同计算出这个周期里, 应该创建或删除的`Pod`数目
	* 根据实际在`Running`状态`Pod`的数目
	* 已经成功退出的`Pod`的数目
	* `parallelism`和`completions`

`Job Controller`实际上控制了, 作业执行的**并行度**, 总共需要完成的**任务数**.

### 常用, 使用`Job`对象的方法

#### 1. 简单粗暴: 外部管理器+`Job`模板

把`Job`的`YAML`文件定义为一个"模板", 用一个外部工具控制这些"模板"生成`Job`.

```
apiVersion: batch/v1
kind: Job
metadata:
  name: process-item-$ITEM
  labels:
    jobgroup: jobexample
spec:
  template:
    metadata:
      name: jobexample
      labels:
        jobgroup: jobexample
    spec:
      containers:
      - name: c
        image: busybox
        command: ["sh", "-c", "echo Processing item $ITEM && sleep 5"]
      restartPolicy: Never
```

创建`Job`时, 替换掉`$ITEM`变量. 可以通过`shell`替换, 如:

```
$ mkdir ./jobs
$ for i in apple banana cherry
do
  cat job-tmpl.yaml | sed "s/\$ITEM/$i/" > ./jobs/job-$i.yaml
done
``` 

后续通过`kubectl create`创建即可.

在这种模式下使用`Job`对象,`completions`和`parallelism`这两个字段都应该使用默认值`1`. 作业`Pod`的并行控制, 完全交由外部工具来进行管理(`KubeFlow`).

#### 2. 拥有固定任务数目的并行`Job`

只关心最后是否有指定数目(`spec.completions`)个任务成功退出. 至于执行时的并行度是多少, 并不关心.

使用工作队列.

```
apiVersion: batch/v1
kind: Job
metadata:
  name: job-wq-1
spec:
  completions: 8
  parallelism: 2
  template:
    metadata:
      name: job-wq-1
    spec:
      containers:
      - name: c
        image: myrepo/job-wq-1
        env:
        - name: BROKER_URL
          value: amqp://guest:guest@rabbitmq-service:5672
        - name: QUEUE
          value: job1
      restartPolicy: OnFailure
```

创建这个`Job`后, 会以并发度为`2`的方式, 每两个`Pod`一组, 创建出`8`个`Pod`. 每个 `Pod`都会去连接`BROKER_URL`, 从`RabbitMQ`里读取任务, 然后各自进行处理.

`Pod`里的执行逻辑用伪代码表示.

```
/* job-wq-1的伪代码 */
queue := newQueue($BROKER_URL, $QUEUE)
task := queue.Pop()
process(task)
exit
```

作为用户只关心最终一共有`8`个计算任务启动并且退出.

#### 3. 指定并行度(`parallelism`), 但不设置固定的`completions`的值

必须自己决定什么时候启动新的`Pod`, 什么时候`Job`才算执行完成.

```
apiVersion: batch/v1
kind: Job
metadata:
  name: job-wq-2
spec:
  parallelism: 2
  template:
    metadata:
      name: job-wq-2
    spec:
      containers:
      - name: c
        image: gcr.io/myproject/job-wq-2
        env:
        - name: BROKER_URL
          value: amqp://guest:guest@rabbitmq-service:5672
        - name: QUEUE
          value: job2
      restartPolicy: OnFailure
 ```
 
`Pod`的逻辑.
 
```
/* job-wq-2的伪代码 */
for !queue.IsEmpty($BROKER_URL, $QUEUE) {
  task := queue.Pop()
  process(task)
}
print("Queue empty, exiting")
exit
```

### CronJob

下面这个`CronJob`意思是"从当前开始, 每分钟执行一次".

```
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: hello
spec:
  schedule: "*/1 * * * *"
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: hello
            image: busybox
            args:
            - /bin/sh
            - -c
            - date; echo Hello from the Kubernetes cluster
          restartPolicy: OnFailure
```

关键词`jobTemplate`. `CronJob`是一个`Job`对象的控制器(`Controller`).

`CronJob`是依据`schedule`字段表示的, `"*/1 * * * *"`.

五个部分分别是:

* 分钟
* 小时
* 日
* 月
* 星期

```
$ kubectl create -f ./cronjob.yaml
cronjob "hello" created

# 一分钟后
$ kubectl get jobs
NAME               DESIRED   SUCCESSFUL   AGE
hello-4111706356   1         1         2s
```

`CronJob`会记录下这次`Job`执行的时间.

```
$ kubectl get cronjob hello
NAME      SCHEDULE      SUSPEND   ACTIVE    LAST-SCHEDULE
hello     */1 * * * *   False     0         Thu, 6 Sep 2018 14:34:00 -070
```

如果一个`Job`还没有执行完, 另外一个新`Job`就产生了. 可以定义`spec.concurrencyPolicy`策略:

* `concurrencyPolicy=Allow`, 意味着这些`Job`可以同时存在.
* `concurrencyPolicy=Forbid`, 不会创建新的`Pod`, **该创建周期被跳过**.
* `concurrencyPolicy=Replace`, 新产生的会替换旧的, 没有执行完成的`Job`.

某一次`Job`创建失败, 这次创建就会被标记为`"miss"`. 当在指定的时间窗口内, `miss`的数目达到`100`时, 那么`CronJob`会停止再创建这个`Job`.

时间窗口可以通过`spec.startingDeadlineSeconds`指定, 即多少时间内`miss`达到了`100`次, 那么这个`Job`就不会被创建执行了.

## 扩展