package main

import (
	"encoding/json"
	"log"

	"github.com/pion/webrtc/v3"
)

func setupPeerA(room *Room) {
	log.Println("Setting up Peer A for room:", room.ID)
	var err error
	room.peerA_pc, err = newPeerConnection()
	if err != nil {
		log.Println("Failed to create Peer A PC:", err)
		return
	}
	room.peerA_pc.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c != nil {
			ice := c.ToJSON()
			room.iceChanA <- &ice
		}
	})
	room.peerA_dc, err = room.peerA_pc.CreateDataChannel("chat", nil)
	if err != nil {
		log.Println("Failed to create Peer A DC:", err)
		return
	}
	room.peerA_dc.OnOpen(func() {
		log.Printf("✅ Server-side Data Channel for Peer A OPENED")
		_ = room.peerA_ws.WriteJSON(Message{Type: "status", Payload: "Connection established! You can now chat."})
	})
	room.peerA_dc.OnMessage(func(msg webrtc.DataChannelMessage) {
		log.Printf("Message from Peer B DC relayed to Peer A WS")
		var relayedMsg Message
		if err := json.Unmarshal(msg.Data, &relayedMsg); err == nil {
			_ = room.peerA_ws.WriteJSON(relayedMsg)
		} else {
			log.Println("Error unmarshalling relayed message for Peer A:", err)
		}
	})
	offer, _ := room.peerA_pc.CreateOffer(nil)
	_ = room.peerA_pc.SetLocalDescription(offer)
	room.offerChan <- offer
	answer := <-room.answerChan
	_ = room.peerA_pc.SetRemoteDescription(answer)
	go func() {
		for ice := range room.iceChanB {
			_ = room.peerA_pc.AddICECandidate(*ice)
		}
	}()
}
