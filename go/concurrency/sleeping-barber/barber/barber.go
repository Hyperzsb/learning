package barber

import "fmt"

type Barber struct {
	Id         int
	IsSleeping bool
}

func NewBarber(id int) Barber {
	return Barber{Id: id, IsSleeping: false}
}

func (b *Barber) Sleeping() {
	fmt.Printf("Barber %d is sleeping\n", b.Id)
	b.IsSleeping = true
}

func (b *Barber) WakenUp(customer int) {
	fmt.Printf("Barber %d is waken up by customer %d\n", b.Id, customer)
	b.IsSleeping = false
}

func (b *Barber) Cutting(customer int) {
	fmt.Printf("Barber %d is having customer %d hair cut\n", b.Id, customer)
}
