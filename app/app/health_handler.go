package app

import (
	"net/http"
)

// For Kubernetes readiness and liveness probes.
// Immediately responding with an HTTP 200 status.

// swagger:route GET /healthz get liveness
// If the server is ready, return HTTP 200 code
// responses:
//  200: HealthRespOk
func HandleHealth(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("."))
}
