<?php
//工厂方法模式

interface People
{
    public function say();
}

class Man implements People
{
    public function say()
    {
        echo 'this is a man';
    }
}

class Women implements People
{
    public function say()
    {
        echo 'this is a women';
    }
}


abstract class CreatePeople
{
    abstract public function create();
}

class FactoryMan extends CreatePeople
{
    public function create()
    {
        return new Man();
    }

}

class FactoryWomen extends CreatePeople
{
    public function create()
    {
        return new Women();
    }
}

class  Client
{
    // 具体生产对象并执行对象方法测试
    public function test() {
        $factory = new FactoryMan();
        $people = $factory->create();
        $people->say();

        $factory = new FactoryWomen();
        $people = $factory->create();
        $people->say();
    }
}

// 执行
$demo = new Client();
$demo->test();
