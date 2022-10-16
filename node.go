package main

import (
	"encoding/json"
	"github.com/cfoxon/jrc"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

//Node Types
const (
	FULL = iota
	LIGHT
	HISTORY
)

//Features
const (
	LiveState = iota
	FullTransactionHistory
	AccountHistory
	MarketHistory
	NftHistory
)

const (
	AccountHistoryEndpoint = "accountHistory"
	MarketHistoryEndpoint  = "marketHistory"
	NftHistoryEndpoint     = "nftHistory"
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
		if x.Features[AccountHistory] {
			query = AccountHistoryEndpoint + "?account=null"
		} else if x.Features[MarketHistory] {
			query = MarketHistoryEndpoint + "?symbol=BEE"
		} else if x.Features[NftHistory] {
			query = NftHistoryEndpoint + "?id=1"
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
	return false
}

type ConfigNode struct {
	Name     string   `json:"name"`
	Address  string   `json:"address"`
	Type     string   `json:"type"`
	Active   bool     `json:"active"`
	Features []string `json:"features"`
}

func LoadNodes() []Node {
	file, err := os.Open("nodes.json")
	if err != nil {
		log.Fatal("Error loading nodes file")
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal("couldn't close node file")
		}
	}(file)
	var readNodes = make([]ConfigNode, 100)
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&readNodes)
	if err != nil {
		log.Fatal("can't decode nodes JSON: ", err)
	}
	var newNodes = make([]Node, 100)
	for _, node := range readNodes {
		if !node.Active {
			continue
		}
		newNode := Node{}
		newNode.Name = node.Name
		newNode.Address = node.Address
		if node.Type == "Full" {
			newNode.Type = FULL
		} else if node.Type == "LIGHT" {
			newNode.Type = LIGHT
		} else if node.Type == "HISTORY" {
			newNode.Type = HISTORY
		} else {
			log.Fatal("Invalid type detected on node: ", node.Name)
		}
		newNode.Active = node.Active
		newNode.Features = map[int]bool{ // Map literal
			LiveState:              false,
			FullTransactionHistory: false,
			AccountHistory:         false,
			MarketHistory:          false,
			NftHistory:             false,
		}
		for _, feature := range node.Features {
			if feature == "LiveState" {
				newNode.Features[LiveState] = true
			}
			if feature == "FullTransactionHistory" {
				newNode.Features[FullTransactionHistory] = true
			}
			if feature == "AccountHistory" {
				newNode.Features[AccountHistory] = true
			}
			if feature == "MarketHistory" {
				newNode.Features[MarketHistory] = true
			}
			if feature == "NftHistory" {
				newNode.Features[NftHistory] = true
			}
		}
		newNodes = append(newNodes, newNode)
	}
	return newNodes
}
