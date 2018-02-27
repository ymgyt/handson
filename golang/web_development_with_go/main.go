package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/YmgchiYt/golang-handson/web_development_with_go/controllers"
	"github.com/YmgchiYt/golang-handson/web_development_with_go/middleware"
	"github.com/YmgchiYt/golang-handson/web_development_with_go/models"
	"github.com/YmgchiYt/golang-handson/web_development_with_go/views"
	"github.com/gorilla/mux"
	"github.com/k0kubun/pp"
)

const (
	//host = "postgres"
	host = "localhost"
	port = 5432
	user = "postgres"
	//user     = "me"
	password = ""
	dbname   = "goweb_dev"
)

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1> Not Found !!!</h1>")
}

var (
	homeView    *views.View
	contactView *views.View
	faqView     *views.View
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%q dbname=%s sslmode=disable", host, port, user, password, dbname)
	services, err := models.NewServices(psqlInfo)
	if err != nil {
		pp.Println(err)
		panic(err)
	}
	defer services.Close()
	services.AutoMigrate()

	r := mux.NewRouter()
	staticCtl := controllers.NewStatic()
	usersCtl := controllers.NewUsers(services.User)
	galleriesCtl := controllers.NewGalleries(services.Gallery, r)

	requireUserMw := middleware.RequireUser{
		UserService: services.User,
	}
	newGallery := requireUserMw.Apply(galleriesCtl.New)
	createGallery := requireUserMw.ApplyFn(galleriesCtl.Create)

	r.Handle("/", staticCtl.Home).Methods("GET")
	r.Handle("/contact", staticCtl.Contact).Methods("GET")
	r.HandleFunc("/signup", usersCtl.New).Methods("GET")
	r.HandleFunc("/signup", usersCtl.Create).Methods("POST")
	r.Handle("/login", usersCtl.LoginView).Methods("GET")
	r.HandleFunc("/login", usersCtl.Login).Methods("POST")
	r.HandleFunc("/cookietest", usersCtl.CookieTest).Methods("GET")
	// Gallery routes
	r.HandleFunc("/galleries", requireUserMw.ApplyFn(galleriesCtl.Index)).Methods("GET")
	r.Handle("/galleries/new", newGallery).Methods("GET")
	r.HandleFunc("/galleries", createGallery).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}", galleriesCtl.Show).Methods("GET").Name(controllers.ShowGallery)
	r.HandleFunc("/galleries/{id:[0-9]+}/edit", requireUserMw.ApplyFn(galleriesCtl.Edit)).Methods("GET")
	r.HandleFunc("/galleries/{id:[0-9]+}/update", requireUserMw.ApplyFn(galleriesCtl.Update)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/delete", requireUserMw.ApplyFn(galleriesCtl.Delete)).Methods("POST")

	r.NotFoundHandler = http.HandlerFunc(notFound)
	log.Fatal(http.ListenAndServe(":3000", r))
}
