package main

import (
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"github.com/logansua/nfl_app/file"
	"github.com/logansua/nfl_app/pagination"
	"github.com/logansua/nfl_app/utils"
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

	Name string `valid:"email" json:"name"`
}

var (
	db  *gorm.DB
	err error
)

func CreatePlayer(w http.ResponseWriter, r *http.Request) {
	var player Player
	json.NewDecoder(r.Body).Decode(&player)

	valid, errors := govalidator.ValidateStruct(&player)
	if !valid {
		panic(errors)
	}

	db.Create(&player)

	utils.JsonResponse(w, player)
}

func GetPlayers(w http.ResponseWriter, r *http.Request) {
	var players []Player

	query := r.URL.Query()

	var paging pagination.Pagination

	paging.ParseParams(query)

	db.
		Offset(paging.Offset).
		Limit(paging.Limit).
		Find(&players)

	utils.JsonResponse(w, players)
}

func GetPlayer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var player Player
	db.First(&player, params["id"])

	utils.JsonResponse(w, player)
}

func DeletePlayer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var player Player
	db.First(&player, params["id"])
	db.Delete(&player)

	w.WriteHeader(204)
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
		log.Fatal("failed to connect database")
	}

	defer db.Close()

	router.HandleFunc("/players", CreatePlayer).Methods(http.MethodPost)
	router.HandleFunc("/players", GetPlayers).Methods(http.MethodGet)
	router.HandleFunc("/players/{id}", GetPlayer).Methods(http.MethodGet)
	router.HandleFunc("/players/{id}", DeletePlayer).Methods(http.MethodDelete)
	router.HandleFunc("/players/avatar", file.UploadFileHandler()).Methods(http.MethodPost)

	log.Print(fmt.Sprintf("Server started on %s:%s", os.Getenv("APP_HOST"), os.Getenv("APP_PORT")))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("APP_PORT")), router))
}
