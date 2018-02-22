package controllers

import (
	"github.com/blockcypher/gobcy"
)

type BlockcypherController struct {
	client gobcy.API
}

func MakeBlockCypherController(APIToken string, isTestNet bool) BlockcypherController {
	net := "main"
	if isTestNet {
		net = "test3"
	}
	return BlockcypherController{
		client: gobcy.API{APIToken, "btc", net},
	}
}

func (controller BlockcypherController) GetConfirmedBalance(address string) (int, error) {
	balance, err := controller.client.GetAddr(address, map[string]string{"confirmations": "1"})
	if err != nil {
		return 0, err
	}

	return balance.FinalBalance, nil
}
