#PHP - IOC/DI

---

####IOC

假设应用程序有储存需求,若直接在高层的应用程序中调用低层模块API,导致应用程序对低层模块产生依赖.

		/**
		 * 高层
		 */
		class Business
		{
		    private $writer;
		
		    public function __construct()
		    {
		        $this->writer = new FloppyWriter();
		    }
		
		    public function save()
		    {
		        $this->writer->saveToFloppy();
		    }
		}
		
		/**
		 * 低层，软盘存储
		 */
		class FloppyWriter
		{
		    public function saveToFloppy()
		    {
		        echo __METHOD__;
		    }
		}
		
		$biz = new Business();
		$biz->save(); // FloppyWriter::saveToFloppy
		
假设程序要移植到另一个平台,而该平台使用USB磁盘作为存储介质,则这个程序无法直接重用,必须加以修改才行.本例由于低层变化导致高层也跟着变化,不好的设计.

**关键点**

* 控制反转 Inversion of Control
* 依赖关系的转移
* 依赖抽象而非实践

程序不应该依赖于具体的实现,而是要依赖抽像的接口.请看代码演示:

		/**
		 * 接口
		 */
		interface IDeviceWriter
		{
		    public function saveToDevice();
		}
		
		/**
		 * 高层
		 */
		class Business
		{
		    /**
		     * @var IDeviceWriter
		     */
		    private $writer;
		
		    /**
		     * @param IDeviceWriter $writer
		     */
		    public function setWriter($writer)
		    {
		        $this->writer = $writer;
		    }
		
		    public function save()
		    {
		        $this->writer->saveToDevice();
		    }
		}
		
		/**
		 * 低层，软盘存储
		 */
		class FloppyWriter implements IDeviceWriter
		{
		
		    public function saveToDevice()
		    {
		        echo __METHOD__;
		    }
		}
		
		/**
		 * 低层，USB盘存储
		 */
		class UsbDiskWriter implements IDeviceWriter
		{
		
		    public function saveToDevice()
		    {
		        echo __METHOD__;
		    }
		}
		
		$biz = new Business();
		$biz->setWriter(new UsbDiskWriter());
		$biz->save(); // UsbDiskWriter::saveToDevice
		
		$biz->setWriter(new FloppyWriter());
		$biz->save(); // FloppyWriter::saveToDevice
		
控制权从实际的`FloppyWriter`转移到了抽象的`IDeviceWriter`接口上,让`Business`依赖于`IDeviceWriter`接口,且`FloppyWriter`、`UsbDiskWriter`也依赖于`IDeviceWriter`接口.

**IOC(控制反转)和DI(依赖注入)**

这就是`IoC`,面对`变化`,`高层不用修改一行代码`,不再依赖低层,而是`依赖注入`,这就引出了`DI`.

**实用的注入方式**

* `Setter injection` 使用setter方法,上面的代码演示的是Setter方式的注入.
* `Constructor injection` 实用构造函数.
* `Property Injection` 直接设置属性.

####依赖注入容器 Dependency Injection Container

**作用**

* 管理应用程序中的`全局`对象(包括实例化、处理依赖关系).
* 可以掩饰加载对象(仅用到时才创建对象).
		
		比如我自己的框架里面的mongo连接初始化,可以在用到的时候在初始化,而不是以开始就全局初始化.
		
* 促进编写可重用,可测试和松耦合的代码.

		我现在的框架里面测试代码覆盖率不够,只能覆盖到operation操作.而且还需要考虑怎么把数据层解耦测试.


**引发的问题**

如果这个组件有很多依赖,我们需要创建多个参数的`setter`方法来传递依赖关系,或者建立一个多个参数的构造函数来传递它们,另外在使用组件前还要每次都创建依赖,这让我们的代码像这样不易维护

		/创建依赖实例或从注册表中查找
		$connection = new Connection();
		$session = new Session();
		$fileSystem = new FileSystem();
		$filter = new Filter();
		$selector = new Selector();
		
		//把实例作为参数传递给构造函数
		$some = new SomeComponent($connection, $session, $fileSystem, $filter, $selector);
		
		// ... 或者使用setter
		
		$some->setConnection($connection);
		$some->setSession($session);
		$some->setFileSystem($fileSystem);
		$some->setFilter($filter);
		$some->setSelector($selector);
		
这也是我个人现在维护框架遇见的问题.

假设我们必须在应用的不同地方使用和创建这些对象.如果当你永远不需要任何依赖实例时,你需要去删掉构造函数的参数,或者去删掉注入的setter.为了解决这样的问题,我们再次回到全局注册表创建组件.不管怎么样,在创建对象之前,它增加了一个新的抽象层:

		class SomeComponent
		{
		
		    // ...
		
		    /**
		     * Define a factory method to create SomeComponent instances injecting its dependencies
		     */
		    public static function factory()
		    {
		
		        $connection = new Connection();
		        $session = new Session();
		        $fileSystem = new FileSystem();
		        $filter = new Filter();
		        $selector = new Selector();
		
		        return new self($connection, $session, $fileSystem, $filter, $selector);
		    }
		
		}
		
这里是用工厂函数封装了创建对象注入各种依赖示例的复杂性,原来调用的时候需要初始化一次注入,现在统一封装后只需要调用工厂函数即可.

**一个优雅的解决方案**

为依赖实例提供一个`容器`.这个容器担任`全局的注册表`,就像我们刚才看到的那样.使用依赖实例的容器作为一个桥梁来获取依赖实例,使我们能够降低我们的组件的复杂性:

		class SomeComponent
		{
		
		    protected $_di;
		
		    public function __construct($di)
		    {
		        $this->_di = $di;
		    }
		
		    public function someDbTask()
		    {
		
		        // 获得数据库连接实例
		        // 总是返回一个新的连接
		        $connection = $this->_di->get('db');
		
		    }
		
		    public function someOtherDbTask()
		    {
		
		        // 获得共享连接实例
		        // 每次请求都返回相同的连接实例
		        $connection = $this->_di->getShared('db');
		
		        // 这个方法也需要一个输入过滤的依赖服务
		        $filter = $this->_di->get('filter');
		
		    }
		
		}
		
		$di = new Phalcon\DI();
		
		//在容器中注册一个db服务
		$di->set('db', function() {
		    return new Connection(array(
		        "host" => "localhost",
		        "username" => "root",
		        "password" => "secret",
		        "dbname" => "invo"
		    ));
		});
		
		//在容器中注册一个filter服务
		$di->set('filter', function() {
		    return new Filter();
		});
		
		//在容器中注册一个session服务
		$di->set('session', function() {
		    return new Session();
		});
		
		//把传递服务的容器作为唯一参数传递给组件
		$some = new SomeComponent($di);
		
		$some->someTask();
				


