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
	SEND_TRANSACTION       = "/v1/transaction/send"
	GET_TRANSACTION        = "/v1/transaction/get"
	GET_TRANSACTION_STATUS = "/v1/transaction/status"
	GET_HASH_FOR_SIGN      = "/v1/transaction/getHashForSign"
	GET_HASH_FOR_SIGN2     = "/v1/transaction/getHashForSign2"

	//Account
	GET_ACCOUNT_BRIEF   = "/v1/account/brief"
	GET_ACCOUNT         = "/v1/account/info"
	GET_TRANSFER_CREDIT = "/v1/balance/GetTransferCredit"

	//Contract
	GET_CONTRACT_ABI  = "/v1/contract/abi"
	GET_CONTRACT_CODE = "/v1/contract/code"

	// Common query
	QUERY_DB_VALUE = "/v1/common/queryDB"
	JSON_TO_BIN    = "/v1/common/jsontobin"

	//node
	GET_GEN_BLK_TIME = "/v1/node/generateblocktime"
	GET_CONN_COUNT   = "/v1/node/connectioncount"

	//delegate
	GET_ALL_DELEFATE = "/v1/delegate/getall"

	//global
	GET_GLOBAL_STAKED             = "/v1/global/stakedbalance"
	GET_FORECAST_RESOURCE_BALANCE = "/v1/resource/forecastresource"

	//p2p
	GET_ALL_PEERINFO           = "/v1/p2p/getpeers"
	CONNECT_PEER_BY_ADDRESS    = "/v1/p2p/connectpeer"
	DISCONNECT_PEER_BY_ADDRESS = "/v1/p2p/disconnectpeer"
	GET_PEER_STATE_BY_ADDRESS  = "/v1/p2p/getpeerstate"

	//MutlSign
	Proposal_Review = "/v1/proposal/review"
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
	/*Route{
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
	},*/
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
		GET_TRANSACTION_STATUS,
		GetTransactionStatus,
	},
	Route{
		"",
		"POST",
		GET_ACCOUNT_BRIEF,
		GetAccountBrief,
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
		QUERY_DB_VALUE,
		QueryDBValue,
	},
	Route{
		"",
		"POST",
		JSON_TO_BIN,
		JsonToBin,
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
	Route{
		"",
		"POST",
		GET_ALL_DELEFATE,
		GetAllDelegates,
	},
	Route{
		"",
		"GET",
		GET_GLOBAL_STAKED,
		GetGlobalStakedBalance,
	},
	Route{
		"",
		"POST",
		GET_FORECAST_RESOURCE_BALANCE,
		GetForecastResBalance,
	},
	Route{
		"",
		"POST",
		CONNECT_PEER_BY_ADDRESS,
		ConnectPeerbyAddress,
	},
	Route{
		"",
		"POST",
		DISCONNECT_PEER_BY_ADDRESS,
		DisConnectPeerbyAddress,
	},
	Route{
		"",
		"POST",
		GET_PEER_STATE_BY_ADDRESS,
		GetPeerStatebyAddress,
	},
	Route{
		"",
		"POST",
		GET_ALL_PEERINFO,
		GetPeers,
	},
	//Multi Sign
	Route{
		"",
		"POST",
		Proposal_Review,
		ReviewProposal,
	},
	Route{
		"",
		"POST",
		GET_HASH_FOR_SIGN,
		GetHashForSign,
	},
	Route{
		"",
		"POST",
		GET_HASH_FOR_SIGN2,
		GetHashForSign2,
	},
}
