package utils

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"time"
)

var MSKLoc *time.Location

func init() {
	MSKLoc, _ = time.LoadLocation("Europe/Moscow")
}

func RunFileServer(filename string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		file, err := os.Open(filename)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		_, err = io.Copy(w, file)
		if err != nil {
			panic(fmt.Errorf("copying error: %w", err))
		}
	}))
}
