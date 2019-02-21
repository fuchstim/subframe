package networking

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strconv"
	"subframe/server/database"
	"subframe/server/logger"
)

var nlog = logger.Logger{Prefix: "networking/NodeConnector"}

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
		nlog.Info("Sending StorageNode GET Request to " + address + "/storage" + queryString + "...")
		resp, err = http.Get(address + "/storage" + queryString)

	} else {
		//There is data to be POSTed, send POST Request
		nlog.Info("Sending StorageNode POST Request to " + address + "/storage" + queryString + "...")
		resp, err = http.Post(address+"/storage"+queryString, "raw", bytes.NewBufferString(data))
	}
	if err != nil {
		nlog.Error("Error sending request: " + err.Error())
		return ""
	}
	defer resp.Body.Close()

	nlog.Info("Reading response...")
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		nlog.Error("Error reading response: " + err.Error())
		return ""
	}

	nlog.Info("Read response.")
	return string(body)
}

func sendCoordinatorNodeRequest(address string, queryString string) (response string) {
	//TODO: Send Request, get response; if in coordinator network send request via socket
	nlog.Info("Sending CoordinatorNode HTTP Request to " + address + "/coordinator" + queryString + "...")
	resp, err := http.Get(address + "/coordinator" + queryString)
	if err != nil {
		nlog.Error("Error sending request: " + err.Error())
		return ""
	}
	defer resp.Body.Close()

	nlog.Info("Reading response...")
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		nlog.Error("Error reading response: " + err.Error())
		return ""
	}

	nlog.Info("Read response")
	return string(body)
}

//Ping returns the current Ping to the specified address
func Ping(address string) (ping int) {
	//TODO: Get Ping of Node
	nlog.Info("Pinging Node " + address)

	ping = 123

	nlog.Info("Ping test for " + address + " returned: " + strconv.Itoa(ping))
	return ping
}

//GetMessageStatus queries the CoordinatorNetwork for the status of the specified message
func GetMessageStatus(messageID string) (status int) {
	nlog.Info("Getting Status for Message " + messageID + " from CoordinatorNetwork...")
	//If Message is not present in local database, no need to check status
	if !database.CheckMessageStorage(messageID) {
		nlog.Error("Message " + messageID + " does not appear to be stored on this Node.")
		return
	}

	//Get Status from up to three different coordinator nodes
	nlog.Info("Getting CoordinatorNodes...")
	coordinatorNodes := database.GetRandomCoordinatorNodes(3)
	if len(coordinatorNodes) == 0 {
		nlog.Error("Did not receive any CoordinatorNodes.")
		return
	}
	nlog.Info("Got " + strconv.Itoa(len(coordinatorNodes)) + " CoordinatorNodes.")
	newStatus := make([]string, len(coordinatorNodes))
	for index, value := range coordinatorNodes {
		newStatus[index] = sendCoordinatorNodeRequest(value, "/status/"+messageID)
	}

	nlog.Info("Got status from " + strconv.Itoa(len(coordinatorNodes)) + " Nodes. Checking...")
	for _, value := range newStatus {
		if value != newStatus[0] {
			//TODO: Network is out of sync; handle appropriately
			nlog.Error("Status do not match. CoordinatorNetwork appears out of sync.")
			return -1
		}
	}

	//Network is in sync, return status
	nlog.Info("New Status appear valid. Returning.")
	status, err := strconv.Atoi(newStatus[0])
	if err == nil {
		return status
	}
	nlog.Error("Error returning new Status: " + err.Error())
	return -1
}
