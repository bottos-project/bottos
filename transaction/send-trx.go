package transaction

import (
	"fmt"
)

func Add() error {
	return nil
}

func CallSendTrx(account_name string, balance uint64) (string, error) {

	fmt.Println("receive transaction")
	Add()

	return "get trx", nil
}
