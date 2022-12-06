package api

import (
	. "backend/api/images"
	. "backend/api/users"
	"net/http"

	"github.com/gorilla/mux"
)

func FileServer(r *mux.Router) {
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))))
}

func InitRouters(r *mux.Router) {
	r.HandleFunc("/sign-in", SignIn).Methods(http.MethodPost)
	r.HandleFunc("/images", UploadImage).Methods(http.MethodPost)
}
