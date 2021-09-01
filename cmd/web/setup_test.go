package main

import (
	"net/http"
	"os"
	"testing"
)

// This code will run before every testcase and setups an env for all the testcases
func TestMain(m *testing.M) {

	os.Exit(m.Run())
}

// creating custom http handler for testcases which requires http.Handler
type myHandler struct{}

func (mh *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	
}