package main

import (
	"log"
	"net/http"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	log.Printf("ERROR: %v", err)
	app.writeJSON(w, http.StatusInternalServerError, envelope{
		"error": "the server encountered a problem and could not process your request",
	}, nil)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.writeJSON(w, http.StatusNotFound, envelope{
		"error": "the requested resource could not be found",
	}, nil)
}

func (app *application) badRequest(w http.ResponseWriter, msg string) {
	app.writeJSON(w, http.StatusBadRequest, envelope{"error": msg}, nil)
}
