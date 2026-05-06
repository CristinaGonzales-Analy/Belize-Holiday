package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", app.health)

	// 1. Holidays in the current month
	mux.HandleFunc("GET /v1/holidays/current-month", app.currentMonthHandler)

	// 2. All occasions for 2026
	mux.HandleFunc("GET /v1/holidays/occasions", app.occasionsHandler)

	// 3. All dates for 2026
	mux.HandleFunc("GET /v1/holidays/dates", app.datesHandler)

	// 4. All days of week
	mux.HandleFunc("GET /v1/holidays/days", app.daysHandler)

	// 5. Is today a holiday
	mux.HandleFunc("GET /v1/holidays/today", app.todayHandler)

	// 6. Next holiday after today
	mux.HandleFunc("GET /v1/holidays/next", app.nextHolidayHandler)

	// 7. Holidays this month
	mux.HandleFunc("GET /v1/holidays/this-month", app.thisMonthHandler)

	// 8. Holidays next month
	mux.HandleFunc("GET /v1/holidays/next-month", app.nextMonthHandler)

	// 9. All holidays for a given year
	mux.HandleFunc("GET /v1/holidays/year", app.byYearHandler)

	return mux
}

func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
