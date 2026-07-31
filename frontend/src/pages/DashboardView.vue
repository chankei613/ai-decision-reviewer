<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useStatsStore } from '@/stores/stats'
import { useI18n } from '@/i18n'
import { levelColor } from '@/levelColors'

const { t } = useI18n()
const store = useStatsStore()

onMounted(() => store.load())

const levelBars = computed(() => {
  const entries = Object.entries(store.pendingByLevel)
  const max = Math.max(1, ...entries.map(([, v]) => v))
  return entries
    .sort((a, b) => b[1] - a[1])
    .map(([level, count]) => ({ level, count, pct: (count / max) * 100, color: levelColor(level) }))
})

function fmtSeconds(s: number): string {
  if (s === 0) return '—'
  if (s < 60) return `${Math.round(s)}s`
  if (s < 3600) return `${Math.round(s / 60)}m`
  return `${(s / 3600).toFixed(1)}h`
}
</script>

<template>
  <div class="p-6 space-y-6 overflow-y-auto h-full">
    <h2 class="text-sm font-semibold">{{ t('dashboard.title') }}</h2>

    <div v-if="store.error" class="text-sm border rounded px-3 py-2 border-red-300 text-red-600">
      {{ t('error.prefix') }}{{ store.error }}
      <button @click="store.load" class="ml-2 underline">{{ t('error.retry') }}</button>
    </div>
    <div v-else-if="store.loading" class="text-sm text-muted-foreground">{{ t('loading') }}</div>

    <template v-else>
      <div class="grid grid-cols-2 gap-3 max-w-md">
        <div class="border border-border rounded-lg p-4">
          <div class="text-xs text-muted-foreground">{{ t('dashboard.pending') }}</div>
          <div class="text-2xl font-semibold tabular-nums mt-1">{{ store.pending }}</div>
        </div>
        <div class="border border-border rounded-lg p-4">
          <div class="text-xs text-muted-foreground">{{ t('dashboard.avgResolution') }}</div>
          <div class="text-2xl font-semibold tabular-nums mt-1">{{ fmtSeconds(store.avgResolutionSeconds) }}</div>
        </div>
      </div>

      <div v-if="levelBars.length">
        <h3 class="text-xs font-semibold text-muted-foreground mb-2">{{ t('dashboard.byLevel') }}</h3>
        <div class="space-y-1.5 max-w-md">
          <div v-for="row in levelBars" :key="row.level" class="flex items-center gap-2 text-sm">
            <span class="w-32 shrink-0 text-xs">{{ t('level.' + row.level) }}</span>
            <div class="flex-1 h-2 bg-gray-100 rounded-full overflow-hidden">
              <div class="h-full rounded-full" :style="{ width: row.pct + '%', background: row.color }" />
            </div>
            <span class="w-8 text-right text-xs tabular-nums text-muted-foreground">{{ row.count }}</span>
          </div>
        </div>
      </div>
      <div v-else class="text-sm text-muted-foreground">{{ t('dashboard.empty') }}</div>
    </template>
  </div>
</template>
