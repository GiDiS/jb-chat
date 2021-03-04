package public

import "net/http"

type CrossOriginServer struct{}

func (s *CrossOriginServer) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if origin := req.Header.Get("Origin"); origin != "" {
			rw.Header().Set("Access-Control-Allow-Origin", origin)
			rw.Header().Set("Access-Control-Allow-Credentials", "true")
			rw.Header().Set("Access-Control-Allow-Methods", "POST, PUT, PATCH, GET, DELETE, OPTIONS")
			rw.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding")
		}
		if req.Method == "OPTIONS" {
			//rw.WriteHeader(204)
			return
		}
		rw.Header().Del("Origin")
		next.ServeHTTP(rw, req)
	})
}
