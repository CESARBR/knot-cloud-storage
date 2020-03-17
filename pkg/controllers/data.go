package controllers

import (
	"encoding/json"
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

type errorMessage struct {
	error   bool
	message string
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

// GetAll handles incoming data listing requests.
func (d *DataController) GetAll(w http.ResponseWriter, r *http.Request) {
	query, errURL := getURLQueryParams(r)
	if errURL.error {
		d.respondWithError(w, http.StatusUnprocessableEntity, errURL.message)
		return
	}

	token := r.Header.Get("auth_token")
	things, err := d.DataInteractor.List(token, query)
	if err != nil {
		d.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	d.respondWithJSON(w, http.StatusOK, things)
}

// GetByID handles incmoning data retrievel by id requests.
func (d *DataController) GetByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	query, errURL := getURLQueryParams(r)
	query.SensorID = params["id"]

	if errURL.error {
		d.respondWithError(w, http.StatusUnprocessableEntity, errURL.message)
		return
	}

	token := r.Header.Get("auth_token")
	thing, err := d.DataInteractor.List(token, query)
	if err != nil {
		d.respondWithError(w, http.StatusBadRequest, "Invalid Thing ID")
		return
	}
	d.respondWithJSON(w, http.StatusOK, thing)
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

func getURLQueryParams(r *http.Request) (query *entities.Query, errorStatus errorMessage) {
	var err error
	order := 1
	skip := 0
	take := 10
	startDate, err := time.Parse("2006-1-02", "2006-1-02")
	errorStatus = checkError(err, "error when trying to parse time")
	finishDate := time.Now()

	for k, v := range r.URL.Query() {
		switch k {
		case "order":
			order, err = strconv.Atoi(v[0])
			errorStatus = checkError(err, "order must be in the following format: 1 or -1")
		case "skip":
			skip, err = strconv.Atoi(v[0])
			errorStatus = checkError(err, "skip must be an integer")
		case "take":
			take, err = strconv.Atoi(v[0])
			errorStatus = checkError(err, "take must be an integer (maximum 100)")
			if take > maxItemsAllowedToRequest {
				take = maxItemsAllowedToRequest
			}
		case "startDate":
			startDate, err = time.Parse("2006-01-02 15:04:05", v[0])
			errorStatus = checkError(err, "date must be in the following format: YYYY-MM-DD HH:MM:SS")
		case "finishDate":
			finishDate, err = time.Parse("2006-01-02 15:04:05", v[0])
			errorStatus = checkError(err, "date must be in the following format: YYYY-MM-DD HH:MM:SS")
		}
	}

	return &entities.Query{
		Order:      order,
		Skip:       skip,
		Take:       take,
		StartDate:  startDate,
		FinishDate: finishDate,
	}, errorStatus
}

func checkError(err error, text string) (errorStatus errorMessage) {
	if err != nil {
		errorStatus.error = true
		errorStatus.message = text
	}

	return errorStatus
}
