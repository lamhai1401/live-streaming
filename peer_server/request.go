package main

import (
	"encoding/json"
	"net/http"
)

func Respond(w http.ResponseWriter, data map[string]interface{}) {
	/*This is for response*/
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

type Session struct {
	SDP  string `json:"sdp"`
	Type string `json:"type"`
}

type BroadcastRequestBody struct {
	IDRoom  string  `json:"idRoom"`
	IDPeer  string  `json:"idPeer"` // for detecting peer slave
	Session Session `json:"session"`
}
