package api

import (
	json2 "encoding/json"
	"github.com/deCodedLife/gorest/database"
	"github.com/deCodedLife/gorest/rest"
	. "github.com/deCodedLife/gorest/tool"
	"net/http"
	"net/url"
	"strings"
)

type IRequest struct {
	Object string `json:"object"`
}

func getObjectRef(object string, variables url.Values) (interface{}, error) {
	var uriParams = make(map[string]interface{})
	schemes, err := rest.GetSchemas()
	if err != nil {
		return nil, err
	}
	for _, schema := range schemes {
		if schema.Table != object {
			continue
		}
		var relatedObjects []database.SchemaParam

		for _, param := range schema.Params {
			var valueExists bool

			if param.TakeFrom != "" {
				relatedObjects = append(relatedObjects, param)
			}

			for variable := range variables {
				if variable == param.Article {

					value := variables.Get(variable)

					if value == "" {
						break
					}

					valueExists = true
					break
				}
			}

			if valueExists == false {
				continue
			}

			uriParams[param.Article] = variables.Get(param.Article)
		}
		data, err := schema.SELECT(uriParams)
		object := data[0]

		if err != nil {
			return nil, err
		}

		for _, param := range relatedObjects {
			relatedObject := strings.Split(param.TakeFrom, "/")[0]
			relatedField := strings.Split(param.TakeFrom, "/")[1]

			for _, scheme := range schemes {
				if scheme.Table != relatedObject {
					continue
				}

				request := make(map[string]interface{})
				request[relatedField] = object[param.Article]
				relatedList, err := scheme.SELECT(request)
				return request, nil
				related := relatedList[0]

				if err != nil {
					return nil, err
				}

				object[param.Article] = related[relatedField]
			}
		}
		return object, nil
	}
	return nil, nil
}

func HandleObjectRef(w http.ResponseWriter, r *http.Request) {

	var userRequest IRequest

	err := json2.NewDecoder(r.Body).Decode(&userRequest)
	HandleError(err, CustomError{}.WebError(w, http.StatusNotAcceptable, err))

	variables := r.URL.Query()

	object, err := getObjectRef(userRequest.Object, variables)
	HandleError(err, CustomError{}.WebError(w, http.StatusInternalServerError, err))

	SendData(w, http.StatusOK, object)

}
