package networking

import (
	"io/ioutil"
	"net/http"
)

//NODE_STORAGE specifies that the request is to be sent to a StorageNode
var NODE_STORAGE = 1

//NODE_COORDINATOR specifies that the request is to be sent to a CoordinatorNode
var NODE_COORDINATOR = 2

//SendNodeRequest sends a synchronous request to the specified node
func SendNodeRequest(nodeType int, address string, queryString string, data string) (response string) {
	switch nodeType {
	case NODE_STORAGE:
		return sendStorageNodeRequest(address, queryString, data)
	case NODE_COORDINATOR:
		return sendCoordinatorNodeRequest(address, queryString, data)
	}
	return ""
}

func sendStorageNodeRequest(address string, queryString string, data string) (response string) {
	//TODO: Send Request, get response
	resp, err := http.Get(address + "/storage" + queryString)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	return string(body)
}

func sendCoordinatorNodeRequest(address string, queryString string, data string) (response string) {
	//TODO: Send Request, get response
	return ""
}

func Ping(address string) (ping int) {
	//TODO: Get Ping of Node
	return 123
}
