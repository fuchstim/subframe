package settings

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"subframe/server/logger"
	. "subframe/status"
)

var log = logger.Logger{Prefix: "settings/Main"}

//BootstrapNode is used for Bootstrapping the local instance
var BootstrapNode = ""

//DataPath is used to store message and database files
var DataPath = "./data"

//RemoteAddress is used to access the local instance remotely
var RemoteAddress = "localhost:9123"

//LocalAddress is the IP and Port the StorageNode instance listens on
var LocalAddress = "0.0.0.0:9123"

//DiskSpace is the maximum space used for message storage
var DiskSpace = 5000

//MaxWorkers is the maximum number of workers to spawn
var MaxWorkers = 10

//QueueMaxLength is the maximum length of a queue before a new worker is spawned, if the current worker count does not exceed MaxWorkers
var QueueMaxLength = 10

//MessageMaxSize defines the maximum size of an individual message file
var MessageMaxSize = 100

//MessageMinCheckDelay defines the minimum time in hours between individual checks of the message status
var MessageMinCheckDelay = 12

//MessageMaxStoreTime defines the maximum time a message is stored locally, in days
var MessageMaxStoreTime = 7

//ColorizedOutput defines whether realtime logs should be colorized
var ColorizedLogs = false

//Read reads settings from local storage and overwrites them with command-line-arguments
func Read() {
	log.Info(InProgress, "Reading Settings...")
	jsonstring, err := ioutil.ReadFile(DataPath + "/settings.json")
	if err == nil {
		data := make(map[string]interface{})
		err := json.Unmarshal(jsonstring, &data)
		if err == nil {
			RemoteAddress, _ = data["RemoteAddress"].(string)

			LocalAddress, _ = data["LocalAddress"].(string)

			tmp, ok := data["DiskSpace"].(float64)
			if ok {
				DiskSpace = int(tmp)
			}

			tmp, ok = data["MaxWorkers"].(float64)
			if ok {
				MaxWorkers = int(tmp)
			}

			tmp, ok = data["QueueMaxLength"].(float64)
			if ok {
				QueueMaxLength = int(tmp)
			}

			tmp, ok = data["MessageMaxSize"].(float64)
			if ok {
				MessageMaxSize = int(tmp)
			}

			tmp, ok = data["MessageMinCheckDelay"].(float64)
			if ok {
				MessageMinCheckDelay = int(tmp)
			}

			tmp, ok = data["MessageMaxStoreTime"].(float64)
			if ok {
				MessageMaxStoreTime = int(tmp)
			}

			ColorizedLogs, _ = data["ColorizedLogs"].(bool)
		} else {
			log.Warn(SettingsReadError, "Failed to read settings from file ("+err.Error()+"). Falling back to defaults or using command line arguments...")
		}
	}

	parseCommandLineArgs()
	logger.ColorizedLogs = ColorizedLogs
	log.Info(OK, "Successfully read Settings.")
	Write()
}

//Write writes settings to local storage
func Write() {
	//Write settings to disk
	log.Info(InProgress, "Writing settings...")
	data := make(map[string]interface{})
	data["RemoteAddress"] = RemoteAddress
	data["LocalAddress"] = LocalAddress
	data["DiskSpace"] = DiskSpace
	data["MaxWorkers"] = MaxWorkers
	data["QueueMaxLength"] = QueueMaxLength
	data["MessageMaxSize"] = MessageMaxSize
	data["MessageMinCheckDelay"] = MessageMinCheckDelay
	data["MessageMaxStoreTime"] = MessageMaxStoreTime
	data["ColorizedLogs"] = ColorizedLogs

	jsonstring, err := json.MarshalIndent(data, "", "\t")
	f, err := os.Create(DataPath + "/settings.json")
	if err != nil {
		log.Fatal(SettingsWriteError, err.Error())
	}
	f.Write(jsonstring)
	f.Close()
	log.Info(OK, "Wrote settings to "+DataPath+"/settings.json")
}

func parseCommandLineArgs() {
	log.Info(InProgress, "Parsing Commandline Arguments...")
	flag.StringVar(&BootstrapNode, "bootstrap-node", BootstrapNode, "If set, SuBFraMe will reinitialize the local Node Database and sync it with the BootstrapNode")
	flag.StringVar(&DataPath, "data-dir", DataPath, "The SuBFraMe data directory, messages, databases and settings will be stored here")
	flag.StringVar(&RemoteAddress, "remote-address", RemoteAddress, "The remote address of this SuBFraMe Instance")
	flag.StringVar(&LocalAddress, "local-address", LocalAddress, "The IP and Port the Node Interface will listen on")
	flag.IntVar(&DiskSpace, "disk-space", DiskSpace, "The maximum space SuBFraMe will use to store Messages in MB")
	flag.IntVar(&MaxWorkers, "max-workers", MaxWorkers, "The maximum number of worker threads")
	flag.IntVar(&QueueMaxLength, "max-queue-length", QueueMaxLength, "The maximum size a queue can have before a new worker is spawned, before exceeding max-workers")
	flag.IntVar(&MessageMaxSize, "message-max-size", MessageMaxSize, "The maximum size of an individual message file, in MB")
	flag.IntVar(&MessageMinCheckDelay, "message-min-check-delay", MessageMinCheckDelay, "The minimum time in hours between individual checks of the same message against the coordinator network")
	flag.IntVar(&MessageMaxStoreTime, "message-max-store-time", MessageMaxStoreTime, "The maximum time a message is stored locally, in days")
	flag.BoolVar(&ColorizedLogs, "colorized-output", ColorizedLogs, "Turns on or off colorized realtime logs")
	flag.Parse()
	log.Info(OK, "Parsed Commandline Arguments.")
}
