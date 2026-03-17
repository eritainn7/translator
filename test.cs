public class Car {
    private string name;
    private int born;

    public Car(string name, string born) {
        this.name = name;
        this.born = born;
    }

    public string getName() {
        return name;
    }

    public int getBorn() {
        return born;
    }
}

public class Program {
    public static void main() {
        Car[] cars = {
            new Car("Bugati", 2008),
            new Car("BMV", 2012),
            new Car("Hundai", 2013)
        }

        for (int i = 0; i < cars.Length; i++) {
            cars[i].getName();
        }
    }
}