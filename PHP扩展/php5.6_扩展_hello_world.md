#php5.6 扩展 加密解密

---

因为需要在`php7`实现加密和解密,需要写成php扩展格式.先在`php5.6`实现.

		cd php_source_code/php-5.6.15/ext
		./ext_skel --extname=marmot
		
		ext ls marmot
		CREDITS      EXPERIMENTAL config.m4    config.w32   marmot.c     marmot.php   php_marmot.h tests


**进行编译测试**

		ext pwd
		/Users/chloroplast1983/Sites/php_source_code/php-5.6.15/ext
		cd marmot
		sudo phpize
		Configuring for:
		PHP Api Version:         20131106
		Zend Module Api No:      20131226
		Zend Extension Api No:   220131226
		./configure
		...
		make && make install
		...
		Installing shared extensions:     /usr/local/Cellar/php56/5.6.14/lib/php/extensions/no-debug-non-zts-20131226/
		
		
		ls /usr/local/Cellar/php56/5.6.14/lib/php/extensions/no-debug-non-zts-20131226
		php --ini(查看扩展路径)
		
		sudo vi /usr/local/etc/php/5.6/conf.d/ext-marmot.ini
		[marmot]
		extension_dir="/usr/local/Cellar/php56/5.6.14/lib/php/extensions/no-debug-non-zts-20131226/"
		extension = "marmot.so"
		
		重启php
		sudo apachectl restart
		
		php -d enable_dl=On marmot.php
		Functions available in the test extension:
		confirm_marmot_compiled
		hello_world
		
		Congratulations! You have successfully modified ext/marmot/config.m4. Module marmot is now compiled into PHP.
		
		
上面的扩展编译路径具体在`phpinfo`中`extension_dir`查看.

		php -i | grep "extension_dir"
		extension_dir => /usr/local/Cellar/php56/5.6.14/lib/php/extensions/no-debug-non-zts-20131226 => /usr/local/Cellar/php56/5.6.14/lib/php/extensions/no-debug-non-zts-20131226

**修改config.m4**

我们需要动态编译,去掉前面的`dnl`.`dnl`是注释.

		PHP_ARG_WITH(marmot, for marmot support,
		dnl Make sure that the comment is aligned:
		[  --with-marmot             Include marmot support])
		
**修改`php_marmot.h`**


