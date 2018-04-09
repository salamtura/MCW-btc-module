package controllers

import (
	"encoding/json"
	"log"
	"math"
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

	latestTransaction := model.Transaction{}

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

func (controller *ExchangeController) CreateTransactionEntry(ethereumAddress string) (*model.Transaction, bool, error) {
	transaction := new(model.Transaction)

	err := controller.database.Where("ethereum_address = ?", ethereumAddress).Order("id desc", false).First(transaction).Error

	if err == nil && transaction.Status == model.TRANSACTON_STATUS_NEW && transaction.BitcoinAddress != "" {
		return transaction, false, nil
	}

	controller.indexMutex.Lock()
	index := controller.currentChildIndex
	index += 1

	address, err := helpers.DeriveAddress(controller.xpub, index)
	controller.currentChildIndex = index
	controller.indexMutex.Unlock()

	if err != nil {
		return nil, false, err
	}

	transaction = &model.Transaction{
		EthereumAddress: ethereumAddress,
		BitcoinAddress:  address,
		Index:           index,
		Status:          model.TRANSACTON_STATUS_NEW,
	}

	err = controller.database.Create(transaction).Error

	return transaction, true, err
}

func (controller *ExchangeController) BuyTokens(transaction *model.Transaction) {
	rate, err := controller.getExchangeRate()

	if err != nil {
		transaction.Error = err.Error()
		transaction.Status = model.TRANSACTION_STATUS_ERROR
		controller.database.Save(transaction)
		return
	}

	receivedBTC, err := controller.MonitoringController.waitForTransfer(transaction.BitcoinAddress)

	if err != nil {
		transaction.Error = err.Error()
		transaction.Status = model.TRANSACTION_STATUS_ERROR
		controller.database.Save(transaction)
		return
	}

	transaction.AmountTransferred = receivedBTC
	controller.database.Save(transaction)

	receivedEth := receivedBTC * rate

	exchangeRate, err := controller.TokenManagementController.GetTokenExchangeRate()

	if err != nil {
		transaction.Error = err.Error()
		transaction.Status = model.TRANSACTION_STATUS_ERROR
		controller.database.Save(transaction)
		return
	}

	bigWei := big.NewFloat(0).Mul(big.NewFloat(receivedEth), big.NewFloat(math.Pow(10, 18)))

	tokensAmount, _ := big.NewFloat(0).Mul(bigWei, exchangeRate).Int(nil)

	tokensLeft, err := controller.TokenManagementController.GetTokensLeft()

	if err != nil {
		transaction.Error = err.Error()
		transaction.Status = model.TRANSACTION_STATUS_ERROR
		controller.database.Save(transaction)
		return
	}

	tokensToTransfer := tokensAmount

	if tokensAmount.Cmp(tokensLeft) == 1 {
		tokensToTransfer = tokensLeft
	}

	if err := controller.TokenManagementController.MintTokens(common.HexToAddress(transaction.EthereumAddress), tokensToTransfer); err != nil {
		transaction.Error = err.Error()
		transaction.Status = model.TRANSACTION_STATUS_ERROR
		controller.database.Save(transaction)
		return
	}

	transaction.Status = model.TRANSACTION_STATUS_SUCCESS
	controller.database.Save(transaction)
}

func (controller ExchangeController) ResumeMonitoring() {
	unfinishedTransactions := new([]model.Transaction)

	if err := controller.database.Where("status = ?", model.TRANSACTON_STATUS_NEW).Find(unfinishedTransactions).Error; err != nil {
		log.Println(err)
	}

	for _, transaction := range *unfinishedTransactions {
		go controller.BuyTokens(&transaction)
	}
}
