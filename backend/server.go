package main

import (
	"net/http"
	"os"
	"time"

	"github.com/netscrn/gocookieauth/data/database"
	"github.com/netscrn/gocookieauth/web"
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

	h := web.SetUpMainHandler(db)
	s := http.Server{
		Addr:         ":8080",
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      h,
	}

	err = s.ListenAndServeTLS("cert.pem", "key.pem")
	if err != nil {
		panic(err)
	}
}
