package interfaces

import (
	"fmt"

	"github.com/nstoker/bookish-palm-domain/src/domain"
	"github.com/nstoker/bookish-palm-domain/src/usecases"
)

// DbHandler interface
type DbHandler interface {
	Execute(statement string)
	Query(statement string) Row
}

// Row Interface
type Row interface {
	Scan(dest ...interface{})
	Next() bool
}

// DbRepo structure
type DbRepo struct {
	dbHandlers map[string]DbHandler
	dbHandler  DbHandler
}

// DbUserRepo DbRepo
type DbUserRepo DbRepo

// DbCustomerRepo DbRepo
type DbCustomerRepo DbRepo

// DbOrderRepo DbRepo
type DbOrderRepo DbRepo

// DbItemRepo DbRepo
type DbItemRepo DbRepo

// NewDbUserRepo create new user db repo
func NewDbUserRepo(dbHandlers map[string]DbHandler) *DbUserRepo {
	dbUserRepo := new(DbUserRepo)
	dbUserRepo.dbHandlers = dbHandlers
	dbUserRepo.dbHandler = dbHandlers["DbUserRepo"]
	return dbUserRepo
}

// Store usecases user
func (repo *DbUserRepo) Store(user usecases.User) {
	isAdmin := "no"
	if user.IsAdmin {
		isAdmin = "yes"
	}

	repo.dbHandler.Execute(fmt.Sprintf(`INSERT INTO users (id, customer_id, is_admin)
										VALUES ('%d', '%d', '%v')`,
		user.ID, user.Customer.ID, isAdmin))
	customerRepo := NewDbCustomerRepo(repo.dbHandlers)
	customerRepo.Store(user.Customer)
}

// FindByID finds a user by id
func (repo *DbUserRepo) FindByID(id int) usecases.User {
	sql := fmt.Sprintf(
		`SELECT is_admin, customer_id
		FROM users
		WHERE id='%d' LIMIT 1`,
		id)
	row := repo.dbHandler.Query(sql)
	var isAdmin string
	var customerID int
	row.Next()
	row.Scan(&isAdmin, &customerID)
	customerRepo := NewDbCustomerRepo(repo.dbHandlers)
	u := usecases.User{ID: id, Customer: customerRepo.FindByID(customerID)}
	u.IsAdmin = isAdmin == "yes"
	return u
}

// NewDbCustomerRepo returns a new customer repo
func NewDbCustomerRepo(dbHandlers map[string]DbHandler) *DbCustomerRepo {
	dbCustomerRepo := new(DbCustomerRepo)
	dbCustomerRepo.dbHandlers = dbHandlers
	dbCustomerRepo.dbHandler = dbHandlers["DbCustomerRepo"]

	return dbCustomerRepo
}

// Store a record
func (repo *DbCustomerRepo) Store(customer domain.Customer) {
	repo.dbHandler.Execute(fmt.Sprintf(
		`INSERT INTO customers (id, name)
		VALUES('%d', '%v')`,
		customer.ID,
		customer.Name))
}

// FindByID find a record by id
func (repo *DbCustomerRepo) FindByID(id int) domain.Customer {
	row := repo.dbHandler.Query(fmt.Sprintf(`
	SELECT name
	FROM customers
	WHERE id = '%d' LIMIT 1`,
		id))

	var name string
	row.Next()
	row.Scan(&name)
	return domain.Customer{ID: id, Name: name}
}

// NewDbOrderRepo get a new db order repo
func NewDbOrderRepo(dbHandlers map[string]DbHandler) *DbOrderRepo {
	dbOrderRepo := new(DbOrderRepo)
	dbOrderRepo.dbHandlers = dbHandlers
	dbOrderRepo.dbHandler = dbHandlers["DbOrderRepo"]

	return dbOrderRepo
}

// Store stores a record
func (repo *DbOrderRepo) Store(order domain.Order) {
	repo.dbHandler.Execute(
		fmt.Sprintf(`
		INSERT INTO orders (id, customer_id)
		VALUES ('%d', '%v')`,
			order.ID, order.Customer.ID))

	for _, item := range order.Items {
		repo.dbHandler.Execute(
			fmt.Sprintf(`
			INSERT INTO items2orders (item_id, order_id)
			VALUES ('%d', '%d')`,
				item.ID, order.ID))
	}
}

// FindByID find record by ID
func (repo *DbOrderRepo) FindByID(id int) domain.Order {
	row := repo.dbHandler.Query(
		fmt.Sprintf(`
		SELECT customer_id
		FROM orders
		WHERE id = '%d' LIMIT 1`, id))

	var customerID int
	row.Next()
	row.Scan(&customerID)
	customerRepo := NewDbCustomerRepo((repo.dbHandlers))
	order := domain.Order{ID: id, Customer: customerRepo.FindByID(customerID)}
	var itemID int
	itemRepo := NewDbItemRepo(repo.dbHandlers)
	row = repo.dbHandler.Query(
		fmt.Sprintf(
			`SELECT item_id
			FROM items2orders
			WHERE order_id = '%d'`, order.ID))

	for row.Next() {
		row.Scan(&itemID)
		order.Add(itemRepo.FindByID(itemID))
	}

	return order
}

// NewDbItemRepo creates a new item repo
func NewDbItemRepo(dbHandlers map[string]DbHandler) *DbItemRepo {
	dbItemRepo := new(DbItemRepo)
	dbItemRepo.dbHandlers = dbHandlers
	dbItemRepo.dbHandler = dbHandlers["DbItemRepo"]

	return dbItemRepo
}

// Store stores an item
func (repo *DbItemRepo) Store(item domain.Item) {
	available := "no"
	if item.Available {
		available = "yes"
	}

	repo.dbHandler.Execute(
		fmt.Sprintf(
			`INSERT INTO items (id, name, value, available)
			VALUES ('%d', '%v', '%f', '%v')`,
			item.ID, item.Name, item.Value, available))
}

// FindByID finds item by id
func (repo *DbItemRepo) FindByID(id int) domain.Item {
	row := repo.dbHandler.Query(
		fmt.Sprintf(`
		SELECT name, value, available
		FROM items WHERE id = '%d' LIMIT 1`, id))

	var name string
	var value float64
	var available string
	row.Next()
	row.Scan(&name, &value, &available)
	item := domain.Item{ID: id, Name: name, Value: value}
	item.Available = available == "yes"

	return item
}
