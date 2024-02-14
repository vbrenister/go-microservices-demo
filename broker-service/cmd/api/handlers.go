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

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := apicommon.JsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.WriteJSON(w, http.StatusAccepted, payload)
}

func (app *Config) HandleSumbission(w http.ResponseWriter, r *http.Request) {
	var restPayload RequestPayload

	err := app.ReadJSON(w, r, &restPayload)
	if err != nil {
		app.ErroJSON(w, err)
		return
	}

	switch restPayload.Action {
	case "auth":
		app.authenticate(w, restPayload.Auth)
	default:
		app.ErroJSON(w, errors.New("unknown action"))
	}

}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	authServiceUrl := os.Getenv("AUTH_SERVICE_URL")

	jsonData, _ := json.MarshalIndent(a, "", "\t")

	response, err := http.Post(fmt.Sprintf("%s/authenticate", authServiceUrl), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		app.ErroJSON(w, err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		app.ErroJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.ErroJSON(w, errors.New("error calling auth service"))
		return
	}

	var authResponse apicommon.JsonResponse

	err = json.NewDecoder(response.Body).Decode(&authResponse)
	if err != nil {
		app.ErroJSON(w, err)
		return
	}

	if authResponse.Error {
		app.ErroJSON(w, errors.New(authResponse.Message))
		return
	}

	var payload apicommon.JsonResponse
	payload.Error = false
	payload.Message = "Logged in"
	payload.Data = authResponse.Data

	app.WriteJSON(w, http.StatusAccepted, payload)
}
