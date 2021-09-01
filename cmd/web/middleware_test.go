package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {

	var myhandler myHandler

	h := NoSurf(&myhandler)

	switch v := h.(type) {

	case http.Handler:
		// do nothing
	default:
		t.Error(fmt.Sprintf("type is not http Handler but is %T", v))
	}

}

func TestSession(t *testing.T) {

	var myhandler myHandler

	h := SessionLoad(&myhandler)

	switch v := h.(type) {

	case http.Handler:
		// do nothing
	default:
		t.Error(fmt.Sprintf("type is not http Handler but is %T", v))
	}

}
