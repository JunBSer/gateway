package grpc_gateway

import (
	"encoding/json"
	"github.com/JunBSer/gateway/internal/metadata"
	"github.com/JunBSer/gateway/pkg/logger"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	authServiceAddr string
	logger          logger.Logger
}

func IsUserMatchesTheID(r *http.Request, userID string) bool {
	body := map[string]interface{}{}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return false
	}

	id, ok := body["user_id"]
	if !ok {
		return true
	}

	strID, ok := id.(string)
	if !ok {
		return false
	}

	return strID == userID
}

func (s *GatewayServer) authMiddleware(cfg *metadata.EndpointConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cleanPath := strings.SplitN(r.URL.Path, "?", 2)[0]

			ep, found := cfg.MatchEndpoint(r.Method, cleanPath)
			if !found {
				s.respondError(w, r, "Unknown endpoint", http.StatusNotFound)
				return
			}

			switch ep.Level {
			case metadata.AuthNone:
				next.ServeHTTP(w, r)
				return

			case metadata.AuthUser, metadata.AuthAdmin:

				token := strings.TrimSpace(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer"))
				if token == "" {
					s.respondError(w, r, "Authorization required", http.StatusUnauthorized)
					return
				}

				ctx2, valid := s.validateToken(r.Context(), token)
				if !valid {
					s.respondError(w, r, "Invalid token", http.StatusUnauthorized)
					return
				}

				if ep.Level == metadata.AuthAdmin {
					if isAdmin, _ := ctx2.Value(IsAdminKey).(bool); !isAdmin {
						s.respondError(w, r, "Admin access required", http.StatusForbidden)
						return
					}
				}

				if ep.Level == metadata.AuthUser {
					if usID, _ := ctx2.Value(UsIDKey).(string); usID != "" {
						if !IsUserMatchesTheID(r, usID) {
							s.respondError(w, r, "You not permitted to use this userID", http.StatusForbidden)
							return
						}
					}
				}

				next.ServeHTTP(w, r.WithContext(ctx2))
				return
			default:

				s.respondError(w, r, "Unknown authorization level", http.StatusInternalServerError)
				return
			}
		})
	}
}
