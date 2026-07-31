import { defineStore } from 'pinia'
import { ListKeys, IssueKey, RevokeKey, GetAPIURL, GetAppVersion } from '../../wailsjs/go/main/App'
import { db } from '../../wailsjs/go/models'

export const useKeysStore = defineStore('keys', {
  state: () => ({
    keys: [] as db.AgentKey[],
    apiURL: '',
    appVersion: '',
    lastIssuedKey: null as { name: string; apiKey: string } | null,
    loading: false,
    error: null as string | null,
  }),
  actions: {
    async load() {
      this.loading = true
      this.error = null
      try {
        const [keys, url, version] = await Promise.all([ListKeys(), GetAPIURL(), GetAppVersion()])
        this.keys = keys ?? []
        this.apiURL = url
        this.appVersion = version
      } catch (e) {
        this.error = String(e)
      } finally {
        this.loading = false
      }
    },
    async issue(name: string) {
      this.error = null
      try {
        const result = await IssueKey(name)
        this.lastIssuedKey = { name: result.name, apiKey: result.api_key }
        await this.load()
      } catch (e) {
        this.error = String(e)
      }
    },
    async revoke(id: string) {
      this.error = null
      try {
        await RevokeKey(id)
        await this.load()
      } catch (e) {
        this.error = String(e)
      }
    },
    dismissIssuedKey() {
      this.lastIssuedKey = null
    },
  },
})
