package object;

public class C8 {

    public static void main(String args[]) {
        StaticTest a = new StaticTest();
        StaticTest b = new StaticTest();

        System.out.println("i in a is :"+a.i);
        System.out.println("i in b is :"+b.i);
        a.i++;
        System.out.println("i in b is :"+b.i);
    }
}
