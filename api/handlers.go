package api

import (
	"github.com/gorilla/mux"
	"net/http"

	. "backend/api/exel"
	. "backend/api/images"
	. "backend/api/users"
)

func FileServer(r *mux.Router) {
	//http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/")))
	r.PathPrefix("/assets/").Handler(http.FileServer(http.Dir("assets/"))).Host("api.klimsystems.ru")
	r.PathPrefix("/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/")))).Host("api.klimsystems.ru")
}

func InitRouters(r *mux.Router) {
	r.HandleFunc("/sign-in", SignIn).Methods(http.MethodPost).Host("api.klimsystems.ru")
	r.HandleFunc("/images", UploadImage).Methods(http.MethodPost).Host("api.klimsystems.ru")
	r.HandleFunc("/exel", UploadTables).Methods(http.MethodPost).Host("api.klimsystems.ru")
}
