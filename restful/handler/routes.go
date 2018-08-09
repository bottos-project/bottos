package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

const (
	GET_GEN_BLK_TIME = "/v1/node/generateblocktime"
	GET_CONN_COUNT   = "/v1/node/connectioncount"

	GET_BLK_TXS_BY_HEIGHT = "/v1/block/transactions/height/:height"
	GET_BLK_BY_HEIGHT     = "/v1/block/details/height/:height"
	GET_BLK_BY_HASH       = "/v1/block/details/hash/:hash"
	GET_BLK_HEIGHT        = "/v1/block/height"
	//GET_BLK_HASH          = "/v1/block/hash/:height"
	GET_BLK_INFO   = "/v1/block/header"
	GET_BLK_DETAIL = "/v1/block/detail"

	SEND_TRANSACTION = "/v1/transaction/send"
	GET_TRANSACTION  = "/v1/transaction/get"

	GET_ACCOUNT = "/v1/account/get"

	GET_KEY_VALUE = "/v1/account/get"

	GET_ABI = "/v1/contract/getabi"

	GET_TRANSFER_CREDIT = "/v1/balance/GetTransferCredit"

	GET_TX                = "/v1/transaction/:hash"
	GET_STORAGE           = "/v1/storage/:hash/:key"
	GET_BALANCE           = "/v1/balance/:addr"
	GET_CONTRACT_STATE    = "/v1/contract/:hash"
	GET_SMTCOCE_EVT_TXS   = "/v1/smartcode/event/transactions/:height"
	GET_SMTCOCE_EVTS      = "/v1/smartcode/event/txhash/:hash"
	GET_BLK_HGT_BY_TXHASH = "/v1/block/height/txhash/:hash"
	GET_MERKLE_PROOF      = "/v1/merkleproof/:hash"

	POST_RAW_TX = "/api/v1/transaction"
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
		"GET",
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
		GET_ABI,
		GetAbi,
	},
	Route{
		"",
		"POST",
		GET_TRANSFER_CREDIT,
		GetTransferCredit,
	},
}
