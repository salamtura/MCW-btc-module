package controllers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeBlockCypherController(t *testing.T) {
	controller := MakeBlockCypherController("token", true)

	assert.Equal(t, "token", controller.client.Token)
	assert.Equal(t, "btc", controller.client.Coin)
	assert.Equal(t, "test3", controller.client.Chain)
}
