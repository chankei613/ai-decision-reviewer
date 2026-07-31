# AI Decision Reviewer — 仕様書

> 作成: 2026-07-31
> ステータス: 設計フェーズ

---

## 1. 製品概要

**「AIが勝手に判断した」リスクを軽減するレビューUI** — concept.md 6章の「Decision Queue」の実体。
AIが実行中に自分で「ここは人間の判断が要る」とフラグを立てたものを、人間がインボックス形式で
レビューし、承認・差し戻し・フィードバックを行う。

### 解決する問題

- 自律実行するAIが、低信頼度・破壊的操作・予算超過などの局面でも人間の確認なしに進んでしまう
- 逆に、全てを人間承認必須にすると自律実行の価値が失われる（Autonomy Matrixが解く問題と表裏）
- 「判断が必要な場面だけ」を人間に集約する場所が無い

### ソリューション

任意のAIシステムが `POST /api/v1/decisions` でエスカレーションを送ると、インボックスに積まれる。
人間は一覧から確認し、Approve/Reject＋任意のフィードバックで応答する。AIは
`GET /api/v1/decisions/:id` をポーリングして結果を受け取り、処理を再開/中断する。

---

## 2. コアコンセプト

### エスカレーションレベル

concept.md 6章の5段階ラダーのうち、**人間のアクションを要する3段階**を本製品の対象にする
（Silent/Loggedは通常動作でExecution Ledger側の領分、本製品はキューに積む価値のあるものだけ扱う）。

| レベル | 意味 | 例 |
|---|---|---|
| `interrupt` | 判断要求。処理は一時停止して待つ | 破壊的操作の実行可否、曖昧な仕様の解釈 |
| `urgent` | 早急な対応が必要 | 予算80%超過、連続失敗 |
| `emergency_stop` | 緊急停止。スコープ違反等 | 予算超過、許可されていない操作の試行 |

### DecisionItem（1件のエスカレーション）

```typescript
interface DecisionItem {
  id: string
  received_at: string          // サーバー側で付与

  source: string                 // 送信元システム（自由記述。"comet-taskAI" | "ai-scheduler" | 任意）
  agent_id: string
  subject: string                 // 対象タスクの自由記述ID（例: "task#4821"）

  level: "interrupt" | "urgent" | "emergency_stop"
  reason: string                  // 機械可読な理由コード（例: "budget_exceeded", "destructive_action", "low_confidence"）
  summary: string                 // 人間が読む前提の1〜2文
  context: Record<string, unknown> // 判断材料（自由形式。confidence内訳・提案アクション等）

  status: "pending" | "approved" | "rejected"
  resolution: {
    decision: "approve" | "reject"
    feedback: string
    resolved_at: string
    resolved_by: string
  } | null
}
```

**一度だけ解決できる**: `resolution` が既に埋まっているアイテムへの再度のapprove/rejectは409で拒否する
（監査ログとしての一貫性を保つため。Execution Ledgerの「追記専用」思想を踏襲するが、
本製品は「未解決→解決」の1回だけの状態遷移を許す点で異なる）。

---

## 3. 機能一覧

### Phase 1 (MVP: Ingestion + 解決API)

| 機能 | 説明 |
|------|------|
| アイテム追加API | `POST /api/v1/decisions`（APIキー認証） |
| 一覧・フィルタ | status・level・source・agent_idで絞り込み + ページネーション |
| 承認/差し戻しAPI | `POST /api/v1/decisions/:id/approve` \| `/reject`（feedback任意） |
| 二重解決防止 | 解決済みアイテムへの再解決は409 |
| APIキー管理 | ブートストラップ発行（0件時のみ未認証）・失効 |

### Phase 2 (集計・リアルタイム)

| 機能 | 説明 |
|------|------|
| 集計API | pending件数・レベル別件数・平均解決時間 |
| SSEストリーム | 新規追加・解決イベントをリアルタイム配信 |

### Phase 3 (UI)

| 機能 | 説明 |
|------|------|
| インボックスビュー | pending一覧。レベル別の色分け・並び替え・SSEでリアルタイム更新 |
| 詳細ドロワー | summary/context全表示 + Approve/Reject/Feedbackフォーム |
| 解決済みビュー | 履歴一覧・フィルタ |
| ダッシュボード | pending件数・レベル別内訳・平均解決時間 |
| Help画面 | インボックスの使い方・レベルの意味・承認/差し戻しの流れ |

---

## 4. UX フロー

```
AIが実行中に判断が必要と判断
 └── POST /api/v1/decisions（level/reason/summary/context）
      └── インボックスに積まれる（SSEで即座にUIへ反映）
           └── 人間がレビュー
                ├── Approve（+feedback任意） → AIは処理を再開
                └── Reject（+feedback任意）  → AIは処理を中断/別経路へ

AI側
 └── POST後、GET /api/v1/decisions/:id をポーリング
      └── status が pending でなくなったら resolution を読んで続行判断
```

---

## 5. データストア

SQLite（ローカル、`~/.ai-decision-reviewer/decisions.db`）

```sql
decision_items (
  id, received_at, source, agent_id, subject,
  level, reason, summary, context JSON,
  status,
  resolution_decision, resolution_feedback, resolution_resolved_at, resolution_resolved_by
)
agent_keys (id, name, api_key_hash, created_at, revoked_at)
```

`context` はJSONカラム（GORM `serializer:json`、既存製品と同じパターン）。
`resolution_*` はnull許容カラムとしてDecisionItemにフラットに持たせる（別テーブルに分けない。
1:1で常に一緒に読み書きするため、JOINのコストを避ける）。
