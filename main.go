package main

import (
	"log"
	"net/http"

	"github.com/chankei613/ai-decision-reviewer/internal/api"
	"github.com/chankei613/ai-decision-reviewer/internal/db"
)

func main() {
	conn, err := db.Init("ai-decision-reviewer.db")
	if err != nil {
		log.Fatalf("db init failed: %v", err)
	}

	router := api.NewRouter(conn)
	log.Println("ai-decision-reviewer backend listening on :8423")
	if err := http.ListenAndServe(":8423", router); err != nil {
		log.Fatal(err)
	}
}
