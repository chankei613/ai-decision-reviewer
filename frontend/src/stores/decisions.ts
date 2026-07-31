import { defineStore } from 'pinia'
import {
  ListDecisions,
  GetDecision,
  ApproveDecision,
  RejectDecision,
} from '../../wailsjs/go/main/App'
import { db, api } from '../../wailsjs/go/models'

export interface DecisionFiltersState {
  status: string
  level: string
  source: string
  agentId: string
}

function emptyFilters(): DecisionFiltersState {
  return { status: '', level: '', source: '', agentId: '' }
}

export const useDecisionsStore = defineStore('decisions', {
  state: () => ({
    filters: emptyFilters(),
    items: [] as db.DecisionItem[],
    total: 0,
    limit: 50,
    offset: 0,
    selected: null as db.DecisionItem | null,
    loading: false,
    resolving: false,
    error: null as string | null,
  }),
  actions: {
    async search() {
      this.loading = true
      this.error = null
      try {
        const result: api.ListDecisionsResult = await ListDecisions(
          this.filters.status,
          this.filters.level,
          this.filters.source,
          this.filters.agentId,
          this.limit,
          this.offset,
        )
        this.items = result.items ?? []
        this.total = result.total
      } catch (e) {
        this.error = String(e)
      } finally {
        this.loading = false
      }
    },
    async selectItem(id: string) {
      this.error = null
      try {
        this.selected = await GetDecision(id)
      } catch (e) {
        this.error = String(e)
      }
    },
    clearSelection() {
      this.selected = null
    },
    async approve(id: string, feedback: string) {
      this.resolving = true
      this.error = null
      try {
        await ApproveDecision(id, feedback)
        this.clearSelection()
        await this.search()
      } catch (e) {
        this.error = String(e)
      } finally {
        this.resolving = false
      }
    },
    async reject(id: string, feedback: string) {
      this.resolving = true
      this.error = null
      try {
        await RejectDecision(id, feedback)
        this.clearSelection()
        await this.search()
      } catch (e) {
        this.error = String(e)
      } finally {
        this.resolving = false
      }
    },
    resetFilters() {
      this.filters = emptyFilters()
      this.offset = 0
    },
    nextPage() {
      if (this.offset + this.limit < this.total) {
        this.offset += this.limit
        this.search()
      }
    },
    prevPage() {
      if (this.offset > 0) {
        this.offset = Math.max(0, this.offset - this.limit)
        this.search()
      }
    },
  },
})
