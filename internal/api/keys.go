package api

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/chankei613/ai-decision-reviewer/internal/db"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type IssueKeyResult struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	APIKey string `json:"api_key"` // 発行時にしか見せない生のキー
}

func generateAPIKey() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

// IssueKey は新しいAPIキーを発行する（HTTP・ネイティブバインディング共用）。
func (s *Server) IssueKey(name string) (IssueKeyResult, error) {
	if name == "" {
		return IssueKeyResult{}, errNameRequired
	}

	raw, err := generateAPIKey()
	if err != nil {
		return IssueKeyResult{}, err
	}

	ak := db.AgentKey{
		ID:         uuid.NewString(),
		Name:       name,
		APIKeyHash: HashAPIKey(raw),
		CreatedAt:  time.Now(),
	}
	if err := s.DB.Create(&ak).Error; err != nil {
		return IssueKeyResult{}, err
	}

	return IssueKeyResult{ID: ak.ID, Name: ak.Name, APIKey: raw}, nil
}

func (s *Server) ListKeys() ([]db.AgentKey, error) {
	var keys []db.AgentKey
	err := s.DB.Order("created_at asc").Find(&keys).Error
	return keys, err
}

func (s *Server) RevokeKey(id string) error {
	now := time.Now()
	res := s.DB.Model(&db.AgentKey{}).Where("id = ? AND revoked_at IS NULL", id).Update("revoked_at", now)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errKeyNotFound
	}
	return nil
}

var (
	errNameRequired = &apiError{"name is required"}
	errKeyNotFound  = &apiError{"key not found or already revoked"}
)

type apiError struct{ msg string }

func (e *apiError) Error() string { return e.msg }

// ─── HTTPハンドラー ────────────────────────────────────────────────────

type issueKeyRequest struct {
	Name string `json:"name"`
}

func (s *Server) httpIssueKey(w http.ResponseWriter, r *http.Request) {
	var body issueKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	result, err := s.IssueKey(body.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusCreated, result)
}

func (s *Server) httpListKeys(w http.ResponseWriter, r *http.Request) {
	keys, err := s.ListKeys()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, keys)
}

func (s *Server) httpRevokeKey(w http.ResponseWriter, r *http.Request) {
	if err := s.RevokeKey(chi.URLParam(r, "id")); err != nil {
		if err == errKeyNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	// レスポンスヘッダー送信後のエンコードエラーはクライアント切断等でリカバリ不能なため意図的に無視する
	_ = json.NewEncoder(w).Encode(v)
}
