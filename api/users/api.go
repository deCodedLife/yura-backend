package api

import . "github.com/deCodedLife/gorest/database"

func selectUser(dbSchema Schema, username string) ([]map[string]interface{}, error) {
	var userGetRequest = map[string]interface{}{
		"username": username,
	}

	return dbSchema.SELECT(userGetRequest)
}
