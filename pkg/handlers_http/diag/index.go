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
	metrics bool
	pprof   bool
}

func NewRootHandlers(metrics, pprof bool) *rootHandlers {
	return &rootHandlers{
		metrics: metrics,
		pprof:   pprof,
	}
}

func (a *rootHandlers) RegisterHandlers(r *mux.Router) {
	r.HandleFunc("/healthz", a.HealthzHandler()).Methods("GET")
	r.HandleFunc("/readyz", a.ReadyzHandler()).Methods("GET")
	if a.metrics {
		registerMetrics(r)
	}
	if a.pprof {
		registerPprof(r)
	}
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
