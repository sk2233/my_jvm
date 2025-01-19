public class GetClassTest {

    public static void main(String[] args) {
        System.out.println(Object.class.getName()); // java.lang.Object
        System.out.println(int[].class.getName()); // [I
        System.out.println(int[][].class.getName()); // [[I
        System.out.println(Object[].class.getName()); // [Ljava.lang.Object;
        System.out.println(Object[][].class.getName()); // [[Ljava.lang.Object;
        System.out.println(Runnable.class.getName()); // java.lang.Runnable
        System.out.println("abc".getClass().getName()); // java.lang.String
    }

}