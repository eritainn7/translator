public class Program
{   
    public string Str()
    {
        string name = "";
        for (int i = 0; i < 10; i++)
        {
            name += i;
        }
        return "Alex";
    }

    public int fact(int x)
    {
        if (x < 2) return 1;
        return x * fact(x-1);
    }

    public static void main()
    {
        for (int i=0; i < 10; i++) Console.WriteLine(fact(i));
    }
}