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
	queue        = "storage-thing-data"
	exchange     = "data.published"
	exchangeType = "fanout"
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
	err := h.amqp.OnMessage(msgChan, queue, exchange, exchangeType, "")
	if err != nil {
		return fmt.Errorf("fail to subscribe in message queue: %w", err)
	}

	return nil
}

func (h *Handler) onMsgReceived(msgChan chan network.InMsg) {
	for {
		msg := <-msgChan
		h.logger.Infof("Exchange: %s, routing key: %s", msg.Exchange, msg.RoutingKey)
		h.logger.Infof("Message received: %s", string(msg.Body))

		err := h.handleMessages(msg)
		if err != nil {
			h.logger.Error(err)
			continue
		}
	}
}

func (h *Handler) handleMessages(msg network.InMsg) error {
	token, ok := msg.Headers["Authorization"].(string)
	if !ok {
		return errors.New("authorization token not provided")
	}

	return h.handlePublishData(token, msg.Body)
}

func (h *Handler) handlePublishData(token string, body []byte) error {
	msg := network.DataPublish{}
	err := json.Unmarshal(body, &msg)
	if err != nil {
		return fmt.Errorf("message body parsing error: %w", err)
	}

	err = h.dataInteractor.Save(token, msg.ID, msg.Data)
	if err != nil {
		return err
	}

	h.logger.Info("data successfully saved")
	return nil
}
