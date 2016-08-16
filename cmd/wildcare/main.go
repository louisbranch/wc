package main

import (
	"log"

	"github.com/larissavoigt/wildcare/internal/http"
	"github.com/larissavoigt/wildcare/internal/mysql"
)

func main() {

	db, err := mysql.Open("wildcare")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	h := &http.Handler{
		UserService:    &mysql.UserService{db},
		SessionService: &mysql.SessionService{db},
	}

	log.Println("running on localhost:8080")

	err = h.ListenAndServe(":8080")

	if err != nil {
		log.Fatal(err)
	}
}
