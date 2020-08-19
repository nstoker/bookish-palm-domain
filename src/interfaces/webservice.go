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
	orderID, _ := strconv.Atoi(req.FormValue("orderId"))
	items, err := handler.OrderInteractor.Items(userID, orderID)
	if err != nil {
		logrus.Warningf(err.Error())
	}
	for _, item := range items {
		io.WriteString(res, fmt.Sprintf("item id: %d\n", item.ID))
		io.WriteString(res, fmt.Sprintf("item name: %s\n", item.Name))
		io.WriteString(res, fmt.Sprintf("item value: %f\n", item.Value))

	}
}
