package main

import (
	"encoding/json"

	"github.com/pion/webrtc/v3"
	"github.com/sirupsen/logrus"
)

func setupPeerB(room *Room) {
	logrus.WithField("room_id", room.ID).Info("Setting up peer B")
	var err error
	room.peerB_pc, err = newPeerConnection()
	if err != nil {
		logrus.WithError(err).Error("Failed to create peer B PC")
		return
	}
	room.peerB_pc.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c != nil {
			ice := c.ToJSON()
			room.iceChanB <- &ice
		}
	})
	room.peerB_pc.OnDataChannel(func(dc *webrtc.DataChannel) {
		logrus.WithField("label", dc.Label()).Info("Peer B received data channel")
		room.mu.Lock()
		room.peerB_dc = dc
		room.mu.Unlock()
		dc.OnOpen(func() {
			logrus.Info("✅ Server-side Data Channel for Peer B OPENED")
			_ = room.peerB_ws.WriteJSON(Message{Type: "status", Payload: "Connection established! You can now chat."})
		})
		dc.OnMessage(func(msg webrtc.DataChannelMessage) {
			logrus.Info("Message from Peer A DC relayed to Peer B WS")
			var relayedMsg Message
			if err := json.Unmarshal(msg.Data, &relayedMsg); err == nil {
				_ = room.peerB_ws.WriteJSON(relayedMsg)
			} else {
				logrus.WithError(err).Error("Error Unmarshalling relayed message for peer B")
			}
		})
	})
	offer := <-room.offerChan
	if err := room.peerB_pc.SetRemoteDescription(offer); err != nil {
		logrus.WithError(err).Error("Failed to set Remote Description for peer B")
	}
	answer, err := room.peerB_pc.CreateAnswer(nil)
	if err != nil {
		logrus.WithError(err).Error("Failed to create answer for peer B")
	}
	if err := room.peerB_pc.SetLocalDescription(answer); err != nil {
		logrus.WithError(err).Error("Failed to set Local Description for peer B")
	}
	room.answerChan <- answer
	go func() {
		for ice := range room.iceChanA {
			_ = room.peerB_pc.AddICECandidate(*ice)
		}
	}()
}
