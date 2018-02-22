package routing

import (
	"MCW-btc-module/controllers"

	"github.com/ethereum/go-ethereum/common"
	"github.com/labstack/echo"

	"net/http"
)

type ExchangeRouter struct {
	*controllers.ExchangeController
}

func MakeExchangeRouter(controller *controllers.ExchangeController) ExchangeRouter {
	return ExchangeRouter{
		ExchangeController: controller,
	}
}

func (router ExchangeRouter) Register(group *echo.Group) {
	group.GET("/:address", router.buyTokens)
}

func (router ExchangeRouter) buyTokens(context echo.Context) error {
	btcAddressChan := make(chan string)
	errChan := make(chan error)

	go router.ExchangeController.BuyTokens(common.HexToAddress(context.Param("address")), btcAddressChan, errChan)

	address := <-btcAddressChan
	err := <-errChan

	if err != nil {
		return err
	}

	return context.String(http.StatusOK, address)
}
