package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type WebAPIError struct {
	ErrorString string `json:"error"`
}

func NewWebAPIError(errorString string) WebAPIError {
	return WebAPIError{
		ErrorString: errorString,
	}
}

func WebWrite(w http.ResponseWriter, data interface{}) {
	bs, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Error in web write: %s", err)
		bs, _ = json.Marshal(NewWebAPIError("internal server error"))
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bs)
}
