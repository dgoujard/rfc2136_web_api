package app

import (
	"fmt"
	"net/http"
)

func (app *App)AuthCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			unauthorized(w, "Dns API Auth")
			return
		}
		if username == app.username && password == app.password {
			next.ServeHTTP(w, r)
			return
		}
		unauthorized(w, "Dns API Auth")
	})
}

func unauthorized(w http.ResponseWriter, realm string) {
	w.Header().Add("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, realm))
	w.WriteHeader(http.StatusUnauthorized)
}