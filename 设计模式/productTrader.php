<?php

class Specification
{
    private $critera;

    public function setCritera(string $critera)
    {
        $this->critera = $critera;
    }

    public function getCritera(): string
    {
        return $this->critera;
    }

    public function isSatisfiedBy(Product $product) : bool
    {
        return $product->getCritera() == $this->getCritera();
    }
}

abstract class Product
{
    public abstract function getCritera() : string;
}

class ConcreteProductA extends Product
{
    public function getCritera() : string
    {
        return "SpecForConcreateProductA";
    }
}

class ConcreteProductB extends Product
{
    public function getCritera() : string
    {
        return "SpecForConcreateProductB";
    }
}

abstract class ProductCreator
{
    public abstract function create(Specification $spec) : Product;
}

class ConcreteProductCreator extends ProductCreator
{
    public function create(Specification $spec) : Product
    {
        if ($spec->getCritera() == "SpecForConcreateProductA") {
            return new ConcreteProductA();
        }
        else if ($spec->getCritera() == "SpecForConcreateProductB") {
            return new ConcreteProductB();
        }
    }
}

class ProductTrader
{
    private $dict;

    public function createFor(Specification $spec) : Product
    {
        $creator = $this->lookUpCreator($spec);
        $product = $creator->create($spec);

        return $product;
    }

    private function encryptObj($obj)
    {
        return md5(serialize($obj));
    }

    public function lookUpCreator(Specification $spec)
    {
        return $this->dict[$this->encryptObj($spec)];
    }

    public function addCreator(Specification $spec, ProductCreator $creator)
    {
        $this->dict[$this->encryptObj($spec)] = $creator;
    }

    public function removeCreatro(Specification $spec)
    {
        unset($this->dict[$this->encryptObj($spec)]);
    }
}

class Client
{
    public function testCase1()
    {
        $specOne = new Specification();
        $specOne->setCritera("SpecForConcreateProductA");

        $specTwo = new Specification();
        $specTwo->setCritera("SpecForConcreateProductB");

        $creator = new ConcreteProductCreator();

        $trader = new ProductTrader();
        $trader->addCreator($specOne, $creator);
        $trader->addCreator($specTwo, $creator);

        $specThree = new Specification();
        $specThree->setCritera("SpecForConcreateProductA");

        $product = $trader->createFor($specThree);
        var_dump($product);
    }
}

$client = new Client();
$client->testCase1();
