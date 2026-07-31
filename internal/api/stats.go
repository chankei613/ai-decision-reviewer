package api

import (
	"net/http"
	"time"

	"github.com/chankei613/ai-decision-reviewer/internal/db"
)

type StatsResult struct {
	Pending              int64            `json:"pending"`
	PendingByLevel       map[string]int64 `json:"pending_by_level"`
	AvgResolutionSeconds float64          `json:"avg_resolution_seconds"`
}

func (s *Server) Stats() (StatsResult, error) {
	result := StatsResult{PendingByLevel: map[string]int64{}}

	if err := s.DB.Model(&db.DecisionItem{}).Where("status = ?", db.StatusPending).Count(&result.Pending).Error; err != nil {
		return StatsResult{}, err
	}

	type levelCount struct {
		Level db.EscalationLevel
		Count int64
	}
	var levelCounts []levelCount
	if err := s.DB.Model(&db.DecisionItem{}).
		Where("status = ?", db.StatusPending).
		Select("level, count(*) as count").
		Group("level").
		Scan(&levelCounts).Error; err != nil {
		return StatsResult{}, err
	}
	for _, lc := range levelCounts {
		result.PendingByLevel[string(lc.Level)] = lc.Count
	}

	var resolved []db.DecisionItem
	if err := s.DB.Where("status != ? AND resolution_resolved_at IS NOT NULL", db.StatusPending).
		Find(&resolved).Error; err != nil {
		return StatsResult{}, err
	}
	if len(resolved) > 0 {
		var total time.Duration
		for _, item := range resolved {
			total += item.ResolutionResolvedAt.Sub(item.ReceivedAt)
		}
		result.AvgResolutionSeconds = total.Seconds() / float64(len(resolved))
	}

	return result, nil
}

func (s *Server) httpStats(w http.ResponseWriter, r *http.Request) {
	result, err := s.Stats()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, result)
}
