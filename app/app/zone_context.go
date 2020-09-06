package app

import (
	"context"
	"github.com/go-chi/chi"
	"net/http"
)

func (app *App)ZoneCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		zoneName := chi.URLParam(r, "zoneName")
		var found = false
		for _, item := range app.zones {
			if item == zoneName {
				found = true
			}
		}
		if !found {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		ctx := context.WithValue(r.Context(), "zoneName", zoneName)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
