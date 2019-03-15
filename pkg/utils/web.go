package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type WebAPIError struct {
	ErrorString string `json:"error"`
	ErrorCode   int    `json:"code"`
}

func NewWebAPIError(code int, errorString string) WebAPIError {
	return WebAPIError{
		ErrorString: errorString,
		ErrorCode:   code,
	}
}

func WebWrite(w http.ResponseWriter, data interface{}) {
	bs, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Error in web write: %s\n", err)
		WebWriteError(w, 500, "internal server error")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bs)
}

func WebWriteError(w http.ResponseWriter, code int, errorString string) {
	msg := NewWebAPIError(code, errorString)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	WebWrite(w, msg)
}
