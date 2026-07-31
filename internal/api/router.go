package api

import (
	"net/http"

	"github.com/chankei613/ai-decision-reviewer/internal/events"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

// Server は全ロジックの実体。HTTPハンドラーとWailsネイティブバインディングの
// 両方がこの同じ Server のメソッドを呼ぶことで、UIとAPIの挙動がズレないようにする。
type Server struct {
	DB     *gorm.DB
	Events *events.Broker
}

func New(conn *gorm.DB) *Server {
	return &Server{DB: conn, Events: events.NewBroker()}
}

func (s *Server) Router() http.Handler {
	r := chi.NewRouter()

	r.Route("/api/v1/keys", func(r chi.Router) {
		r.Use(APIKeyAuth(s.DB, "/api/v1/keys"))
		r.Post("/", s.httpIssueKey)
		r.Get("/", s.httpListKeys)
		r.Delete("/{id}", s.httpRevokeKey)
	})

	r.Route("/api/v1/decisions", func(r chi.Router) {
		r.Use(APIKeyAuth(s.DB))
		r.Post("/", s.httpCreateDecision)
		r.Get("/", s.httpListDecisions)
		r.Get("/{id}", s.httpGetDecision)
		r.Post("/{id}/approve", s.httpApproveDecision)
		r.Post("/{id}/reject", s.httpRejectDecision)
	})

	r.Route("/api/v1/stats", func(r chi.Router) {
		r.Use(APIKeyAuth(s.DB))
		r.Get("/", s.httpStats)
	})

	r.Route("/api/v1/events", func(r chi.Router) {
		r.Use(APIKeyAuth(s.DB))
		r.Get("/", s.httpStreamEvents)
	})

	return r
}

// NewRouter はcmd/adrserve（単体HTTPサーバー）向けの簡易コンストラクタ。
func NewRouter(conn *gorm.DB) http.Handler {
	return New(conn).Router()
}
