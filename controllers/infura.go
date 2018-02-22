package controllers

import (
	"encoding/json"
	"errors"

	"MCW-btc-module/helpers"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const (
	mainnetInfuraEndpoint = "https://mainnet.infura.io/"
	testnetInfuraEndpoint = "https://rinkeby.infura.io/"
)

type InfuraController struct {
	infuraEndpoint string
}

func MakeInfuraController(accessToken string, isTestnet bool) InfuraController {

	endpoint := mainnetInfuraEndpoint

	if isTestnet {
		endpoint = testnetInfuraEndpoint
	}
	endpoint += accessToken
	return InfuraController{
		endpoint,
	}
}

type (
	infuraResponse struct {
		ID      int64        `json:"id"`
		JsonRPC string       `json:"jsonrpc"`
		Result  string       `json:"result"`
		Error   *infuraError `json:"error"`
	}

	infuraError struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	requestPayload struct {
		JsonRPC string        `json:"jsonrpc"`
		ID      int           `json:"id"`
		Method  string        `json:"method"`
		Params  []interface{} `json:"params"`
	}

	methodParameter struct {
		To   common.Address `json:"to"`
		Data string         `json:"data"`
	}
)

func makeRequestPayload(methodName string, params []interface{}) requestPayload {
	return requestPayload{
		JsonRPC: "2.0",
		Method:  methodName,
		Params:  params,
	}
}

func makeCallRequestPayload(targetContract common.Address, data []byte) requestPayload {
	return makeRequestPayload(
		"eth_call",
		[]interface{}{
			methodParameter{
				To:   targetContract,
				Data: hexutil.Bytes(data).String(),
			},
			"latest",
		},
	)
}

func makeEstimateGasRequestPayload(targetContract common.Address, data []byte) requestPayload {
	return makeRequestPayload(
		"eth_estimateGas",
		[]interface{}{
			methodParameter{
				To:   targetContract,
				Data: hexutil.Bytes(data).String(),
			},
		},
	)
}

func (controller InfuraController) callInfura(payload requestPayload) (*infuraResponse, error) {

	marshaledPayload, err := json.Marshal(payload)

	if err != nil {
		return nil, err
	}

	_, responseBody, err := helpers.Post(controller.infuraEndpoint, helpers.Headers{"Content-Type": "application/json"}, marshaledPayload)

	if err != nil {
		return nil, err
	}

	response := new(infuraResponse)

	if err := json.Unmarshal(responseBody, response); err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, errors.New(response.Error.Message)
	}

	return response, nil
}
func (controller InfuraController) retrieveParameter(name string, abi abi.ABI, targetContract common.Address) (*common.Address, error) {
	packedData, err := abi.Pack(name)

	if err != nil {
		return nil, err
	}
	payload := makeCallRequestPayload(targetContract, packedData)

	response, err := controller.callInfura(payload)

	if err != nil {
		return nil, err
	}

	resultAddress := common.HexToAddress(response.Result)

	return &resultAddress, nil
}
