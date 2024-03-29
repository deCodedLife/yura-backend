package api

import (
	"encoding/json"
	"errors"
	. "github.com/deCodedLife/gorest/tool"
	"net/http"
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

func CreateDirectory(w http.ResponseWriter, r *http.Request) {

	var request FilesRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	HandleError(err, CustomError{}.WebError(w, http.StatusNotAcceptable, err))

	err = MkDir(request.Path)
	HandleError(err, CustomError{}.WebError(w, http.StatusInternalServerError, err))

	SendData(w, http.StatusOK, nil)

}

func DeleteFile(w http.ResponseWriter, r *http.Request) {

	var request FilesRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	HandleError(err, CustomError{}.WebError(w, http.StatusNotAcceptable, err))

	err = RemoveFile(request.Path)
	HandleError(err, CustomError{}.WebError(w, http.StatusInternalServerError, err))

	SendData(w, http.StatusOK, nil)

}

func UploadFile(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(1024 << 20)
	HandleError(err, CustomError{}.WebError(w, http.StatusNotAcceptable, err))

	filesList, err := HandleFile(r.MultipartForm.File["file"], FileConfigs{
		FileType: "",
		SavePath: r.MultipartForm.Value["path"],
	})

	HandleError(err, CustomError{}.WebError(w, http.StatusNotAcceptable, err))

	SendData(w, http.StatusOK, filesList)

}

func UploadImage(w http.ResponseWriter, r *http.Request) {

	defer func() {
		recover()
	}()

	err := r.ParseMultipartForm(5 << 20)
	HandleError(err, CustomError{}.WebError(w, http.StatusNotAcceptable, err))

	files := r.MultipartForm.File["image"]

	filesList, err := HandleFile(files, FileConfigs{
		FileType:    "image",
		FileSubType: nil,
		TypeError:   errors.New("only images expected"),
	})

	HandleError(err, CustomError{}.WebError(w, http.StatusNotAcceptable, err))
	SendData(w, http.StatusOK, filesList)

}
