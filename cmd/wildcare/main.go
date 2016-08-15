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
		UserService:    mysql.UserService{db},
		SessionService: mysql.SessionService{db},
	}

	h.ListenAndServe(":8080")
}
