package main

import (
	"encoding/json"
	"fmt"
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
	db.AutoMigrate((&model.Item{}))


	tx := db.Begin()
	defer func() {
			if r := recover(); r != nil {
					tx.Rollback()
					http.Error(w, "Error saving order: "+ fmt.Sprint(r), http.StatusInternalServerError)
					return
			}
	}()

	// Efficiently upsert items (avoid redundant lookups)
	itemMap := make(map[string]model.Item)
	for _, item := range order.Item {
			if existingItem, ok := itemMap[item.Name]; ok {
					order.Item = append(order.Item, existingItem) // Use existing item
			} else {
					var newItem model.Item
					if err := tx.FirstOrCreate(&newItem, model.Item{Name: item.Name}).Error; err != nil {
							tx.Rollback()
							http.Error(w, "Error upserting item: "+err.Error(), http.StatusInternalServerError)
							return
					}
					itemMap[item.Name] = newItem
					order.Item = append(order.Item, newItem)
			}
	}

	if err := db.Create(&order).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}



	// w.WriteHeader(http.StatusOK)

	response := struct {
		Message string `json:"message"`
		OrderID uint   `json:"order_id"`
	} {
		Message: "Pedido Salvo",
		OrderID: order.ID,
	}
	json.NewEncoder(w).Encode(response)

}

func buscarPedidosHandler(w http.ResponseWriter, r *http.Request) {

	dsn := "host=localhost user=postgres password=postgres dbname=talosorder port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var orders []model.Order
	if err := db.Preload("Item").Find(&orders).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)

}

func main() {
	http.HandleFunc("/order", salvarPedidoHander)
	http.HandleFunc("/orders", buscarPedidosHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
