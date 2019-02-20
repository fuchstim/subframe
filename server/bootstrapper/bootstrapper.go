package bootstrapper

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"subframe/server/database"
	"subframe/server/jobqueue"
	"subframe/server/networking"
	"subframe/server/settings"
)

//Bootstrap bootstraps the local node, if settings.BootstrapNode is set
func Bootstrap() {
	bootstrapNode := settings.BootstrapNode
	if bootstrapNode == "" {
		return
	}

	if database.ClearNodeTables() != 200 {
		log.Fatal("Could not clear databases before bootstrapping. Exiting...")
	}
	pullStorageNodes()
	pullCoordinatorNodes()
}

func pullStorageNodes() {
	response := networking.SendNodeRequest(networking.NODE_STORAGE, settings.BootstrapNode, "/control/get-storage-nodes", "")
	var storageNodes []string
	err := json.Unmarshal([]byte(response), &storageNodes)

	if response == "" || err != nil {
		log.Fatal("No StorageNodes received. Exiting...")
	}

	task := func(data interface{}) {
		storageNodes, ok := data.([]string)
		if !ok {
			return
		}
		for _, node := range storageNodes {
			ping := networking.Ping(node)
			database.AddStorageNode(node, ping)
			fmt.Println("Added StorageNode " + node + " with Ping " + strconv.Itoa(ping) + " to Database")
		}
	}
	job := jobqueue.Job{
		Task: task,
		Data: storageNodes,
	}
	select {
	case jobqueue.Queue <- job:
	}
}

func pullCoordinatorNodes() {
	response := networking.SendNodeRequest(networking.NODE_STORAGE, settings.BootstrapNode, "/control/get-coordinator-nodes", "")
	var coordinatorNodes []string
	err := json.Unmarshal([]byte(response), &coordinatorNodes)

	if response == "" || err != nil {
		log.Fatal("No CoordinatorNodes received. Exiting...")
	}

	task := func(data interface{}) {
		coordinatorNodes, ok := data.([]string)
		if !ok {
			return
		}
		for _, node := range coordinatorNodes {
			ping := networking.Ping(node)
			database.AddCoordinatorNode(node, ping)
			fmt.Println("Added CoordinatorNode " + node + " with Ping " + strconv.Itoa(ping) + " to Database")
		}
	}
	job := jobqueue.Job{
		Task: task,
		Data: coordinatorNodes,
	}
	select {
	case jobqueue.Queue <- job:
	}
}
