package grpc_gateway

import (
	"github.com/JunBSer/gateway/internal/metadata"
	"net/http"
	"strings"
)

func (s *GatewayServer) registerSwagger(mux *http.ServeMux) {
	mux.HandleFunc("/docs/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs"+strings.TrimPrefix(r.URL.Path, "/docs"))
	})

	s.Endpoints.AddEndpoint("GET", "/docs", "Docs.Index", metadata.AuthNone)
	s.Endpoints.AddEndpoint("GET", "/docs/swagger.json", "Docs.Swagger", metadata.AuthNone)
}

//http.Handle("/docs/auth/swagger.json",
//serveSwagger(filepath.Join("docs", "auth", "swagger.json")))
//
//http.Handle("/docs/booking/swagger.json",
//serveSwagger(filepath.Join("docs", "booking", "swagger.json")))
