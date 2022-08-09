package object;

//04, -5
public class DataOnly {
    int i;
    double d;
    boolean b;

    public DataOnly()
    {
        i = 0;
        d = 0.0;
        b = false;
    }

    public static void main(String args[]) {
         DataOnly d = new DataOnly();
         d.i = 1;
         d.d = 2.1;
         d.b = true;

         System.out.println("i :" + d.i);
         System.out.println("d :" + d.d);
         System.out.println("b :" + d.b);
    }
}
