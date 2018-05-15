// Package handlers provides request handlers.
package handlers

import (
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
	"reflect"
)

func getIdFromPath(w http.ResponseWriter, r *http.Request) (int64, error) {
	idString := mux.Vars(r)["id"]
	if idString == "" {
		return -1, errors.New("user id cannot be empty.")
	}

	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		return -1, err
	}

	return id, nil
}

//****big bug with golang date format parsing ***
//https://stackoverflow.com/questions/14106541/go-parsing-date-time-strings-which-are-not-standard-formats
func ConvertFormDate(value string) reflect.Value {
	if v, err := time.Parse("01/02/2006", value); err == nil {
		return reflect.ValueOf(v)
	} else if v, err := time.Parse("2006-01-02 00:00:00 +0000 UTC", value); err == nil {
		return reflect.ValueOf(v)
	}
	return reflect.Value{} // this is the same as the private const invalidType
}