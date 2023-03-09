package api

import (
	"encoding/json"
	"errors"
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

func UploadFile(w http.ResponseWriter, r *http.Request) {

	defer func() {
		recover()
	}()

	err := r.ParseMultipartForm(5 << 20)
	HandleError(err, CustomError{}.WebError(w, http.StatusNotAcceptable, err))

	files := r.MultipartForm.File["file"]

	filesList, err := HandleFile(files, FileConfigs{
		FileType:    "",
		FileSubType: nil,
		TypeError:   errors.New("something went wrong"),
		SavePath:    r.MultipartForm.Value["path"][0] + "/",
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
		SavePath:    "",
	})

	HandleError(err, CustomError{}.WebError(w, http.StatusNotAcceptable, err))
	SendData(w, http.StatusOK, filesList)

}
