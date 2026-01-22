package router

import (
	"net/http"
)

func MainRouter() *http.ServeMux {
	tRouter := teacherRouter()
	sRouter := studentRouter()
	tRouter.Handle("/", sRouter)
	return tRouter
}
