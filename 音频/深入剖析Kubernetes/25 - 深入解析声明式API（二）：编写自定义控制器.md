# 25 | 深入解析声明式API（二）：编写自定义控制器

## 笔记

为`24`章的`Network`自定义`API`对象编写一个自定义控制器(`Custom Controller`).

**基于声明式 API 的业务功能实现, 往往需要通过扩至器模式来"监视" API 对象的变化 (比如, 创建或者删除 Network), 然后以此来决定实际要执行的具体工作**.

编写自定义控制器代码包括:

* `main`函数
* 自定义控制器的定义
* 控制器里的业务逻辑

### 编写`main`函数

`main`函数的主要工作就是, 定义并初始化一个自定义控制器(`Custom Controller`), 然后启动它.

```
func main() {
  ...
  
  cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
  ...
  kubeClient, err := kubernetes.NewForConfig(cfg)
  ...
  networkClient, err := clientset.NewForConfig(cfg)
  ...
  
  networkInformerFactory := informers.NewSharedInformerFactory(networkClient, ...)
  
  controller := NewController(kubeClient, networkClient,
  networkInformerFactory.Samplecrd().V1().Networks())
  
  go networkInformerFactory.Start(stopCh)
 
  if err = controller.Run(2, stopCh); err != nil {
    glog.Fatalf("Error running controller: %s", err.Error())
  }
}
```

通过三步完成了初始化并启动一个自定义控制器的工作

#### 第一步

`main`函数根据`Master`配置, 创建一个`Kubernetes`的`client`(kubeClient)和`Network`对象的`client`(networkClient).

如果没有提供`Master`配置, `main`函数会使用一种**InClusterConfig**的方式来创建这个`client`, 这个方式会假设你的自定义控制器是以`Pod`的方式运行在`Kubernetes`集群里的.

#### 第二步

`main`函数为`Network`对象创建一个叫做`InformerFactory`(`networkInformerFactory`), 并使用它生成一个`Network`对象的`Informer`传递给控制器.

#### 第三步

`main`函数启动上述的`Informer`, 然后执行`controller.Run`, 启动自定义控制器.

### 自定义控制器工作原理

![](./img/25_01.png)

**控制器, 从 Kubernetes 的 APIServer 里获取它所关心的对象, 即自定义的 Network 对象**

该操作依靠`Informer`(通知器)完成的. `Informer`与`API`对象是一一对应的, 传递给自定义控制器的, 正式一个`Network`对象的`Informer`.

`Network Informer`使用`networkClient`跟`APIServer`建立了连接. `Informer`使用的`Reflector`包负责维护这个连接.

`Reflector`使用的是`ListAndWatch`方法来"获取"并"监听"这些`Network`对象实例的变化. 一旦`APIServer`端有新的`Network`实例被创建, 删除或者更新, `Reflector`都会收到"事件通知". 该事件及它对应的`API`对象这个组合, 被称为增量`Delta`, 被放进一个`Delta FIFO Queue`(增量先进先出队列)中.

`Informe`会不断地从这个`Delta FIFO Queue`里读取(`Pop`)增量. 每拿到一个增量, `Informer`就会判断这个增量里的事件类型, 然后创建或者更新本地对象的缓存(在`Kubernetes`里叫做`Store`).

* `Added`(添加对象), `Informer`会通过一个叫做`Indexer`的库把这个增量的`API`对象保存在本地缓存汇总, 创建索引.
* `Deleted`(删除对象), `Informer`会从本地缓存中删除这个对象.

**同步本地缓存**是**Informer**的一个职责.

`Informer`的第二个职责, 根据这些事件的类型, 触发事先注册好的`ResourceEventHandler`. 这些`Handler`, 需要在创建控制器的时候注册给它对应的`Informer`.

控制器定义

```
func NewController(
  kubeclientset kubernetes.Interface,
  networkclientset clientset.Interface,
  networkInformer informers.NetworkInformer) *Controller {
  ...
  controller := &Controller{
    kubeclientset:    kubeclientset,
    networkclientset: networkclientset,
    networksLister:   networkInformer.Lister(),
    networksSynced:   networkInformer.Informer().HasSynced,
    workqueue:        workqueue.NewNamedRateLimitingQueue(...,  "Networks"),
    ...
  }
    networkInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
    AddFunc: controller.enqueueNetwork,
    UpdateFunc: func(old, new interface{}) {
      oldNetwork := old.(*samplecrdv1.Network)
      newNetwork := new.(*samplecrdv1.Network)
      if oldNetwork.ResourceVersion == newNetwork.ResourceVersion {
        return
      }
      controller.enqueueNetwork(new)
    },
    DeleteFunc: controller.enqueueNetworkForDelete,
 return controller
}
```

用前面`main`函数里创建的两个`client`和前面创建的`Informer`, 初始化了自定义控制器.

为`networkInformer`注册了三个`Handler`(`AddFunc`, `UpdateFunc`和`DeleteFunc`), 分别对应`API`对象的"添加""更新"和"删除"事件. 具体的处理操作, 将该事件对应的`API`对象加入到工作队列中.

实际入队的是`API`对象的`Key`.(`<namespace>/<name>`).

后面的控制循环, 会不断地从这个工作队列拿到这些`Key`, 然后开始执行真正的控制逻辑.

**所谓的`Informer`, 其实就是一个带有本地缓存和索引机制的, 可以注册 EventHandler 的 client**, 是自定义控制器跟`APIServer`进行数据同步的重要组件.

#### Informer 原理

* `Informer`通过`ListAndWatch`的方法, 把`APIServer`中的API对象缓存在了本地, 并负责更新和维护这个缓存.
* `ListAndWatich`, 通过`LIST API`获取所有最新版本的`API`对象. 再通过`WATCH API`来"监听"所有这些`API`对象的变化.
* 监听到的时间变化, `Informer`实时地更新本地缓存, 调用这些事件对应的`EventHandler`.

过程中, 每经过`resyncPeriod`指定的时间, `Informer`维护的本地缓存, 都会使用最近一次`LIST`返回的结果强制更新一次, 从而保证缓存的有效性. 在`Kubernetes`中, 这个缓存强制更新的操作较`resync`.

`resync`也会触发`Informer`注册的"更新"时间, 此时因为新, 旧两个`Network`对象的`ResouceVersion`是一样的. 这种情况下, `Informer`就不需要对这个更新事件在做进一步的处理了(上述代码也判断了新, 旧两个`Network`对象的版本(ResourceVersion)是否发生了变化, 然后才开始进行的入队操作.

### 控制循环(`Control Loop`)

是`main`函数最后调用`controller.Run()`启动的"控制循环".

```
func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
 ...
  if ok := cache.WaitForCacheSync(stopCh, c.networksSynced); !ok {
    return fmt.Errorf("failed to wait for caches to sync")
  }
  
  ...
  for i := 0; i < threadiness; i++ {
    go wait.Until(c.runWorker, time.Second, stopCh)
  }
  
  ...
  return nil
}
```

逻辑

* 等待`Informer`完成一次本地缓存的数据同步操作
* 直接通过`goroutine`启动一个(或者并发启动多个)”无限循环"的任务

无限循环任务的每一个循环周期, 执行的是真正的业务逻辑.

```
func (c *Controller) runWorker() {
  for c.processNextWorkItem() {
  }
}

func (c *Controller) processNextWorkItem() bool {
  obj, shutdown := c.workqueue.Get()
  
  ...
  
  err := func(obj interface{}) error {
    ...
    if err := c.syncHandler(key); err != nil {
     return fmt.Errorf("error syncing '%s': %s", key, err.Error())
    }
    
    c.workqueue.Forget(obj)
    ...
    return nil
  }(obj)
  
  ...
  
  return true
}

func (c *Controller) syncHandler(key string) error {

  namespace, name, err := cache.SplitMetaNamespaceKey(key)
  ...
  
  network, err := c.networksLister.Networks(namespace).Get(name)
  if err != nil {
    if errors.IsNotFound(err) {
      glog.Warningf("Network does not exist in local cache: %s/%s, will delete it from Neutron ...",
      namespace, name)
      
      glog.Warningf("Network: %s/%s does not exist in local cache, will delete it from Neutron ...",
    namespace, name)
    
     // FIX ME: call Neutron API to delete this network by name.
     //
     // neutron.Delete(namespace, name)
     
     return nil
  }
    ...
    
    return err
  }
  
  glog.Infof("[Neutron] Try to process network: %#v ...", network)
  
  // FIX ME: Do diff().
  //
  // actualNetwork, exists := neutron.Get(namespace, name)
  //
  // if !exists {
  //   neutron.Create(namespace, name)
  // } else if !reflect.DeepEqual(actualNetwork, network) {
  //   neutron.Update(namespace, name)
  // }
  
  return nil
}
```

在执行周期里(`processNextWorkItem`), **首先**从工作队列里出队(`workqueue.Get`)了一个成员, 也就是一个`Key`.

**然后**, 在`syncHandler`方法中, 我使用这个`key`, 尝试从`Informer`维护的缓存中拿到了它所对应的`Network`对象. 使用了`networksLister`获取这个`Key`对应的`Network`对象(实际在访问本地缓存的索引). 如果拿不到对象, 则代表这个对象的`Key`是通过前面的"删除"事件添加进工作队列的. 尽管有这个`Key`, 但是对应的`Network`对象已经被删除了.

**如果能够获取到对应的 Network 对象, 就可以执行控制器模式里的对比"期望状态"和"实际状态"的逻辑了**.

自定义控制器拿到的这个对象, 是`APIServer`里保存的期望状态(已经被缓存到本地了).

"实际状态"来自于实际的集群. 控制循环通过`Neutron API`来查询实际的网络情况.

通过`Neutron`来查询`Network`对应对应的真实网络是否存在

* 不存在, “期望状态"和"实际状态"不一致. 需要调用`Neutron API`来创建真实的网络.
* 存在, 读取这个真实网络的信息, 然后检查是否信息一致, 是否需要更新.

**通过对比, "期望状态"和”实际状态"的差异, 完成一次调谐(`Reconcile`)的过程**

### 总结

`Informer`一个自带缓存和索引机制, 可以触发`Handler`的客户端库. 这个本地缓存在`Kubernetes`中一般被称为`Store`, 索引一般被称为`Index`.

`Informer`使用了`Reflector`包, 一个可以通过`ListAndWatch`机制获取并监视`API`对象变化的客户端封装.

`Reflector`和`Informer`之间, 用到了一个"增量先进先出队列"进行协同. 

`Informer`与编写的控制循环之间, 使用了一个工作队列来进行协同.

作为开发者, 你就只需要关注如何拿到"实际状态", 然后如何拿它去跟"期望状态"做对比, 从而决定接下来要做的业务逻辑即可.

以上内容, 就是 Kubernetes API 编程范式的核心思想.

## 扩展