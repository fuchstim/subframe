package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"subframe/server/settings"
)

var storageDB sql.DB
var coordinatorDB sql.DB

//Init initializes and / or opens the required SuBFraMe Databases
func Init() {
	//Initialize and / or open sqlite databases
	println("Initializing database connections...")
	storageDBPath := settings.DataPath + "/databases/storage.db"
	storageDB, err := sql.Open("sqlite3", storageDBPath)

	storageDB.Close()

	println("Initialized database connections.")
}

//Close closes all Database connections
func Close() {
	println("Closing Database connections...")

	println("Closed database connections.")
}

//LogMessageStorage logs to the StorageNode Database that a message has been received and stored locally
func LogMessageStorage(id string) (status int) {

	return http.StatusOK
}

//CheckMessageStorage checks whether a message is is present in the local database
func CheckMessageStorage(id string) (hasMessage bool) {
	//Check whether database has message id
	return false
}

//CheckMessageStatusStorage checks the status of a locally stored message against the Coordinator Network and handles it respectively
func CheckMessageStatusStorage(id string) {
	//Check status of message against coordinator network, then delete or keep message and log time of last check
}

//CheckDueMessageStatusStorage checks the status of all messages checked more that settings.MessageMinCheckDelay ago, removes them if they exceed settings.MessageMaxStoreTime or have been received
func CheckDueMessageStatusStorage() {
	//Run CheckMessageStatusStorage on all messages which have been checked more than settings.MessageMinCheckDelay ago,
	//remove all messages which have been received before now - settings.MessageMaxStoreTime
}
