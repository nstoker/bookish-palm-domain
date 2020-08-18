package usecases

import (
	"fmt"

	"github.com/nstoker/bookish-palm-domain/src/domain"
)

// UserRepository interface
type UserRepository interface {
	Store(user User)
	FindByID(id int) User
}

// User structure
type User struct {
	ID       int
	IsAdmin  bool
	Customer domain.Customer
}

// Item structure
type Item struct {
	ID    int
	Name  string
	Value float64
}

// Logger interface
type Logger interface {
	Log(message string) error
}

// OrderInteractor structure
type OrderInteractor struct {
	UserRepository  UserRepository
	OrderRepository domain.OrderRepository
	ItemRepository  domain.ItemRepository
	Logger          Logger
}

// Items not sure what this does yet
func (interactor *OrderInteractor) Items(userID, orderID int) ([]Item, error) {
	var items []Item
	user := interactor.UserRepository.FindByID(userID)
	order := interactor.OrderRepository.FindByID(orderID)
	if user.Customer.ID != order.Customer.ID {
		message := "User #%i (customer #%i) is not allowed to see items in order ##%i (of customer #%i)"
		err := fmt.Errorf(message, user.ID, user.Customer.ID, order.ID, order.Customer.ID)
		interactor.Logger.Log(err.Error())
		items = make([]Item, 0)
		return items, err
	}
	items = make([]Item, len(order.Items))
	for i, item := range order.Items {
		items[i] = Item{item.ID, item.Name, item.Value}
	}

	return items, nil
}

// Add item to order
func (interactor *OrderInteractor) Add(userID, orderID, itemID int) error {
	var message string
	user := interactor.UserRepository.FindByID(userID)
	order := interactor.OrderRepository.FindByID(orderID)
	if user.Customer.ID != order.Customer.ID {
		message = "User #%i (customer #%i) is not allowed to add items to order #%i (of customer #%i)"
		err := fmt.Errorf(message,
			user.ID,
			user.Customer.ID,
			order.ID,
			order.Customer.ID)
		interactor.Logger.Log(err.Error())
		return err
	}

	item := interactor.ItemRepository.FindByID(itemID)
	if domainErr := order.Add(item); domainErr != nil {
		message = "Could not add item #%i to order #%i (of customer #%i) as user #%i because a business rule was violated '%s'"
		err := fmt.Errorf(message, item.ID, order.ID, order.Customer.ID, user.ID, domainErr.Error())
		interactor.Logger.Log(err.Error())
		return err
	}

	interactor.OrderRepository.Store(order)
	interactor.Logger.Log(fmt.Sprintf(
		"User added item '%s' (#%d) to order #%d",
		item.Name, item.ID, order.ID))
	return nil
}

// AdminOrderInteractor structure
type AdminOrderInteractor struct {
	OrderInteractor
}

// Add item
func (interactor *AdminOrderInteractor) Add(userID, orderID, itemID int) error {
	var message string
	user := interactor.UserRepository.FindByID(userID)
	order := interactor.OrderRepository.FindByID(orderID)
	if !user.IsAdmin {
		message = "User #%i (customer #%i) is not allowed to add items to order #%i (of customer #%i), because he is not an administrator"
		err := fmt.Errorf(message, user.ID, user.Customer.ID, order.ID, order.Customer.ID)
		interactor.Logger.Log(err.Error())
		return err
	}

	item := interactor.ItemRepository.FindByID(itemID)
	if domainErr := order.Add(item); domainErr != nil {
		message = "Could not add item #%i to order #%i (of customer #%i) as user#%i because a business rule was violated: '%s'"
		err := fmt.Errorf(message, item.ID, order.ID, order.Customer.ID, user.ID, domainErr)
		interactor.Logger.Log(err.Error())
		return err
	}
	interactor.OrderRepository.Store(order)
	interactor.Logger.Log(fmt.Sprintf(
		"Admin added item '%s' (#%d) to order #%d",
		item.Name,
		item.ID,
		order.ID,
	))

	return nil
}
