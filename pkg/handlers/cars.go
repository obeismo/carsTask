package handlers

import (
	"em_test/pkg/database"
	"em_test/pkg/model"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func AddNewCarHandler(w http.ResponseWriter, r *http.Request) {
	newCar := model.Car{}
	validate := validator.New()

	err := json.NewDecoder(r.Body).Decode(&newCar)
	if err != nil {
		newErrorResponse(w, http.StatusBadRequest, "Failed to decode JSON")
		return
	}

	err = validate.Struct(newCar)
	if err != nil {
		newErrorResponse(w, http.StatusBadRequest, "Invalid JSON data")
		return
	}

	carID, err := database.AddNewCar(newCar)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, "The car with that license plate number is already in the db")
		return
	}

	newResponse(w, http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"car_id": strconv.Itoa(carID)})
}

func GetCarsHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	filter := make(map[string]string)
	filter["reg_num"] = queryParams.Get("reg_num")
	filter["mark"] = queryParams.Get("mark")
	filter["model"] = queryParams.Get("model")
	filter["year"] = queryParams.Get("year")

	offset, err := strconv.Atoi(queryParams.Get("page"))
	if err != nil {
		newErrorResponse(w, http.StatusBadRequest, "Invalid page number")
		return
	}

	limit, err := strconv.Atoi(queryParams.Get("page_size"))
	if err != nil {
		newErrorResponse(w, http.StatusBadRequest, "Invalid page size")
		return
	}

	cars, err := database.GetCarsByFilter(filter, offset, limit)
	if err != nil {
		logrus.Errorf("failet to fetch cars: %s", err.Error())
		newErrorResponse(w, http.StatusInternalServerError, "Failed to fetch cars")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cars)
}

func UpdateCarHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		newErrorResponse(w, http.StatusBadRequest, "Invalid car ID")
		return
	}

	var updateFields map[string]string
	err = json.NewDecoder(r.Body).Decode(&updateFields)
	if err != nil {
		newErrorResponse(w, http.StatusBadRequest, "Failed to decode JSON request")
		return
	}

	err = database.UpdateCarByID(id, updateFields)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, "Failed to update car")
		return
	}

	newResponse(w, http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "successfully updated"})
}

func DeleteCarHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		newErrorResponse(w, http.StatusBadRequest, "Invalid car ID")
		return
	}

	err = database.DeleteCarByID(id)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, "There's no car with that ID in the db")
		return
	}

	newResponse(w, http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": fmt.Sprintf("Vehicle with id = %d has been successfully deleted", id)})
}
