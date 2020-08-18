package main

import (
	"net/http"

	"github.com/nstoker/bookish-palm-domain/src/infrastructure"
	"github.com/nstoker/bookish-palm-domain/src/interfaces"
	"github.com/nstoker/bookish-palm-domain/src/usecases"
	"github.com/sirupsen/logrus"
)

func main() {
	dbHandler := infrastructure.NewSqliteHandler("./production.sqlite")

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

	logrus.Infof("Starting up")
	http.ListenAndServe(":8080", nil)
}
