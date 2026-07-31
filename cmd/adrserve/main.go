// cmd/adrserve はAI Decision Reviewer APIをlocalhostで提供する単体サーバー。
// フロントエンドのdev時（`npm run dev` + このサーバー）や、Wailsを介さずヘッドレスで
// 動かしたい場合に使う。Wailsアプリ本体（ルートのmain.go）も同じ internal/api.NewRouter を
// 使い回すため、挙動はズレない。
//
//	go run ./cmd/adrserve -addr :8423 -db ai-decision-reviewer.db
package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/chankei613/ai-decision-reviewer/internal/api"
	"github.com/chankei613/ai-decision-reviewer/internal/db"
)

func main() {
	addr := flag.String("addr", ":8423", "待ち受けアドレス")
	dbPath := flag.String("db", "ai-decision-reviewer.db", "SQLiteファイル")
	flag.Parse()

	conn, err := db.Init(*dbPath)
	if err != nil {
		log.Fatalf("db init failed: %v", err)
	}

	router := api.NewRouter(conn)
	log.Printf("ai-decision-reviewer backend listening on %s", *addr)
	if err := http.ListenAndServe(*addr, router); err != nil {
		log.Fatal(err)
	}
}
