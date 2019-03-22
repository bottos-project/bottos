package signature

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/bottos-project/bottos/api"
	"github.com/bottos-project/bottos/common"
	berr "github.com/bottos-project/bottos/common/errors"
	"github.com/bottos-project/bottos/config"
	"github.com/bottos-project/crypto-go/crypto"
	log "github.com/cihub/seelog"
)

func SignWithKey(digest []byte, privateKey []byte) ([]byte, error) {
	signdata, err := crypto.Sign(digest, privateKey)
	if err != nil {
		return nil, errors.New("crypto signature failed")
	}
	return signdata, err
}

func SignWithWallet(digest []byte, account string, url string) ([]byte, error) {
	values := map[string]interface{}{
		"account_name": account,
		"hash":         common.BytesToHex(digest),
	}
	if account != config.BOTTOS_CONTRACT_NAME {
		values["type"] = "delegate"
	}
	jsonValue, _ := json.Marshal(values)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil || resp.StatusCode != 200 {
		return nil, fmt.Errorf("wallet signature failed: %v", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("wallet signature failed: %v", err)
	}
	var respStruct api.SignDataResponse
	err = json.Unmarshal(body, &respStruct)
	if err != nil || respStruct.Errcode != uint32(berr.ErrNoError) {
		return nil, fmt.Errorf("wallet signature failed: %v, errorcode: %v", err, respStruct.Errcode)
	}
	signdata, err := common.HexToBytes(respStruct.Result.SignValue)
	return signdata, err
}

func SignHash(digest []byte, account string, url string) ([]byte, error) {
	values := map[string]interface{}{
		"account_name": account,
		"hash":         common.BytesToHex(digest),
	}
	jsonValue, _ := json.Marshal(values)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil || resp.StatusCode != 200 {
		return nil, fmt.Errorf("wallet signature failed: %v", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("wallet signature failed: %v", err)
	}
	var respStruct api.SignDataResponse
	err = json.Unmarshal(body, &respStruct)
	if err != nil || respStruct.Errcode != uint32(berr.ErrNoError) {
		return nil, fmt.Errorf("wallet signature failed: %v", err)
	}
	signdata, err := common.HexToBytes(respStruct.Result.SignValue)
	return signdata, err
}

func SignByDelegate(digest []byte, pubkey string) ([]byte, error) {
	if config.BtoConfig.Delegate.Signature.Type == "key" {
		prikey, err := config.GetDelegateSignKey(pubkey)
		if err != nil {
			return nil, errors.New("crypto signature failed")
		}
		signdata, err := SignWithKey(digest, prikey)
		if err != nil {
			log.Errorf("COMMON sign with key, error %v, data %x, digest %x, pubkey %v", err, signdata, digest, pubkey)
		}
		return signdata, err
	} else if config.BtoConfig.Delegate.Signature.Type == "wallet" {
		signdata, err := SignWithWallet(digest, config.BtoConfig.Delegate.Account, config.BtoConfig.Delegate.Signature.URL)
		if err != nil {
			log.Errorf("COMMON sign with wallet, error %v, data %x, digest %x, url %v, account %v", err, signdata, digest, config.BtoConfig.Delegate.Signature.URL, config.BtoConfig.Delegate.Account)
		}
		return signdata, err
	}

	return nil, errors.New("crypto signature failed")
}

func Sign(digest []byte, prikey []byte) ([]byte, error) {
	return SignWithKey(digest, prikey)
}

func VerifySign(pubkey []byte, digest []byte, signdata []byte) bool {
	result := crypto.VerifySign(pubkey, digest, signdata)
	log.Debugf("COMMON verify signature, result %v, signdata %x, pubkey %x, digest %x", result, signdata, pubkey, digest)
	return result
}
