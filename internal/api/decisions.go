package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/chankei613/ai-decision-reviewer/internal/db"
	"github.com/chankei613/ai-decision-reviewer/internal/events"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

var (
	errDecisionNotFound = &apiError{"decision not found"}
	errAlreadyResolved  = &apiError{"decision already resolved"}
	errInvalidLevel     = &apiError{"level must be one of: interrupt, urgent, emergency_stop"}
)

func validLevel(l db.EscalationLevel) bool {
	switch l {
	case db.LevelInterrupt, db.LevelUrgent, db.LevelEmergencyStop:
		return true
	default:
		return false
	}
}

type CreateDecisionInput struct {
	Source  string             `json:"source"`
	AgentID string             `json:"agent_id"`
	Subject string             `json:"subject"`
	Level   db.EscalationLevel `json:"level"`
	Reason  string             `json:"reason"`
	Summary string             `json:"summary"`
	Context map[string]any     `json:"context"`
}

// CreateDecision はAIが送るエスカレーションを1件追加する（HTTP・ネイティブバインディング共用）。
func (s *Server) CreateDecision(in CreateDecisionInput) (db.DecisionItem, error) {
	if in.AgentID == "" || in.Summary == "" {
		return db.DecisionItem{}, &apiError{"agent_id and summary are required"}
	}
	if !validLevel(in.Level) {
		return db.DecisionItem{}, errInvalidLevel
	}

	item := db.DecisionItem{
		ID:         uuid.NewString(),
		ReceivedAt: time.Now(),
		Source:     in.Source,
		AgentID:    in.AgentID,
		Subject:    in.Subject,
		Level:      in.Level,
		Reason:     in.Reason,
		Summary:    in.Summary,
		Context:    in.Context,
		Status:     db.StatusPending,
	}
	if item.Context == nil {
		item.Context = map[string]any{}
	}
	if err := s.DB.Create(&item).Error; err != nil {
		return db.DecisionItem{}, err
	}

	s.Events.Publish(events.Event{Type: events.EventCreated, DecisionID: item.ID, At: time.Now()})
	return item, nil
}

type ListDecisionsFilters struct {
	Status  string
	Level   string
	Source  string
	AgentID string
}

func FiltersFromQuery(q url.Values) ListDecisionsFilters {
	return ListDecisionsFilters{
		Status:  q.Get("status"),
		Level:   q.Get("level"),
		Source:  q.Get("source"),
		AgentID: q.Get("agent_id"),
	}
}

type ListDecisionsResult struct {
	Items []db.DecisionItem `json:"items"`
	Total int64             `json:"total"`
}

func (s *Server) ListDecisions(f ListDecisionsFilters, limit, offset int) (ListDecisionsResult, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}

	q := s.DB.Model(&db.DecisionItem{})
	if f.Status != "" {
		q = q.Where("status = ?", f.Status)
	}
	if f.Level != "" {
		q = q.Where("level = ?", f.Level)
	}
	if f.Source != "" {
		q = q.Where("source = ?", f.Source)
	}
	if f.AgentID != "" {
		q = q.Where("agent_id = ?", f.AgentID)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return ListDecisionsResult{}, err
	}

	var items []db.DecisionItem
	if err := q.Order("received_at desc").Limit(limit).Offset(offset).Find(&items).Error; err != nil {
		return ListDecisionsResult{}, err
	}
	return ListDecisionsResult{Items: items, Total: total}, nil
}

func (s *Server) GetDecision(id string) (db.DecisionItem, error) {
	var item db.DecisionItem
	if err := s.DB.First(&item, "id = ?", id).Error; err != nil {
		return db.DecisionItem{}, errDecisionNotFound
	}
	return item, nil
}

// resolve は承認/差し戻しの共通処理。既に解決済みなら errAlreadyResolved を返す
// （監査ログとしての一貫性を守るため、二重解決は許さない）。
func (s *Server) resolve(id, decision, feedback string) (db.DecisionItem, error) {
	item, err := s.GetDecision(id)
	if err != nil {
		return db.DecisionItem{}, err
	}
	if item.Status != db.StatusPending {
		return db.DecisionItem{}, errAlreadyResolved
	}

	now := time.Now()
	item.ResolutionDecision = decision
	item.ResolutionFeedback = feedback
	item.ResolutionResolvedAt = &now
	item.ResolutionResolvedBy = "human"
	if decision == "approve" {
		item.Status = db.StatusApproved
	} else {
		item.Status = db.StatusRejected
	}

	if err := s.DB.Save(&item).Error; err != nil {
		return db.DecisionItem{}, err
	}

	s.Events.Publish(events.Event{Type: events.EventResolved, DecisionID: item.ID, At: now})
	return item, nil
}

func (s *Server) ApproveDecision(id, feedback string) (db.DecisionItem, error) {
	return s.resolve(id, "approve", feedback)
}

func (s *Server) RejectDecision(id, feedback string) (db.DecisionItem, error) {
	return s.resolve(id, "reject", feedback)
}

// ─── HTTPハンドラー ────────────────────────────────────────────────────

func (s *Server) httpCreateDecision(w http.ResponseWriter, r *http.Request) {
	var body CreateDecisionInput
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	item, err := s.CreateDecision(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (s *Server) httpListDecisions(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))

	result, err := s.ListDecisions(FiltersFromQuery(q), limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (s *Server) httpGetDecision(w http.ResponseWriter, r *http.Request) {
	item, err := s.GetDecision(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

type resolveRequest struct {
	Feedback string `json:"feedback"`
}

func (s *Server) httpApproveDecision(w http.ResponseWriter, r *http.Request) {
	var body resolveRequest
	_ = json.NewDecoder(r.Body).Decode(&body)

	item, err := s.ApproveDecision(chi.URLParam(r, "id"), body.Feedback)
	if err != nil {
		s.writeResolveError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (s *Server) httpRejectDecision(w http.ResponseWriter, r *http.Request) {
	var body resolveRequest
	_ = json.NewDecoder(r.Body).Decode(&body)

	item, err := s.RejectDecision(chi.URLParam(r, "id"), body.Feedback)
	if err != nil {
		s.writeResolveError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (s *Server) writeResolveError(w http.ResponseWriter, err error) {
	switch err {
	case errDecisionNotFound:
		http.Error(w, err.Error(), http.StatusNotFound)
	case errAlreadyResolved:
		http.Error(w, err.Error(), http.StatusConflict)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
