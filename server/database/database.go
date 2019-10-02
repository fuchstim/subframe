package database

import (
	"database/sql"
	"net/http"
	"strconv"
	"subframe/server/logger"
	"subframe/server/settings"

	//Importing SQLite Driver
	_ "github.com/mattn/go-sqlite3"
)

var storageDB *sql.DB
var coordinatorDB *sql.DB
var log = logger.Logger{Prefix: "database/Main"}

//Init initializes and / or opens the required SuBFraMe Databases
func Init() {
	//Initialize and / or open sqlite databases
	log.Info("Opening Database Files...")
	storageDBPath := settings.DataPath + "/databases/storage.db"
	coordinatorDBPath := settings.DataPath + "/databases/coordinator.db"

	var err error
	storageDB, err = sql.Open("sqlite3", storageDBPath)
	if err != nil {
		log.Fatal("Error opening StorageDatabase: " + err.Error())
		return
	}

	coordinatorDB, err = sql.Open("sqlite3", coordinatorDBPath)
	if err != nil {
		log.Fatal("Error opening CoordinatorDatabase: " + err.Error())
		return
	}

	log.Info("Successfully openened Database Files. Initializing Table structure if not present...")

	//Create Tables for storageDatabase
	log.Info("Creating tables for StorageDatabase...")
	statement := `
	CREATE TABLE IF NOT EXISTS messages(
		id varchar(255) not null primary key, 
		verified tinyint not null default 0,
		expiresOn timestamp not null, 
		lastCheck timestamp
	);
	`
	_, err = storageDB.Exec(statement)
	if err != nil {
		log.Fatal("Failed to create Tables for StorageDatabase: " + err.Error())
		return
	}

	log.Info("Created Tables for StorageDatabase.")

	//Create Tables for coordinatorDatabase
	log.Info("Creating tables for CoordinatorDatabase...")
	statement = `
	CREATE TABLE IF NOT EXISTS storageNodes(
		address varchar(255) not null primary key, 
		lastPing int not null
	);
	CREATE TABLE IF NOT EXISTS coordinatorNodes(
		address varchar(255) not null primary key, 
		lastPing int not null
	);
	CREATE TABLE IF NOT EXISTS messages(
		id varchar(255) not null, 
		storageNode varchar(255) not null, 
		reportedOn timestamp not null, 
		verified tinyint not null default 0
	);
	`
	_, err = coordinatorDB.Exec(statement)
	if err != nil {
		log.Fatal("Failed to create Tables for CoordinatorDatabase: " + err.Error())
		return
	}

	log.Info("Created Tables for CoordinatorDatabase.")
	log.Info("Initialized database connections.")
}

//Close closes all Database connections
func Close() {
	log.Info("Closing Database connections...")
	storageDB.Close()
	coordinatorDB.Close()
	log.Info("Closed database connections.")
}

//LogMessageStorage logs to the StorageNode Database that a message has been received and stored locally
func LogMessageStorage(id string) (status int) {
	log.Info("Logging new Message " + id + "...")
	if CheckMessageStorage(id) {
		log.Error("Message " + id + " already present in Database.")
		return http.StatusConflict
	}

	query := "INSERT INTO messages(id, expiresOn) VALUES (?, date('now', '+' || ? || ' days'))"
	stmt, err := storageDB.Prepare(query)
	if err != nil {
		log.Error("Error logging Message " + id + " to Database: " + err.Error())
		return http.StatusInternalServerError
	}
	defer stmt.Close()
	_, err = stmt.Exec(id, settings.MessageMaxStoreTime)
	if err != nil {
		log.Error("Error loggin Message " + id + " to Database: " + err.Error())
		return http.StatusInternalServerError
	}
	log.Info("Successfully logged Message " + id + " to Database.")
	return http.StatusOK
}

//CheckMessageStorage checks whether a message is is present in the local database
func CheckMessageStorage(id string) (hasMessage bool) {
	log.Info("Checking whether Message " + id + " is in Database...")
	query := "SELECT id FROM messages WHERE id=?"
	stmt, err := storageDB.Prepare(query)
	if err != nil {
		log.Error("Error: " + err.Error())
		return false
	}
	defer stmt.Close()

	var res string
	err = stmt.QueryRow(id).Scan(&res)
	if err != nil {
		log.Info("Message " + id + " does not appear to be present in database.")
		return false
	}

	log.Info("Message " + id + " is present in database.")
	return true
}

//CheckMessageStatusStorage checks the status of a locally stored message against the Coordinator Network and handles it respectively
func CheckMessageStatusStorage(id string) {
	//TODO: Check status of message against coordinator network, then delete or keep message and log time of last check
}

//CheckDueMessageStatusStorage checks the status of all messages checked more that settings.MessageMinCheckDelay ago, removes them if they exceed settings.MessageMaxStoreTime or have been received
func CheckDueMessageStatusStorage() {
	//TODO: Run CheckMessageStatusStorage on all messages which have been checked more than settings.MessageMinCheckDelay ago,
	//remove all messages which have been received before now - settings.MessageMaxStoreTime
}

//AddStorageNode adds a StorageNode to the local database
func AddStorageNode(address string, ping int) (status int) {
	log.Info("Adding StorageNode " + address + " to database...")
	query := "INSERT INTO storageNodes(address, lastPing) VALUES (?,?)"
	stmt, err := coordinatorDB.Prepare(query)
	if err != nil {
		log.Error("Error adding StorageNode " + address + " to database: " + err.Error())
		return http.StatusInternalServerError
	}
	defer stmt.Close()
	_, err = stmt.Exec(address, ping)
	if err != nil {
		log.Error("Error adding StorageNode " + address + " to database: " + err.Error())
		return http.StatusInternalServerError
	}
	log.Info("Added StorageNode " + address + " to Database.")
	return http.StatusOK
}

//GetStorageNodes returns known StorageNodes
func GetStorageNodes(limit int) (storageNodes []string) {
	log.Info("Exporting " + strconv.Itoa(limit) + " StorageNodes...")
	var nodes []string
	query := "SELECT address FROM storageNodes LIMIT " + strconv.Itoa(limit)
	rows, err := coordinatorDB.Query(query)
	if err != nil {
		log.Error("Error exporting StorageNodes: " + err.Error())
		return nodes
	}
	defer rows.Close()
	for rows.Next() {
		var address string
		err = rows.Scan(&address)
		if err != nil {
			continue
		}
		nodes = append(nodes, address)
	}
	log.Info("Returning " + strconv.Itoa(len(nodes)) + " StorageNodes.")
	return nodes
}

//AddCoordinatorNode adds a CoordinatorNode to the local database
func AddCoordinatorNode(address string, ping int) (status int) {
	log.Info("Adding CoordinatorNode " + address + " to database...")
	query := "INSERT INTO coordinatorNodes(address, lastPing) VALUES (?,?)"
	stmt, err := coordinatorDB.Prepare(query)
	if err != nil {
		log.Error("Error adding CoordinatorNode " + address + " to database: " + err.Error())
		return http.StatusInternalServerError
	}
	defer stmt.Close()
	_, err = stmt.Exec(address, ping)
	if err != nil {
		log.Error("Error adding CoordinatorNode " + address + " to database: " + err.Error())
		return http.StatusInternalServerError
	}
	log.Info("Added CoordinatorNode " + address + " to Database.")
	return http.StatusOK
}

//GetCoordinatorNodes returns known CoordinatorNodes
func GetCoordinatorNodes() (storageNodes []string) {
	log.Info("Exporting CoordinatorNodes...")
	var nodes []string
	query := "SELECT address FROM coordinatorNodes"
	rows, err := coordinatorDB.Query(query)
	if err != nil {
		log.Error("Error exporting CoordinatorNodes: " + err.Error())
		return nodes
	}
	defer rows.Close()
	for rows.Next() {
		var address string
		err = rows.Scan(&address)
		if err != nil {
			continue
		}
		nodes = append(nodes, address)
	}
	log.Info("Returning " + strconv.Itoa(len(nodes)) + " CoordinatorNodes.")
	return nodes
}

//GetRandomCoordinatorNodes returns max <number> random CoordinatorNodes
func GetRandomCoordinatorNodes(max int) (nodes []string) {
	log.Info("Getting " + strconv.Itoa(max) + " random CoordinatorNodes...")
	//TODO: Return random CoordinatorNodes

	result := []string{"node1", "node2", "node3"}

	log.Info("Returning " + strconv.Itoa(len(result)) + " CoordinatorNodes.")
	return result
}

//ClearNodeTables removes all elements from storageNodes and coordinatorNodes tables, for bootstrapping
func ClearNodeTables() (status int) {
	log.Info("Clearing Node Tables...")
	query := "DELETE FROM storageNodes; DELETE FROM coordinatorNodes"
	_, err := coordinatorDB.Exec(query)
	if err != nil {
		log.Error("Error clearing Node Tables: " + err.Error())
		return http.StatusInternalServerError
	}
	log.Info("Cleared Node Tables.")
	return http.StatusOK
}

//UpdateMessageStatusStorage updates the status of a message in the local database
func UpdateMessageStatusStorage(messageID string, status int) {
	log.Info("Updating Status of Message " + messageID)
	query := "UPDATE messages SET verified=? WHERE id=?"
	stmt, err := storageDB.Prepare(query)
	if err != nil {
		log.Error("Error updating status of message " + messageID + ": " + err.Error())
		return
	}
	defer stmt.Close()

	log.Info("Updated status of Message " + messageID + ". New status: " + strconv.Itoa(status))
	stmt.Exec(status, messageID)
}
