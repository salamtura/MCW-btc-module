package controllers

import (
	"crypto/ecdsa"
	"errors"
	"math/big"
	"strings"

	"MCW-btc-module/gocontracts"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type TokenManagementController struct {
	InfuraController
	crowdsaleContractABI     abi.ABI
	crowdsaleContractAddress common.Address
	crowdsaleOwnerAddress    common.Address
	OwnerPrivateKey          *ecdsa.PrivateKey
}

func MakeTokenManagementController(infuraController InfuraController, crowdsaleAddress string, crowdsaleOwnerAddress string, crowdsaleOwnerPrivateKey string) (*TokenManagementController, error) {
	crowdsaleABI, err := abi.JSON(strings.NewReader(gocontracts.MocrowCoinCrowdsaleABI))

	if err != nil {
		return nil, err
	}

	key, err := crypto.HexToECDSA(crowdsaleOwnerPrivateKey)

	if err != nil {
		return nil, err
	}

	return &TokenManagementController{
		InfuraController:         infuraController,
		crowdsaleContractABI:     crowdsaleABI,
		crowdsaleContractAddress: common.HexToAddress(crowdsaleAddress),
		crowdsaleOwnerAddress:    common.HexToAddress(crowdsaleOwnerAddress),
		OwnerPrivateKey:          key,
	}, nil
}

func (controller TokenManagementController) getCrowdaleParameter(parameter string) (*common.Address, error) {
	return controller.InfuraController.retrieveParameter(parameter, controller.crowdsaleContractABI, controller.crowdsaleContractAddress)
}

func (controller TokenManagementController) isPreICO() (bool, error) {
	parameter, err := controller.getCrowdaleParameter("isPreIco")
	if err != nil {
		return false, err
	}

	return parameter.Big().Int64() == 1, nil
}

func (controller TokenManagementController) preICOTokensRemaining() (*big.Int, error) {
	parameter, err := controller.getCrowdaleParameter("tokensRemainingPreIco")
	if err != nil {
		return nil, err
	}

	return parameter.Big(), nil
}

func (controller TokenManagementController) preICOExchangeRate() (*big.Int, error) {
	parameter, err := controller.getCrowdaleParameter("TOKEN_RATE_PRE_ICO")
	if err != nil {
		return nil, err
	}

	return parameter.Big(), nil
}

func (controller TokenManagementController) isICO() (bool, error) {
	parameter, err := controller.getCrowdaleParameter("isIco")
	if err != nil {
		return false, err
	}

	return parameter.Big().Int64() == 1, nil
}

func (controller TokenManagementController) icoTokensRemaining() (*big.Int, error) {
	parameter, err := controller.getCrowdaleParameter("tokensRemainingIco")
	if err != nil {
		return nil, err
	}

	return parameter.Big(), nil
}

func (controller TokenManagementController) icoExchangeRate() (*big.Int, error) {
	parameter, err := controller.getCrowdaleParameter("TOKEN_RATE_ICO")
	if err != nil {
		return nil, err
	}

	return parameter.Big(), nil
}

func (controller TokenManagementController) GetTokensLeft() (*big.Int, error) {
	isPreICO, err := controller.isPreICO()

	if err != nil {
		return nil, err
	}

	isICO, err := controller.isICO()
	if err != nil {
		return nil, err
	}

	tokensLeft := big.NewInt(0)

	if isPreICO {
		tokensLeft, err = controller.preICOTokensRemaining()
	} else if isICO {
		tokensLeft, err = controller.icoTokensRemaining()
	}

	return tokensLeft, err
}

func (controller TokenManagementController) GetTokenExchangeRate() (*big.Float, error) {
	isPreICO, err := controller.isPreICO()

	if err != nil {
		return nil, err
	}

	isICO, err := controller.isICO()
	if err != nil {
		return nil, err
	}

	rate := big.NewInt(0)

	if isPreICO {
		rate, err = controller.preICOExchangeRate()
	} else if isICO {
		rate, err = controller.icoExchangeRate()
	}

	return big.NewFloat(0).SetInt(rate), err
}

func (controller TokenManagementController) MintTokens(receiver common.Address, weiAmount *big.Int) error {
	packedData := make([]byte, 0)

	isPreICO, err := controller.isPreICO()

	if err != nil {
		return err
	}

	isICO, err := controller.isICO()
	if err != nil {
		return err
	}

	if isPreICO {
		packedData, err = controller.crowdsaleContractABI.Pack("sellTokensForBTCPreIco", receiver, weiAmount)
	} else if isICO {
		packedData, err = controller.crowdsaleContractABI.Pack("sellTokensForBTCIco", receiver, weiAmount)
	} else {
		err = errors.New("it's neither pre-ico, nor ico")
	}

	if err != nil {
		return err
	}

	gasPriceChan := make(chan common.Address)
	gasLimitChan := make(chan common.Address)
	nonceChan := make(chan common.Address)

	go func() {
		response, err := controller.InfuraController.callInfura(makeRequestPayload(
			"eth_gasPrice",
			[]interface{}{},
		))

		if err != nil {
			gasPriceChan <- common.Address{}
			return
		}

		gasPriceChan <- common.HexToAddress(response.Result)
	}()

	go func() {
		response, err := controller.InfuraController.callInfura(makeEstimateGasRequestPayload(
			controller.crowdsaleContractAddress,
			packedData,
		))

		if err != nil {
			gasLimitChan <- common.BigToAddress(big.NewInt(500000))
			return
		}

		gasLimitChan <- common.HexToAddress(response.Result)
	}()

	go func() {
		response, err := controller.InfuraController.callInfura(makeRequestPayload(
			"eth_getTransactionCount",
			[]interface{}{controller.crowdsaleOwnerAddress, "pending"},
		))

		if err != nil {
			nonceChan <- common.Address{}
			return
		}

		nonceChan <- common.HexToAddress(response.Result)
	}()

	gasPrice := <-gasPriceChan
	gasLimit := <-gasLimitChan
	nonce := <-nonceChan

	// We do not check nonce here since it actually can be zero.
	if (gasPrice == common.Address{}) {
		return errors.New("unable to fetch gas price. aborting")
	}

	transaction := types.NewTransaction(
		nonce.Big().Uint64(),
		controller.crowdsaleContractAddress,
		big.NewInt(0),
		gasLimit.Big().Uint64(),
		gasPrice.Big(),
		packedData,
	)

	signer := types.HomesteadSigner{}

	signature, err := crypto.Sign(signer.Hash(transaction).Bytes(), controller.OwnerPrivateKey)

	if err != nil {
		return err
	}

	signedTransaction, err := transaction.WithSignature(signer, signature)

	if err != nil {
		return err
	}

	_, err = controller.InfuraController.callInfura(makeRequestPayload(
		"eth_sendRawTransaction",
		[]interface{}{hexutil.Bytes(types.Transactions{signedTransaction}.GetRlp(0)).String()},
	))

	if err != nil {
		return err
	}

	return nil
}
