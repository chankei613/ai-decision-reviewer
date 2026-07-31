# 開発計画

**予測期間:** 2〜3週間相当（ロードマップ見積もり）。過去実績（A〜E）に基づき短縮を狙う。

| Phase | 内容 |
|---|---|
| Phase 0 | プロジェクト立ち上げ（Go初期化・docs・GitHub repo） |
| Phase 1 | データモデル・Ingestion/CRUD API（APIキー認証・二重解決防止） |
| Phase 2 | 集計API・SSEストリーム |
| Phase 3 | Wails + Vue3 UI（インボックス・詳細ドロワー・解決済みビュー・ダッシュボード・Help） |
| Phase 4 | 仕上げ・署名・配布・LP |

## 優先順位の根拠

concept.md 9章の指示通り、Decision QueueはUI・分析より前に「エスカレーションが正しく積まれ、
承認/差し戻しが正しく反映される」ことが最重要。Phase 1（Ingestion）とPhase 2（SSE）を
UIより先に固め、`curl`ベースの手動テストで正しさを検証してからUIに進む
（execution-ledgerと同じ優先順位判断）。
