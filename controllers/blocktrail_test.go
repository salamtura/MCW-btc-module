package controllers

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"
)

func TestMakeBlocktrailController(t *testing.T) {
	controller := MakeBlocktrailController("key", true)

	assert.Equal(t, blocktrailTestnetEndpoint+"key", controller.endpoint)

	controller = MakeBlocktrailController("key2", false)

	assert.Equal(t, blocktrailMainnetEndpoint+"key2", controller.endpoint)
}

func TestBlocktrailController_GetConfirmedBalance(t *testing.T) {
	controller := MakeBlocktrailController("key", true)

	testAddress := "testaddress"

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet,
		strings.Replace(controller.endpoint, "%address%", testAddress, 1),
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(
				http.StatusOK,
				map[string]interface{}{
					"data": []map[string]interface{}{
						{
							"hash":          "902912aeafe06a03ca95c70cad2e709c89e9b4f4a99aa6a0ae386408ae131b0f",
							"time":          "2014-09-05T17:08:04+0000",
							"confirmations": 279,
							"is_coinbase":   false,
							"value":         15000,
							"index":         0,
							"address":       "1NcXPMRaanz43b1kokpPuYDdk6GGDvxT2T",
							"type":          "pubkeyhash",
							"script":        "DUP HASH160 0x14 0xed12908714ffd43142bf9832692017e8ad54e9a8 EQUALVERIFY CHECKSIG",
							"script_hex":    "76a914ed12908714ffd43142bf9832692017e8ad54e9a888ac",
						}},
					"current_page": 1,
					"per_page":     20,
					"total":        4,
				})
		},
	)

	value, err := controller.GetConfirmedBalance(testAddress)

	if assert.NoError(t, err) {
		assert.Equal(t, 15000, value)
	}
}

func TestBlocktrailController_GetConfirmedBalance_Fail(t *testing.T) {
	controller := MakeBlocktrailController("key", true)

	testAddress := "testaddress"

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet,
		strings.Replace(controller.endpoint, "%address%", testAddress, 1),
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, "wrong json"), nil
		},
	)

	_, err := controller.GetConfirmedBalance(testAddress)

	assert.Error(t, err)

	httpmock.RegisterResponder(http.MethodGet,
		strings.Replace(controller.endpoint, "%address%", testAddress, 1),
		func(request *http.Request) (*http.Response, error) {
			return nil, errors.New("blocktrail malfunction")
		},
	)

	_, err = controller.GetConfirmedBalance(testAddress)

	assert.Error(t, err)
}
