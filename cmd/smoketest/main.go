// cmd/smoketest はAI Decision ReviewerのAPIを一時DBで自前起動し、
// ブートストラップ鍵発行 → decision追加 → 一覧/取得 → 承認 → 二重解決の拒否 →
// SSEでイベントが届くこと、の一連が通しで動くことを確認する。
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"time"

	"github.com/chankei613/ai-decision-reviewer/internal/api"
	"github.com/chankei613/ai-decision-reviewer/internal/db"
)

func main() {
	dbPath := "smoketest.db"
	_ = os.Remove(dbPath)
	defer func() { _ = os.Remove(dbPath) }()

	conn, err := db.Init(dbPath)
	if err != nil {
		log.Fatalf("db init: %v", err)
	}

	srv := httptest.NewServer(api.NewRouter(conn))
	defer srv.Close()

	// 1. bootstrap key issuance
	issueBody, _ := json.Marshal(map[string]string{"name": "smoketest"})
	resp, err := http.Post(srv.URL+"/api/v1/keys", "application/json", bytes.NewReader(issueBody))
	if err != nil {
		log.Fatal(err)
	}
	var issued api.IssueKeyResult
	if err := json.NewDecoder(resp.Body).Decode(&issued); err != nil {
		log.Fatal(err)
	}
	_ = resp.Body.Close()
	if issued.APIKey == "" {
		log.Fatal("FAIL: bootstrap key issuance returned empty key")
	}
	fmt.Println("PASS: bootstrap key issued")

	// 2. second unauthenticated key issuance must be rejected
	resp, err = http.Post(srv.URL+"/api/v1/keys", "application/json", bytes.NewReader(issueBody))
	if err != nil {
		log.Fatal(err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		log.Fatalf("FAIL: expected 401 for unauthenticated 2nd key issuance, got %d", resp.StatusCode)
	}
	fmt.Println("PASS: bootstrap closes after first key (2nd unauthenticated request -> 401)")

	// 3. start an SSE subscriber before creating the decision, so we can observe the "created" event
	sseReq, _ := http.NewRequest(http.MethodGet, srv.URL+"/api/v1/events", nil)
	sseReq.Header.Set("Authorization", "Bearer "+issued.APIKey)
	sseResp, err := http.DefaultClient.Do(sseReq)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = sseResp.Body.Close() }()
	if ct := sseResp.Header.Get("Content-Type"); ct != "text/event-stream" {
		log.Fatalf("FAIL: expected text/event-stream, got %q", ct)
	}
	sseLines := make(chan string, 16)
	go func() {
		scanner := bufio.NewScanner(sseResp.Body)
		for scanner.Scan() {
			sseLines <- scanner.Text()
		}
	}()

	// 4. create a decision
	createBody, _ := json.Marshal(api.CreateDecisionInput{
		Source:  "smoketest",
		AgentID: "claude-01",
		Subject: "task#1",
		Level:   db.LevelInterrupt,
		Reason:  "destructive_action",
		Summary: "About to run a migration, please confirm",
	})
	req, _ := http.NewRequest(http.MethodPost, srv.URL+"/api/v1/decisions", bytes.NewReader(createBody))
	req.Header.Set("Authorization", "Bearer "+issued.APIKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	var item db.DecisionItem
	if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
		log.Fatal(err)
	}
	_ = resp.Body.Close()
	if item.ID == "" || item.Status != db.StatusPending {
		log.Fatalf("FAIL: unexpected decision on create: %+v", item)
	}
	fmt.Println("PASS: decision created")

	// 5. an invalid level must be rejected
	badBody, _ := json.Marshal(api.CreateDecisionInput{AgentID: "x", Summary: "y", Level: "not-a-level"})
	req, _ = http.NewRequest(http.MethodPost, srv.URL+"/api/v1/decisions", bytes.NewReader(badBody))
	req.Header.Set("Authorization", "Bearer "+issued.APIKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		log.Fatalf("FAIL: expected 400 for invalid level, got %d", resp.StatusCode)
	}
	fmt.Println("PASS: invalid level rejected")

	// 6. list with status=pending filter finds it
	req, _ = http.NewRequest(http.MethodGet, srv.URL+"/api/v1/decisions?status=pending", nil)
	req.Header.Set("Authorization", "Bearer "+issued.APIKey)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	var listResult api.ListDecisionsResult
	if err := json.NewDecoder(resp.Body).Decode(&listResult); err != nil {
		log.Fatal(err)
	}
	_ = resp.Body.Close()
	if listResult.Total != 1 {
		log.Fatalf("FAIL: expected 1 pending decision, got %d", listResult.Total)
	}
	fmt.Println("PASS: filtered list finds the pending decision")

	// 7. approve it
	approveBody, _ := json.Marshal(map[string]string{"feedback": "looks safe"})
	req, _ = http.NewRequest(http.MethodPost, srv.URL+"/api/v1/decisions/"+item.ID+"/approve", bytes.NewReader(approveBody))
	req.Header.Set("Authorization", "Bearer "+issued.APIKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	var resolved db.DecisionItem
	if err := json.NewDecoder(resp.Body).Decode(&resolved); err != nil {
		log.Fatal(err)
	}
	_ = resp.Body.Close()
	if resolved.Status != db.StatusApproved || resolved.ResolutionFeedback != "looks safe" {
		log.Fatalf("FAIL: unexpected resolution: %+v", resolved)
	}
	fmt.Println("PASS: decision approved")

	// 8. double-resolution must be rejected with 409
	req, _ = http.NewRequest(http.MethodPost, srv.URL+"/api/v1/decisions/"+item.ID+"/reject", bytes.NewReader(approveBody))
	req.Header.Set("Authorization", "Bearer "+issued.APIKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusConflict {
		log.Fatalf("FAIL: expected 409 for double resolution, got %d", resp.StatusCode)
	}
	fmt.Println("PASS: double resolution rejected with 409")

	// 9. stats reflects the resolved item
	req, _ = http.NewRequest(http.MethodGet, srv.URL+"/api/v1/stats", nil)
	req.Header.Set("Authorization", "Bearer "+issued.APIKey)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	var stats api.StatsResult
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		log.Fatal(err)
	}
	_ = resp.Body.Close()
	if stats.Pending != 0 {
		log.Fatalf("FAIL: expected 0 pending after resolution, got %d", stats.Pending)
	}
	fmt.Println("PASS: stats reflects resolution")

	// 10. confirm the SSE subscriber saw both the created and resolved events
	seenCreated, seenResolved := false, false
	deadline := time.After(3 * time.Second)
	for !seenCreated || !seenResolved {
		select {
		case line := <-sseLines:
			if strings.Contains(line, "decision:created") {
				seenCreated = true
			}
			if strings.Contains(line, "decision:resolved") {
				seenResolved = true
			}
		case <-deadline:
			log.Fatalf("FAIL: SSE timeout waiting for events (created=%v resolved=%v)", seenCreated, seenResolved)
		}
	}
	fmt.Println("PASS: SSE delivered created and resolved events")

	fmt.Println("SMOKE TEST OK")
}
