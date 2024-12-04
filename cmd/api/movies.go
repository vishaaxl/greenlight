package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"greenlight.vishaaxl.net/internal/data"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string   `json:"title" validate:"required"`
		Year    int      `json:"year" validate:"required,min=1888"`
		Runtime int32    `json:"runtime" validate:"required"`
		Genres  []string `json:"genres" validate:"required,unique"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// validate request body
	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		errors := err.(validator.ValidationErrors)
		app.badRequestResponse(w, r, errors)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) getMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIdParam(r)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Casablanca",
		Runtime:   102,
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"data": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
