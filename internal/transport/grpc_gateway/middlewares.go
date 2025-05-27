package grpc_gateway

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

func (s *GatewayServer) applyMiddlewares(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		mw := middlewares[i]

		h = mw(h)
	}
	return h
}

func (s *GatewayServer) loggingMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s.Logger.Info(r.Context(), "Request started",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("remote", r.RemoteAddr))

			start := time.Now()
			defer func() {
				s.Logger.Info(r.Context(), "Request completed",
					zap.Duration("duration", time.Since(start)))
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func (s *GatewayServer) corsMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding")

			if r.Method == "OPTIONS" {
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
