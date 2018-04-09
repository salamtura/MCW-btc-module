package controllers

import (
	"net/http"
	"strings"
	"testing"

	"MCW-btc-module/model"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"
)

func TestMakeExchangeController(t *testing.T) {
	db, err := gorm.Open("sqlite3", "test.db")

	db.DropTableIfExists(model.BTCTransaction{})
	db.AutoMigrate(model.BTCTransaction{})

	assert.NoError(t, err)

	controller, err := MakeExchangeController(
		MonitoringController{},
		TokenManagementController{},
		db,
		"tpubDAbGZM7PHnNp75QbARzDM7id7zpBcxKH9EX7VFE2Pr15EuWQEdRzSZSB4fhnHBxeLyzZB6QnhewQQhdkRHx6wCow3iTj6BXfwGsj8RevWoC",
	)

	assert.NoError(t, err)
	assert.NotNil(t, controller)
}

func TestMakeExchangeControllerError(t *testing.T) {
	db, err := gorm.Open("sqlite3", "test.db")

	assert.NoError(t, err)

	controller, err := MakeExchangeController(
		MonitoringController{},
		TokenManagementController{},
		db,
		"pubDAbGZM7PHnNp75QbARzDM7id7zpBcxKH9EX7VFE2Pr15EuWQEdRzSZSB4fhnHBxeLyzZB6QnhewQQhdkRHx6wCow3iTj6BXfwGsj8RevWoC",
	)

	assert.Error(t, err)
	assert.Nil(t, controller)
}

func TestGetExchangeRate(t *testing.T) {
	db, err := gorm.Open("sqlite3", "test.db")

	assert.NoError(t, err)

	controller, err := MakeExchangeController(
		MonitoringController{},
		TokenManagementController{},
		db,
		"tpubDAbGZM7PHnNp75QbARzDM7id7zpBcxKH9EX7VFE2Pr15EuWQEdRzSZSB4fhnHBxeLyzZB6QnhewQQhdkRHx6wCow3iTj6BXfwGsj8RevWoC",
	)

	if assert.NoError(t, err) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder(
			http.MethodGet,
			"https://shapeshift.io/rate/btc_eth",
			func(request *http.Request) (*http.Response, error) {

				return httpmock.NewStringResponse(
					http.StatusOK,
					`{"rate": "10"}`,
				), nil

			},
		)

		rate, err := controller.getExchangeRate()

		assert.NoError(t, err)
		assert.Equal(t, float64(10), rate)
	}
}

func TestExchangeController_BuyTokens(t *testing.T) {
	db, err := gorm.Open("sqlite3", "test.db")

	assert.NoError(t, err)

	controller, err := MakeExchangeController(
		MonitoringController{},
		TokenManagementController{},
		db,
		"tpubDAbGZM7PHnNp75QbARzDM7id7zpBcxKH9EX7VFE2Pr15EuWQEdRzSZSB4fhnHBxeLyzZB6QnhewQQhdkRHx6wCow3iTj6BXfwGsj8RevWoC",
	)

	if assert.NoError(t, err) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder(
			http.MethodGet,
			"https://shapeshift.io/rate/btc_eth",
			func(request *http.Request) (*http.Response, error) {

				return httpmock.NewStringResponse(
					http.StatusOK,
					`{"rate": "10"}`,
				), nil

			},
		)
		httpmock.RegisterResponder(
			http.MethodPost,
			controller.infuraEndpoint,
			func(request *http.Request) (*http.Response, error) {

				return httpmock.NewStringResponse(
					http.StatusOK,
					`{"jsonrpc": "2.0", "result": "0x1"}`,
				), nil

			},
		)

		transaction, isNew, err := controller.CreateTransactionEntry("")

		httpmock.RegisterResponder(http.MethodGet,
			strings.Replace(controller.endpoint, "%address%", transaction.BitcoinAddress, 1),
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

		if assert.NoError(t, err) {
			assert.True(t, isNew)
			controller.BuyTokens(transaction)

		}
	}
}
