package interfaces

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/nstoker/bookish-palm-domain/src/usecases"
	"github.com/sirupsen/logrus"
)

// OrderInteractor interface
type OrderInteractor interface {
	Items(userID, orderID int) ([]usecases.Item, error)
	Add(userID, orderID, itemID int) error
}

// WebserviceHandler structure
type WebserviceHandler struct {
	OrderInteractor OrderInteractor
}

// ShowOrder show an order
func (handler WebserviceHandler) ShowOrder(res http.ResponseWriter, req *http.Request) {
	userID, _ := strconv.Atoi(req.FormValue("userId"))
	logrus.Infof("UserID: %d", userID)
	orderID, _ := strconv.Atoi(req.FormValue("orderId"))
	logrus.Infof("OrderID: %d", orderID)
	logrus.Infof("Handler %+v", handler.OrderInteractor)
	logrus.Infof("--> %+v", handler.OrderInteractor)
	logrus.Infof("  --> %+v", handler.OrderInteractor.Items)
	items, err := handler.OrderInteractor.Items(userID, orderID)
	if err != nil {
		logrus.Fatalf(err.Error())
	}
	logrus.Infof("items %+v", items)
	for _, item := range items {
		io.WriteString(res, fmt.Sprintf("item id: %d\n", item.ID))
		io.WriteString(res, fmt.Sprintf("item name: %s\n", item.Name))
		io.WriteString(res, fmt.Sprintf("item value: %f\n", item.Value))

	}
}
