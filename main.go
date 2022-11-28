package main

import (
	. "backend/api"
	"github.com/deCodedLife/gorest/rest"
	. "github.com/deCodedLife/gorest/tool"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	Handlers := rest.Construct()

	r := mux.NewRouter()

	for _, api := range Handlers {
		r.HandleFunc("/"+api.Path, api.Handler).Methods(api.Method)
	}

	FileServer(r)
	InitRouters(r)

	err := http.ListenAndServe(":8080", r)
	HandleError(err, CustomError{}.Unexpected(err))
}
