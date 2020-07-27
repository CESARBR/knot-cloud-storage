package server

import (
	"encoding/json"
	"fmt"

	"github.com/CESARBR/knot-cloud-storage/pkg/interactor"
	"github.com/CESARBR/knot-cloud-storage/pkg/logging"
	"github.com/CESARBR/knot-cloud-storage/pkg/network"
	"github.com/pkg/errors"
)

const (
	dataQueue                 = "storage-thing-data"
	deviceQueue               = "storage-thing-events"
	dataPublishedExchange     = "data.published"
	dataPublishedExchangeType = "fanout"
	deviceExchange            = "device"
	deviceExchangeType        = "direct"
	deviceUnregisteredKey     = "device.unregistered"
)

// Handler handle messages received from a service
type Handler struct {
	logger         logging.Logger
	amqp           *network.Amqp
	dataInteractor interactor.Interactor
}

// NewHandler creates a new Handler instance with the necessary dependencies
func NewHandler(logger logging.Logger, amqp *network.Amqp, dataInteractor interactor.Interactor) *Handler {
	return &Handler{logger, amqp, dataInteractor}
}

// Start starts to listen messages
func (h *Handler) Start(started chan bool) {
	h.logger.Debug("Msg handler started")
	msgChan := make(chan network.InMsg)
	err := h.subscribeToMessages(msgChan)
	if err != nil {
		h.logger.Error(err)
		started <- false
		return
	}

	go h.onMsgReceived(msgChan)

	started <- true
}

// Stop stops to listen for messages
func (h *Handler) Stop() {
	h.logger.Debug("Handler stopped")
}

func (h *Handler) subscribeToMessages(msgChan chan network.InMsg) error {
	var err error
	subscribe := func(msgChan chan network.InMsg, queue, exchange, kind, key string) {
		if err != nil {
			return
		}
		err = h.amqp.OnMessage(msgChan, queue, exchange, kind, key)
	}

	subscribe(msgChan, dataQueue, dataPublishedExchange, dataPublishedExchangeType, "")
	subscribe(msgChan, deviceQueue, deviceExchange, deviceExchangeType, deviceUnregisteredKey)

	return nil
}

func (h *Handler) onMsgReceived(msgChan chan network.InMsg) {
	for {
		var err error
		msg := <-msgChan
		h.logger.Infof("Exchange: %s, routing key: %s", msg.Exchange, msg.RoutingKey)
		h.logger.Infof("Message received: %s", string(msg.Body))

		token, ok := msg.Headers["Authorization"].(string)
		if !ok {
			err = errors.New("authorization token not provided")
			h.logger.Error(err)
			continue
		}

		switch msg.RoutingKey {
		case deviceUnregisteredKey:
			err = h.handleUnregisteredDevice(msg, token)
		case "": // handling broadcasted data events
			if msg.Exchange == dataPublishedExchange {
				err = h.handlePublishedData(msg, token)
			}
		}

		if err != nil {
			h.logger.Error(err)
			continue
		}
	}
}

func (h *Handler) handlePublishedData(msg network.InMsg, token string) error {
	dataMsg := network.DataPublish{}
	err := json.Unmarshal(msg.Body, &dataMsg)
	if err != nil {
		return fmt.Errorf("message body parsing error: %w", err)
	}

	err = h.dataInteractor.Save(token, dataMsg.ID, dataMsg.Data)
	if err != nil {
		return err
	}

	h.logger.Info("data successfully saved")
	return nil
}

func (h *Handler) handleUnregisteredDevice(msg network.InMsg, token string) error {
	device := network.DeviceUnregistered{}
	err := json.Unmarshal(msg.Body, &device)
	if err != nil {
		return fmt.Errorf("message body parsing error: %w", err)
	}

	err = h.dataInteractor.Delete(token, device.ID)
	if err != nil {
		return err
	}

	h.logger.Info("data successfully deleted")
	return nil
}
