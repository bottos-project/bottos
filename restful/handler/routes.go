package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

const (
	//block
	GET_BLK_INFO   = "/v1/block/height"
	GET_BLK_DETAIL = "/v1/block/detail"

	// Transaction
	SEND_TRANSACTION = "/v1/transaction/send"
	GET_TRANSACTION  = "/v1/transaction/get"
	GET_TRANSACTION_STATUS = "/v1/transaction/status"

	//Account
	GET_ACCOUNT_BRIEF   = "/v1/account/brief"
	GET_ACCOUNT = "/v1/account/info"
	GET_TRANSFER_CREDIT = "/v1/balance/GetTransferCredit"

	//Contract
	GET_CONTRACT_ABI = "/v1/contract/abi"
	GET_CONTRACT_CODE = "/v1/contract/code"

	// Common query
	GET_KEY_VALUE = "/v1/common/query"

	//node
	GET_GEN_BLK_TIME = "/v1/node/generateblocktime"
	GET_CONN_COUNT   = "/v1/node/connectioncount"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"TodoIndex",
		"GET",
		"/todos",
		TodoIndex,
	},
	Route{
		"TodoShow",
		"GET",
		"/todos/{todoId}",
		TodoShow,
	},
	Route{
		"",
		"GET",
		GET_GEN_BLK_TIME,
		GetGenerateBlockTime,
	},


	Route{
		"",
		"GET",
		GET_BLK_INFO,
		GetInfo,
	},
	Route{
		"",
		"POST",
		GET_BLK_DETAIL,
		GetBlock,
	},
	Route{
		"",
		"POST",
		SEND_TRANSACTION,
		SendTransaction,
	},
	Route{
		"",
		"POST",
		GET_TRANSACTION,
		GetTransaction,
	},
	Route{
		"",
		"POST",
		GET_ACCOUNT,
		GetAccount,
	},
	Route{
		"",
		"POST",
		GET_KEY_VALUE,
		GetKeyValue,
	},
	Route{
		"",
		"POST",
		GET_CONTRACT_ABI,
		GetContractAbi,
	},
	Route{
		"",
		"POST",
		GET_CONTRACT_CODE,
		GetContractCode,
	},
	Route{
		"",
		"POST",
		GET_TRANSFER_CREDIT,
		GetTransferCredit,
	},
}
