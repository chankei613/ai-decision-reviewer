import { defineStore } from 'pinia'
import { Stats } from '../../wailsjs/go/main/App'

export const useStatsStore = defineStore('stats', {
  state: () => ({
    pending: 0,
    pendingByLevel: {} as Record<string, number>,
    avgResolutionSeconds: 0,
    loading: false,
    error: null as string | null,
  }),
  actions: {
    async load() {
      this.loading = true
      this.error = null
      try {
        const result = await Stats()
        this.pending = result.pending
        this.pendingByLevel = result.pending_by_level ?? {}
        this.avgResolutionSeconds = result.avg_resolution_seconds
      } catch (e) {
        this.error = String(e)
      } finally {
        this.loading = false
      }
    },
  },
})
