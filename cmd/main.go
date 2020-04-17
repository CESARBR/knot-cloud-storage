package main

import (
	"os"

	"github.com/CESARBR/knot-cloud-storage/internal/config"
	"github.com/CESARBR/knot-cloud-storage/pkg/controllers"
	"github.com/CESARBR/knot-cloud-storage/pkg/data"
	"github.com/CESARBR/knot-cloud-storage/pkg/interactor"
	"github.com/CESARBR/knot-cloud-storage/pkg/logging"
	"github.com/CESARBR/knot-cloud-storage/pkg/network"
	"github.com/CESARBR/knot-cloud-storage/pkg/server"
	"github.com/CESARBR/knot-cloud-storage/pkg/things"
)

func main() {
	config := config.Load()

	logrus := logging.NewLogrus(config.Logger.Level)
	logger := logrus.Get("Main")
	logger.Info("Starting KNoT Cloud Storage")

	quit := make(chan bool, 1)

	mongo := data.NewMongoDB(config.MongoDB.Host, config.MongoDB.Name)
	database, err := mongo.Connect()
	if err != nil {
		logger.Error("Failed to connect with the database: ", err)
		os.Exit(1)
	}

	amqpStartChan := make(chan bool, 1)
	amqp := network.NewAmqp(config.RabbitMQ.URL, logrus.Get("AMQP"))

	thingsService := things.New(config.Users.Host, uint16(config.Users.Port), logger)
	dataStore := data.NewStore(database, logrus.Get("Storage"))
	dataInteractor := interactor.NewDataInteractor(thingsService, dataStore, logrus.Get("Interactor"))
	dataController := controllers.NewDataController(dataInteractor, logrus.Get("Controller"))

	handlerStartChan := make(chan bool, 1)
	handler := server.NewHandler(logrus.Get("Handler"), amqp, dataInteractor)

	serverStartChan := make(chan bool, 1)
	server := server.NewServer(config.Server.Port, logrus.Get("Server"), dataController)

	go amqp.Start(amqpStartChan)
	go server.Start(serverStartChan)

	// Main loop
	for {
		select {
		case started := <-serverStartChan:
			if started {
				logger.Info("Server started")
			}
		case started := <-amqpStartChan:
			if started {
				logger.Info("AMQP connection started")
				go handler.Start(handlerStartChan)
			}
		case started := <-handlerStartChan:
			if started {
				logger.Info("Handler started")
			} else {
				quit <- true
			}
		case <-quit:
			handler.Stop()
			amqp.Stop()
			server.Stop()
			os.Exit(0)
		}
	}
}
