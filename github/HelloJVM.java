public class HelloJVM {
    public static void main(String[] args) {
        write0("2233");
        System.out.println("Hello World");
    }
    private static native void write0(String x);
}
