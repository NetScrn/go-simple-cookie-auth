package web

import (
	"database/sql"
	"net/http"

	"github.com/netscrn/gocookieauth/data/users"
	"github.com/netscrn/gocookieauth/web/controllers"
	"github.com/netscrn/gocookieauth/web/middleware"
)

func SetUpMainHandler(db *sql.DB) http.Handler {
	ur := users.NewUserRepo(db)
	uc := controllers.NewUsersController(ur)

	m := http.NewServeMux()
	m.HandleFunc("/user", uc.CreateUser)
	return middleware.CORS(middleware.CommonHeaders(m))
}
