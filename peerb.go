package main

import (
	"encoding/json"
	"log"

	"github.com/pion/webrtc/v3"
)

func setupPeerB(room *Room) {
	log.Println("Setting up Peer B for room:", room.ID)
	var err error
	room.peerB_pc, err = newPeerConnection()
	if err != nil {
		log.Println("Failed to create Peer B PC:", err)
		return
	}
	room.peerB_pc.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c != nil {
			ice := c.ToJSON()
			room.iceChanB <- &ice
		}
	})
	room.peerB_pc.OnDataChannel(func(dc *webrtc.DataChannel) {
		log.Println("Peer B received Data Channel:", dc.Label())
		room.mu.Lock()
		room.peerB_dc = dc
		room.mu.Unlock()
		dc.OnOpen(func() {
			log.Printf("✅ Server-side Data Channel for Peer B OPENED")
			_ = room.peerB_ws.WriteJSON(Message{Type: "status", Payload: "Connection established! You can now chat."})
		})
		dc.OnMessage(func(msg webrtc.DataChannelMessage) {
			log.Printf("Message from Peer A DC relayed to Peer B WS")
			var relayedMsg Message
			if err := json.Unmarshal(msg.Data, &relayedMsg); err == nil {
				_ = room.peerB_ws.WriteJSON(relayedMsg)
			} else {
				log.Println("Error unmarshalling relayed message for Peer B:", err)
			}
		})
	})
	offer := <-room.offerChan
	_ = room.peerB_pc.SetRemoteDescription(offer)
	answer, _ := room.peerB_pc.CreateAnswer(nil)
	_ = room.peerB_pc.SetLocalDescription(answer)
	room.answerChan <- answer
	go func() {
		for ice := range room.iceChanA {
			_ = room.peerB_pc.AddICECandidate(*ice)
		}
	}()
}
