package object;

import java.util.*;

/** The first Thinking in Java example program
 * Displays a string and today's date
 * @aythor Bruce Eckel
 * @author www.MindView.net
 * @version 4.0
 */
public class C12HelloDate {

    /** Entry point to class & application
     *
     * @param args array of string arguments
     * @throws exceptions No exceptions thrown
     */
    public static void main(String args[]) {
        System.out.println("Hello, it's: ");
        System.out.println(new Date());
    }
}
