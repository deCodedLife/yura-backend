package api

import (
	"net/http"

	"github.com/gorilla/mux"

	. "backend/api/conditioners"
)

func FileServer(r *mux.Router) {
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))))
}

func InitRouters(r *mux.Router) {
	//r.HandleFunc("/auth", )
	r.HandleFunc("/conditioners/popular", GetPopular).Methods(http.MethodGet)
}
