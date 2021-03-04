package public

import (
	"github.com/gorilla/mux"
	"net/http"
)

type HandlersProvider interface {
	RegisterHandlers(r *mux.Router)
}

type RootHandlers interface {
	HandlersProvider
}

type rootHandlers struct {
}

func NewRootHandlers() *rootHandlers {
	return &rootHandlers{}
}

func (a *rootHandlers) RegisterHandlers(r *mux.Router) {

	cors := CrossOriginServer{}
	r.Use(cors.Middleware)

	uiFilesApi := NewLocalFilesApi("/ui", "ui/build/")
	uiFilesApi.RegisterHandlers(r)

	//uiFilesApi := NewEmbedFilesApi("/ui", "ui/build/")
	//uiFilesApi.RegisterHandlers(r)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/ui", http.StatusMovedPermanently)
	})

}
