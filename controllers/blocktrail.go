package controllers

import (
	"encoding/json"
	"strings"

	"MCW-btc-module/helpers"
)

type BlocktrailController struct {
	endpoint string
}

const blocktrailTestnetEndpoint = "https://api.blocktrail.com/v1/tbtc/address/%address%/unspent-outputs?api_key="
const blocktrailMainnetEndpoint = "https://api.blocktrail.com/v1/btc/address/%address%/unspent-outputs?api_key="

type (
	BlocktrailUnspentInput struct {
		Value         int `json:"value"`
		Confirmations int `json:"confirmations"`
	}
	blocktrailResponse struct {
		Data []BlocktrailUnspentInput `json:"data"`
	}
)

func MakeBlocktrailController(APIKey string, isTestnet bool) BlocktrailController {
	endpoint := blocktrailMainnetEndpoint
	if isTestnet {
		endpoint = blocktrailTestnetEndpoint
	}
	return BlocktrailController{
		endpoint: endpoint + APIKey,
	}
}

func (controller BlocktrailController) GetConfirmedBalance(address string) (int, error) {
	_, response, err := helpers.Get(strings.Replace(controller.endpoint, "%address%", address, 1), helpers.Headers{})

	if err != nil {
		return 0, err
	}

	responseJSON := new(blocktrailResponse)

	if err := json.Unmarshal(response, responseJSON); err != nil {
		return 0, err
	}

	value := 0

	for _, input := range responseJSON.Data {
		if input.Confirmations >= 6 {
			value += input.Value
		}
	}

	return value, nil
}
