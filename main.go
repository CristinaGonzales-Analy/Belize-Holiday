package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

type application struct {
	db *sql.DB
}

func main() {
	dsn := "postgres://holidays_user:password@localhost:5432/belize_holidays?sslmode=disable"

	db, err := openDB(dsn)
	if err != nil {
		log.Fatalf("Unable to open Database: %v", err)
	}
	defer db.Close()

	app := &application{db: db}

	srv := &http.Server{
		Addr:    ":4000",
		Handler: app.enableCORS(app.routes()),
	}

	log.Printf("Belize Holidays API starting on %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}
	return db, nil
}

func (app *application) health(w http.ResponseWriter, r *http.Request) {
	app.writeJSON(w, http.StatusOK, envelope{"status": "available"}, nil)
}
