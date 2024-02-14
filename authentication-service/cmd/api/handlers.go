package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

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

	err = app.logRequest("auth", fmt.Sprintf("Logged in %s", user.Email))
	if err != nil {
		app.ErroJSON(w, err, http.StatusInternalServerError)
		return
	}

	payload := apicommon.JsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in %s", user.Email),
		Data:    user,
	}

	app.WriteJSON(w, http.StatusAccepted, payload)
}

func (app *Config) logRequest(name string, data string) error {
	logServiceUrl := os.Getenv("LOG_SERVICE_URL")

	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	jsonData, err := json.MarshalIndent(entry, "", "\t")
	if err != nil {
		return err
	}
	_, err = http.Post(fmt.Sprintf("%s/log", logServiceUrl), "application/json", bytes.NewBuffer(jsonData))
	return err
}
