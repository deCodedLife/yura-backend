package api

import (
	"github.com/gorilla/mux"
	"net/http"

	. "backend/api/exel"
	. "backend/api/files"
	. "backend/api/users"
)

func InitRouters(r *mux.Router) {
	r.HandleFunc("/api/sign-in", SignIn).Methods(http.MethodPost)
	r.HandleFunc("/api/images", UploadImage).Methods(http.MethodPost)
	r.HandleFunc("/api/files", UploadFile).Methods(http.MethodPost)
	r.HandleFunc("/api/deleteFile", DeleteFile).Methods(http.MethodPost)
	r.HandleFunc("/api/createFolder", CreateDirectory).Methods(http.MethodPost)
	r.HandleFunc("/api/exel", UploadTables).Methods(http.MethodPost)
	r.HandleFunc("/api/fm/ls", LS).Methods(http.MethodPost)
}
