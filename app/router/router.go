package router

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"io"
	"net/http"
	"strings"

	"github.com/dgoujard/rfc2136_web_api/app/app"
	"github.com/dgoujard/rfc2136_web_api/app/requestlog"
	"github.com/dgoujard/rfc2136_web_api/app/router/middleware"
	"github.com/markbates/pkger"
)

func New(a *app.App) *chi.Mux {
	l := a.Logger()

	r := chi.NewRouter()
	rcors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(rcors.Handler)
	// Routes for healthz
	r.Get("/healthz", app.HandleHealth)

	// Routes for APIs
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(a.AuthCheck)
		r.Use(middleware.ContentTypeJson)

		r.Method("GET", "/zones", requestlog.NewHandler(a.HandleZoneList, l))

		r.Route("/zones/{zoneName}", func(r chi.Router) {
			r.Use(a.ZoneCtx)
			r.Method("GET", "/", requestlog.NewHandler(a.HandleZoneRecordsList, l))
			r.Method("POST", "/", requestlog.NewHandler(a.HandleCreateZoneRecords, l))
			r.Method("GET", "/{recordId}", requestlog.NewHandler(a.HandleGetZoneRecordsId, l))
			r.Method("DELETE", "/{recordId}", requestlog.NewHandler(a.HandleDeleteZoneRecordsId, l))
		})
	})

	FileServer(r, "/ui", pkger.Dir("/public/ui"))
	FileServer(r, "/swagger", pkger.Dir("/public/swagger"))

	r.Mount("/", RootRouter())

	return r
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}


func RootRouter() chi.Router {
	r := chi.NewRouter()

	// catch any remaining routes and serve them the index.html
	// let react-router deal with them
	// TODO: handle 404 in react router http://knowbody.github.io/react-router-docs/guides/NotFound.html
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		f, err := pkger.Open("/public/ui/index.html")
		if err != nil {
			http.Error(w, err.Error(), 422)
		}
		defer f.Close()

		if _, err := io.Copy(w, f); err != nil {
			http.Error(w, err.Error(), 422)
		}
	})

	return r
}