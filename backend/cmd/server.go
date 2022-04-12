package main

import (
	"github.com/netscrn/gocookieauth/internal/data/database"
	"github.com/netscrn/gocookieauth/internal/web"
	"net/http"
	"os"
	"time"
)

func main() {
	env, exists := os.LookupEnv("SIMPLE_AUTH_ENV")
	if !exists {
		env = "development"
	}

	db, err := database.SetUpdDB(env)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	h := web.SetUpMainHandler(db, env)
	s := http.Server{
		Addr:         ":3001",
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      h,
	}
	defer s.Close()

	err = s.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
