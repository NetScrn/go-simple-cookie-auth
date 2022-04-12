package controllers

import (
	"fmt"
	auth "github.com/netscrn/gocookieauth/internal/web/middleware/authentication"
	"net/http"
)

func AuthenticatedOnly(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	u, ok := auth.UserFromContext(ctx)
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	resBody := []byte(fmt.Sprintf(`{"message": "Hello %s"}`, u.Email))

	w.WriteHeader(http.StatusOK)
	w.Write(resBody)
}
