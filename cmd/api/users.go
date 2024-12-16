package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"greenlight.vishaaxl.net/internal/data"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Password string `json:"password" validate:"required,min=6"`
		Name     string `json:"name" validate:"required,min=2"`
		Email    string `json:"email" validate:"required,email"`
	}

	err := app.readJSON(w, r, &input)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		errors := err.(validator.ValidationErrors)
		app.badRequestResponse(w, r, errors)
		return
	}

	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			app.badRequestResponse(w, r, fmt.Errorf("email already exists"))
			return
		default:
			app.serverErrorResponse(w, r, err)
		}
	}

	err = app.mailer.Send(user.Email, "Welcome to greenlight", "user_welcome.gohtml", user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Write a JSON response containing the user data along with a 201 Created status
	// code.
	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
