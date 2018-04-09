package routing

import (
	"errors"

	"MCW-btc-module/controllers"

	"github.com/ethereum/go-ethereum/common"
	"github.com/labstack/echo"

	"net/http"
)

type ExchangeRouter struct {
	*controllers.ExchangeController
	*controllers.WhitelistController
}

func MakeExchangeRouter(exchangeController *controllers.ExchangeController, whitelistController *controllers.WhitelistController) ExchangeRouter {
	return ExchangeRouter{
		ExchangeController:  exchangeController,
		WhitelistController: whitelistController,
	}
}

func (router ExchangeRouter) Register(group *echo.Group) {
	group.GET("/:address", router.buyTokens)
}

func (router ExchangeRouter) buyTokens(context echo.Context) error {
	address := common.HexToAddress(context.Param("address"))

	if address == common.HexToAddress("") {
		return errors.New("invalid address")
	}

	whitelisted, err := router.IsWhitelisted(address)

	if err != nil {
		return err
	}

	if !whitelisted {
		return errors.New("address is not whitelisted")
	}

	transaction, isNew, err := router.ExchangeController.CreateTransactionEntry(address.String())
	if err != nil {
		return err
	}

	if isNew {
		go router.ExchangeController.BuyTokens(transaction)

		if err != nil {
			return err
		}
	}

	return context.String(http.StatusOK, transaction.BitcoinAddress)
}

