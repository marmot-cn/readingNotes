### self 和 static

---

* `self`: 调用的就是本身代码片段这个类.
* `static`: 调用的是从堆内存中提取出来,访问的是当前实例化的那个类(static就是调用的当前调用的类).

**self**

```php
root@bb0aa005831b:/var/www/html# cat test4.php
<?php

class Boo {

    protected static $str = "This is class Boo";

    public static function get_info(){

        echo get_called_class(), PHP_EOL;
        echo self::$str, PHP_EOL;
    }


}
class Foo extends Boo{

    protected static $str = "This is class Foo";

}


Foo::get_info();

输出:
Foo
This is class Boo
```

**static**

```php
<?php

class Boo {

    protected static $str = "This is class Boo";

    public static function get_info(){

        echo get_called_class(), PHP_EOL;
        echo static::$str, PHP_EOL;
    }


}
class Foo extends Boo{

    protected static $str = "This is class Foo";

}


Foo::get_info();

输出:
Foo
This is class Foo
```

