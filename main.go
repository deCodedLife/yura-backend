package main

import (
	. "backend/api"
	"github.com/deCodedLife/gorest/rest"
	. "github.com/deCodedLife/gorest/tool"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
)

func main() {

	if _, err := os.Stat("assets"); os.IsNotExist(err) {
		err := os.Mkdir("assets", 777)
		log.Println(err.Error())
		panic(err)
	}

	Handlers := rest.Construct()

	r := mux.NewRouter().StrictSlash(true)

	for _, api := range Handlers {
		r.HandleFunc("/api/"+api.Path, api.Handler).Methods(api.Method)
	}

	FileServer(r)
	InitRouters(r)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}).Handler(r)

	//err := http.ListenAndServe(":8080", r) // Test server
	err := http.ListenAndServeTLS(":443", "certificate.crt", "private.key", c)
	HandleError(err, CustomError{}.Unexpected(err))
}
