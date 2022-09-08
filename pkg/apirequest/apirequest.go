package apirequest

import (
	"io"
	"net/http"
)

// HTTPRequest requests the given url
func HTTPRequest(url string) ([]byte, error) {
	req, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	bs, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	return bs, nil
}
