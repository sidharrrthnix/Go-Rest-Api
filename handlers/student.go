package handlers

import (
	"fmt"
	"net/http"
)

func StudentsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, Students!")
}
