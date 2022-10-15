package main

import (
	"encoding/json"
	"fmt"
	"github.com/cfoxon/jrc"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type QueryParams struct {
	Jsonrpc string          `json:"jsonrpc"`
	Id      int             `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

func Router(r *gin.Engine) {
	r.POST("/", handleMain)
	//r.POST("/history")
	r.POST("/contracts", handleContracts)
	r.POST("/blockchain", handleBlockchain)
}

func handleMain(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		// Handle error
	}
	query := QueryParams{}
	json.Unmarshal(jsonData, &query)
	if strings.HasPrefix(query.Method, "contracts.") || strings.HasPrefix(query.Method, "blockchain.") {
		rpcClient, _ := jrc.NewServer("https://engine.rishipanthee.com")
		jr2query := jrc.RpcRequest{Method: query.Method, JsonRpc: "2.0", Id: query.Id, Params: query.Params}
		resp, _ := rpcClient.Exec(jr2query)
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
		fmt.Println("https://enginehistory.rishipanthee.com/" + endpoint + paramsString)
		resp, err := http.Get("https://enginehistory.rishipanthee.com/" + endpoint + paramsString)
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
	rpcClient, _ := jrc.NewServer("https://engine.rishipanthee.com")
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
	rpcClient, _ := jrc.NewServer("https://engine.rishipanthee.com")
	jr2query := jrc.RpcRequest{Method: "blockchain." + query.Method, JsonRpc: "2.0", Id: query.Id, Params: query.Params}
	resp, _ := rpcClient.Exec(jr2query)
	c.JSON(
		http.StatusOK,
		resp,
	)
}
