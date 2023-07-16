package learning.controller;

import learning.bean.Car;
import learning.service.CarService;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

@Slf4j
@RestController
public class CarController {
    private CarService carService;

    @Autowired
    public void setCarService(CarService carService) {
        this.carService = carService;
    }

    @GetMapping("/car")
    public Car getCar() {
        log.info("GET /car");
        return carService.getCar();
    }
}
