package core

import (
	"context"
	"errors"
	"fmt"

	"github.com/bottos-project/bottos/core/trx"
	"github.com/bottos-project/bottos/core/common/types"
)

type CoreService struct {
	txRepo txRepository
	txp *trx.TxPool
}

type txRepository interface {
	CallSendTrx(account_name string, balance uint64) (string, error)
}

func NewCoreSrvice(txpool *trx.TxPool) *CoreService {
	return &CoreService{txp: txpool}
}

func (h *CoreService) SendTrx(ctx context.Context, req *TxRequest, rsp *TxResponse) error {
	if req == nil {

		return errors.New("Missing storage request")
	}
	fmt.Println(req.AccountName)

	//id, err := trx.CallSendTrx(req.AccountName, 111)
	//if err != nil {
	//	return errors.New("get PUTURL failed")
	//}

	h.txp.Add(&types.Transaction{Id: "111", AccountName: req.AccountName})

	fmt.Println("success")
	rsp.Id = req.AccountName
	return nil
}
