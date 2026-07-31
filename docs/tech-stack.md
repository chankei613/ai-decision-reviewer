# 技術選定

**決定日:** 2026-07-31
**ステータス:** 確定

---

## 決定

| レイヤー | 採用 | 理由 |
|---|---|---|
| Desktop基盤 | Wails v2（v2.12.0固定） | A/B/K/C/D/Eの6製品すべてで実績あり |
| Backend | Go 1.22+ | 同上 |
| Frontend | Vue 3 + Vite + Pinia | 同上 |
| DB | SQLite + GORM | 同上 |
| CI | GitHub Actions | Go 1.23 / macos-14 を最初から採用 |
| 配布 | `.app` / `.exe` シングルバイナリ + コード署名・公証 | `build-release.sh`で環境変数ベースの認証情報を使う |

## 前2製品（execution-ledger, context-bundle-builder）から必ず引き継ぐこと

Phase 4で毎回同じ問題を踏まないよう、Phase 0の時点で以下を確定事項として扱う（後回しにしない）：

1. **`.golangci.yml` はv2形式で最初から書く**（`version: "2"` + `formatters.enable: [gofmt]`）
2. **`frontend/wailsjs/` は `.gitignore` に入れない**。コミット対象にする（CIのtypecheckがクリーンチェックアウトで壊れるのを防ぐ）
3. **CIのgo-testジョブは `./internal/... ./cmd/...` にスコープを限定する**。ルートに`//go:embed all:frontend/dist`があるため`go build ./...`はfrontend/distが無い環境で失敗する
4. **`build-release.sh`・`wails.json`に個人のメールアドレス・実名を直書きしない**。`APPLE_ID`/`TEAM_ID`/`DEVELOPER_NAME`/`APP_PASSWORD`は全て環境変数
5. **release.ymlはharness-managerのv0.1.1構成をそのままコピー**（Go1.23・macos-14・notarytool・windows拡張子対応）
6. **Phase 3のUI実装に、アプリ内Help画面を最初から含める**（2026-07-31にユーザーがリリース必須要件化。後付けにしない）

## Apple Developer認証情報について

TEAM_ID: `D5R956CRBE` / Developer ID Application: `keisuke haraguchi`。
このMacのKeychainに証明書があるため`codesign`はローカルで即実行できる。
公証用のApple app-specific passwordは、execution-ledger/context-bundle-builderで発行したものが
同一Apple ID（chankei613@gmail.com）宛のため使い回せる可能性が高いが、値そのものはセッションを跨ぐと
記憶していないため、必要になったらユーザーに確認する。

## 却下した案

- **Webhook配信**: AI側への通知はポーリング（`GET /api/v1/decisions/:id`）とSSEに留める。
  Webhookは受信側にHTTPサーバーが要り、ローカルファーストの個人開発ツールという方針に対してオーバースペック
