package object;

//11
public class AllTheColorOfTheRainbow {
    int anIntegerRepresentingColors;

    void changeTheHueOfTheColor(int newHue) {
        anIntegerRepresentingColors = newHue;
    }

    public static void main(String args[])
    {
        AllTheColorOfTheRainbow a = new AllTheColorOfTheRainbow();

        a.changeTheHueOfTheColor(10);

        System.out.println("anIntegerRepresentingColors is "+a.anIntegerRepresentingColors);
    }
}
