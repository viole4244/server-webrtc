package main

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload,omitempty"`
	RoomID  string      `json:"roomId,omitempty"`
	Sender  string      `json:"sender,omitempty"`
}

type Room struct {
	ID         string
	peerA_pc   *webrtc.PeerConnection
	peerB_pc   *webrtc.PeerConnection
	peerA_dc   *webrtc.DataChannel
	peerB_dc   *webrtc.DataChannel
	peerA_ws   *websocket.Conn
	peerB_ws   *websocket.Conn
	offerChan  chan webrtc.SessionDescription
	answerChan chan webrtc.SessionDescription
	iceChanA   chan *webrtc.ICECandidateInit
	iceChanB   chan *webrtc.ICECandidateInit
	mu         sync.Mutex
}

var (
	rooms      = make(map[string]*Room)
	roomsMutex = &sync.Mutex{}

	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)
