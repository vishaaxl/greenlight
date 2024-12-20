package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthCheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.getMovieHandler)

	router.HandlerFunc(http.MethodPost, "/v1/movies", app.createMovieHandler)
	router.HandlerFunc(http.MethodGet, "/v1/movies", app.listMoviesHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.deleteMovieHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.updateMovieHandler)

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/users/activate", app.activateUserHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationHandler)

	return app.recoverPanic(app.rateLimit(app.authenticate(router)))
}
