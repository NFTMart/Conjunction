Proof Of Concept

Config: Hardcode in 1 endpoint for node, and one for history.

Main: Take in HTTP requests in JSONRPC2.0 specs. if url path is `/` If method starts with contracts or blockchain, send the requests to a node processor. If the requests start with history send them to a history processor. Else return method not found json spec error. If url path is `/contracts` scrape off /contracts and add contracts. to the start of the method and send to node processor. If url path is `/blockchain` scrape off /blockchain and add blockchain. to the start of the method and send to node processor. If url path is one of the history ones, forward the raw request to history node.


Node processor: Send raw call to node for it to process. Node needs to be running https://github.com/hive-engine/hivesmartcontracts/pull/5 this pr at a minimum to work. 

History processor: Modify call to change params to url params. Then send the call to the history node.


JS version of modify params: 
```js
let paramsString = "";
for (let key in params) {
    if (paramsString != "") {
        paramsString += "&";
    }
    paramsString += key + "=" + encodeURIComponent(params[key]);
}
```