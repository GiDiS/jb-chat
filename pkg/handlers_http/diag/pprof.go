package diag

import (
	"github.com/gorilla/mux"
	"net/http/pprof"
)

func registerPprof(r *mux.Router) {
	// Регистрация pprof-обработчиков

	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)
	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/*", pprof.Index)
}
