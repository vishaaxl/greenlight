package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

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

	err = app.models.Permissions.AddForUser(user.ID, "movies:read")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.background(func() {
		data := map[string]interface{}{
			"activationToken": token.PlainText,
			"Name":            user.Name,
		}

		err = app.mailer.Send(user.Email, "Welcome to greenlight", "user_welcome.gohtml", data)
		if err != nil {
			app.logger.PrintError(err, nil)
		}
	})

	// Write a JSON response containing the user data along with a 201 Created status
	// code.
	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TokenPlainText string `json:"token" validate:"required,min=26,max=26"`
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

	user, err := app.models.Users.GetForToken(data.ScopeActivation, input.TokenPlainText)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.errorResponse(w, r, http.StatusBadRequest, "invalid or expired token")
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	user.Activated = true
	err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	err = app.models.Tokens.DeleteAllForUser(data.ScopeActivation, user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
