package main

import (
	"net/http"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	logrus.SetOutput(&lumberjack.Logger{
		Filename:   "roomapp.log",
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	})
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/ws", handleConnections)

	logrus.Info("🚀 Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logrus.WithError(err).Fatal("Listenandserve Failed")
	}
}
