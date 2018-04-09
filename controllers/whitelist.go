package controllers

import (
	"strings"

	"MCW-btc-module/gocontracts"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/ethereum/go-ethereum/common"
)

type WhitelistController struct {
	InfuraController
	crowdsaleAddress common.Address
	crowdsaleABI     abi.ABI
}

func MakeWhitelistController(
	infuraController InfuraController,
	crowdsaleAddress string,
) (*WhitelistController, error) {

	whitelistAbi, err := abi.JSON(strings.NewReader(gocontracts.MocrowCoinCrowdsaleABI))

	if err != nil {
		return nil, err
	}

	return &WhitelistController{
		infuraController,
		common.HexToAddress(crowdsaleAddress),
		whitelistAbi,
	}, nil
}

func (controller WhitelistController) IsWhitelisted(address common.Address) (bool, error) {
	packedData, err := controller.crowdsaleABI.Pack("isWhitelisted", address)

	if err != nil {
		return false, err
	}

	payload := makeCallRequestPayload(controller.crowdsaleAddress, packedData)

	response, err := controller.InfuraController.callInfura(payload)

	if err != nil {
		return false, err
	}

	return common.HexToAddress(response.Result).Big().Int64() == 1, nil
}
