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

func (d *DataController) respondWithError(w http.ResponseWriter, code int, msg string) {
	d.respondWithJSON(w, code, map[string]string{"message": msg})
}

func (d *DataController) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if payload != nil {
		_, err := w.Write(response)
		if err != nil {
			d.logger.Error(err)
		}
	}
}

// List handles incoming data list requests
func (d *DataController) List(w http.ResponseWriter, r *http.Request) {
	query, err := getQueryParams(r)
	if err != nil {
		d.respondWithError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	token := r.Header.Get("auth_token")
	data, err := d.DataInteractor.List(token, query)
	if err != nil {
		d.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	d.respondWithJSON(w, http.StatusOK, data)
}

// Save handles incoming data insertion requests.
func (d *DataController) Save(w http.ResponseWriter, r *http.Request) {
	var thing entities.Data
	if err := json.NewDecoder(r.Body).Decode(&thing); err != nil {
		d.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	token := r.Header.Get("auth_token")
	data := []entities.Payload{thing.Payload}
	if err := d.DataInteractor.Save(token, thing.From, data, time.Now()); err != nil {
		d.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer r.Body.Close()
	d.respondWithJSON(w, http.StatusCreated, thing)
}

func (d *DataController) DeleteByDeviceID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	deviceId := params["deviceId"]

	err := d.DataInteractor.Delete(deviceId)
	if err != nil {
		d.respondWithError(w, http.StatusUnprocessableEntity, "Invalid information")
		return
	}

	d.respondWithJSON(w, http.StatusOK, nil)
}

func (d *DataController) DeleteAll(w http.ResponseWriter, r *http.Request) {
	err := d.DataInteractor.Delete("")
	if err != nil {
		d.respondWithError(w, http.StatusUnprocessableEntity, "Invalid information")
		return
	}

	d.respondWithJSON(w, http.StatusOK, nil)
}

func getQueryParams(r *http.Request) (query *entities.Query, err error) {
	var startDate time.Time
	finishDate := time.Now()
	params := mux.Vars(r)
	order := 1
	skip := 0
	take := 10

	for k, v := range r.URL.Query() {
		switch k {
		case "order":
			order, err = strconv.Atoi(v[0])
			if err != nil {
				return nil, errors.New("order must be in the following format: 1 or -1")
			}
		case "skip":
			skip, err = strconv.Atoi(v[0])
			if err != nil {
				return nil, errors.New("skip must be an integer")
			}
		case "take":
			take, err = strconv.Atoi(v[0])
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
