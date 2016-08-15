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

/*

func main() {
	http.Handle("/signup", session.Middleware(signup))
	http.Handle("/login", session.Middleware(login))

	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		session.Destroy(w, r)
		http.Redirect(w, r, "/", http.StatusFound)
	})

	http.Handle("/", session.Middleware(index))

	http.ListenAndServe(":8080", nil)
}

*/
