package learning.bean;

import lombok.AllArgsConstructor;
import lombok.Data;

@Data
@AllArgsConstructor
public class Car {
    private String manufacturer;
    private String type;
    private String name;
    private Integer price;
}
