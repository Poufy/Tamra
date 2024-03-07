package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
)

type DocsRouter struct {
	logger logrus.FieldLogger
}

func NewDocsRouter(logger logrus.FieldLogger) *DocsRouter {
	return &DocsRouter{logger: logger}
}

func (router *DocsRouter) GetRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/swagger.json", router.swaggerJSON)
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/api/v1/docs/swagger.json"),
	))
	return r
}

func swaggerhttpRoute() chi.Router {
	r := chi.NewRouter()
	r.Get("/*", httpSwagger.Handler(
		httpSwagger.URL("/api/v1/docs/swagger.json"),
	))
	return r
}

func (router *DocsRouter) swaggerJSON(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "docs/swagger.json")
}

// func docsServeRoute() chi.Router {
// 	r := chi.NewRouter()
// 	// Currently we only have the swagger documentation as static content
// 	// Use the relative path to the docs directory since it is not in the same directory as the main.go file
// 	fileServer := http.FileServer(http.Dir("../../docs"))

// 	// Strip the /api/v1/docs prefix from the URL before serving the files, since the resulting URL should be the file name
// 	r.Handle("/*", http.StripPrefix("/api/v1/docs", fileServer))
// 	return r
// }
