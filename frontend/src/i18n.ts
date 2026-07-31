import { ref } from 'vue'

export type Locale = 'en' | 'ja'

// localStorage に保存して再起動後も維持する
const saved = localStorage.getItem('locale') as Locale | null
const locale = ref<Locale>(saved === 'en' || saved === 'ja' ? saved : 'ja')

const messages: Record<Locale, Record<string, string>> = {
  en: {
    'app.subtitle': 'AI Decision Reviewer',
    'lang.toggle': 'JA',
    'nav.inbox': 'Inbox',
    'nav.resolved': 'Resolved',
    'nav.dashboard': 'Dashboard',
    'nav.help': 'Help',
    'nav.settings': 'Settings',

    'error.prefix': 'Error: ',
    'error.retry': 'Retry',
    'loading': 'Loading…',

    'level.interrupt': 'Interrupt',
    'level.urgent': 'Urgent',
    'level.emergency_stop': 'Emergency stop',

    'inbox.title': 'Inbox',
    'inbox.empty': 'Nothing waiting for review. New escalations will appear here in real time.',
    'inbox.filter.level': 'Level',
    'inbox.filter.all': 'All levels',

    'detail.title': 'Decision detail',
    'detail.close': 'Close',
    'detail.summary': 'Summary',
    'detail.reason': 'Reason',
    'detail.context': 'Context',
    'detail.context.empty': 'No additional context',
    'detail.feedback.placeholder': 'Optional feedback for the AI…',
    'detail.approve': 'Approve',
    'detail.reject': 'Reject',
    'detail.alreadyResolved': 'This item was already resolved.',

    'resolved.title': 'Resolved',
    'resolved.empty': 'Nothing resolved yet.',
    'resolved.filter.status': 'Status',
    'resolved.filter.all': 'All',
    'resolved.status.approved': 'Approved',
    'resolved.status.rejected': 'Rejected',

    'dashboard.title': 'Dashboard',
    'dashboard.pending': 'Pending',
    'dashboard.avgResolution': 'Avg. resolution time',
    'dashboard.byLevel': 'Pending by level',
    'dashboard.empty': 'No pending items — nothing to show yet.',

    'settings.title': 'Settings',
    'settings.api.title': 'API endpoint',
    'settings.api.desc': 'AI systems POST escalations here and poll for the resolution.',
    'settings.keys.title': 'API keys',
    'settings.keys.name': 'Key name',
    'settings.keys.issue': 'Issue key',
    'settings.keys.issued': 'Key issued — copy it now, it will not be shown again',
    'settings.keys.copy': 'Copy',
    'settings.keys.revoke': 'Revoke',
    'settings.keys.revoked': 'Revoked',
    'settings.keys.empty': 'No keys issued yet. Issue one to allow AI systems to send escalations.',
    'settings.version': 'Version',
    'settings.quit': 'Quit',
    'settings.quit.confirm': 'Quit the app? Ingestion will stop until you reopen it.',

    'help.title': 'Help',
    'help.intro': 'What Decision Queue is for, and how a review normally goes.',

    'help.what.title': 'What this app does',
    'help.what.body': 'Any AI system can POST an escalation here when it hits something it should not decide alone — a destructive action, a low-confidence judgment call, a budget overrun. You review it and either Approve or Reject, optionally with feedback. The AI polls (or listens over SSE) for your decision and continues accordingly.',

    'help.levels.title': 'The three levels',
    'help.levels.interrupt': 'Interrupt — the AI is paused and waiting. Typical for ambiguous instructions or a destructive action that needs a yes/no.',
    'help.levels.urgent': 'Urgent — needs attention soon. Typical for a budget nearing its limit or repeated failures.',
    'help.levels.emergency': 'Emergency stop — something went outside the agreed scope. Review these first.',

    'help.flow.title': 'Reviewing an item',
    'help.flow.1': 'Open Inbox — pending items appear there in real time as they arrive.',
    'help.flow.2': 'Click an item to open its detail — read the summary and any context the AI attached.',
    'help.flow.3': 'Approve or Reject. Add feedback if it helps the AI understand why.',
    'help.flow.4': 'Once resolved, an item cannot be resolved again — it moves to Resolved as a permanent record.',

    'help.stuck.title': 'Common snags',
    'help.stuck.1': 'Approve/Reject button does nothing → the item was likely already resolved (e.g. from another client); check Resolved.',
    'help.stuck.2': 'Nothing shows up in Inbox → confirm the sending AI system has a valid API key and is POSTing to the address shown in Settings.',
    'help.stuck.3': 'Want to see history → check Resolved, which keeps every past decision with its feedback.',
  },
  ja: {
    'app.subtitle': 'AI Decision Reviewer',
    'lang.toggle': 'EN',
    'nav.inbox': 'インボックス',
    'nav.resolved': '解決済み',
    'nav.dashboard': 'ダッシュボード',
    'nav.help': 'ヘルプ',
    'nav.settings': '設定',

    'error.prefix': 'エラー: ',
    'error.retry': '再試行',
    'loading': '読み込み中…',

    'level.interrupt': '判断要求',
    'level.urgent': '至急',
    'level.emergency_stop': '緊急停止',

    'inbox.title': 'インボックス',
    'inbox.empty': 'レビュー待ちのアイテムはありません。新しいエスカレーションはここにリアルタイムで表示されます。',
    'inbox.filter.level': 'レベル',
    'inbox.filter.all': 'すべてのレベル',

    'detail.title': '判断の詳細',
    'detail.close': '閉じる',
    'detail.summary': 'サマリ',
    'detail.reason': '理由',
    'detail.context': 'コンテキスト',
    'detail.context.empty': '追加情報はありません',
    'detail.feedback.placeholder': 'AIへのフィードバック（任意）…',
    'detail.approve': '承認',
    'detail.reject': '差し戻し',
    'detail.alreadyResolved': 'このアイテムは既に解決済みです。',

    'resolved.title': '解決済み',
    'resolved.empty': 'まだ解決したアイテムはありません。',
    'resolved.filter.status': 'ステータス',
    'resolved.filter.all': 'すべて',
    'resolved.status.approved': '承認済み',
    'resolved.status.rejected': '差し戻し済み',

    'dashboard.title': 'ダッシュボード',
    'dashboard.pending': '未解決件数',
    'dashboard.avgResolution': '平均解決時間',
    'dashboard.byLevel': 'レベル別の未解決件数',
    'dashboard.empty': '未解決のアイテムがないため、表示するものがありません。',

    'settings.title': '設定',
    'settings.api.title': 'APIエンドポイント',
    'settings.api.desc': 'AIシステムはここへエスカレーションをPOSTし、解決結果をポーリングします。',
    'settings.keys.title': 'APIキー',
    'settings.keys.name': 'キー名',
    'settings.keys.issue': 'キーを発行',
    'settings.keys.issued': 'キーを発行しました — この場では二度と表示されないので今すぐコピーしてください',
    'settings.keys.copy': 'コピー',
    'settings.keys.revoke': '失効',
    'settings.keys.revoked': '失効済み',
    'settings.keys.empty': 'まだキーがありません。AIシステムがエスカレーションを送れるよう発行してください。',
    'settings.version': 'バージョン',
    'settings.quit': '終了',
    'settings.quit.confirm': 'アプリを終了しますか？再度開くまでIngestionは停止します。',

    'help.title': 'ヘルプ',
    'help.intro': 'Decision Queueが何のためにあるか、レビューの流れをまとめました。',

    'help.what.title': 'このアプリでできること',
    'help.what.body': '任意のAIシステムは、自分だけでは決めるべきでない場面（破壊的操作・低信頼度の判断・予算超過など）に遭遇すると、ここへエスカレーションをPOSTします。人間はそれをレビューし、承認または差し戻しを行い、任意でフィードバックを添えます。AIはポーリングまたはSSEでその結果を受け取り、処理を続けます。',

    'help.levels.title': '3つのレベル',
    'help.levels.interrupt': '判断要求 — AIは一時停止して待っています。曖昧な指示や、yes/noの確認が必要な破壊的操作が典型例です。',
    'help.levels.urgent': '至急 — 早めの対応が必要です。予算が上限に近づいている、失敗が繰り返されている等が典型例です。',
    'help.levels.emergency': '緊急停止 — 合意したスコープの外に出てしまった状態です。最優先でレビューしてください。',

    'help.flow.title': 'アイテムのレビュー方法',
    'help.flow.1': 'インボックスを開く — 届いたアイテムはリアルタイムで表示されます。',
    'help.flow.2': 'アイテムをクリックして詳細を開く — サマリとAIが添付したコンテキストを確認します。',
    'help.flow.3': '承認または差し戻しを行う。AIが理解しやすいようフィードバックを添えるとよいです。',
    'help.flow.4': '一度解決したアイテムは再度解決できません — 解決済みへ永続的な記録として移動します。',

    'help.stuck.title': 'よくある詰まりどころ',
    'help.stuck.1': '承認/差し戻しボタンが反応しない → 別のクライアントなどで既に解決済みの可能性があります。解決済みビューを確認してください。',
    'help.stuck.2': 'インボックスに何も表示されない → 送信元のAIシステムが有効なAPIキーを持ち、設定画面に表示されているアドレスへPOSTしているか確認してください。',
    'help.stuck.3': '過去の履歴を見たい → 解決済みビューに、フィードバックも含めた全ての過去の判断が残っています。',
  },
}

export function useI18n() {
  function t(key: string, params?: Record<string, string | number>): string {
    let msg = messages[locale.value][key] ?? messages.en[key] ?? key
    if (params) {
      for (const [k, v] of Object.entries(params)) {
        msg = msg.replace(`{${k}}`, String(v))
      }
    }
    return msg
  }

  function toggleLocale() {
    locale.value = locale.value === 'en' ? 'ja' : 'en'
    localStorage.setItem('locale', locale.value)
  }

  return { t, locale, toggleLocale }
}
