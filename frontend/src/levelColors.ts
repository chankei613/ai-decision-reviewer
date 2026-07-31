// レベル配色は dataviz スキルの検証済みパレット（status palette）をそのまま使う。
// 色だけで意味を伝えない: 呼び出し側は必ずアイコン+ラベルとセットで使うこと。
export type Level = 'interrupt' | 'urgent' | 'emergency_stop'

export const LEVEL_COLORS: Record<Level, string> = {
  interrupt: '#ec835a',
  urgent: '#fab219',
  emergency_stop: '#d03b3b',
}

const ICON_BY_LEVEL: Record<Level, string> = {
  interrupt: '?',
  urgent: '!',
  emergency_stop: '⛔',
}

export function levelColor(level: string): string {
  return LEVEL_COLORS[level as Level] ?? '#8a8a86'
}

export function levelIcon(level: string): string {
  return ICON_BY_LEVEL[level as Level] ?? '·'
}
