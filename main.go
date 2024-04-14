package main

import (
	"encoding/json"
	"go-api/model"
	"log"
	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func salvarPedidoHander(w http.ResponseWriter, r *http.Request) {

	var order model.Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}


	dsn := "host=localhost user=postgres password=postgres dbname=talosorder port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	db.AutoMigrate(&model.Order{})


	if err := db.Create(&order).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}

func main() {
	http.HandleFunc("/order", salvarPedidoHander)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
