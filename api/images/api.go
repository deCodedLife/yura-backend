package api

import (
	"errors"
	"net/http"

	. "backend/api/files"
	. "github.com/deCodedLife/gorest/tool"
)

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
