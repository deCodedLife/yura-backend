package main

import (
	. "backend/api"
	"github.com/deCodedLife/gorest/rest"
	. "github.com/deCodedLife/gorest/tool"
	"github.com/gorilla/mux"
	"net/http"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Headers:", "*")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")

		next.ServeHTTP(w, r)
		return
	})
}

func main() {
	Handlers := rest.Construct()

	r := mux.NewRouter()
	r.Use(CORS)

	for _, api := range Handlers {
		r.HandleFunc("/"+api.Path, api.Handler).Methods(api.Method)
	}

	FileServer(r)
	InitRouters(r)

	//err := http.ListenAndServe(":8080", r)
	err := http.ListenAndServeTLS(":443", "certificate.crt", "private.key", r)
	HandleError(err, CustomError{}.Unexpected(err))
}
