package object;

public class C6 {

    int storage(String s) {
        return s.length() * 2;
    }

    public static void main(String args[]) {
        String s = "title";

        C6 c = new C6();

        System.out.println("storage is :" + c.storage(s));
    }
}
