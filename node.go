package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/cfoxon/jrc"
)

// Node Types
const (
	FULL = iota
	LIGHT
	HISTORY
)

// Features
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

var featureNodes map[int][]Node

func (x *Node) cleanup() {
	if x.Address[len(x.Address)-1] != '/' {
		x.Address += "/"
	}
}

// HealthCheck Does a health check on the node and sets active to correct state depending on results of the check
func (x *Node) HealthCheck() bool {
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
		body, err := io.ReadAll(resp.Body)
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

func LoadNodes(fileName string) []Node {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal("Error loading nodes file")
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal("couldn't close node file")
		}
	}(file)
	var readNodes = make([]ConfigNode, 0)
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&readNodes)
	if err != nil {
		log.Fatal("can't decode nodes JSON: ", err)
	}
	var liveStateNodes = make([]Node, 0)
	var fullTransactionHistoryNodes = make([]Node, 0)
	var accountHistoryNodes = make([]Node, 0)
	var marketHistoryNodes = make([]Node, 0)
	var nftHistroyNodes = make([]Node, 0)

	var newNodes = make([]Node, 0)

	featureNodes = map[int][]Node{
		LiveState:              nil,
		FullTransactionHistory: nil,
		AccountHistory:         nil,
		MarketHistory:          nil,
		NftHistory:             nil,
	}

	for _, node := range readNodes {
		if !node.Active {
			continue
		}
		newNode := Node{}
		newNode.Name = node.Name
		newNode.Address = node.Address
		if node.Type == "FULL" {
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
				liveStateNodes = append(liveStateNodes, newNode)
				featureNodes[LiveState] = liveStateNodes
			}
			if feature == "FullTransactionHistory" {
				newNode.Features[FullTransactionHistory] = true
				fullTransactionHistoryNodes = append(fullTransactionHistoryNodes, newNode)
				featureNodes[FullTransactionHistory] = fullTransactionHistoryNodes
			}
			if feature == "AccountHistory" {
				newNode.Features[AccountHistory] = true
				accountHistoryNodes = append(accountHistoryNodes, newNode)
				featureNodes[AccountHistory] = accountHistoryNodes
			}
			if feature == "MarketHistory" {
				newNode.Features[MarketHistory] = true
				marketHistoryNodes = append(marketHistoryNodes, newNode)
				featureNodes[MarketHistory] = marketHistoryNodes
			}
			if feature == "NftHistory" {
				newNode.Features[NftHistory] = true
				nftHistroyNodes = append(nftHistroyNodes, newNode)
				featureNodes[NftHistory] = nftHistroyNodes
			}
		}
		newNodes = append(newNodes, newNode)
	}
	return newNodes
}

// Gets an active node address with a a feature support
func GetNodeAddress(feature int) string {
	var nodeToReturn = featureNodes[feature][0].Address
	//Shift nodes
	return nodeToReturn
}
