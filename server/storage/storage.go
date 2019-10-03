package storage

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"subframe/server/database"
	"subframe/server/logger"
	"subframe/server/settings"
	. "subframe/status"
	"subframe/structs/message"
)

var messagesPath string
var databasePath string
var logPath string
var log = logger.Logger{Prefix: "storage/Main"}

//Init initializes the data directory
func Init() {
	log.Info(InProgress, "Initializing Storage Directories...")
	createDirIfNotExist(settings.DataPath)
	log.Info(OK, "Initialized "+settings.DataPath)

	messagesPath = settings.DataPath + "/messages"
	createDirIfNotExist(messagesPath)
	log.Info(OK, "Initialized "+messagesPath)

	databasePath = settings.DataPath + "/databases"
	createDirIfNotExist(databasePath)
	log.Info(OK, "Initialized "+databasePath)

	logPath = settings.DataPath + "/logs"
	createDirIfNotExist(logPath)
	logger.LogPath = logPath

	log.Info(OK, "Initialized "+logPath)
}

//Finish might do something soon
func Finish() {
	log.Info(InProgress, "Finishing Storage...")

	log.Info(OK, "Finished Storage.")
}

//Get loads a message from local disk
func Get(id string) (msg message.Message, status int) {
	//Read message from disk and return
	log.Info("Getting Message " + id + "...")

	if !database.CheckMessageStorage(id) {
		log.Warn("Error getting Message " + id + ": Not in database")
		return message.Message{}, http.StatusNotFound
	}

	dat, err := ioutil.ReadFile(messagesPath + "/" + id)
	if err != nil {
		log.Warn("Error getting Message " + id + ": " + err.Error())
		return message.Message{}, http.StatusNotFound
	}
	log.Info("Got Message " + id)
	return message.Message{
		ID:      id,
		Content: string(dat),
	}, http.StatusOK
}

//Put writes a message to local disk
func Put(msg message.Message) (status int) {
	id := msg.ID
	content := []byte(msg.Content)

	log.Info("Putting Message " + id)

	if database.CheckMessageStorage(id) {
		log.Error("Error storing Message " + id + ": Already in database")
		return http.StatusConflict
	}

	if !checkStorageSpace(len(content)) {
		log.Warn("Could not store Message " + id + ": Insufficient Storage.")
		return http.StatusInsufficientStorage
	}

	if _, err := os.Stat(messagesPath + "/" + id); os.IsNotExist(err) {
		err = ioutil.WriteFile(messagesPath+"/"+id, content, 0600)
		if err != nil {
			log.Error("Error storing Message " + id + ": " + err.Error())
			return http.StatusInternalServerError
		}

		log.Info("Successfully stored Message " + id)
		return http.StatusOK
	}
	log.Error("Error storing Message " + id + ": File exists")
	return http.StatusConflict
}

//Delete removes a message from local disk
func Delete(id string) (status int) {
	//Delete message from disk
	log.Fatal("Method DELETE not yet implemented")
	return http.StatusOK
}

//Creates Directory if it does not yet exist
func createDirIfNotExist(dir string) {
	//TODO: Fix error on windows reporting directories exists when they do not
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		log.Warn("Directory " + dir + " does not exist. Creating...")
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			log.Fatal("Directory " + dir + " could not be created: " + err.Error())
		}
	}
}

//Check whether Size of Data Directory exceeds size limit set in settings.DiskSpace
func checkStorageSpace(size int) bool {
	dirsize, _ := dirSize(messagesPath)
	dirsize += int64(size)
	return dirsize/1024/1024 < int64(settings.DiskSpace)
}

//Get Size of Directory
func dirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}
