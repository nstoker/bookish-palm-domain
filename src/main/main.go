package main

import (
	"net/http"
	"os"

	"github.com/nstoker/bookish-palm-domain/src/infrastructure"
	"github.com/nstoker/bookish-palm-domain/src/interfaces"
	"github.com/nstoker/bookish-palm-domain/src/usecases"
	"github.com/sirupsen/logrus"
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func main() {
	databaseName := "./production.sqlite"
	if !fileExists(databaseName) {
		logrus.Fatalf("Can't find '%s'.", databaseName)
	}
	dbHandler := infrastructure.NewSqliteHandler(databaseName)

	handlers := make(map[string]interfaces.DbHandler)
	handlers["DbUserRepo"] = dbHandler
	handlers["DbCustomerRepo"] = dbHandler
	handlers["DbItemRepo"] = dbHandler
	handlers["DbOrderRepo"] = dbHandler

	orderInteractor := new(usecases.OrderInteractor)
	orderInteractor.UserRepository = interfaces.NewDbUserRepo(handlers)
	orderInteractor.ItemRepository = interfaces.NewDbItemRepo(handlers)
	orderInteractor.OrderRepository = interfaces.NewDbOrderRepo(handlers)

	webserviceHandler := interfaces.WebserviceHandler{}
	webserviceHandler.OrderInteractor = orderInteractor

	http.HandleFunc("/orders", func(res http.ResponseWriter, req *http.Request) {
		webserviceHandler.ShowOrder(res, req)
	})

	logrus.Infof("Starting up, database: %s", databaseName)
	http.ListenAndServe(":8080", nil)
}
