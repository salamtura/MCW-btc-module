package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/big"
	"net/http"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"
)

func TestMakeTokenManagementController(t *testing.T) {
	controller, err := MakeTokenManagementController(InfuraController{"endpoint"}, "", "", "9df9993fcb4d9f4520770ce69f1623bccf3690489205c11e42a78bddc6526123")

	if assert.NoError(t, err) {
		assert.Equal(t, "endpoint", controller.infuraEndpoint)
	}
}

func TestMakeTokenManagementControllerError(t *testing.T) {
	controller, err := MakeTokenManagementController(InfuraController{"endpoint"}, "", "", "9df9cb4d9f4520770ce69f1623bccf3690489205c11e42a78bddc6526123")

	assert.Error(t, err)
	assert.Nil(t, controller)
}

func TestTokenManagementController(t *testing.T) {
	controller, err := MakeTokenManagementController(InfuraController{"endpoint"}, "", "", "9df9993fcb4d9f4520770ce69f1623bccf3690489205c11e42a78bddc6526123")

	if assert.NoError(t, err) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder(
			http.MethodPost,
			controller.infuraEndpoint,
			func(request *http.Request) (*http.Response, error) {
				return httpmock.NewStringResponse(
					http.StatusOK,
					`{"jsonrpc": "2.0", "result": "0xabc"}`,
				), nil

			},
		)

		address, err := controller.getCrowdaleParameter("tokensRemainingIco")

		if assert.NoError(t, err) {
			assert.Equal(t, common.HexToAddress("0xabc"), *address)
		}

		boolValue, err := controller.isPreICO()

		if assert.NoError(t, err) {
			assert.False(t, boolValue)
		}

		boolValue, err = controller.isICO()

		if assert.NoError(t, err) {
			assert.False(t, boolValue)
		}

		bigValue, err := controller.preICOExchangeRate()

		if assert.NoError(t, err) {
			assert.Equal(t, common.HexToAddress("0xabc").Big(), bigValue)
		}

		bigValue, err = controller.preICOTokensRemaining()

		if assert.NoError(t, err) {
			assert.Equal(t, common.HexToAddress("0xabc").Big(), bigValue)
		}

		bigValue, err = controller.icoExchangeRate()

		if assert.NoError(t, err) {
			assert.Equal(t, common.HexToAddress("0xabc").Big(), bigValue)
		}

		bigValue, err = controller.icoTokensRemaining()

		if assert.NoError(t, err) {
			assert.Equal(t, common.HexToAddress("0xabc").Big(), bigValue)
		}

		bigValue, err = controller.GetTokensLeft()

		if assert.NoError(t, err) {
			assert.Equal(t, big.NewInt(0), bigValue)
		}

		bigFloatValue, err := controller.GetTokenExchangeRate()

		if assert.NoError(t, err) {
			assert.Equal(t, big.NewFloat(0), bigFloatValue)
		}

		assert.Equal(t, errors.New("it's neither pre-ico, nor ico"), controller.MintTokens(common.Address{}, big.NewInt(0)))
	}
}

func TestTokenManagementController_Fail(t *testing.T) {
	controller, err := MakeTokenManagementController(InfuraController{"endpoint"}, "", "", "9df9993fcb4d9f4520770ce69f1623bccf3690489205c11e42a78bddc6526123")

	if assert.NoError(t, err) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder(
			http.MethodPost,
			controller.infuraEndpoint,
			func(request *http.Request) (*http.Response, error) {

				return nil, errors.New("mock infura failure")

			},
		)

		address, err := controller.getCrowdaleParameter("tokensRemainingIco")

		assert.Error(t, err)
		assert.Nil(t, address)

		boolValue, err := controller.isPreICO()

		assert.Error(t, err)
		assert.False(t, boolValue)

		boolValue, err = controller.isICO()

		assert.Error(t, err)
		assert.False(t, boolValue)

		bigValue, err := controller.preICOExchangeRate()

		assert.Error(t, err)
		assert.Nil(t, bigValue)

		bigValue, err = controller.preICOTokensRemaining()

		assert.Error(t, err)
		assert.Nil(t, bigValue)

		bigValue, err = controller.icoExchangeRate()

		assert.Error(t, err)
		assert.Nil(t, bigValue)

		bigValue, err = controller.icoTokensRemaining()

		assert.Error(t, err)
		assert.Nil(t, bigValue)

		bigValue, err = controller.GetTokensLeft()

		assert.Error(t, err)
		assert.Nil(t, bigValue)

		bigFloatValue, err := controller.GetTokenExchangeRate()

		assert.Error(t, err)
		assert.Nil(t, bigFloatValue)

		assert.Error(t, controller.MintTokens(common.Address{}, big.NewInt(0)))
	}
}

func TestControllerPreICOStage(t *testing.T) {
	controller, err := MakeTokenManagementController(InfuraController{"endpoint"}, "", "", "9df9993fcb4d9f4520770ce69f1623bccf3690489205c11e42a78bddc6526123")

	if assert.NoError(t, err) {
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

				data := requestBody.Params[0].(map[string]interface{})["data"].(string)

				testData, err := controller.crowdsaleContractABI.Pack("isPreIcoStage")

				if assert.NoError(t, err) {
					if data == hexutil.Bytes(testData).String() {
						return httpmock.NewStringResponse(
							http.StatusOK,
							`{"jsonrpc": "2.0", "result": "`+common.BigToAddress(big.NewInt(1)).String()+`"}`,
						), nil
					}
				}

				return httpmock.NewStringResponse(
					http.StatusOK,
					`{"jsonrpc": "2.0", "result": "0xabc"}`,
				), nil

			},
		)

		boolValue, err := controller.isPreICO()

		if assert.NoError(t, err) {
			assert.True(t, boolValue)
		}

		bigValue, err := controller.GetTokensLeft()

		if assert.NoError(t, err) {
			assert.Equal(t, common.HexToAddress("0xabc").Big(), bigValue)
		}

		bigFloatValue, err := controller.GetTokenExchangeRate()

		if assert.NoError(t, err) {
			assert.Equal(t, big.NewFloat(0).SetInt(common.HexToAddress("0xabc").Big()), bigFloatValue)
		}
	}

}

func TestControllerICOStage(t *testing.T) {
	controller, err := MakeTokenManagementController(InfuraController{"endpoint"}, "", "", "9df9993fcb4d9f4520770ce69f1623bccf3690489205c11e42a78bddc6526123")

	if assert.NoError(t, err) {
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

				data := requestBody.Params[0].(map[string]interface{})["data"].(string)

				testData, err := controller.crowdsaleContractABI.Pack("isIcoStage")

				if assert.NoError(t, err) {
					if data == hexutil.Bytes(testData).String() {
						return httpmock.NewStringResponse(
							http.StatusOK,
							`{"jsonrpc": "2.0", "result": "`+common.BigToAddress(big.NewInt(1)).String()+`"}`,
						), nil
					}
				}

				return httpmock.NewStringResponse(
					http.StatusOK,
					`{"jsonrpc": "2.0", "result": "0xabc"}`,
				), nil

			},
		)

		boolValue, err := controller.isICO()

		if assert.NoError(t, err) {
			assert.True(t, boolValue)
		}

		bigValue, err := controller.GetTokensLeft()

		if assert.NoError(t, err) {
			assert.Equal(t, common.HexToAddress("0xabc").Big(), bigValue)
		}

		bigFloatValue, err := controller.GetTokenExchangeRate()

		if assert.NoError(t, err) {
			assert.Equal(t, big.NewFloat(0).SetInt(common.HexToAddress("0xabc").Big()), bigFloatValue)
		}
	}

}

func TestControllerICOStage_Fail(t *testing.T) {
	controller, err := MakeTokenManagementController(InfuraController{"endpoint"}, "", "", "9df9993fcb4d9f4520770ce69f1623bccf3690489205c11e42a78bddc6526123")

	if assert.NoError(t, err) {
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

				data := requestBody.Params[0].(map[string]interface{})["data"].(string)

				testData, err := controller.crowdsaleContractABI.Pack("isIcoStage")

				if assert.NoError(t, err) {
					if data == hexutil.Bytes(testData).String() {
						return nil, errors.New("something wrong went with infura")
					}
				}

				return httpmock.NewStringResponse(
					http.StatusOK,
					`{"jsonrpc": "2.0", "result": "0xabc"}`,
				), nil

			},
		)

		boolValue, err := controller.isICO()

		assert.Error(t, err)
		assert.False(t, boolValue)

		bigValue, err := controller.GetTokensLeft()

		assert.Error(t, err)
		assert.Nil(t, bigValue)

		bigFloatValue, err := controller.GetTokenExchangeRate()

		assert.Error(t, err)
		assert.Nil(t, bigFloatValue)
	}

}
