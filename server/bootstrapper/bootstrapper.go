package bootstrapper

import (
	"encoding/json"
	"strconv"
	"subframe/server/database"
	"subframe/server/jobqueue"
	"subframe/server/logger"
	"subframe/server/networking"
	"subframe/server/settings"
	. "subframe/status"
	"subframe/structs/node"
)

var log = logger.Logger{Prefix: "bootstrapper/Main"}

//Bootstrap bootstraps the local node, if settings.BootstrapNode is set
func Bootstrap() {
	bootstrapNode := settings.BootstrapNode
	if bootstrapNode == "" {
		log.Info(OK, "No BootstrapNode set. Skipping Bootstrapping.")
		return
	}

	log.Info(InProgress, "Bootstrapping with Node "+settings.BootstrapNode+"...")

	if database.ClearNodeTables() != 200 {
		log.Fatal(DBWriteError, "Could not clear databases before bootstrapping.")
	}
	pullStorageNodes()
	pullCoordinatorNodes()
}

func pullStorageNodes() {
	log.Info(InProgress, "Pulling StorageNodes...")
	status, response := networking.SendNodeRequest(networking.NODE_STORAGE, settings.BootstrapNode, "/control/get-storage-nodes", "")
	var storageNodes []node.Node
	err := json.Unmarshal([]byte(response), &storageNodes)

	if status != OK {
		log.Fatal(status, "Error getting StorageNodes")
	}

	if err != nil {
		log.Fatal(GenericInternalError, "Error getting StorageNodes: "+err.Error())
	}

	task := func(data interface{}) {
		log := logger.Logger{Prefix: "bootstrapper/DatabaseThread-StorageNodes"}
		storageNodes, ok := data.([]node.Node)
		if !ok {
			log.Error(GenericInternalError, "Failed to add StorageNodes to Database")
			return
		}
		for _, node := range storageNodes {
			node.Ping = networking.Ping(node.Address)
			database.AddStorageNode(node)
			log.Info(OK, "Added StorageNode "+node.Address+" with Ping "+strconv.Itoa(node.Ping)+" to Database")
		}
	}
	job := jobqueue.Job{
		Task: task,
		Data: storageNodes,
	}
	select {
	case jobqueue.Queue <- job:
	}
	log.Info(OK, "Pulled StorageNodes.")
}

func pullCoordinatorNodes() {
	log.Info(InProgress, "Pulling CoordinatorNodes...")
	status, response := networking.SendNodeRequest(networking.NODE_STORAGE, settings.BootstrapNode, "/control/get-coordinator-nodes", "")
	var coordinatorNodes []node.Node
	err := json.Unmarshal([]byte(response), &coordinatorNodes)

	if status != OK {
		log.Fatal(status, "Error getting Coordinator Nodes")
	}

	if err != nil {
		log.Fatal(GenericInternalError, "Error getting CoordinatorNodes: "+err.Error())
	}

	task := func(data interface{}) {
		log := logger.Logger{Prefix: "bootstrapper/DatabaseThread-CoordinatorNodes"}
		coordinatorNodes, ok := data.([]node.Node)
		if !ok {
			log.Error(DBWriteError, "Failed to add CoordinatorNodes to Database")
			return
		}
		for _, node := range coordinatorNodes {
			node.Ping = networking.Ping(node.Address)
			database.AddCoordinatorNode(node)
			log.Info(OK, "Added CoordinatorNode "+node.Address+" with Ping "+strconv.Itoa(node.Ping)+" to Database")
		}
	}
	job := jobqueue.Job{
		Task: task,
		Data: coordinatorNodes,
	}
	select {
	case jobqueue.Queue <- job:
	}
	log.Info(OK, "Pulled CoordinatorNodes.")
}
