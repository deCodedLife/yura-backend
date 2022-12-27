package api

import (
	"github.com/gorilla/mux"
	"net/http"

	. "backend/api/exel"
	. "backend/api/images"
	. "backend/api/users"
)

func FileServer(r *mux.Router) {
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))))
}

func InitRouters(r *mux.Router) {
	r.HandleFunc("/sign-in", SignIn).Methods(http.MethodPost)
	r.HandleFunc("/images", UploadImage).Methods(http.MethodPost)
	r.HandleFunc("/exel", UploadTables).Methods(http.MethodPost)
}
