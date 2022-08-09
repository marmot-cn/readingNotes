package object;

public class Incrementable {

    static void increment() {
        StaticTest.i++;
    }

    public static void main (String args[]) {
        Incrementable i = new Incrementable();
        i.increment();
        Incrementable.increment();
        increment();

        System.out.println("i is :" + StaticTest.i);
    }
}
