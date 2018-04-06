package model

type BTCTransaction struct {
	ID                uint    `gorm:"primary_key" json:"id"`
	EthereumAddress   string  `json:"ethereumAddress"`
	BitcoinAddress    string  `json:"bitcoinAddress"`
	AmountTransferred float64 `json:"amountTransferred"`
	Index             uint32  `json:"depth"`
	Error             string  `json:"error"`
	Status            int8    `json:"status"`
}

const TRANSACTION_STATUS_ERROR = -1
const TRANSACTON_STATUS_NEW = 0
const TRANSACTION_STATUS_SUCCESS = 1
