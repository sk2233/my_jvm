public class HelloWorld {

    public static void main(String[] args) {
//         for (String arg : args) {
//             System.out.println(arg);
//         }
        max(22,33);
    }

    public static native int max(int a, int b);

}