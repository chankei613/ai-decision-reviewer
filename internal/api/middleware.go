package api

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/chankei613/ai-decision-reviewer/internal/db"
	"gorm.io/gorm"
)

type ctxKey string

const agentKeyCtxKey ctxKey = "agentKey"

// HashAPIKey は生のAPIキーをSHA-256でハッシュ化する。発行時と検証時で同じ関数を使う。
func HashAPIKey(key string) string {
	sum := sha256.Sum256([]byte(key))
	return hex.EncodeToString(sum[:])
}

func bearerToken(r *http.Request) string {
	authz := r.Header.Get("Authorization")
	if key := strings.TrimPrefix(authz, "Bearer "); key != "" && key != authz {
		return key
	}
	return ""
}

// APIKeyAuth は "Authorization: Bearer <key>" を検証する。
// AgentKeyが0件のときのみ bootstrapPaths に列挙したパスを未認証で通す
// （最初の1件を発行するにはキーが要る、という鶏と卵を解くため）。
func APIKeyAuth(conn *gorm.DB, bootstrapPaths ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := bearerToken(r)

			if key == "" {
				var count int64
				conn.Model(&db.AgentKey{}).Where("revoked_at IS NULL").Count(&count)
				if count == 0 {
					for _, p := range bootstrapPaths {
						if r.URL.Path == p {
							next.ServeHTTP(w, r)
							return
						}
					}
				}
				http.Error(w, "missing bearer token", http.StatusUnauthorized)
				return
			}

			var ak db.AgentKey
			if err := conn.Where("api_key_hash = ? AND revoked_at IS NULL", HashAPIKey(key)).First(&ak).Error; err != nil {
				http.Error(w, "invalid api key", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), agentKeyCtxKey, &ak)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
