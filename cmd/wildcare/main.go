package main

import (
	"log"

	"github.com/larissavoigt/wildcare/internal/mysql"
)

func main() {

	db, err := mysql.Open("wildcare")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	us := &mysql.UserService{DB: db}

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

func index(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if r.Method != "GET" || path != "/" {
		view.NotFound(w)
		return
	}

	content := struct {
		User *user.User
	}{}

	u, ok := session.CurrentUser(r.Context())

	if ok {
		content.User = u
	}

	view.Render(w, "home/index", content)
}

func signup(w http.ResponseWriter, r *http.Request) {
	_, ok := session.CurrentUser(r.Context())

	if ok {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	content := struct{ Error error }{}

	switch r.Method {
	case "GET":
		view.Render(w, "user/signup", content)
	case "POST":
		r.ParseForm()
		email := r.Form.Get("email")
		password := r.Form.Get("password")

		user, err := user.Create(email, password)
		if err != nil {
			content.Error = err
			view.Render(w, "user/signup", content)
			return
		}

		_, err = session.Create(w, user.ID)
		if err != nil {
			content.Error = err
			view.Render(w, "user/signup", content)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	default:
		view.NotFound(w)
	}
}
func login(w http.ResponseWriter, r *http.Request) {
	_, ok := session.CurrentUser(r.Context())

	if ok {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	content := struct{ Error error }{}

	switch r.Method {
	case "GET":
		view.Render(w, "user/login", content)
	case "POST":
		r.ParseForm()
		email := r.Form.Get("email")
		password := r.Form.Get("password")

		u, ok := user.Authenticate(email, password)

		if !ok {
			content.Error = errors.New("Email and password combination doesn't match")
			view.Render(w, "user/login", content)
			return
		}

		_, err := session.Create(w, u.ID)
		if err != nil {
			content.Error = errors.New("Failed to create session")
			view.Render(w, "user/login", content)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	default:
		view.NotFound(w)
	}
}
*/
