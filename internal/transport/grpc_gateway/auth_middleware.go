package grpc_gateway

import (
	"github.com/JunBSer/gateway/pkg/logger"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	authServiceAddr string
	logger          logger.Logger
}

func (s *GatewayServer) authMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s.logger.Debug(r.Context(), "Auth middleware start")

			defer func() {
				s.logger.Debug(r.Context(), "Auth middleware exit")
			}()

			if _, ok := s.endpoints.PublicPaths[r.URL.Path]; ok {
				s.logger.Debug(r.Context(), "Public path access")
				next.ServeHTTP(w, r)
				return
			}

			token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
			if token == "" {
				s.respondError(w, r, "Authorization required", http.StatusUnauthorized)
				next.ServeHTTP(w, r)
				return
			}

			newCtx, isValid := s.validateToken(r.Context(), token)
			if !isValid {
				s.respondError(w, r, "Invalid token", http.StatusUnauthorized)
				next.ServeHTTP(w, r)
				return
			}

			r = r.WithContext(newCtx)
			next.ServeHTTP(w, r)
		})
	}
}
