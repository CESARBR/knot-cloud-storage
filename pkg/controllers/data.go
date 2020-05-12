package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/CESARBR/knot-cloud-storage/pkg/entities"
	"github.com/CESARBR/knot-cloud-storage/pkg/interactor"
	"github.com/CESARBR/knot-cloud-storage/pkg/logging"
	"github.com/gorilla/mux"
)

const maxItemsAllowedToRequest = 100

// DataController handles data operations in the storage.
type DataController struct {
	DataInteractor interactor.Interactor

	logger logging.Logger
}

// NewDataController constructs the DataController.
func NewDataController(dataInteractor interactor.Interactor, logger logging.Logger) *DataController {
	return &DataController{dataInteractor, logger}
}

func (d *DataController) writeResponse(w http.ResponseWriter, statusCode int, msg interface{}) {
	w.WriteHeader(statusCode)
	if msg == nil {
		return
	}

	js, err := json.Marshal(msg)
	if err != nil {
		d.logger.Errorf("unable to marshal json: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	if err != nil {
		d.logger.Errorf("unable to write to connection HTTP: %s", err)
		return
	}
}

// List handles incoming data list requests
func (d *DataController) List(w http.ResponseWriter, r *http.Request) {
	d.logger.Debug("get data request received")
	query, err := getQueryParams(r)
	if err != nil {
		d.logger.Errorf("failed to get query params: %s", err)
		d.writeResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	d.logger.Debug("query params received %v", query)

	token := r.Header.Get("auth_token")
	data, err := d.DataInteractor.List(token, query)
	if err != nil {
		d.logger.Errorf("failed to get data: %s", err)
		d.writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	d.logger.Info("data successfully obtained")
	d.writeResponse(w, http.StatusOK, data)
}

// Save handles incoming data insertion requests.
func (d *DataController) Save(w http.ResponseWriter, r *http.Request) {
	d.logger.Debug("save data request received")
	var data entities.Data
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		d.logger.Errorf("failed to decode payload: %s", err)
		d.writeResponse(w, http.StatusUnprocessableEntity, "Invalid request payload")
		return
	}
	d.logger.Debug("payload received: %v", data)

	token := r.Header.Get("auth_token")
	payloads := []entities.Payload{data.Payload}
	if err := d.DataInteractor.Save(token, data.From, payloads); err != nil {
		d.logger.Errorf("failed to save data: %s", err)
		d.writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer r.Body.Close()
	d.logger.Info("successfully saved data")
	d.writeResponse(w, http.StatusCreated, nil)
}

// DeleteByDeviceID handles request to delete data by device ID
func (d *DataController) DeleteByDeviceID(w http.ResponseWriter, r *http.Request) {
	d.logger.Debug("delete data request received")
	params := mux.Vars(r)
	deviceID := params["deviceId"]
	token := r.Header.Get("auth_token")

	err := d.DataInteractor.Delete(token, deviceID)
	if err != nil {
		d.logger.Errorf("failed to delete data: %s", err)
		d.writeResponse(w, http.StatusUnprocessableEntity, "Invalid information")
		return
	}

	d.logger.Info("data with %s successfully deleted", deviceID)
	d.writeResponse(w, http.StatusOK, nil)
}

// DeleteAll handles request to delete all the data associated with the user
func (d *DataController) DeleteAll(w http.ResponseWriter, r *http.Request) {
	d.logger.Debug("delete all data request received")
	token := r.Header.Get("auth_token")
	err := d.DataInteractor.Delete(token, "")
	if err != nil {
		d.logger.Errorf("failed to delete data: %s", err)
		d.writeResponse(w, http.StatusUnprocessableEntity, "Invalid information")
		return
	}

	d.logger.Info("successfully deleted all data")
	d.writeResponse(w, http.StatusOK, nil)
}

func getQueryParams(r *http.Request) (query *entities.Query, err error) {
	var startDate time.Time
	finishDate := time.Now()
	params := mux.Vars(r)
	order := 1
	var skip int64 = 0
	var take int64 = 10

	for k, v := range r.URL.Query() {
		switch k {
		case "order":
			order, err = strconv.Atoi(v[0])
			if err != nil {
				return nil, errors.New("order must be in the following format: 1 or -1")
			}
		case "skip":
			skip, err = strconv.ParseInt(v[0], 10, 64)
			if err != nil {
				return nil, errors.New("skip must be an integer")
			}
		case "take":
			take, err = strconv.ParseInt(v[0], 10, 64)
			if err != nil {
				return nil, errors.New("take must be an integer (maximum 100)")
			}

			if take > maxItemsAllowedToRequest {
				take = maxItemsAllowedToRequest
			}
		case "startDate":
			startDate, err = time.Parse("2006-01-02 15:04:05", v[0])
			if err != nil {
				return nil, errors.New("date must be in the following format: YYYY-MM-DD HH:MM:SS")
			}

		case "finishDate":
			finishDate, err = time.Parse("2006-01-02 15:04:05", v[0])
			if err != nil {
				return nil, errors.New("date must be in the following format: YYYY-MM-DD HH:MM:SS")
			}
		}
	}

	return &entities.Query{
		ThingID:    params["deviceId"],
		SensorID:   params["sensorId"],
		Order:      order,
		Skip:       skip,
		Take:       take,
		StartDate:  startDate,
		FinishDate: finishDate,
	}, nil
}
