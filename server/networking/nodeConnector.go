package networking

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strconv"
	"subframe/server/database"
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
		return sendCoordinatorNodeRequest(address, queryString)
	}
	return ""
}

func sendStorageNodeRequest(address string, queryString string, data string) (response string) {
	var resp *http.Response
	var err error
	if data == "" {
		//There is no data to be POSTed, send GET Request
		resp, err = http.Get(address + "/storage" + queryString)

	} else {
		//There is data to be POSTed, send POST Request
		resp, err = http.Post(address+"/storage"+queryString, "raw", bytes.NewBufferString(data))
	}
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

func sendCoordinatorNodeRequest(address string, queryString string) (response string) {
	//TODO: Send Request, get response; if in coordinator network send request via socket
	resp, err := http.Get(address + "/coordinator" + queryString)
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

//Ping returns the current Ping to the specified address
func Ping(address string) (ping int) {
	//TODO: Get Ping of Node
	return 123
}

//GetMessageStatus queries the CoordinatorNetwork for the status of the specified message
func GetMessageStatus(messageID string) (status int) {
	//If Message is not present in local database, no need to check status
	if !database.CheckMessageStorage(messageID) {
		return
	}

	//Get Status from up to three different coordinator nodes
	coordinatorNodes := database.GetRandomCoordinatorNodes(3)
	newStatus := make([]string, len(coordinatorNodes))
	for index, value := range coordinatorNodes {
		newStatus[index] = sendCoordinatorNodeRequest(value, "/status/"+messageID)
	}

	for _, value := range newStatus {
		if value != newStatus[0] {
			//TODO: Network is out of sync; handle appropriately
			return -1
		}
	}

	//Network is in sync, return status
	status, err := strconv.Atoi(newStatus[0])
	if err == nil {
		return status
	}
	return -1
}
