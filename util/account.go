package util

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Account struct {
	ecdsa.PrivateKey
}

func NewAccount(privKey *ecdsa.PrivateKey) *Account {
	return &Account{*privKey}
}

func NewAccountFromRaw(privKeyStr string) (*Account, error) {
	privKey, err := crypto.HexToECDSA(privKeyStr)
	if err != nil {
		return nil, err
	}
	return NewAccount(privKey), nil
}

func (a *Account) Address() common.Address {
	return crypto.PubkeyToAddress(a.PublicKey)
}
