import { describe, it, expect, beforeEach } from 'vitest'
import { useI18n } from '@/i18n'

describe('i18n', () => {
  beforeEach(() => {
    window.localStorage.clear()
  })

  it('defaults to ja and resolves known keys', () => {
    const { t, locale } = useI18n()
    expect(locale.value).toBe('ja')
    expect(t('nav.inbox')).toBe('インボックス')
  })

  it('toggleLocale switches between ja and en and persists', () => {
    const { t, locale, toggleLocale } = useI18n()
    toggleLocale()
    expect(locale.value).toBe('en')
    expect(t('nav.inbox')).toBe('Inbox')
    expect(window.localStorage.getItem('locale')).toBe('en')
  })

  it('falls back to the key itself for unknown keys', () => {
    const { t } = useI18n()
    expect(t('does.not.exist')).toBe('does.not.exist')
  })
})
