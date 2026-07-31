// Package db はAI Decision ReviewerのGORMモデルとSQLite初期化を提供する。
// concept.md 6章のエスカレーションラダーのうち、人間のアクションを要する3段階
// （interrupt/urgent/emergency_stop）を対象にする。docs/spec.md 参照。
package db

import "time"

type EscalationLevel string

const (
	LevelInterrupt     EscalationLevel = "interrupt"
	LevelUrgent        EscalationLevel = "urgent"
	LevelEmergencyStop EscalationLevel = "emergency_stop"
)

type DecisionStatus string

const (
	StatusPending  DecisionStatus = "pending"
	StatusApproved DecisionStatus = "approved"
	StatusRejected DecisionStatus = "rejected"
)

// DecisionItem は1件のエスカレーション。未解決の間はResolution系フィールドが空、
// 解決は一度だけ許される（二重解決はAPI側で409にする）。
type DecisionItem struct {
	ID         string    `gorm:"primaryKey" json:"id"`
	ReceivedAt time.Time `gorm:"index" json:"received_at"`

	Source  string `gorm:"index" json:"source"`
	AgentID string `gorm:"index" json:"agent_id"`
	Subject string `gorm:"index" json:"subject"`

	Level   EscalationLevel `gorm:"index" json:"level"`
	Reason  string          `json:"reason"`
	Summary string          `json:"summary"`
	Context map[string]any  `gorm:"serializer:json" json:"context"`

	Status DecisionStatus `gorm:"index" json:"status"`

	ResolutionDecision   string     `json:"resolution_decision,omitempty"`
	ResolutionFeedback   string     `json:"resolution_feedback,omitempty"`
	ResolutionResolvedAt *time.Time `json:"resolution_resolved_at,omitempty"`
	ResolutionResolvedBy string     `json:"resolution_resolved_by,omitempty"`
}

// AgentKey — Ingestion/解決APIを叩くためのAPIキー。ハッシュのみ保存する。
type AgentKey struct {
	ID         string     `gorm:"primaryKey" json:"id"`
	Name       string     `json:"name"`
	APIKeyHash string     `json:"-"`
	CreatedAt  time.Time  `json:"created_at"`
	RevokedAt  *time.Time `json:"revoked_at"`
}
