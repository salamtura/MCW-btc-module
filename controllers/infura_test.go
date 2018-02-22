package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"MCW-btc-module/gocontracts"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"
)

func TestMakeCallRequestPayload(t *testing.T) {
	testAddress := common.HexToAddress("0x123")
	payload := makeCallRequestPayload(testAddress, common.Hex2Bytes("aa"))

	parameter := payload.Params[0].(methodParameter)

	assert.Equal(t, "2.0", payload.JsonRPC)
	assert.Equal(t, "eth_call", payload.Method)
	assert.Equal(t, testAddress, parameter.To)
	assert.Equal(t, "0xaa", parameter.Data)
}

func TestMakeEstimateGasRequestPayload(t *testing.T) {
	testAddress := common.HexToAddress("0x123")
	payload := makeEstimateGasRequestPayload(testAddress, common.Hex2Bytes("aa"))

	parameter := payload.Params[0].(methodParameter)

	assert.Equal(t, "2.0", payload.JsonRPC)
	assert.Equal(t, "eth_estimateGas", payload.Method)
	assert.Equal(t, testAddress, parameter.To)
	assert.Equal(t, "0xaa", parameter.Data)
}

func TestInfuraController(t *testing.T) {

	controller := MakeInfuraController("test", false)

	assert.Equal(t, mainnetInfuraEndpoint+"test", controller.infuraEndpoint)

	controller = MakeInfuraController("test", true)

	assert.Equal(t, testnetInfuraEndpoint+"test", controller.infuraEndpoint)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodPost,
		controller.infuraEndpoint,
		func(request *http.Request) (*http.Response, error) {
			defer request.Body.Close()

			requestBodyBytes, err := ioutil.ReadAll(request.Body)

			if err != nil {
				return httpmock.NewStringResponse(http.StatusInternalServerError, "failure"), nil
			}

			requestBody := new(requestPayload)

			if err := json.Unmarshal(requestBodyBytes, requestBody); err != nil {
				return httpmock.NewStringResponse(http.StatusInternalServerError, "failure"), nil
			}

			paramZero := requestBody.Params[0].(map[string]interface{})

			if requestBody.Method == "eth_call" {
				if len(paramZero["to"].(string)[2:]) != 40 || paramZero["data"].(string)[0:2] != "0x" {
					return httpmock.NewStringResponse(
						http.StatusOK,
						`{"jsonrpc": "2.0", "error": {"code": 100, "message": "invalid argument 0"}}`,
					), nil
				}
			}

			return httpmock.NewStringResponse(
				http.StatusOK,
				`{"jsonrpc": "2.0", "result": "0xabc"}`,
			), nil

		},
	)

	payload := makeCallRequestPayload(common.Address{}, []byte{})
	response, err := controller.callInfura(payload)

	if assert.NoError(t, err) {
		assert.Equal(t, "2.0", response.JsonRPC)
		assert.Equal(t, "0xabc", response.Result)
	}

	tokenContractAddress := common.HexToAddress("0xc780504e46526ce52e54b2f08c8c191fdd5743e8")

	tokenABI, err := abi.JSON(strings.NewReader(gocontracts.MocrowCoinABI))

	assert.NoError(t, err)

	param, err := controller.retrieveParameter("totalSupply", tokenABI, tokenContractAddress)

	if assert.NoError(t, err) {
		assert.Equal(t, common.HexToAddress("0xabc"), *param)
	}

	param, err = controller.retrieveParameter("totalSupply", abi.ABI{}, tokenContractAddress)

	if assert.Error(t, err) {
		assert.Nil(t, param)
	}

	param, err = controller.retrieveParameter("totalSupplies", tokenABI, tokenContractAddress)

	if assert.Error(t, err) {
		assert.Nil(t, param)
	}

	// TEST FOR UNLIKELY ERRORS

	httpmock.RegisterResponder(http.MethodPost,
		controller.infuraEndpoint,
		func(request *http.Request) (*http.Response, error) {
			return nil, errors.New("test_error")
		},
	)

	response, err = controller.callInfura(payload)

	if assert.Error(t, err) {
		assert.Nil(t, response)
	}

	param, err = controller.retrieveParameter("totalSupply", tokenABI, tokenContractAddress)

	if assert.Error(t, err) {
		assert.Nil(t, param)
	}

	httpmock.RegisterResponder(http.MethodPost,
		controller.infuraEndpoint,
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, "{C}"), nil
		},
	)

	response, err = controller.callInfura(payload)

	if assert.Error(t, err) {
		assert.Nil(t, response)
	}

	param, err = controller.retrieveParameter("totalSupply", tokenABI, tokenContractAddress)

	if assert.Error(t, err) {
		assert.Nil(t, param)
	}

	httpmock.RegisterResponder(http.MethodPost,
		controller.infuraEndpoint,
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				`{"jsonrpc": "2.0", "error": {"code": 100, "message": "invalid argument 0"}}`,
			), nil
		},
	)

	response, err = controller.callInfura(payload)

	if assert.Error(t, err) {
		assert.Nil(t, response)
	}

	param, err = controller.retrieveParameter("totalSupply", tokenABI, tokenContractAddress)

	if assert.Error(t, err) {
		assert.Nil(t, param)
	}

}
