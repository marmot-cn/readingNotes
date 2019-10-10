#Composer

###简介
Composer 是 PHP 的一个依赖管理工具。它允许你申明项目所依赖的代码库，它会在你的项目中为你安装他们.

更多参考 [官方手册][id]
[id]:http://docs.phpcomposer.com/

###安装

本机安装命令

		curl -sS https://getcomposer.org/installer | php -- --install-dir=/usr/local/bin/ --filename=composer

###声明依赖关系

---
创建一个composer.json,其中描述了项目的依赖关系

	{
    	"require": {
        	"monolog/monolog": "1.0.*"
   	 	}
	}
require 需要一个 包名称 （例如 monolog/monolog） 映射到 包版本 （例如 1.0.*）

[composer.json架构][id]
[id]:http://docs.phpcomposer.com/04-schema.html

#####包名称
包名称由供应商名称和其项目名称构成。通常容易产生相同的项目名称，而供应商名称的存在则很好的解决了命名冲突的问题。它允许两个不同的人创建同样名为 json 的库，而之后它们将被命名为 igorw/json 和 seldaek/json.

phpunit/phpunit,供应商名称与项目名称相同.

#####包版本
1. 确切版本号`1.0.2`
2. 范围`>`,`>=`,`<`,`<=`,`!=`
3. 通配符`1.0.*` 等同于 `>=1.0`,`<1.1`
4. 赋值运算符  `~1.2` 指定最低版本，但允许版本号的最后一位数字上升.*`~1.2` 相当于 `>=1.2`,`<2.0`，而 `~1.2.3` 相当于 `>=1.2.3`,`<1.3`*

#####稳定性
默认情况下只有稳定的发行版才会被考虑在内。如果你也想获得 RC、beta、alpha 或 dev 版本，你可以使用 稳定标志。你可以对所有的包做 最小稳定性 设置，而不是每个依赖逐一设置。

###安装依赖包

---
获取定义的依赖到你的本地项目，只需要调用 composer.phar 运行 install 命令。

		composer install
		
		Loading composer repositories with package information
		Installing dependencies (including require-dev)
  		- Installing monolog/monolog (1.0.2)
    	  Downloading: 100%
   		  Downloading: 100%
    	  Downloading: 100%
    	- Installing monolog/monolog (1.0.2)
          Cloning b704c49a3051536f67f2d39f13568f74615b9922

第三方的代码到一个指定的目录`vendor`. install 命令将创建一个 composer.lock 文件到你项目的根目录中。

**composer.lock - 锁文件**

Composer 将把安装时确切的版本号列表写入 composer.lock 文件。这将锁定改项目的特定版本.

请提交你应用程序的 composer.lock(包括 composer.json)到你的版本库中

如果不存在 composer.lock 文件,Composer 将读取 composer.json 并创建锁文件.

任何人建立项目都将下载与指定版本完全相同的依赖.即使从那时起你的依赖已经发布了许多新的版本.*在json文件中根据包版本的范围指定,可能会有新的版本*


**Packagist**

一个公共的资源库
[packagist website][id]
[id]:http://packagist.org

**自动加载**

Composer 生成了一个 vendor/autoload.php 文件。你可以简单的引入这个文件，你会得到一个免费的自动加载支持。

		require 'vendor/autoload.php';
		
**使用monolog**
		
		require 'vendor/autoload.php';
		$log = new Monolog\Logger('name');
		$log->pushHandler(new Monolog\Handler\StreamHandler('app.log', Monolog\Logger::WARNING));

		$log->addWarning('Foo');
		
**生成app.log**
		
		[2015-05-17 12:18:29] name.WARNING: Foo [] []


		

 