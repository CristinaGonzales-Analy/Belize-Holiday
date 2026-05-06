package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Holiday struct {
	ID       int64  `json:"id"`
	Day      string `json:"day"`
	Date     string `json:"date"`
	Month    int    `json:"month"`
	Year     int    `json:"year"`
	Occasion string `json:"occasion"`
}

// current month
func (app *application) currentMonthHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	app.holidaysByMonth(w, r, int(now.Month()), now.Year())
}

// Occassions
func (app *application) occasionsHandler(w http.ResponseWriter, r *http.Request) {
	query := `SELECT occasion FROM holidays WHERE year = 2026 ORDER BY id`

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	rows, err := app.db.QueryContext(ctx, query)
	if err != nil {
		app.serverError(w, err)
		return
	}
	defer rows.Close()

	var occasions []string
	for rows.Next() {
		var o string
		if err := rows.Scan(&o); err != nil {
			app.serverError(w, err)
			return
		}
		occasions = append(occasions, o)
	}
	if err := rows.Err(); err != nil {
		app.serverError(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{
		"year":      2026,
		"occasions": occasions,
	}, nil)
}

// All dates
func (app *application) datesHandler(w http.ResponseWriter, r *http.Request) {
	query := `SELECT date FROM holidays WHERE year = 2026 ORDER BY id`

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	rows, err := app.db.QueryContext(ctx, query)
	if err != nil {
		app.serverError(w, err)
		return
	}
	defer rows.Close()

	var dates []string
	for rows.Next() {
		var d string
		if err := rows.Scan(&d); err != nil {
			app.serverError(w, err)
			return
		}
		dates = append(dates, d)
	}
	if err := rows.Err(); err != nil {
		app.serverError(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{
		"year":  2026,
		"dates": dates,
	}, nil)
}

// All days
func (app *application) daysHandler(w http.ResponseWriter, r *http.Request) {
	query := `SELECT day FROM holidays WHERE year = 2026 ORDER BY id`

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	rows, err := app.db.QueryContext(ctx, query)
	if err != nil {
		app.serverError(w, err)
		return
	}
	defer rows.Close()

	var days []string
	for rows.Next() {
		var d string
		if err := rows.Scan(&d); err != nil {
			app.serverError(w, err)
			return
		}
		days = append(days, d)
	}
	if err := rows.Err(); err != nil {
		app.serverError(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{
		"year": 2026,
		"days": days,
	}, nil)
}

// Check if today is holiday
func (app *application) todayHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	query := `SELECT id, day, date, month, year, occasion
	          FROM holidays
	          WHERE month = $1 AND year = $2`

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	holidays, err := app.queryHolidays(ctx, query, int(now.Month()), now.Year())
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Match by day-of-month
	todayNum := strconv.Itoa(now.Day())
	for _, h := range holidays {
		if startsWithNum(h.Date, todayNum) {
			app.writeJSON(w, http.StatusOK, envelope{
				"is_holiday": true,
				"occasion":   h.Occasion,
				"message":    fmt.Sprintf("Today is %s.", h.Occasion),
			}, nil)
			return
		}
	}

	app.writeJSON(w, http.StatusOK, envelope{
		"is_holiday": false,
		"message":    "Unlucky you! Today is not a holiday",
	}, nil)
}

// checks if starsWithNUm
func startsWithNum(date, num string) bool {
	if len(date) < len(num) {
		return false
	}
	return date[:len(num)] == num
}

// Next holiday after today
func (app *application) nextHolidayHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	query := `SELECT id, day, date, month, year, occasion
	          FROM holidays
	          WHERE (year > $1)
	             OR (year = $1 AND month > $2)
	             OR (year = $1 AND month = $2 AND CAST(REGEXP_REPLACE(date, '[^0-9].*', '') AS INTEGER) > $3)
	          ORDER BY year, month, CAST(REGEXP_REPLACE(date, '[^0-9].*', '') AS INTEGER)
	          LIMIT 1`

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var h Holiday
	err := app.db.QueryRowContext(ctx, query, now.Year(), int(now.Month()), now.Day()).
		Scan(&h.ID, &h.Day, &h.Date, &h.Month, &h.Year, &h.Occasion)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.writeJSON(w, http.StatusOK, envelope{
				"found":   false,
				"message": "No more holidays found after today.",
			}, nil)
			return
		}
		app.serverError(w, err)
		return
	}

	// Calculate days away
	holidayTime := time.Date(h.Year, time.Month(h.Month), dayNum(h.Date), 0, 0, 0, 0, time.UTC)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	daysAway := int(holidayTime.Sub(today).Hours() / 24)

	app.writeJSON(w, http.StatusOK, envelope{
		"found":     true,
		"days_away": daysAway,
		"holiday":   h,
	}, nil)
}

// Extracts integer
func dayNum(date string) int {
	n := 0
	for _, c := range date {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		} else {
			break
		}
	}
	return n
}

// this month
func (app *application) thisMonthHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	app.holidaysByMonth(w, r, int(now.Month()), now.Year())
}

// Next month
func (app *application) nextMonthHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	month := int(now.Month()) + 1
	year := now.Year()
	if month > 12 {
		month = 1
		year++
	}
	app.holidaysByMonth(w, r, month, year)
}

// Month Holidays
func (app *application) holidaysByMonth(w http.ResponseWriter, r *http.Request, month, year int) {
	query := `SELECT id, day, date, month, year, occasion
	          FROM holidays
	          WHERE month = $1 AND year = $2
	          ORDER BY CAST(REGEXP_REPLACE(date, '[^0-9].*', '') AS INTEGER)`

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	holidays, err := app.queryHolidays(ctx, query, month, year)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{
		"month":    time.Month(month).String(),
		"year":     year,
		"count":    len(holidays),
		"holidays": holidays,
	}, nil)
}

// By year
func (app *application) byYearHandler(w http.ResponseWriter, r *http.Request) {
	yearStr := r.URL.Query().Get("year")
	if yearStr == "" {
		app.badRequest(w, "missing required query parameter: year")
		return
	}
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		app.badRequest(w, "year must be an integer")
		return
	}
	if year != 2026 {
		app.writeJSON(w, http.StatusNotFound, envelope{
			"error":   fmt.Sprintf("no holiday data available for year %d — only 2026 is supported", year),
			"support": []int{2026},
		}, nil)
		return
	}

	query := `SELECT id, day, date, month, year, occasion
	          FROM holidays
	          WHERE year = $1
	          ORDER BY month, CAST(REGEXP_REPLACE(date, '[^0-9].*', '') AS INTEGER)`

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	holidays, err := app.queryHolidays(ctx, query, year)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{
		"year":     year,
		"count":    len(holidays),
		"holidays": holidays,
	}, nil)
}
