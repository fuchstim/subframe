package bootstrapper

import (
	"encoding/json"
	"strconv"
	"subframe/server/database"
	"subframe/server/jobqueue"
	"subframe/server/logger"
	"subframe/server/networking"
	"subframe/server/settings"
)

var log = logger.Logger{Prefix: "bootstrapper/Main"}

//Bootstrap bootstraps the local node, if settings.BootstrapNode is set
func Bootstrap() {
	bootstrapNode := settings.BootstrapNode
	if bootstrapNode == "" {
		log.Info("No BootstrapNode set. Skipping Bootstrapping.")
		return
	}

	log.Info("Bootstrapping with Node " + settings.BootstrapNode + "...")

	if database.ClearNodeTables() != 200 {
		log.Fatal("Could not clear databases before bootstrapping.")
	}
	pullStorageNodes()
	pullCoordinatorNodes()
}

func pullStorageNodes() {
	log.Info("Pulling StorageNodes...")
	response := networking.SendNodeRequest(networking.NODE_STORAGE, settings.BootstrapNode, "/control/get-storage-nodes", "")
	var storageNodes []string
	err := json.Unmarshal([]byte(response), &storageNodes)

	if response == "" || err != nil {
		if response == "" {
			log.Fatal("Error getting StorageNodes: Empty Response")
		} else {
			log.Fatal("Error getting StorageNodes: " + err.Error())
		}
	}

	task := func(data interface{}) {
		log := logger.Logger{Prefix: "bootstrapper/DatabaseThread-StorageNodes"}
		storageNodes, ok := data.([]string)
		if !ok {
			log.Error("Failed to add StorageNodes to Database")
			return
		}
		for _, node := range storageNodes {
			ping := networking.Ping(node)
			database.AddStorageNode(node, ping)
			log.Info("Added StorageNode " + node + " with Ping " + strconv.Itoa(ping) + " to Database")
		}
	}
	job := jobqueue.Job{
		Task: task,
		Data: storageNodes,
	}
	select {
	case jobqueue.Queue <- job:
	}
	log.Info("Pulled StorageNodes.")
}

func pullCoordinatorNodes() {
	log.Info("Pulling CoordinatorNodes...")
	response := networking.SendNodeRequest(networking.NODE_STORAGE, settings.BootstrapNode, "/control/get-coordinator-nodes", "")
	var coordinatorNodes []string
	err := json.Unmarshal([]byte(response), &coordinatorNodes)

	if response == "" || err != nil {
		if response == "" {
			log.Fatal("Error getting CoordinatorNodes: Empty Response")
		} else {
			log.Fatal("Error getting CoordinatorNodes: " + err.Error())
		}
	}

	task := func(data interface{}) {
		log := logger.Logger{Prefix: "bootstrapper/DatabaseThread-CoordinatorNodes"}
		coordinatorNodes, ok := data.([]string)
		if !ok {
			log.Error("Failed to add CoordinatorNodes to Database")
			return
		}
		for _, node := range coordinatorNodes {
			ping := networking.Ping(node)
			database.AddCoordinatorNode(node, ping)
			log.Info("Added CoordinatorNode " + node + " with Ping " + strconv.Itoa(ping) + " to Database")
		}
	}
	job := jobqueue.Job{
		Task: task,
		Data: coordinatorNodes,
	}
	select {
	case jobqueue.Queue <- job:
	}
	log.Info("Pulled CoordinatorNodes.")
}
