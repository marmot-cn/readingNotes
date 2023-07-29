<?php 

interface TV{
  public function open();
  public function watch();
}

interface PC{
  public function work();
  public function play();
}

class HaierTv implements TV
{
  public function open()
  {
      echo "Open Haier TV <br>";
  }

  public function watch()
  {
      echo "I'm watching TV <br>";
  }
}

class HaierPc implements PC
{
  public function work()
  {
      echo "I'm working on a Haier computer <br>";
  }
  public function play()
  {
      echo "Haier computers can be used to play games <br>";
  }
}


class LenovoTv implements TV
{
  public function open()
  {
      echo "Open Lenovo TV <br>";
  }

  public function watch()
  {
      echo "I'm watching TV <br>";
  }
}

class LenovoPc implements PC
{
  public function work()
  {
      echo "I'm working on a Lenovo computer <br>";
  }
  public function play()
  {
      echo "Lenovo computers can be used to play games <br>";
  }
}

interface Factory{
  public function createPc();
  public function createTv();
}

class LenovoFactory implements Factory
{
  public static function createTv()
  {
      return new LenovoTv();
  }
  public static function createPc()
  {
      return new LenovoPc();
  }
}

class HaierFactory implements Factory
{
  public static function createTv()
  {
      return new HaierTv();
  }
  public static function createPc()
  {
      return new HaierPc();
  }
}

$factory = new LenovoFactory();
$tv = factory->createTv()
$tv->open();
$tv->watch();
$pc = factory->createPc()
$pc->work();
$pc->play();

$factory = new HaierFactory();
$tv = factory->createTv()
$tv->open();
$tv->watch();
$pc = factory->createPc()
$pc->work();
$pc->play();
