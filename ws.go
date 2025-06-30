package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}
	defer ws.Close()

	log.Println("🔗 Client connected")

	var username string
	var associatedRoomID string
	var isPeerA bool

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Read error: %v. Cleaning up.", err)
			cleanupDisconnectedUser(ws, associatedRoomID, username)
			break
		}

		roomsMutex.Lock()
		room, roomExists := rooms[msg.RoomID]
		roomsMutex.Unlock()

		switch msg.Type {
		case "register":
			name, ok := msg.Payload.(string)
			if ok && name != "" {
				username = name
				log.Printf("User %s registered.", username)
				_ = ws.WriteJSON(Message{Type: "registered"})
			}

		case "create":
			if roomExists {
				_ = ws.WriteJSON(Message{Type: "error", Payload: "Room already exists"})
				continue
			}
			newRoom := &Room{
				ID:         msg.RoomID,
				peerA_ws:   ws,
				offerChan:  make(chan webrtc.SessionDescription),
				answerChan: make(chan webrtc.SessionDescription),
				iceChanA:   make(chan *webrtc.ICECandidateInit, 5),
				iceChanB:   make(chan *webrtc.ICECandidateInit, 5),
			}
			roomsMutex.Lock()
			rooms[msg.RoomID] = newRoom
			roomsMutex.Unlock()

			associatedRoomID = msg.RoomID
			isPeerA = true
			_ = ws.WriteJSON(Message{Type: "status", Payload: "Room created. Waiting for the other peer..."})
			go setupPeerA(newRoom)

		case "join":
			if !roomExists || room.peerB_ws != nil {
				_ = ws.WriteJSON(Message{Type: "error", Payload: "Room not found or is full"})
				continue
			}
			room.mu.Lock()
			room.peerB_ws = ws
			room.mu.Unlock()

			associatedRoomID = msg.RoomID
			isPeerA = false
			_ = ws.WriteJSON(Message{Type: "status", Payload: "Joined room. Establishing connection..."})
			_ = room.peerA_ws.WriteJSON(Message{Type: "status", Payload: username + " joined. Establishing connection..."})
			go setupPeerB(room)

		case "chat", "file-start", "file-chunk", "file-end":
			if !roomExists {
				continue
			}
			msg.Sender = username
			messageBytes, err := json.Marshal(msg)
			if err != nil {
				log.Println("Error marshalling message:", err)
				continue
			}
			room.mu.Lock()
			if isPeerA && room.peerA_dc != nil && room.peerA_dc.ReadyState() == webrtc.DataChannelStateOpen {
				_ = room.peerA_dc.SendText(string(messageBytes))
			} else if !isPeerA && room.peerB_dc != nil && room.peerB_dc.ReadyState() == webrtc.DataChannelStateOpen {
				_ = room.peerB_dc.SendText(string(messageBytes))
			}
			room.mu.Unlock()
		}
	}
}

func cleanupDisconnectedUser(ws *websocket.Conn, roomID, username string) {
	if roomID == "" {
		return
	}
	roomsMutex.Lock()
	defer roomsMutex.Unlock()

	if room, ok := rooms[roomID]; ok {
		var otherPeerWS *websocket.Conn
		if ws == room.peerA_ws {
			otherPeerWS = room.peerB_ws
		} else {
			otherPeerWS = room.peerA_ws
		}
		if otherPeerWS != nil {
			_ = otherPeerWS.WriteJSON(Message{Type: "peer-disconnect", Payload: username + " has disconnected."})
		}
		delete(rooms, roomID)
	}
}

func newPeerConnection() (*webrtc.PeerConnection, error) {
	return webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{URLs: []string{"stun:stun.l.google.com:19302"}},
		},
	})
}
