package main

import (
	"logger/data"
	"net/http"

	"github.com/vbrenister/apicommon"
)

type JSONPaylod struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	var requestPayload JSONPaylod
	err := app.ReadJSON(w, r, &requestPayload)

	if err != nil {
		app.ErroJSON(w, err)
		return
	}

	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err = app.Models.LogEntry.Insert(event)
	if err != nil {
		app.ErroJSON(w, err)
		return
	}

	resp := apicommon.JsonResponse{
		Error:   false,
		Message: "Logged",
	}

	app.WriteJSON(w, http.StatusCreated, resp)
}
