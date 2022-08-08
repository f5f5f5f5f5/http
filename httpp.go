package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Id      uint16
	Name    string
	Company string
	Mark    string
}

var users = []User{}

func home_page(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/home.html", "templates/head.html", "templates/footer.html")
	tmpl.ExecuteTemplate(w, "home", nil)
}

func aboutme_page(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/aboutme.html", "templates/head.html", "templates/footer.html")
	tmpl.ExecuteTemplate(w, "aboutme", nil)
}

func poll_page(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/poll.html", "templates/head.html", "templates/footer.html")
	tmpl.ExecuteTemplate(w, "poll", nil)
}

func contacts_page(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/contacts.html", "templates/head.html", "templates/footer.html")
	tmpl.ExecuteTemplate(w, "contacts", nil)
}

func results(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	company := r.FormValue("company")
	mark := r.FormValue("mark")

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/http")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	input, err := db.Query(fmt.Sprintf("INSERT INTO users (name, company, mark) VALUES('%s', '%s', '%s')", name, company, mark))

	if err != nil {
		panic(err)
	}

	defer input.Close()

	http.Redirect(w, r, "/contacts/", http.StatusSeeOther)
}

func login_page(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/login.html", "templates/head.html", "templates/footer.html")
	tmpl.ExecuteTemplate(w, "login", nil)
}

func admin_page(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/admin.html", "templates/head.html", "templates/footer.html")

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/http")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	out, err := db.Query("SELECT * FROM users")

	if err != nil {
		panic(err)
	}

	users = []User{}

	for out.Next() {
		var user User
		err = out.Scan(&user.Id, &user.Name, &user.Company, &user.Mark)
		if err != nil {
			panic(err)
		}

		users = append(users, user)
	}

	tmpl.ExecuteTemplate(w, "admin", users)
}

func checkin(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	password := r.FormValue("password")
	if login == "admin" && password == "admin" {
		http.Redirect(w, r, "/admin/", http.StatusSeeOther)
	}
}

func HandleRequest() {
	r := mux.NewRouter()
	r.HandleFunc("/", home_page)
	r.HandleFunc("/home/", home_page)
	r.HandleFunc("/about/", aboutme_page)
	r.HandleFunc("/poll/", poll_page)
	r.HandleFunc("/contacts/", contacts_page)
	r.HandleFunc("/results/", results).Methods("POST")
	r.HandleFunc("/login/", login_page)
	r.HandleFunc("/checkin/", checkin).Methods("POST")
	r.HandleFunc("/admin/", admin_page)
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
func main() {
	HandleRequest()
}
