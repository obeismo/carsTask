package main

import (
	"fmt"
	"net/http"

	_ "em_test/cmd/docs"
	"em_test/pkg/database"
	"em_test/pkg/handlers"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatal("Error loading .env file")
	}

	db, err := database.InitDB()

	if err != nil {
		logrus.Fatalf("Failed to initialize db: %s", err.Error())
		return
	}
	defer db.Close()

	router := mux.NewRouter()

	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8000/swagger/doc.json"),
	))

	router.HandleFunc("/cars", handlers.GetCarsHandler).Methods("GET")
	router.HandleFunc("/cars/{id}", handlers.DeleteCarHandler).Methods("DELETE")
	router.HandleFunc("/cars/{id}", handlers.UpdateCarHandler).Methods("PUT")
	router.HandleFunc("/cars", handlers.AddNewCarHandler).Methods("POST")

	fmt.Println("The server is running on http://localhost:8000")
	fmt.Println("Swagger on http://localhost:8000/swagger/index.html")
	http.ListenAndServe(":8000", router)
}
