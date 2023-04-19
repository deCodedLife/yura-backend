package api

import (
	"errors"
	"fmt"
	"github.com/deCodedLife/gorest/database"
	"github.com/deCodedLife/gorest/rest"
	"github.com/xuri/excelize/v2"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"

	. "backend/api/files"
	. "github.com/deCodedLife/gorest/tool"
)

func toChar(i int) string {
	return strconv.QuoteRune(rune('A' - 1 + i))
}

func columnName(i int) string {
	if i > 26 {
		subColumn := math.RoundToEven(math.Floor(float64(i / 26)))
		column := i - int(26*int(subColumn))
		return fmt.Sprintf("%s%s", toChar(int(subColumn)), toChar(column))
	}

	return toChar(i)
}

func rawObjects(filesList []string, schemas []database.Schema) ([]map[string]interface{}, error) {
	var output []map[string]interface{}
	usedSchemas, err := rest.GetSchemas()

	if err != nil {
		return nil, err
	}

	for index, fileName := range filesList {
		schema := schemas[index]

		file, err := excelize.OpenFile("assets/" + fileName)
		if err != nil {
			return nil, err
		}

		defer file.Close()
		defer os.Remove("assets/" + fileName)

		dataSheet := file.GetSheetList()[0]
		rowsList, err := file.GetRows(dataSheet)
		if err != nil {
			return nil, err
		}

		for index, row := range rowsList {

			object := make(map[string]interface{})

			// Skip titles
			if index == 0 {
				continue
			}

			for rowIndex, rowContent := range row {

				var content interface{}
				content = rowContent

				if schema.Params[rowIndex].Article == "id" {
					continue
				}

				//if schema.Params[rowIndex].DisplayType == "image" {
				//	content, err = DownloadFile(rowContent)
				//
				//	if err != nil {
				//		return nil, err
				//	}
				//}

				if rowContent == "" && schema.Params[rowIndex].Null == "NO" {
					if schema.Params[rowIndex].Default == "" {
						errorMessage := fmt.Sprintf("отсутствует важный параметр %s столбец %s строка %d",
							schema.Params[rowIndex].Title, columnName(rowIndex), index+1)
						return nil, errors.New(errorMessage)
					}
					continue
				}

				dependedParams := schema.Params[rowIndex].TakeFrom

				// Check is Dependent
				if dependedParams != "" {
					var objectSchema database.Schema

					dependedObject := strings.Split(dependedParams, "/")[0]
					dependedArticle := strings.Split(dependedParams, "/")[1]

					for _, schema := range usedSchemas {
						if dependedObject != schema.Table {
							continue
						}
						objectSchema = schema
						break
					}

					dependedValue, err := objectSchema.SELECT(map[string]interface{}{
						dependedArticle: content,
					})

					if err != nil {
						return nil, err
					}
					if len(dependedValue) < 1 {
						notFoundError := fmt.Sprintf("Can't find %s/%s as %s", schema.Params[rowIndex].Article, dependedArticle, content)
						notFoundError = fmt.Sprintf("%s. At row %d %s", notFoundError, rowIndex, rowContent)
						return nil, errors.New(notFoundError)
					}
					content = dependedValue[0]["id"]
				}

				object[schema.Params[rowIndex].Article] = content
			}

			output = append(output, object)
		}
	}

	return output, nil
}

func fromSchemasList(textSchemas []string) ([]database.Schema, error) {
	var outSchemes []database.Schema
	restSchemas, err := rest.GetSchemas()

	if err != nil {
		return nil, err
	}

	for _, textSchema := range textSchemas {
		for _, schema := range restSchemas {
			if schema.Table == textSchema {
				outSchemes = append(outSchemes, schema)
			}
		}
	}

	return outSchemes, nil
}

func UploadTables(w http.ResponseWriter, r *http.Request) {

	defer func() {
		err := recover()
		log.Println(err)
		SendData(w, http.StatusOK, err)
	}()

	err := r.ParseMultipartForm(5 << 20)
	HandleError(err, CustomError{}.WebError(w, http.StatusNotAcceptable, err))

	files := r.MultipartForm.File["tables"]
	objects := r.MultipartForm.Value["objects"]

	if len(files) != len(objects) {
		err := errors.New("mismatch. files count should be equal objects count")
		HandleError(err, CustomError{}.WebError(w, http.StatusNotAcceptable, err))
	}

	fileList, err := HandleFile(files, FileConfigs{
		FileSubType: []string{"vnd.openxmlformats-officedocument.spreadsheetml.sheet", "application/vnd.ms-excel"},
		FileType:    "application",
		TypeError:   errors.New("принимаются только exel таблицы"),
	})
	HandleError(err, CustomError{}.WebError(w, http.StatusNotAcceptable, err))

	expectedSchemas, err := fromSchemasList(objects)
	HandleError(err, CustomError{}.WebError(w, http.StatusInternalServerError, err))

	if len(expectedSchemas) == 0 {
		err := errors.New("объекты не передались")
		HandleError(err, CustomError{}.WebError(w, http.StatusNotAcceptable, err))
	}

	rawObjects, err := rawObjects(fileList, expectedSchemas)
	HandleError(err, CustomError{}.WebError(w, http.StatusNotAcceptable, err))

	for _, schema := range expectedSchemas {
		for _, object := range rawObjects {
			_, err := schema.INSERT(object)
			HandleError(err, CustomError{}.WebError(w, http.StatusInternalServerError, err))
		}
	}

	SendData(w, http.StatusOK, err)
}
