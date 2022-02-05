package authentication

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"log"
	"net/http"

	"github.com/netscrn/gocookieauth/data/sessions"
	"github.com/netscrn/gocookieauth/data/users"
)

func Authenticate(h http.Handler, um users.Manger, tm sessions.TokenManager) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		us, err := r.Cookie("u_session")
		if err != nil {
			if err != http.ErrNoCookie {
				log.Printf("Authenticate - error reading session cookie: %v\n", err)
			}
			h.ServeHTTP(w, r)
			return
		}
		reqCsrfEncoded := r.Header.Get("X-CSRF-Token")
		if reqCsrfEncoded == "" {
			log.Printf("Authenticate - warning! request with session %s, without csrf: %v\n", us.Value, err)
			h.ServeHTTP(w, r)
			return
		}

		computedCsrf := sha256.Sum256([]byte(us.Value))
		reqCsrf, err := base64.StdEncoding.DecodeString(reqCsrfEncoded)
		if err != nil {
			log.Printf("Authenticate - warning! error decoding csrf: %v\n", err)
			h.ServeHTTP(w, r)
			return
		}

		csrfCompRes := subtle.ConstantTimeCompare(computedCsrf[:], reqCsrf)
		if csrfCompRes == 0 {
			log.Printf("Authenticate - warning! request with session %s, with invalid csrf: %v\n", us.Value, err)
			h.ServeHTTP(w, r)
			return
		}

		t, err := tm.Read(r.Context(), us.Value)
		if err != nil {
			if err != sessions.ErrTokenNotFound {
				log.Printf("Authenticate - can't fetch token with id %s: %v\n", us.Value, err)
			}
			h.ServeHTTP(w, r)
			return
		}

		u, err := um.GetUserByID(r.Context(), t.UserID)
		if err != nil {
			log.Printf("Authenticate - can't fetch user with id %d, from token_id %s: %v\n", t.UserID, us.Value, err)
			h.ServeHTTP(w, r)
			return
		}

		r = r.WithContext(ContextWithToken(r.Context(), t))
		r = r.WithContext(ContextWithUser(r.Context(), u))
		h.ServeHTTP(w, r)
	})
}
