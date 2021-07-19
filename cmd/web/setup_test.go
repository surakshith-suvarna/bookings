package main

import (
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	os.Exit(m.Run())
}

//An object which satisfies http handler interface
type myHandler struct{}

func (mh *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
