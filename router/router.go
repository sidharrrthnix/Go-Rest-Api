package router

import (
	"net/http"

	"http-api.com/handlers"
)

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.RootHandler)

	// Single teacher operations (with ID) - MUST come first for specificity
	mux.HandleFunc("GET /teachers/{id}", handlers.GetTeacherByIdHandler)
	mux.HandleFunc("PUT /teachers/{id}", handlers.UpdateTeacherHandler)
	mux.HandleFunc("PATCH /teachers/{id}", handlers.PatchTeacherHandler)
	mux.HandleFunc("DELETE /teachers/{id}", handlers.DeleteTeacherHandler)

	// Batch operations (no trailing slash to avoid conflicts)
	mux.HandleFunc("GET /teachers", handlers.GetAllTeachersHandler)
	mux.HandleFunc("POST /teachers", handlers.AddTeacherHandler)
	mux.HandleFunc("PATCH /teachers", handlers.PatchTeachersHandler)
	mux.HandleFunc("DELETE /teachers", handlers.DeleteTeachersHandler)

	mux.HandleFunc("/students", handlers.StudentsHandler)
	mux.HandleFunc("/students/", handlers.StudentsHandler)

	return mux
}
