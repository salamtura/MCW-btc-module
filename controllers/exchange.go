package controllers

import (
	"encoding/json"
	"math"
	"fmt"
	"math/big"
	"sync"

	"MCW-btc-module/helpers"
	"MCW-btc-module/model"

	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jinzhu/gorm"
)

type ExchangeController struct {
	MonitoringController
	TokenManagementController
	database          *gorm.DB
	xpub              *hdkeychain.ExtendedKey
	currentChildIndex uint32
	indexMutex        *sync.Mutex
}

func MakeExchangeController(
	monitoringController MonitoringController,
	tokenManagementController TokenManagementController,
	database *gorm.DB,
	xpubString string,
) (*ExchangeController, error) {
	xpub, err := hdkeychain.NewKeyFromString(xpubString)
	if err != nil {
		return nil, err
	}

	latestTransaction := model.BTCTransaction{}

	database.Order("index desc").First(&latestTransaction)

	return &ExchangeController{
		MonitoringController:      monitoringController,
		TokenManagementController: tokenManagementController,
		database:                  database,
		xpub:                      xpub,
		currentChildIndex:         latestTransaction.Index,
		indexMutex:                &sync.Mutex{},
	}, nil
}

func (controller ExchangeController) getExchangeRate() (float64, error) {
	type ShapeshiftResponse struct {
		Rate big.Float `json:"rate, string"`
	}
	_, responseBytes, err := helpers.Get("https://shapeshift.io/rate/btc_eth", helpers.Headers{})

	if err != nil {
		return 0, err
	}

	response := new(ShapeshiftResponse)

	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return 0, err
	}

	value, _ := response.Rate.Float64()

	return value, nil
}

func (controller *ExchangeController) BuyTokens(ethereumAddress common.Address, bitcoinAddressChannel chan<- string, errorChannel chan<- error) {
	controller.indexMutex.Lock()
	index := controller.currentChildIndex
	index += 1

	address, err := helpers.DeriveAddress(controller.xpub, index)
	controller.currentChildIndex = index
	controller.indexMutex.Unlock()

	if err != nil {
		bitcoinAddressChannel <- address
		errorChannel <- err
		return
	}

	transaction := &model.BTCTransaction{
		EthereumAddress: ethereumAddress.String(),
		BitcoinAddress:  address,
		Index:           index,
		Status:          model.TRANSACTON_STATUS_NEW,
	}

	err = controller.database.Create(transaction).Error
	fmt.Println(err)
	bitcoinAddressChannel <- address
	errorChannel <- err

	if err != nil {
		return
	}

	rate, err := controller.getExchangeRate()

	if err != nil {
		transaction.Error = err.Error()
		transaction.Status = model.TRANSACTION_STATUS_ERROR
		controller.database.Save(transaction)
		return
	}

	receivedBTC, err := controller.MonitoringController.waitForTransfer(address)

	if err != nil {
		transaction.Error = err.Error()
		transaction.Status = model.TRANSACTION_STATUS_ERROR
		controller.database.Save(transaction)
		return
	}

	transaction.AmountTransferred = receivedBTC
	controller.database.Save(transaction)

	receivedEth := receivedBTC * rate

	//exchangeRate, err := controller.TokenManagementController.GetTokenExchangeRate()
	//
	//if err != nil {
	//	transaction.Error = err.Error()
	//	transaction.Status = model.TRANSACTION_STATUS_ERROR
	//	controller.database.Save(transaction)
	//	return
	//}

	bigWei := big.NewFloat(0).Mul(big.NewFloat(receivedEth), big.NewFloat(math.Pow(10, 18)))

	//tokensAmount, _ := big.NewFloat(0).Mul(bigWei, exchangeRate).Int(nil)
	//
	//tokensLeft, err := controller.TokenManagementController.GetTokensLeft()
	//
	//if err != nil {
	//	transaction.Error = err.Error()
	//	transaction.Status = model.TRANSACTION_STATUS_ERROR
	//	controller.database.Save(transaction)
	//	return
	//}
	//
	//tokensToTransfer := tokensAmount
	//
	//if tokensAmount.Cmp(tokensLeft) == 1 {
	//	tokensToTransfer = tokensLeft
	//}

	bigIntWei, _ := bigWei.Int(nil)

	if err := controller.TokenManagementController.MintTokens(ethereumAddress, bigIntWei); err != nil {
		transaction.Error = err.Error()
		transaction.Status = model.TRANSACTION_STATUS_ERROR
		controller.database.Save(transaction)
		return
	}

	transaction.Status = model.TRANSACTION_STATUS_SUCCESS
	controller.database.Save(transaction)
}
