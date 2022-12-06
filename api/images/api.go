package api

import (
	"crypto/sha512"
	"errors"
	"fmt"
	. "github.com/deCodedLife/gorest/tool"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func UploadImage(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(5 << 20)
	HandleError(err, CustomError{}.WebError(w, http.StatusNotAcceptable, err))

	files := r.MultipartForm.File["image"]
	var output []string

	for _, handler := range files {

		file, err := handler.Open()
		HandleError(err, CustomError{}.Unexpected(err))

		defer file.Close()

		buff := make([]byte, 512)
		_, err = file.Read(buff)
		HandleError(err, CustomError{}.Unexpected(err))

		filetype := http.DetectContentType(buff)
		if len(strings.Split(filetype, "image")) == 1 {
			err := errors.New("only images expected")
			HandleError(err, CustomError{}.WebError(w, http.StatusNotAcceptable, err))
		}

		_, err = file.Seek(0, io.SeekStart)
		HandleError(err, CustomError{}.Unexpected(err))

		chipper := sha512.New()
		chipper.Write([]byte(fmt.Sprintf("%d_%d", time.Now().UnixNano(), filepath.Base(handler.Filename))))

		outputFileName := fmt.Sprintf("%x%s", chipper.Sum(nil), filepath.Ext(handler.Filename))

		serverFile, err := os.Create(fmt.Sprintf("./assets/%s", outputFileName))
		defer serverFile.Close()

		_, err = io.Copy(serverFile, file)
		HandleError(err, CustomError{}.Unexpected(err))

		output = append(output, outputFileName)

	}

	SendData(w, 200, output)

}
