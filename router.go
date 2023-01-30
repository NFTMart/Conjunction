package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/cfoxon/jrc"
	"github.com/gin-gonic/gin"
)

type QueryParams struct {
	Jsonrpc string          `json:"jsonrpc"`
	Id      int             `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

type QueryResponseError struct {
	Jsonrpc string                `json:"jsonrpc"`
	Id      int                   `json:"id"`
	Error   QueryResponseInternal `json:"error"`
}

type QueryResponseInternal struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func Router(r *gin.Engine) {
	r.POST("/", handleMain)
	r.GET("/accountHistory", handleHistory)
	r.GET("/nftHistory", handleHistory)
	r.GET("/marketHistory", handleHistory)
	r.POST("/contracts", handleContracts)
	r.POST("/blockchain", handleBlockchain)
}

func handleMain(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errInternal := QueryResponseInternal{-32603, "Error Processing Request"}
		errRes := QueryResponseError{"2.0", -1, errInternal}
		c.JSON(
			http.StatusOK,
			errRes,
		)
	}
	query := QueryParams{}
	json.Unmarshal(jsonData, &query)
	if strings.HasPrefix(query.Method, "contracts.") {
		rpcClient, _ := jrc.NewServer(GetNodeAddress(LiveState))
		jr2query := jrc.RpcRequest{Method: query.Method, JsonRpc: "2.0", Id: query.Id, Params: query.Params}
		resp, err := rpcClient.Exec(jr2query)
		if err != nil {
			fmt.Print(err)
			c.JSON(
				http.StatusBadGateway,
				gin.H{
					"error": err,
				},
			)
			return
		}
		c.JSON(
			http.StatusOK,
			resp,
		)
	} else if strings.HasPrefix(query.Method, "blockchain.") {
		rpcClient, _ := jrc.NewServer(GetNodeAddress(FullTransactionHistory))
		jr2query := jrc.RpcRequest{Method: query.Method, JsonRpc: "2.0", Id: query.Id, Params: query.Params}
		resp, _ := rpcClient.Exec(jr2query)
		if err != nil {
			fmt.Print(err)
			c.JSON(
				http.StatusBadGateway,
				gin.H{
					"error": err,
				},
			)
			return
		}
		c.JSON(
			http.StatusOK,
			resp,
		)
	} else if strings.HasPrefix(query.Method, "history.") {
		var params map[string]interface{}
		json.Unmarshal(query.Params, &params)
		var paramsString = ""
		for key, value := range params {
			if paramsString != "" {
				paramsString += "&"
			} else {
				paramsString += "?"
			}
			paramsString += key + "=" + url.QueryEscape(value.(string))
		}

		var endpoint = strings.Split(query.Method, ".")[1]
		resp, err := http.Get(GetNodeAddress(AccountHistory) + endpoint + paramsString) //TODO: parse out exactly which histroy endpoint they were wanting
		if err != nil {
			log.Fatalln(err)
		}
		body, _ := ioutil.ReadAll(resp.Body)
		var badBad interface{}
		json.Unmarshal(body, &badBad)
		c.JSON(
			http.StatusOK,
			badBad,
		)
	}
}

func handleContracts(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		// Handle error
	}
	query := QueryParams{}
	json.Unmarshal(jsonData, &query)
	rpcClient, _ := jrc.NewServer(GetNodeAddress(LiveState))
	jr2query := jrc.RpcRequest{Method: "contracts." + query.Method, JsonRpc: "2.0", Id: query.Id, Params: query.Params}
	resp, _ := rpcClient.Exec(jr2query)
	c.JSON(
		http.StatusOK,
		resp,
	)
}

func handleBlockchain(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		// Handle error
	}
	query := QueryParams{}
	json.Unmarshal(jsonData, &query)
	rpcClient, _ := jrc.NewServer(GetNodeAddress(FullTransactionHistory))
	jr2query := jrc.RpcRequest{Method: "blockchain." + query.Method, JsonRpc: "2.0", Id: query.Id, Params: query.Params}
	resp, _ := rpcClient.Exec(jr2query)
	c.JSON(
		http.StatusOK,
		resp,
	)
}

func handleHistory(c *gin.Context) {
	var paramsString = ""
	var query = c.Request.URL.Query()
	for key, value := range query {
		if paramsString != "" {
			paramsString += "&"
		} else {
			paramsString += "?"
		}
		paramsString += key + "=" + strings.Join(value[:], ",")
	}
	resp, err := http.Get(GetNodeAddress(AccountHistory) + c.Request.URL.Path[1:] + paramsString)
	if err != nil {
		log.Fatalln(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var badBad interface{}
	json.Unmarshal(body, &badBad)
	c.JSON(
		http.StatusOK,
		badBad,
	)
}
