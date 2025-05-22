package grpc_gateway

import (
	"net/http"
	"strings"
)

func (s *GatewayServer) registerSwagger(mux *http.ServeMux) {
	mux.HandleFunc("/docs/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs"+strings.TrimPrefix(r.URL.Path, "/docs"))
	})

	s.endpoints.PublicPaths["/docs"] = struct{}{}
	s.endpoints.PublicPaths["/docs/swagger.json"] = struct{}{}
}

//http.Handle("/docs/auth/swagger.json",
//serveSwagger(filepath.Join("docs", "auth", "swagger.json")))
//
//http.Handle("/docs/booking/swagger.json",
//serveSwagger(filepath.Join("docs", "booking", "swagger.json")))
