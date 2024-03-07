package routes

import (
	"Tamra/internal/app/tamra/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type UserRouter struct {
	userHandler   *handlers.UserHandler
	authMiddlware func(http.Handler) http.Handler
	logger        logrus.FieldLogger
}

func NewUserRouter(userHandler *handlers.UserHandler, authMiddleware func(http.Handler) http.Handler, logger logrus.FieldLogger) *UserRouter {
	return &UserRouter{userHandler: userHandler, authMiddlware: authMiddleware, logger: logger}
}

// ? Is this an overkill? Should we just return a chi.Router instead of a chi.Router wrapped in a Router struct?
// ? Should we just pass the handlers to GetRouter instead of creating a new Router struct?
func (router *UserRouter) GetRouter() chi.Router {
	r := chi.NewRouter()
	// Use the authMiddleware for all routes in the user route
	// Middleware checks if the token is valid and if it is, it will call the next handler in the chain
	// It will also append the UUID of the user to the request context so we can use it in the handler
	r.Use(router.authMiddlware)
	r.Post("/", router.userHandler.CreateUser)
	r.Get("/", router.userHandler.GetUsers)
	r.Get("/{id}", router.userHandler.GetUser)
	r.Patch("/{id}", router.userHandler.UpdateUser)
	return r
}
