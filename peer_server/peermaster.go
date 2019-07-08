package main

import (
	"fmt"
	"math/rand"
	"time"

	uuid "github.com/nu7hatch/gouuid"
	"github.com/pion/rtcp"
	"github.com/pion/webrtc"
)

type PeerSlave struct {
	Connection *webrtc.PeerConnection `json:"connection"`
	Track      *webrtc.Track          `json:"track"`
}

type PeerMaster struct {
	RoomID string `json:"roomId"`
	SDP    webrtc.SessionDescription
	Slaves map[string]*PeerSlave
}

func (pm *PeerMaster) GetAttribute() (string, webrtc.SessionDescription) {
	return pm.RoomID, pm.SDP
}

func (pm *PeerMaster) addPeerSlave() (*webrtc.SessionDescription, string, error) {
	uid, err := uuid.NewV4()

	peerConnection, err := webrtc.NewPeerConnection(config)

	if err != nil {
		return nil, "", err
	}

	// create Trach the we send video back to
	outputTrack, err := peerConnection.NewTrack(webrtc.DefaultPayloadTypeVP8, rand.Uint32(), "video", "pion")
	if err != nil {
		return nil, "", err
	}

	// Add this newly created track to the PeerConnection
	if _, err = peerConnection.AddTrack(outputTrack); err != nil {
		return nil, "", err
	}

	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("ICE Connection State has changed: %s\n", connectionState.String())
	})

	// Create answer
	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		return nil, "", err
	}

	// Sets the LocalDescription, and starts our UDP listeners
	err = peerConnection.SetLocalDescription(offer)

	pm.Slaves[uid.String()] = &PeerSlave{
		Connection: peerConnection,
		Track:      outputTrack,
	}

	return &offer, uid.String(), nil
}

func (pm *PeerMaster) RTCConnecting() (*webrtc.SessionDescription, error) {
	peerConnection, err := webrtc.NewPeerConnection(config)

	if err != nil {
		return nil, err
	}

	// create Trach the we send video back to
	outputTrack, err := peerConnection.NewTrack(webrtc.DefaultPayloadTypeVP8, rand.Uint32(), "video", "pion")
	if err != nil {
		return nil, err
	}

	// Add this newly created track to the PeerConnection
	if _, err = peerConnection.AddTrack(outputTrack); err != nil {
		return nil, err
	}

	peerConnection.OnTrack(func(track *webrtc.Track, receiver *webrtc.RTPReceiver) {
		go func() {
			ticker := time.NewTicker(time.Second * 3)
			for range ticker.C {
				errSend := peerConnection.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: track.SSRC()}})
				if errSend != nil {
					fmt.Println(errSend)
				}
			}
		}()

		for {
			// Read RTP packets  sent to Pion
			rtp, readErr := track.ReadRTP()
			if readErr != nil {
				panic(readErr)
			}

			// Replace the SSRC with the SSRC of the outbound track.
			// The only change we are making replacing the SSRC, the RTP packets are unchanged otherwise

			rtp.SSRC = outputTrack.SSRC()
			rtp.PayloadType = webrtc.DefaultPayloadTypeVP8

			if writeErr := outputTrack.WriteRTP(rtp); writeErr != nil {
				panic(writeErr)
			}

			for _, slave := range pm.Slaves {
				rtp.SSRC = slave.Track.SSRC()
				rtp.PayloadType = webrtc.DefaultPayloadTypeVP8

				if writeErr := slave.Track.WriteRTP(rtp); writeErr != nil {
					panic(writeErr)
				}
			}

			// do something with track for peer slave here
		}
	})

	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("ICE Connection State has changed: %s\n", connectionState.String())
	})

	// set remote session
	err = peerConnection.SetRemoteDescription(pm.SDP)

	if err != nil {
		return nil, err
	}

	// Create answer
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		return nil, err
	}

	// Sets the LocalDescription, and starts our UDP listeners
	err = peerConnection.SetLocalDescription(answer)

	if err != nil {
		return nil, err
	}

	return &answer, nil
}
