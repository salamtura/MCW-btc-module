package helpers

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/stretchr/testify/assert"
)

func TestDeriveAddress(t *testing.T) {
	masterKey, err := hdkeychain.NewKeyFromString("tpubDB7iVAmGkzub1fkjb46T5Pqqw6RiyvhmmwT4KnDwgxPwBDAjXKv9SYLRwLcSzryP9pEbytkaVQRs51a5TvTykaAde2czxFKbeStDv1iY8qF")
	testDerivedKey, err := masterKey.Child(uint32(masterKey.Depth()) + 1)

	if assert.NoError(t, err) {
		address, err := DeriveAddress(masterKey, uint32(masterKey.Depth())+1)

		if assert.NoError(t, err) {

			derivedKeyAddress, err := testDerivedKey.Address(&chaincfg.TestNet3Params)

			assert.NoError(t, err)
			assert.Equal(t, derivedKeyAddress.EncodeAddress(), address)
		}
	}

}
