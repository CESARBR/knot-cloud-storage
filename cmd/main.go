package main

import (
	"github.com/CESARBR/knot-cloud-storage/internal/config"
	"github.com/CESARBR/knot-cloud-storage/pkg/logging"
	"github.com/CESARBR/knot-cloud-storage/pkg/server"
)

func main() {
	config := config.Load()
	logrus := logging.NewLogrus(config.Logger.Level)

	logger := logrus.Get("Main")
	logger.Info("Starting KNoT Babeltower")

	server := server.NewServer(config.Server.Port, logrus.Get("Server"))
	server.Start()
}
