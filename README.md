# AI Decision Reviewer

「AIが勝手に判断した」リスクを軽減するレビューUI — comet-taskAI ロードマップ Product G。
concept.mdの「Decision Queue」の実体。

AIが実行中に自分で「ここは人間の判断が要る」とフラグを立てたもの（低信頼度・破壊的操作・
予算超過・繰り返し失敗など）をインボックス形式でレビューし、承認・差し戻し・フィードバックを行う。

詳細は [docs/spec.md](docs/spec.md) を参照。

## 現在のステータス: Phase 1-2（Ingestion/CRUD API + SSE）完了

- [x] Phase 0: プロジェクト立ち上げ
- [x] Phase 1: データモデル・Ingestion/CRUD API（APIキー認証・二重解決防止）
- [x] Phase 2: 集計API・SSEストリーム
- [ ] Phase 3: Wails + Vue3 UI
- [ ] Phase 4: 仕上げ・署名・配布・LP

## 使い方（開発用ヘッドレスサーバー）

```bash
go mod tidy   # 依存解決
go run .      # :8423 でAPIサーバー起動（SQLite: ai-decision-reviewer.db）
go run ./cmd/smoketest  # bootstrap鍵発行 → decision追加 → 承認 → 二重解決の拒否 → SSE配信 の一連を確認する自己完結テスト
```

### APIキー認証

`AgentKey`が0件のときのみ `POST /api/v1/keys` を未認証で許可する（最初の1件を発行するため）。
1件発行された時点で以降は `Authorization: Bearer <key>` が必須になる。

### エスカレーションの送信・解決

```bash
curl -X POST localhost:8423/api/v1/decisions \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "source": "comet-taskAI",
    "agent_id": "claude-01",
    "subject": "task#4821",
    "level": "interrupt",
    "reason": "destructive_action",
    "summary": "本番DBのマイグレーションを実行してよいか確認したい"
  }'

curl -X POST localhost:8423/api/v1/decisions/{id}/approve \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"feedback": "内容確認済み、進めてOK"}'
```

AI側は `GET /api/v1/decisions/{id}` をポーリングするか、`GET /api/v1/events` のSSEを購読して
`status` が `pending` でなくなったタイミングで結果を受け取る。

## API

| メソッド | パス | 用途 |
|---|---|---|
| POST | `/api/v1/keys` | APIキー発行（ブートストラップ時のみ未認証） |
| GET | `/api/v1/keys` | 発行済みキー一覧 |
| DELETE | `/api/v1/keys/{id}` | キー失効 |
| POST | `/api/v1/decisions` | エスカレーション追加（level: interrupt/urgent/emergency_stop） |
| GET | `/api/v1/decisions` | 一覧・フィルタ（status/level/source/agent_id + limit/offset） |
| GET | `/api/v1/decisions/{id}` | 単体取得（AIのポーリング先） |
| POST | `/api/v1/decisions/{id}/approve` | 承認（feedback任意、既に解決済みなら409） |
| POST | `/api/v1/decisions/{id}/reject` | 差し戻し（feedback任意、既に解決済みなら409） |
| GET | `/api/v1/stats` | pending件数・レベル別件数・平均解決時間 |
| GET | `/api/v1/events` | SSEストリーム（decision:created / decision:resolved） |

## ディレクトリ構成

```
internal/db/       GORMモデル（DecisionItem/AgentKey）・SQLite初期化
internal/api/       REST API（keys/decisions/stats/events）+ 認証ミドルウェア
internal/events/    SSE配信用pub/sub（Broker）
cmd/smoketest/      bootstrap→decision追加→承認→二重解決拒否→SSE配信の通しスモークテスト
docs/               設計ドキュメント
```
