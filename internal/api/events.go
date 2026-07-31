package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/chankei613/ai-decision-reviewer/internal/events"
)

// httpStreamEvents はSSEストリーム。新規decision追加・解決の両方を配信する。
// UIのインボックスがこれ1本を購読してリアルタイム更新する。
func (s *Server) httpStreamEvents(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")
	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	stream, unsubscribe := s.Events.Subscribe()
	defer unsubscribe()

	keepalive := time.NewTicker(20 * time.Second)
	defer keepalive.Stop()

	for {
		select {
		case <-r.Context().Done():
			return

		case ev, open := <-stream:
			if !open {
				return
			}
			if err := writeSSE(w, ev); err != nil {
				return
			}
			flusher.Flush()

		case <-keepalive.C:
			if _, err := fmt.Fprint(w, ": keepalive\n\n"); err != nil {
				return
			}
			flusher.Flush()
		}
	}
}

func writeSSE(w http.ResponseWriter, ev events.Event) error {
	payload, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "event: %s\ndata: %s\n\n", ev.Type, payload)
	return err
}
