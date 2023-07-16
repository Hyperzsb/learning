package learning.controller;

import learning.bean.Car;
import learning.service.CarService;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

import javax.annotation.Resource;

@RestController
public class CarController {
    @Resource
    private CarService carService;

    @GetMapping("/car")
    public Car getCar() {
        return carService.getCar();
    }
}
