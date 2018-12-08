package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"github.com/logansua/nfl_app/pagination"
	"log"
	"net/http"
	"os"
	"time"
)

type Model struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

type Player struct {
	Model

	Name string `json:"name"`
}

var db *gorm.DB
var err error

func CreatePlayer(w http.ResponseWriter, r *http.Request) {
	var player Player
	json.NewDecoder(r.Body).Decode(&player)

	db.Create(&player)

	jsonResponse(w, player)
}

func GetPlayers(w http.ResponseWriter, r *http.Request) {
	var players []Player

	query := r.URL.Query()

	var paging pagination.Pagination

	paging.ParseParams(query)

	db.
		//Offset(paging.offset).
		//Limit(paging.limit).
		Find(&players)

	jsonResponse(w, players)
}

func GetPlayer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var player Player
	db.First(&player, params["id"])

	jsonResponse(w, player)
}

func DeletePlayer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var player Player
	db.First(&player, params["id"])
	db.Delete(&player)

	w.WriteHeader(204)
}

func jsonResponse(w http.ResponseWriter, model interface{}) {
	js, err := json.Marshal(model)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	router := mux.NewRouter()

	db, err = gorm.Open("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_NAME"),
		os.Getenv("DATABASE_PASSWORD"),
	))

	if err != nil {
		panic("failed to connect database")
	}

	defer db.Close()

	router.HandleFunc("/players", CreatePlayer).Methods("POST")
	router.HandleFunc("/players", GetPlayers).Methods("GET")
	router.HandleFunc("/players/{id}", GetPlayer).Methods("GET")
	router.HandleFunc("/players/{id}", DeletePlayer).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}
