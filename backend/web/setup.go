package web

import (
	"database/sql"
	"net/http"

	"github.com/netscrn/gocookieauth/data/sessions"
	"github.com/netscrn/gocookieauth/data/users"
	"github.com/netscrn/gocookieauth/web/controllers"
	"github.com/netscrn/gocookieauth/web/middleware"
	auth "github.com/netscrn/gocookieauth/web/middleware/authentication"
)

func SetUpMainHandler(db *sql.DB) http.Handler {
	ur := users.NewUserRepo(db)
	tr := sessions.NewTokensRepo(db)

	uc := controllers.NewUsersController(ur)
	sc := controllers.NewSessionController(ur, tr)

	m := http.NewServeMux()
	m.HandleFunc("/user", uc.CreateUser)
	m.HandleFunc("/login", sc.Login)
	m.HandleFunc("/logout", sc.Logout)
	m.HandleFunc("/auth-only", controllers.AuthenticatedOnly)

	h := middleware.CORS(m)
	h = middleware.CommonHeaders(h)
	h = auth.Authenticate(h, ur, tr)
	return h
}
