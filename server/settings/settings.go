package settings

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
)

//BootstrapNode is used for Bootstrapping the local instance
var BootstrapNode string = ""

//DataPath is used to store message and database files
var DataPath string = "./data"

//RemoteAddress is used to access the local instance remotely
var RemoteAddress string = "localhost:9123"

//LocalAddress is the IP and Port the StorageNode instance listens on
var LocalAddress string = "0.0.0.0:9123"

//DiskSpace is the maximum space used for message storage
var DiskSpace int = 5000

//MaxWorkers is the maximum number of workers to spawn
var MaxWorkers int = 10

//QueueMaxLength is the maximum length of a queue before a new worker is spawned, if the current worker count does not exceed MaxWorkers
var QueueMaxLength int = 10

//MessageMaxSize defines the maximum size of an individual message file
var MessageMaxSize int64 = 100

//MessageMinCheckDelay defines the minimum time in hours between individual checks of the message status
var MessageMinCheckDelay int = 12

//MessageMaxStoreTime defines the maximum time a message is stored locally, in days
var MessageMaxStoreTime int = 7

//Read reads settings from local storage and overwrites them with command-line-arguments
func Read() {
	println("Reading Settings...")

	jsonstring, err := ioutil.ReadFile(DataPath + "/settings.json")
	if err == nil {
		data := make(map[string]interface{})
		err := json.Unmarshal(jsonstring, &data)
		if err == nil {
			RemoteAddress = data["RemoteAddress"].(string)
			LocalAddress = data["LocalAddress"].(string)
			DiskSpace = int(data["DiskSpace"].(float64))
			MaxWorkers = int(data["MaxWorkers"].(float64))
			QueueMaxLength = int(data["QueueMaxLength"].(float64))
			MessageMaxSize = int64(data["MessageMaxSize"].(float64))
			MessageMinCheckDelay = int(data["MessageMinCheckDelay"].(float64))
			MessageMaxStoreTime = int(data["MessageMaxStoreTime"].(float64))
		}
	}

	parseCommandLineArgs()
	println("Read Settings.")
	Write()
}

//Write writes settings to local storage
func Write() {
	//Write settings to disk
	data := make(map[string]interface{})
	data["RemoteAddress"] = RemoteAddress
	data["LocalAddress"] = LocalAddress
	data["DiskSpace"] = DiskSpace
	data["MaxWorkers"] = MaxWorkers
	data["QueueMaxLength"] = QueueMaxLength
	data["MessageMaxSize"] = MessageMaxSize
	data["MessageMinCheckDelay"] = MessageMinCheckDelay
	data["MessageMaxStoreTime"] = MessageMaxStoreTime

	jsonstring, err := json.MarshalIndent(data, "", "\t")
	f, err := os.Create(DataPath + "/settings.json")
	if err != nil {
		log.Panic(err)
	}
	f.Write(jsonstring)
	f.Close()
}

func parseCommandLineArgs() {
	flag.StringVar(&BootstrapNode, "bootstrap-node", BootstrapNode, "If set, SuBFraMe will reinitialize the local Node Database and sync it with the BootstrapNode")
	flag.StringVar(&DataPath, "data-dir", DataPath, "The SuBFraMe data directory, messages, databases and settings will be stored here")
	flag.StringVar(&RemoteAddress, "remote-address", RemoteAddress, "The remote address of this SuBFraMe Instance")
	flag.StringVar(&LocalAddress, "local-address", LocalAddress, "The IP and Port the Node Interface will listen on")
	flag.IntVar(&DiskSpace, "disk-space", DiskSpace, "The maximum space SuBFraMe will use to store Messages in MB")
	flag.IntVar(&MaxWorkers, "max-workers", MaxWorkers, "The maximum number of worker threads")
	flag.IntVar(&QueueMaxLength, "max-queue-length", QueueMaxLength, "The maximum size a queue can have before a new worker is spawned, before exceeding max-workers")
	flag.Int64Var(&MessageMaxSize, "message-max-size", MessageMaxSize, "The maximum size of an individual message file, in MB")
	flag.IntVar(&MessageMinCheckDelay, "message-min-check-delay", MessageMinCheckDelay, "The minimum time in hours between individual checks of the same message against the coordinator network")
	flag.IntVar(&MessageMaxStoreTime, "message-max-store-time", MessageMaxStoreTime, "The maximum time a message is stored locally, in days")
	flag.Parse()
}
