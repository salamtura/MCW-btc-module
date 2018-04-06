package controllers

import (
	"net/http"
	"testing"

	"MCW-btc-module/model"

	"github.com/ethereum/go-ethereum/common"
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
					`{"jsonrpc": "2.0", "result": "0xabc"}`,
				), nil

			},
		)
		btcAddressChan := make(chan string)
		errChan := make(chan error)

		go controller.BuyTokens(common.Address{}, btcAddressChan, errChan)
		address := <-btcAddressChan
		err := <-errChan

		assert.NoError(t, err)
		assert.NotEmpty(t, address)

	}
}
