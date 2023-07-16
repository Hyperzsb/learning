package learning.service.impl;

import learning.bean.Car;
import learning.service.CarService;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

@Service
public class CarServiceImpl implements CarService {
    @Value("${car.manufacturer}")
    private String manufacturer;

    public Car getCar() {
        return new Car(manufacturer, "Type A", "Car A", 10000);
    }
}
