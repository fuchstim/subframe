package database

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"subframe/server/settings"

	//Importing SQLite Driver
	_ "github.com/mattn/go-sqlite3"
)

var storageDB *sql.DB
var coordinatorDB *sql.DB

//Init initializes and / or opens the required SuBFraMe Databases
func Init() {
	//Initialize and / or open sqlite databases
	println("Initializing database connections...")
	storageDBPath := settings.DataPath + "/databases/storage.db"
	coordinatorDBPath := settings.DataPath + "/databases/coordinator.db"

	var err error
	storageDB, err = sql.Open("sqlite3", storageDBPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	coordinatorDB, err = sql.Open("sqlite3", coordinatorDBPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	//Create Tables for storageDatabase
	statement := `
	CREATE TABLE IF NOT EXISTS messages(id varchar(255) not null primary key, receivedOn timestamp not null, expiresOn timestamp not null, lastCheck timestamp);
	`
	_, err = storageDB.Exec(statement)
	if err != nil {
		log.Fatal(err)
		return
	}

	//Create Tables for coordinatorDatabase
	statement = `
	CREATE TABLE IF NOT EXISTS storageNodes(address varchar(255) not null primary key);
	CREATE TABLE IF NOT EXISTS coordinatorNodes(address varchar(255) not null primary key);
	CREATE TABLE IF NOT EXISTS messages(id varchar(255) not null, storageNode varchar(255) not null, reportedOn timestamp not null, verifiedOn timestamp);
	`
	_, err = coordinatorDB.Exec(statement)
	if err != nil {
		log.Fatal(err)
		return
	}

	println("Initialized database connections.")
}

//Close closes all Database connections
func Close() {
	println("Closing Database connections...")
	storageDB.Close()
	coordinatorDB.Close()
	println("Closed database connections.")
}

//LogMessageStorage logs to the StorageNode Database that a message has been received and stored locally
func LogMessageStorage(id string) (status int) {
	if CheckMessageStorage(id) {
		return http.StatusConflict
	}

	query := "INSERT INTO messages(id, receivedOn, expiresOn) VALUES (?, date('now'), date('now', '+' || ? || ' days'))"
	stmt, err := storageDB.Prepare(query)
	if err != nil {
		return http.StatusInternalServerError
	}
	defer stmt.Close()
	_, err = stmt.Exec(id, settings.MessageMaxStoreTime)
	if err != nil {
		return http.StatusInternalServerError
	}
	return http.StatusOK
}

//CheckMessageStorage checks whether a message is is present in the local database
func CheckMessageStorage(id string) (hasMessage bool) {
	query := "SELECT id FROM messages WHERE id=?"
	stmt, err := storageDB.Prepare(query)
	if err != nil {
		return false
	}
	defer stmt.Close()

	var res string
	err = stmt.QueryRow(id).Scan(&res)
	if err != nil {
		return false
	}

	return true
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

//GetStorageNodes returns known StorageNodes
func GetStorageNodes(limit int) (storageNodes []string) {
	var nodes []string
	query := "SELECT address FROM storageNodes LIMIT " + strconv.Itoa(limit)
	rows, err := coordinatorDB.Query(query)
	if err != nil {
		return nodes
	}
	defer rows.Close()
	for rows.Next() {
		var address string
		err = rows.Scan(&address)
		if err != nil {
			return nodes
		}
		nodes = append(nodes, address)
	}
	return nodes
}
