package app

import (
	"encoding/json"
	"net/http"
)

func (app *App) HandleIndex(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	w.WriteHeader(http.StatusOK)

	records , err := app.rfc2136.GetRecordsForZone(app.zones[0])
	if err != nil {
		http.Error(w, err.Error(), 422)
		return
	}
	var b []byte
	b, err = json.Marshal(records)
	if err != nil {
		http.Error(w, err.Error(), 422)
		return
	}
	w.Write(b)
}
