package main

import (
	"net/http"

	"github.com/vbrenister/apicommon"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := apicommon.JsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.WriteJSON(w, http.StatusAccepted, payload)
}
