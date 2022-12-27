package api

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"
	. "github.com/deCodedLife/gorest/database"
	. "github.com/deCodedLife/gorest/rest"
	. "github.com/deCodedLife/gorest/tool"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"
)

type AuthData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func SignIn(w http.ResponseWriter, r *http.Request) {

	var dbSchema Schema
	var authData AuthData

	defer func() {
		recover()
	}()

	byteData, err := ioutil.ReadFile(filepath.Join(SchemaDir, "users.json"))
	HandleError(err, CustomError{}.Unexpected(err))

	err = json.Unmarshal(byteData, &dbSchema)
	HandleError(err, CustomError{}.Unexpected(err))

	err = json.NewDecoder(r.Body).Decode(&authData)
	HandleError(err, CustomError{}.WebError(w, http.StatusNotAcceptable, err))

	data, err := selectUser(dbSchema, authData.Username)
	HandleError(err, CustomError{}.WebError(w, http.StatusInternalServerError, err))

	shaChipper := sha512.New()
	shaChipper.Write([]byte(authData.Password))

	authData.Password = fmt.Sprintf("%x", shaChipper.Sum(nil))
	userPassword := data[0]["password"]

	if authData.Password != userPassword {
		SendData(w, http.StatusNotAcceptable, "username or password is incorrect")
		return
	}

	token := fmt.Sprintf("%d%s%s", time.Now().Day(), authData.Username, userPassword)

	shaChipper = sha512.New()
	shaChipper.Write([]byte(token))

	SendData(w, 200, fmt.Sprintf("%x", shaChipper.Sum(nil)))
}
