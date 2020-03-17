package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	. "github.com/CESARBR/knot-cloud-storage/pkg/entities"
	"github.com/CESARBR/knot-cloud-storage/pkg/interactor"
	"github.com/gorilla/mux"
)

const maxItemsAllowedToRequest = 100

type DataController struct {
	DataInteractor *interactor.DataInteractor
}

type errorMessage struct {
	error   bool
	message string
}

func NewDataController(dataInteractor *interactor.DataInteractor) *DataController {
	return &DataController{dataInteractor}
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"message": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (d *DataController) GetAll(w http.ResponseWriter, r *http.Request) {

	order, skip, take, startDate, finishDate, errUrl := getUrlQueryParams(r)
	if errUrl.error != false {
		respondWithError(w, http.StatusUnprocessableEntity, errUrl.message)
		return
	}

	things, err := d.DataInteractor.GetAll(order, skip, take, startDate, finishDate)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, things)
}

func (d *DataController) GetByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	order, skip, take, startDate, finishDate, errUrl := getUrlQueryParams(r)
	if errUrl.error != false {
		respondWithError(w, http.StatusUnprocessableEntity, errUrl.message)
		return
	}
	thing, err := d.DataInteractor.GetByID(params["id"], order, skip, take, startDate, finishDate)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Thing ID")
		return
	}
	respondWithJson(w, http.StatusOK, thing)
}

func (d *DataController) Save(w http.ResponseWriter, r *http.Request) {
	var thing Data
	if err := json.NewDecoder(r.Body).Decode(&thing); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	thing.Timestamp = time.Now()
	if err := d.DataInteractor.Save(thing); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, thing)
	defer r.Body.Close()

}

func getUrlQueryParams(r *http.Request) (order, skip, take int, startDate, finishDate time.Time, errorStatus errorMessage) {
	var err error = nil
	order = 1
	skip = 0
	take = 10
	startDate, err = time.Parse("2006-1-02", "2006-1-02")
	errorStatus = checkError(err, "error when trying to parse time")
	finishDate = time.Now()

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

	return order, skip, take, startDate, finishDate, errorStatus
}

func checkError(err error, text string) (errorStatus errorMessage) {
	if err != nil {
		errorStatus.error = true
		errorStatus.message = text
	}

	return errorStatus
}
