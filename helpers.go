package main

import (
	"context"
	"encoding/json"
	"net/http"
)

type envelope map[string]any

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	js = append(js, '\n')
	for key, values := range headers {
		for _, v := range values {
			w.Header().Add(key, v)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(js)
	return err
}

func scanHoliday(row interface {
	Scan(...any) error
}) (Holiday, error) {
	var h Holiday
	return h, row.Scan(&h.ID, &h.Day, &h.Date, &h.Month, &h.Year, &h.Occasion)
}

func (app *application) queryHolidays(ctx context.Context, query string, args ...any) ([]Holiday, error) {
	rows, err := app.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var holidays []Holiday
	for rows.Next() {
		h, err := scanHoliday(rows)
		if err != nil {
			return nil, err
		}
		holidays = append(holidays, h)
	}
	return holidays, rows.Err()
}
