package main

import (
	. "backend/api"
	"github.com/deCodedLife/gorest/rest"
	. "github.com/deCodedLife/gorest/tool"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strings"
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

func main() {

	if _, err := os.Stat("assets"); os.IsNotExist(err) {
		err := os.Mkdir("assets", 777)

		if err != nil {
			log.Println(err.Error())
			panic(err)
		}

		_ = os.Mkdir("assets/uploads", 777)
	}

	Handlers := rest.Construct()

	r := mux.NewRouter()

	for _, api := range Handlers {
		r.HandleFunc("/api/"+api.Path, api.Handler).Methods(api.Method)
	}

	InitRouters(r)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix(r.RequestURI, http.FileServer(http.Dir("./public"))).ServeHTTP(w, r)
	})

	r.PathPrefix("/api/assets").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/api/assets/", http.FileServer(http.Dir("./assets"))).ServeHTTP(w, r)
	})

	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prefix := r.RequestURI

		if len(strings.Split(r.RequestURI, ".")) > 1 {
			prefix = "/"
		}

		http.StripPrefix(prefix, http.FileServer(http.Dir("./public"))).ServeHTTP(w, r)
	})

	r.Use(CORS)
	err := http.ListenAndServeTLS(":443", "certificate.crt", "private.key", r)
	HandleError(err, CustomError{}.Unexpected(err))
}
