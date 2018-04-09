package helpers

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
)

func DeriveAddress(key *hdkeychain.ExtendedKey, index uint32, isTestnet bool) (string, error) {
	derivedKey, err := key.Child(index)

	if err != nil {
		return "", err
	}

	params := chaincfg.MainNetParams
	if isTestnet {
		params = chaincfg.TestNet3Params
	}

	address, err := derivedKey.Address(&params)

	if err != nil {
		return "", err
	}

	return address.EncodeAddress(), nil
}
