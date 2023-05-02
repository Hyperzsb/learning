package customer

import "fmt"

type Customer struct {
	Id int
}

func NewCustomer(id int) Customer {
	return Customer{Id: id}
}

func (c Customer) Coming() {
	fmt.Printf("Customer %d is coming\n", c.Id)
}

func (c Customer) Waiting() {
	fmt.Printf("Customer %d is waiting\n", c.Id)
}

func (c Customer) Leaving() {
	fmt.Printf("Customer %d is leaving\n", c.Id)
}
