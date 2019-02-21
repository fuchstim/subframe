package storage

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"subframe/server/settings"
	"subframe/structs/message"
)

var messagesPath string
var databasePath string

//Init intializes the data directory
func Init() {
	createDirIfNotExist(settings.DataPath)
	messagesPath = settings.DataPath + "/messages"
	createDirIfNotExist(messagesPath)
	databasePath = settings.DataPath + "/databases"
	createDirIfNotExist(databasePath)
}

//Finish might do something soon
func Finish() {

}

//Get loads a message from local disk
func Get(id string) (msg message.Message, status int) {
	//Read message from disk and return
	dat, err := ioutil.ReadFile(messagesPath + id)
	if err != nil {
		return message.Message{}, http.StatusNotFound
	}
	return message.Message{
		ID:      id,
		Content: string(dat),
	}, http.StatusOK
}

//Put writes a message to local disk
func Put(msg message.Message) (status int) {
	id := msg.ID
	content := []byte(msg.Content)

	if !checkStorageSpace(len(content)) {
		return http.StatusInsufficientStorage
	}

	if _, err := os.Stat(messagesPath + id); os.IsNotExist(err) {
		err := ioutil.WriteFile(messagesPath+id, content, 0600)
		if err != nil {
			fmt.Println(err)
			return http.StatusInternalServerError
		}
		//Write message to disk
		return http.StatusOK
	}
	return http.StatusConflict
}

//Delete removes a message from local disk
func Delete(id string) (status int) {
	//Delete message from disk
	return http.StatusOK
}

//Creates Directory if it does not yet exist
func createDirIfNotExist(dir string) {
	//TODO: Fix error on windows reporting directories exists when they do not
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		fmt.Println(dir + " does not exist")
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
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
