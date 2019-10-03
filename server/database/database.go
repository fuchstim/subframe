package database

import (
	"database/sql"
	"strconv"
	"subframe/server/logger"
	"subframe/server/settings"
	. "subframe/status"
	"subframe/structs/node"
	"time"

	//Importing SQLite Driver
	_ "github.com/mattn/go-sqlite3"
)

var storageDB *sql.DB
var coordinatorDB *sql.DB
var log = logger.Logger{Prefix: "database/Main"}

//Init initializes and / or opens the required SuBFraMe Databases
func Init() {
	//Initialize and / or open sqlite databases
	log.Info(InProgress, "Opening Database Files...")
	storageDBPath := settings.DataPath + "/databases/storage.db"
	coordinatorDBPath := settings.DataPath + "/databases/coordinator.db"

	var err error
	storageDB, err = sql.Open("sqlite3", storageDBPath)
	if err != nil {
		log.Fatal(DBOpenError, "Error opening StorageDatabase: "+err.Error())
		return
	}

	coordinatorDB, err = sql.Open("sqlite3", coordinatorDBPath)
	if err != nil {
		log.Fatal(DBOpenError, "Error opening CoordinatorDatabase: "+err.Error())
		return
	}

	log.Info(OK, "Successfully openened Database Files. Initializing Table structure if not present...")

	//Create Tables for storageDatabase
	log.Info(InProgress, "Creating tables for StorageDatabase...")
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
		log.Fatal(DBStructureError, "Failed to create Tables for StorageDatabase: "+err.Error())
		return
	}

	log.Info(OK, "Created Tables for StorageDatabase.")

	//Create Tables for coordinatorDatabase
	log.Info(InProgress, "Creating tables for CoordinatorDatabase...")
	statement = `
	CREATE TABLE IF NOT EXISTS storageNodes(
		address varchar(255) not null primary key, 
		lastPing timestamp not null,
		ping int not null
	);
	CREATE TABLE IF NOT EXISTS coordinatorNodes(
		address varchar(255) not null primary key, 
		lastPing timestamp not null,
		ping int not null           
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
		log.Fatal(DBStructureError, "Failed to create Tables for CoordinatorDatabase: "+err.Error())
		return
	}

	log.Info(OK, "Created Tables for CoordinatorDatabase.")
	log.Info(OK, "Initialized database connections.")
}

//Close closes all Database connections
func Close() {
	log.Info(InProgress, "Closing Database connections...")
	storageDB.Close()
	coordinatorDB.Close()
	log.Info(OK, "Closed database connections.")
}

//LogMessageStorage logs to the StorageNode Database that a message has been received and stored locally
func LogMessageStorage(id string) (status int) {
	log.Info(InProgress, "Logging new Message "+id+"...")
	if _, c := CheckMessageStorage(id); c == true {
		log.Error(SNDBIdConflict, "Message "+id+" already present in Database.")
		return SNDBIdConflict
	}

	query := "INSERT INTO messages(id, expiresOn) VALUES (?, date('now', '+' || ? || ' days'))"
	stmt, err := storageDB.Prepare(query)
	if err != nil {
		log.Error(SNDBPrepareError, "Error logging Message "+id+" to Database: "+err.Error())
		return SNDBPrepareError
	}
	defer stmt.Close()
	_, err = stmt.Exec(id, settings.MessageMaxStoreTime)
	if err != nil {
		log.Error(SNDBWriteError, "Error logging Message "+id+" to Database: "+err.Error())
		return SNDBWriteError
	}
	log.Info(OK, "Successfully logged Message "+id+" to Database.")
	return OK
}

//CheckMessageStorage checks whether a message is is present in the local database
func CheckMessageStorage(id string) (status int, hasMessage bool) {
	log.Info(InProgress, "Checking whether Message "+id+" is in Database...")
	query := "SELECT id FROM messages WHERE id=?"
	stmt, err := storageDB.Prepare(query)
	if err != nil {
		log.Error(SNDBReadError, "Error: "+err.Error())
		return SNDBReadError, false
	}
	defer stmt.Close()

	var res string
	err = stmt.QueryRow(id).Scan(&res)
	if err != nil {
		log.Info(OK, "Message "+id+" does not appear to be present in database.")
		return OK, false
	}

	log.Info(OK, "Message "+id+" is present in database.")
	return OK, true
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
func AddStorageNode(n node.Node) (status int) {
	log.Info(InProgress, "Adding StorageNode "+n.Address+" to database...")
	query := "INSERT INTO storageNodes(address, lastPing, ping) VALUES (?,?)"
	stmt, err := coordinatorDB.Prepare(query)
	if err != nil {
		log.Error(CNDBPrepareError, "Error adding StorageNode "+n.Address+" to database: "+err.Error())
		return CNDBPrepareError
	}
	defer stmt.Close()
	_, err = stmt.Exec(n.Address, n.LastPing.Unix(), n.Ping)
	if err != nil {
		log.Error(CNDBWriteError, "Error adding StorageNode "+n.Address+" to database: "+err.Error())
		return CNDBWriteError
	}
	log.Info(OK, "Added StorageNode "+n.Address+" to Database.")
	return OK
}

//GetStorageNodes returns known StorageNodes
func GetStorageNodes(limit int) (status int, storageNodes []node.Node) {
	log.Info(InProgress, "Exporting "+strconv.Itoa(limit)+" StorageNodes...")
	var nodes []node.Node
	query := "SELECT address, lastPing FROM storageNodes LIMIT " + strconv.Itoa(limit)
	rows, err := coordinatorDB.Query(query)
	if err != nil {
		log.Error(CNDBReadError, "Error exporting StorageNodes: "+err.Error())
		return CNDBReadError, nil
	}
	defer rows.Close()
	for rows.Next() {
		var address string
		var lastPing int64
		err = rows.Scan(&address, &lastPing)
		if err != nil {
			continue
		}
		nodes = append(nodes, node.Node{
			Address: address, LastPing: time.Unix(lastPing, 0),
		})
	}
	log.Info(OK, "Returning "+strconv.Itoa(len(nodes))+" StorageNodes.")
	return OK, nodes
}

//AddCoordinatorNode adds a CoordinatorNode to the local database
func AddCoordinatorNode(n node.Node) (status int) {
	log.Info(InProgress, "Adding CoordinatorNode "+n.Address+" to database...")
	query := "INSERT INTO coordinatorNodes(address, lastPing, ping) VALUES (?,?)"
	stmt, err := coordinatorDB.Prepare(query)
	if err != nil {
		log.Error(CNDBPrepareError, "Error adding CoordinatorNode "+n.Address+" to database: "+err.Error())
		return CNDBPrepareError
	}
	defer stmt.Close()
	_, err = stmt.Exec(n.Address, n.LastPing, n.Ping)
	if err != nil {
		log.Error(CNDBWriteError, "Error adding CoordinatorNode "+n.Address+" to database: "+err.Error())
		return CNDBWriteError
	}
	log.Info(OK, "Added CoordinatorNode "+n.Address+" to Database.")
	return OK
}

//GetCoordinatorNodes returns known CoordinatorNodes
func GetCoordinatorNodes() (status int, storageNodes []node.Node) {
	log.Info(InProgress, "Exporting CoordinatorNodes...")
	var nodes []node.Node
	query := "SELECT address, lastPing FROM coordinatorNodes"
	rows, err := coordinatorDB.Query(query)
	if err != nil {
		log.Error(CNDBReadError, "Error exporting CoordinatorNodes: "+err.Error())
		return CNDBReadError, nil
	}
	defer rows.Close()
	for rows.Next() {
		var address string
		var lastPing int64
		err = rows.Scan(&address, &lastPing)
		if err != nil {
			continue
		}
		nodes = append(nodes, node.Node{
			Address: address, LastPing: time.Unix(lastPing, 0),
		})
	}
	log.Info(OK, "Returning "+strconv.Itoa(len(nodes))+" CoordinatorNodes.")
	return OK, nodes
}

//GetRandomCoordinatorNodes returns max <number> random CoordinatorNodes
func GetRandomCoordinatorNodes(max int) (status int, nodes []node.Node) {
	log.Info(InProgress, "Getting "+strconv.Itoa(max)+" random CoordinatorNodes...")
	//TODO: Return random CoordinatorNodes

	result := []node.Node{
		{
			Address:  "test1",
			LastPing: time.Time{},
		},
		{
			Address:  "test2",
			LastPing: time.Time{},
		},
		{
			Address:  "test3",
			LastPing: time.Time{},
		},
	}

	log.Info(OK, "Returning "+strconv.Itoa(len(result))+" CoordinatorNodes.")
	return OK, result
}

//ClearNodeTables removes all elements from storageNodes and coordinatorNodes tables, for bootstrapping
func ClearNodeTables() (status int) {
	log.Info(InProgress, "Clearing Node Tables...")
	query := "DELETE FROM storageNodes; DELETE FROM coordinatorNodes"
	_, err := coordinatorDB.Exec(query)
	if err != nil {
		log.Error(DBWriteError, "Error clearing Node Tables: "+err.Error())
		return DBWriteError
	}
	log.Info(OK, "Cleared Node Tables.")
	return OK
}

//UpdateMessageStatusStorage updates the status of a message in the local database
func UpdateMessageStatusStorage(messageID string, status int) int {
	log.Info(InProgress, "Updating Status of Message "+messageID)
	query := "UPDATE messages SET verified=? WHERE id=?"
	stmt, err := storageDB.Prepare(query)
	if err != nil {
		log.Error(SNDBPrepareError, "Error updating status of message "+messageID+": "+err.Error())
		return SNDBPrepareError
	}
	defer stmt.Close()

	_, err = stmt.Exec(status, messageID)
	if err != nil {
		log.Error(SNDBWriteError, "Failed updating status of message "+messageID+": "+err.Error())
		return SNDBWriteError
	}
	log.Info(OK, "Updated status of Message "+messageID+". New status: "+strconv.Itoa(status))
	return OK
}
