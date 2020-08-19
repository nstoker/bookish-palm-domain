package domain

import (
	"errors"
)

// CustomerRepository details of the customer
type CustomerRepository interface {
	Store(customer Customer)
	FindByID(ID int) Customer
}

// ItemRepository a repository of items
type ItemRepository interface {
	Store(item Item)
	FindByID(id int) Item
}

// OrderRepository an interface to an order repository
type OrderRepository interface {
	Store(order Order)
	FindByID(id int) Order
}

// Customer is a customer
type Customer struct {
	ID   int
	Name string
}

// Item is an item
type Item struct {
	ID        int
	Name      string
	Value     float64
	Available bool
}

// Order is an order
type Order struct {
	ID       int
	Customer Customer
	Items    []Item
}

// Add adds an item to an order
func (order *Order) Add(item Item) error {
	if !item.Available {
		return errors.New("cannot add unavailable items to order")
	}

	if order.value()*item.Value > 250.00 {
		return errors.New("An order may not exceed a total value of $250.00")
	}

	order.Items = append(order.Items, item)
	return nil
}

func (order *Order) value() float64 {
	sum := 0.0
	for i := range order.Items {
		sum = sum + order.Items[i].Value
	}

	return sum
}
