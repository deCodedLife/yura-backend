package main

import (
	. "backend/api"
	"github.com/deCodedLife/gorest/rest"
	. "github.com/deCodedLife/gorest/tool"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Access-Control-Allow-Headers:", "*")
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "*")

		next.ServeHTTP(w, r)
		return
	})
}

func angularHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./public/index.html")
}

func main() {

	if _, err := os.Stat("assets"); os.IsNotExist(err) {
		err := os.Mkdir("assets", 777)
		log.Println(err.Error())
		panic(err)
	}

	Handlers := rest.Construct()

	r := mux.NewRouter()
	r.Host("api.localhost")

	for _, api := range Handlers {
		r.HandleFunc("/api/"+api.Path, api.Handler).Methods(api.Method)
	}

	r.Handle("/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	r.NotFoundHandler = http.HandlerFunc(angularHandler)
	//r.PathPrefix("/assets/").Handler(http.StripPrefix("/public/assets/", http.FileServer(http.Dir("/public/assets/"))))

	FileServer(r)
	InitRouters(r)

	r.Use(CORS)
	//err := http.ListenAndServe(":8080", r) // TEST env
	err := http.ListenAndServeTLS(":443", "certificate.crt", "private.key", r)
	HandleError(err, CustomError{}.Unexpected(err))
}
