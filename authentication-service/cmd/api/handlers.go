package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/vbrenister/apicommon"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var restPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.ReadJSON(w, r, &restPayload)
	if err != nil {
		app.ErroJSON(w, err, http.StatusBadRequest)
		return
	}

	user, err := app.Models.Users.GetByEmail(restPayload.Email)
	if err != nil {
		app.ErroJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	valid, err := user.PasswordMatches(restPayload.Password)
	if err != nil || !valid {
		app.ErroJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	payload := apicommon.JsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in %s", user.Email),
		Data:    user,
	}

	app.WriteJSON(w, http.StatusAccepted, payload)
}
