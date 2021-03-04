package public

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/markbates/pkger"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type StaticFilesApi interface {
	RegisterHandlers(r *mux.Router)
	StaticFileHandler() http.Handler
}

type localFilesApi struct {
	StaticFilesApi
	pathPrefix string
	baseDir    string
}

func NewLocalFilesApi(pathPrefix, baseDir string) *localFilesApi {
	return &localFilesApi{
		baseDir:    strings.TrimRight(baseDir, "/"),
		pathPrefix: pathPrefix,
	}
}

func (a *localFilesApi) RegisterHandlers(r *mux.Router) {
	handler := spaHandler{
		staticPath: a.baseDir,
		indexPath:  "/index.html",
	}
	dirHandler := http.StripPrefix(a.pathPrefix, handler)
	r.PathPrefix(a.pathPrefix).Handler(dirHandler).Methods("GET")
}

type spaHandler struct {
	staticPath string
	indexPath  string
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	path = filepath.Join(h.staticPath, path)
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

type embedFilesApi struct {
	StaticFilesApi
	pathPrefix string
	baseDir    string
}

func NewEmbedFilesApi(pathPrefix, baseDir string) *embedFilesApi {
	pkger.Include("/ui/build")
	return &embedFilesApi{
		baseDir:    strings.TrimRight(baseDir, "/"),
		pathPrefix: pathPrefix,
	}
}

func (a *embedFilesApi) RegisterHandlers(r *mux.Router) {
	r.PathPrefix(a.pathPrefix).Handler(a.StaticFileHandler())
}

func (a *embedFilesApi) StaticFileHandler() http.Handler {
	return http.StripPrefix(a.pathPrefix, a)
}

func (a *embedFilesApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path
	indexPath := "/index.html"
	if !strings.HasPrefix(filePath, "/") {
		filePath = "/" + filePath
	}
	if filePath == "/" {
		filePath = indexPath
	}
	filePath = "/" + a.baseDir + filePath
	err := a.serveFile(w, filePath)
	if err == nil {
		return
	}
	if err == os.ErrNotExist && filePath != indexPath {
		if err := a.serveFile(w, indexPath); err != nil {
			http.NotFound(w, r)
			return
		}
	}
	http.NotFound(w, r)
	return

}

func (a *embedFilesApi) serveFile(w http.ResponseWriter, filePath string) error {

	file, err := pkger.Open(filePath)
	if err == os.ErrNotExist {
		return err
	} else if err != nil {
		http.Error(w, "Something bad happened", http.StatusInternalServerError)
		println(err)
		return nil
	}
	ctype := mime.TypeByExtension(filepath.Ext(filePath))
	if ctype != "" {
		w.Header().Set("Content-Type", ctype)
	}
	stat, err := file.Stat()
	if stat != nil {
		if !stat.ModTime().IsZero() {
			w.Header().Set("Last-Modified", stat.ModTime().UTC().Format(http.TimeFormat))
		}
		if stat.Size() > 0 {
			w.Header().Set("Content-Length", fmt.Sprintf("%d", stat.Size()))
		}
	}
	//
	//if strings.HasSuffix(filePath, ".json") {
	//} else if strings.HasSuffix(filePath, ".js") {
	//	w.Header().Set("Content-Type", "application/javascript")
	//} else if strings.HasSuffix(filePath, ".css") {
	//	w.Header().Set("Content-Type", "text/css")
	//} else if strings.HasSuffix(filePath, ".svg") {
	//	w.Header().Set("Content-Type", "image/svg+xml")
	//}

	_, _ = io.Copy(w, file)
	return nil
}
