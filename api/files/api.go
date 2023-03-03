package api

import (
	"encoding/json"
	"net/http"

	. "github.com/deCodedLife/gorest/tool"
)

const UPLOADS_FOLDER = "assets/uploads"

type FilesRequest struct {
	Path string `json:"path"`
}

func LS(w http.ResponseWriter, r *http.Request) {
	var request FilesRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	HandleError(err, CustomError{}.WebError(w, http.StatusNotAcceptable, err))

	folderContent, err := ListDirectory(UPLOADS_FOLDER + request.Path)
	HandleError(err, CustomError{}.WebError(w, http.StatusInternalServerError, err))

	SendData(w, http.StatusOK, folderContent)
}
