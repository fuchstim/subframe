package networking

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func startStorageNodeAPIService() {
	http.HandleFunc("/", handleRequest)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func handleRequest(responseWriter http.ResponseWriter, request *http.Request) {
	action, messageID := parsePath(request.URL.Path)
	fmt.Fprintf(responseWriter, "You are trying to %s Message %s", action, messageID)
}

func parsePath(path string) (action string, messageID string) {
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		return "", ""
	}
	return parts[1], parts[2]
}
