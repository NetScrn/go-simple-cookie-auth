package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"time"

	"github.com/netscrn/gocookieauth/controllers"
	"github.com/netscrn/gocookieauth/data"
)

func main() {
	db, err := setUpdDB("root", "ss", "simple_auth", "")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	m := setUpMux(db)

	s := http.Server{
		Addr: ":8080",
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler: m,
	}

	err = s.ListenAndServeTLS("cert.pem", "key.pem")
	if err != nil {
		panic(err)
	}
}

func setUpdDB(username, password, dbName, params string) (*sql.DB, error) {
	dbAddress := fmt.Sprintf("%s:%s@/%s?%s", username, password, dbName, params)
	db, err := sql.Open("mysql", dbAddress)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(3 * time.Second)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db, nil
}

func setUpMux(db *sql.DB) *http.ServeMux {
	ur := data.NewUserRepo(db)
	uc := controllers.NewUsersController(ur)

	m := http.NewServeMux()
	m.HandleFunc("/user", uc.CreateUser)
	return m
}