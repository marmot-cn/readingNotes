<?php
//简单工厂方法

interface People
{
    public function say();

}

class Man implements People
{
    public function say()
    {
        echo 'this is a man ';
    }
}

class Women implements People
{
    public function say()
    {
        echo 'this is a women';
    }
}

class SimpleFactory
{
    public static function create($name)
    {
    if ($name == 'man') {
            return new Man();
        } elseif ($name == 'women') {
            return new Women();
        }
    }

}

//具体调用
$man = SimpleFactory::create('man');
$man->say();
$women = SimpleFactory::create('women');
$women->say();
