package api

import (
	"avito/internal/db"
	"avito/internal/errors"
	"avito/internal/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

type Router struct {
	Router *mux.Router
}

func NewRouter(storage *db.PostgresStorage) *Router {
	router := &Router{Router: mux.NewRouter()}
	router.InitRoutes(storage)
	return router
}

func (self Router) InitRoutes(storage *db.PostgresStorage) {
	self.CommonRoutes()
	self.TenderRoutes(storage)
	self.BidRoutes(storage)
}

func (self Router) createRoute(path string, apiFunc apiFunc) *mux.Route {
	return self.Router.HandleFunc(
		"/api"+path,
		func(w http.ResponseWriter, r *http.Request) {
			context := handlers.NewContext(w, r)
			if err := apiFunc(context); err != nil {
				println(err.Error())
				_ = context.RespondWithJson(err.Code, errors.ErrorResponse{Reason: err.Message})
			}
		},
	)
}

type apiFunc func(context *handlers.Context) *errors.AppError
