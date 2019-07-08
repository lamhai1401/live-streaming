package main

import (
	"encoding/json"
	"net/http"

	"github.com/pion/webrtc"
)

var pmList = make(map[string]*PeerMaster)

var handleBroadcast = func(w http.ResponseWriter, r *http.Request) {
	requestBody := BroadcastRequestBody{}
	err := json.NewDecoder(r.Body).Decode(&requestBody)

	defer r.Body.Close()

	if err != nil {
		Respond(w, Message(false, "Invalid request"))
		return
	}

	// create peer master
	pm := PeerMaster{
		SDP:    createSDP(requestBody.Session.SDP, requestBody.Session.Type),
		RoomID: requestBody.IDRoom,
		Slaves: make(map[string]*PeerSlave),
	}

	// save peer master to map
	pmList[pm.RoomID] = &pm

	// peer master rtc connecting
	answer, err := pm.RTCConnecting()
	if err != nil {
		Respond(w, Message(false, "RTC connected failed"))
		return
	}
	session := Session{
		SDP:  answer.SDP,
		Type: answer.Type.String(),
	}

	// response
	resp := Message(true, "Broadcasting connected !!!")
	resp["answer"] = session
	Respond(w, resp)
	return
}

func createSDP(SDP string, Type string) webrtc.SessionDescription {
	sesion := webrtc.SessionDescription{}
	sesion.SDP = SDP

	switch Type {
	case "answer":
		sesion.Type = webrtc.SDPTypeAnswer
	case "offer":
		sesion.Type = webrtc.SDPTypeOffer
	default:
		sesion.Type = webrtc.SDPTypeRollback
	}
	return sesion
}

var handleViewer = func(w http.ResponseWriter, r *http.Request) {
	body := BroadcastRequestBody{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		Respond(w, Message(false, "Invalid request body"))
		return
	}

	pm, ok := pmList[body.IDRoom]
	if !ok {
		Respond(w, Message(false, "Invalid room id"))
		return
	}

	offer, peerID, err := pm.addPeerSlave()

	if err != nil {
		Respond(w, Message(false, "Cannot creaa peer slave"))
		return
	}

	resp := Message(true, "Request success")
	resp["idPeer"] = peerID
	resp["offer"] = Session{
		SDP:  offer.SDP,
		Type: offer.Type.String(),
	}

	Respond(w, resp)
	return
}

var handleViewerAnswer = func(w http.ResponseWriter, r *http.Request) {
	body := BroadcastRequestBody{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		Respond(w, Message(false, "Invalid request body"))
		return
	}

	pm, ok := pmList[body.IDRoom]
	if !ok {
		Respond(w, Message(false, "Invalid room id"))
		return
	}

	peerSlave, ok := pm.Slaves[body.IDPeer]
	if !ok {
		Respond(w, Message(false, "Invalid peer id"))
		return
	}

	// create session
	answer := createSDP(body.Session.SDP, body.Session.Type)

	peerSlave.Connection.SetRemoteDescription(answer)

	Respond(w, Message(true, "Connected !!!"))
	return
}
