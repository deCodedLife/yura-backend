package api

import (
	"encoding/json"
	. "github.com/deCodedLife/gorest/tool"
	"net/http"
	"os"
)

func readMenu() (Menu, error) {

	var menuObject Menu
	menuFile, err := os.ReadFile("menu.json")
	err = json.Unmarshal(menuFile, &menuObject)

	if err != nil {
		return menuObject, err
	}

	return menuObject, nil

}

func HandleMenuRequest(w http.ResponseWriter, r *http.Request) {

	menu, err := readMenu()
	HandleError(err, CustomError{}.WebError(w, http.StatusInternalServerError, err))

	SendData(w, http.StatusOK, menu)

}
