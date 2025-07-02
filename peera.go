package main

import (
	"encoding/json"

	"github.com/pion/webrtc/v3"
	"github.com/sirupsen/logrus"
)

func setupPeerA(room *Room) {
	logrus.WithField("room_id", room.ID).Info("Setting up peer A")
	var err error
	room.peerA_pc, err = newPeerConnection()
	if err != nil {
		logrus.WithError(err).Error("Failed to create peer A PC")
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
		logrus.WithError(err).Error("Failed to create peer A DC")
		return
	}
	room.peerA_dc.OnOpen(func() {
		logrus.Info("✅ Server-side Data Channel for Peer A OPENED")
		_ = room.peerA_ws.WriteJSON(Message{Type: "status", Payload: "Connection established! You can now chat."})
	})
	room.peerA_dc.OnMessage(func(msg webrtc.DataChannelMessage) {
		logrus.Info("Message from Peer B DC relayed to Peer A WS")
		var relayedMsg Message
		if err := json.Unmarshal(msg.Data, &relayedMsg); err == nil {
			_ = room.peerA_ws.WriteJSON(relayedMsg)
		} else {
			logrus.WithError(err).Error("Error Unmarshalling relayed message for peer A")
		}
	})
	offer, err := room.peerA_pc.CreateOffer(nil)
	if err != nil {
		logrus.WithError(err).Error("Failed to create offer for peer A")
	}
	if err := room.peerA_pc.SetLocalDescription(offer); err != nil {
		logrus.WithError(err).Error("Failed to set Local Description for peer A")
	}
	room.offerChan <- offer
	answer := <-room.answerChan
	if err := room.peerA_pc.SetRemoteDescription(answer); err != nil {
		logrus.WithError(err).Error("Failed to set remote description for peer A")
	}
	go func() {
		for ice := range room.iceChanB {
			_ = room.peerA_pc.AddICECandidate(*ice)
		}
	}()
}
