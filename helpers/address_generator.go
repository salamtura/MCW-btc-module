package helpers

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
)

func DeriveAddress(key *hdkeychain.ExtendedKey, index uint32) (string, error) {
	derivedKey, err := key.Child(index)

	if err != nil {
		return "", err
	}
	address, err := derivedKey.Address(&chaincfg.TestNet3Params)

	if err != nil {
		return "", err
	}

	return address.EncodeAddress(), nil
}
