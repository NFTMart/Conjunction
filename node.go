package main

import (
	"github.com/cfoxon/jrc"
	"io/ioutil"
	"net/http"
)

//Node Types
const (
	FULL = iota
	LIGHT
	HISTORY
)

//Features
const (
	LIVE_STATE = iota
	FULL_TRANSACTION_HISTORY
	ACCOUNT_HISTORY
	MARKET_HISTORY
	NFT_HISTORY
)

const (
	ACCOUNT_HISTORY_ENDPOINT = "accountHistory"
	MARKET_HISTORY_ENDPOINT  = "marketHistory"
	NFT_HISTORY_ENDPOINT     = "nftHistory"
)

type Node struct {
	Name     string       `json:"name"`
	Address  string       `json:"address"`
	Type     int          `json:"type"`
	Active   bool         `json:"active"`
	Features map[int]bool `json:"features"`
}

func (x Node) cleanup() {
	if x.Address[len(x.Address)] != '/' {
		x.Address += "/"
	}
}

// HealthCheck Does a health check on the node and sets active to correct state depending on results of the check
func (x Node) HealthCheck() bool {
	x.cleanup()
	if x.Type == FULL || x.Type == LIGHT {
		rpcClient, err := jrc.NewServer(x.Address)
		if err != nil {
			x.Active = false
			return false
		}
		jr2query := jrc.RpcRequest{Method: "blockchain.getStatus", JsonRpc: "2.0", Id: 0}
		resp, err := rpcClient.Exec(jr2query)
		if err != nil {
			x.Active = false
			return false
		}
		if resp.Result == nil {
			x.Active = false
			return false
		}
		x.Active = true
		return true
	}

	if x.Type == HISTORY {
		var query = ""
		if x.Features[ACCOUNT_HISTORY] {
			query = ACCOUNT_HISTORY_ENDPOINT + "?account=null"
		} else if x.Features[MARKET_HISTORY] {
			query = MARKET_HISTORY_ENDPOINT + "?symbol=BEE"
		} else if x.Features[NFT_HISTORY] {
			query = NFT_HISTORY_ENDPOINT + "?id=1"
		} else {
			x.Active = false
			return false
		}
		resp, err := http.Get(x.Address + query)
		if err == nil {
			x.Active = false
			return false
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			x.Active = false
			return false
		}
		if body == nil {
			x.Active = false
			return false
		}
		x.Active = true
		return true
	}
}
