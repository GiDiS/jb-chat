package diag

import (
	"fmt"
	"github.com/GiDiS/jb-chat/pkg/handlers_http/public"
	"github.com/gorilla/mux"
	"net/http"
)

type RootHandlers interface {
	public.HandlersProvider
	HealthzHandler() http.HandlerFunc
	ReadyzHandler() http.HandlerFunc
}

type rootHandlers struct {
}

func NewRootHandlers() *rootHandlers {
	return &rootHandlers{}
}

func (a *rootHandlers) RegisterHandlers(r *mux.Router) {
	r.HandleFunc("/healthz", a.HealthzHandler()).Methods("GET")
	r.HandleFunc("/readyz", a.ReadyzHandler()).Methods("GET")
	//r.Handle("/metrics", promhttp.Handler())
}

func (a *rootHandlers) HealthzHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, http.StatusText(http.StatusOK))
	}
}

func (a *rootHandlers) ReadyzHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, http.StatusText(http.StatusOK))
	}
}
