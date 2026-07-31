// Package events は新規decision追加・解決イベントを、購読中のSSEクライアントへ
// リアルタイム配信するための小さなpub/subを提供する（永続化はinternal/dbが担う）。
package events

import (
	"sync"
	"time"
)

const (
	EventCreated  = "decision:created"
	EventResolved = "decision:resolved"
)

// Event はSSEで配信される1件分のイベント。
type Event struct {
	Type       string    `json:"type"`
	DecisionID string    `json:"decision_id"`
	At         time.Time `json:"at"`
}

// Broker は購読者へのファンアウトのみを担当する。
type Broker struct {
	mu   sync.RWMutex
	subs map[chan Event]struct{}
}

func NewBroker() *Broker {
	return &Broker{subs: make(map[chan Event]struct{})}
}

// Subscribe は新しい購読チャネルと、その解除関数を返す。
// 解除関数は購読者側（SSEハンドラ）が必ず defer で呼ぶこと。
func (b *Broker) Subscribe() (<-chan Event, func()) {
	ch := make(chan Event, 32)

	b.mu.Lock()
	b.subs[ch] = struct{}{}
	b.mu.Unlock()

	var once sync.Once
	unsubscribe := func() {
		once.Do(func() {
			b.mu.Lock()
			delete(b.subs, ch)
			b.mu.Unlock()
			close(ch)
		})
	}
	return ch, unsubscribe
}

// Publish は全購読者に配信する。遅い購読者のバッファが埋まっていたら捨てる（ブロックしない）。
func (b *Broker) Publish(e Event) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for ch := range b.subs {
		select {
		case ch <- e:
		default:
		}
	}
}
