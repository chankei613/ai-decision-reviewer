import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useDecisionsStore } from '@/stores/decisions'
import { ListDecisions, ApproveDecision } from '../../wailsjs/go/main/App'
import { api } from '../../wailsjs/go/models'

describe('decisions store error handling', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.mocked(ListDecisions).mockReset()
    vi.mocked(ApproveDecision).mockReset()
  })

  it('captures a failed search() as store.error and clears loading', async () => {
    vi.mocked(ListDecisions).mockRejectedValueOnce(new Error('network down'))
    const store = useDecisionsStore()

    await store.search()

    expect(store.loading).toBe(false)
    expect(store.error).toContain('network down')
  })

  it('clears the previous error on a successful retry', async () => {
    vi.mocked(ListDecisions).mockRejectedValueOnce(new Error('network down'))
    const store = useDecisionsStore()
    await store.search()
    expect(store.error).not.toBeNull()

    vi.mocked(ListDecisions).mockResolvedValueOnce(api.ListDecisionsResult.createFrom({ items: [], total: 0 }))
    await store.search()

    expect(store.error).toBeNull()
  })

  it('captures a failed approve() (e.g. 409 already-resolved) without throwing', async () => {
    vi.mocked(ApproveDecision).mockRejectedValueOnce(new Error('decision already resolved'))
    const store = useDecisionsStore()

    await store.approve('item-1', 'looks fine')

    expect(store.resolving).toBe(false)
    expect(store.error).toContain('already resolved')
  })
})
