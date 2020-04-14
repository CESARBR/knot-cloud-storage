package main

import (
	"github.com/CESARBR/knot-cloud-storage/internal/config"
	"github.com/CESARBR/knot-cloud-storage/pkg/controllers"
	"github.com/CESARBR/knot-cloud-storage/pkg/data"
	"github.com/CESARBR/knot-cloud-storage/pkg/interactor"
	"github.com/CESARBR/knot-cloud-storage/pkg/logging"
	"github.com/CESARBR/knot-cloud-storage/pkg/server"

	"os"
)

func main() {
	config := config.Load()

	logrus := logging.NewLogrus(config.Logger.Level)
	logger := logrus.Get("Main")
	logger.Info("Starting KNoT Cloud Storage")

	mongo := data.NewMongoDB(config.MongoDB.Host, config.MongoDB.Name)
	database, err := mongo.Connect()
	if err != nil {
		logger.Error("Failed to connect with the database: ", err)
		os.Exit(1)
	}

	dataStore := data.NewDataStore(database, logrus.Get("Storage"))
	dataInteractor := interactor.NewDataInteractor(dataStore, logrus.Get("Interactor"))
	dataController := controllers.NewDataController(dataInteractor)

	server := server.NewServer(config.Server.Port, logrus.Get("Server"), dataController)
	server.Start()
}
