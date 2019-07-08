package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/pion/webrtc"
)

var config = webrtc.Configuration{
	ICEServers: []webrtc.ICEServer{
		{
			URLs: []string{"stun:stun.l.google.com:19302"},
		},
	},
}

func main() {
	router := mux.NewRouter()

	router.
		HandleFunc("/broadcast", handleBroadcast)
	router.
		HandleFunc("/viewing", viewerHandler)

	router.
		HandleFunc("/viewer", handleViewer)
	router.
		HandleFunc("/answer", handleViewerAnswer)
	router.
		HandleFunc("/", broadcastHandler)
	router.
		HandleFunc("/broadcast.js", broadcastJSHandler)
	router.
		HandleFunc("/viewer.js", viewerJSHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	err := http.ListenAndServeTLS(":8080", "server.crt", "server.key", router)
	// err := http.ListenAndServe(fmt.Sprintf(":%s", port), router)

	if err != nil {
		fmt.Print(err)
	}
}

func broadcastHandler(w http.ResponseWriter, r *http.Request) {
	path := fmt.Sprintf("templates/html/broadcast.html")
	fmt.Println(path)
	http.ServeFile(w, r, path)
}

func broadcastJSHandler(w http.ResponseWriter, r *http.Request) {
	path := fmt.Sprintf("templates/js/broadcast.js")
	fmt.Println(path)
	http.ServeFile(w, r, path)
}

func viewerHandler(w http.ResponseWriter, r *http.Request) {
	path := fmt.Sprintf("templates/html/viewer.html")
	fmt.Println(path)
	http.ServeFile(w, r, path)
}

func viewerJSHandler(w http.ResponseWriter, r *http.Request) {
	path := fmt.Sprintf("templates/js/viewer.js")
	fmt.Println(path)
	http.ServeFile(w, r, path)
}
