package account

import (
	"crypto"
)

type Account struct {
	PrivateKey crypto.PrivateKey
	PublicKey  crypto.PublicKey
	Name       string
}

func CreateAccount() *Account {
	acc := Account{}

	return &acc
}
