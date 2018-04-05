package account

import (
	"fmt"
)

type AccountManager struct {
	defAccount string
}

func CreateAccountManager() *AccountManager {
	am := AccountManager{}

	return &am
}

func getAccount() int {

	return 0
}

func SetAccount(aaa string) (int, string) {
	fmt.Println(aaa)
	return 0, "good"
}
