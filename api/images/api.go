package api

import (
	"crypto/sha512"
	"fmt"
	. "github.com/deCodedLife/gorest/tool"
	"io"
	"net/http"
	"os"
	"strings"
	"unicode/utf8"
)

func splitFileName(f string) (string, string) {
	namePieces := strings.Split(f, ".")
	extension := namePieces[len(namePieces)]
	extensionLength := utf8.RuneCountInString(extension)
	fileName := f[:len(namePieces)-extensionLength]
	return fileName, extension
}

func UploadImage(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("image")
	HandleError(err, CustomError{}.WebError(w, http.StatusNotAcceptable, err))

	defer file.Close()

	fileName, extension := splitFileName(handler.Filename)

	chipper := sha512.New()
	chipper.Write([]byte(fileName))

	outputFileName := fmt.Sprintf("%x.%s", chipper.Sum(nil), extension)

	imageFile, err := os.OpenFile(outputFileName, os.O_WRONLY|os.O_CREATE, 0666)
	HandleError(err, CustomError{}.Unexpected(err))

	defer imageFile.Close()

	_, err = io.Copy(imageFile, file)
	HandleError(err, CustomError{}.Unexpected(err))

	SendData(w, 200, outputFileName)
}
