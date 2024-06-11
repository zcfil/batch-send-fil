package auth

import (
	"context"
	"net/http"
	"strings"

	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("auth")

type Handler struct {
	Verify func(ctx context.Context, token string) ([]Permission, error)
	Next   http.HandlerFunc
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	token := r.Header.Get("Authorization")
	if token == "" {
		token = r.FormValue("token")
		if token != "" {
			token = "Bearer " + token
		}
	}

	if token != "" {
		if !strings.HasPrefix(token, "Bearer ") {
			log.Warn("missing Bearer prefix in auth header")
			w.WriteHeader(401)
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")

		allow, err := h.Verify(ctx, token)
		if err != nil {
			log.Warnf("JWT Verification failed (originating from %s): %s", r.RemoteAddr, err)
			w.WriteHeader(401)
			return
		}

		ctx = WithPerm(ctx, allow)
	}

	h.Next(w, r.WithContext(ctx))
}
