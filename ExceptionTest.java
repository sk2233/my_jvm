public class ExceptionTest {

    public static void main(String[] args) {
        test(0);
        test(1);
        test(2);
        test(3);
    }

    private static void test(int x) {
        try {
            if (x == 0) {
                throw new IllegalArgumentException("0!");
            }
            if (x == 1) {
                throw new RuntimeException("1!");
            }
            if (x == 2) {
                throw new Exception("2!");
            }
        } catch (IllegalArgumentException e) {
            System.out.println("IllegalArgumentException");
        } catch (RuntimeException e) {
            System.out.println("RuntimeException");
        } catch (Exception e) {
            System.out.println("Exception");
        } finally { // 会为每个  catch 独立生成 finally 代码进行执行 而不是使用 goto 跳转
            System.out.println(x);
        }
    }

}
