<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import { useDecisionsStore } from '@/stores/decisions'
import { useI18n } from '@/i18n'
import { levelColor, levelIcon } from '@/levelColors'
import DecisionDetailDrawer from '@/components/DecisionDetailDrawer.vue'
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime'

const { t } = useI18n()
const store = useDecisionsStore()

function loadPending() {
  store.filters.status = 'pending'
  store.search()
}

function onLevelChange() {
  store.offset = 0
  loadPending()
}

function fmt(v: any): string {
  const d = new Date(v)
  return isNaN(d.getTime()) ? String(v) : d.toLocaleString(undefined, { month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' })
}

onMounted(() => {
  loadPending()
  EventsOn('decision:created', loadPending)
  EventsOn('decision:resolved', loadPending)
})

onUnmounted(() => {
  EventsOff('decision:created')
  EventsOff('decision:resolved')
})
</script>

<template>
  <div class="flex h-full">
    <aside class="w-56 border-r border-border p-4 space-y-3 shrink-0">
      <h2 class="text-sm font-semibold">{{ t('inbox.title') }}</h2>
      <div class="space-y-1">
        <label class="text-xs text-muted-foreground">{{ t('inbox.filter.level') }}</label>
        <select v-model="store.filters.level" @change="onLevelChange" class="w-full text-sm border border-border rounded px-2 py-1">
          <option value="">{{ t('inbox.filter.all') }}</option>
          <option value="interrupt">{{ t('level.interrupt') }}</option>
          <option value="urgent">{{ t('level.urgent') }}</option>
          <option value="emergency_stop">{{ t('level.emergency_stop') }}</option>
        </select>
      </div>
    </aside>

    <main class="flex-1 overflow-y-auto p-4">
      <div v-if="store.error" class="text-sm border rounded px-3 py-2 mb-3 border-red-300 text-red-600">
        {{ t('error.prefix') }}{{ store.error }}
        <button @click="loadPending" class="ml-2 underline">{{ t('error.retry') }}</button>
      </div>

      <div v-if="store.loading" class="text-sm text-muted-foreground">{{ t('loading') }}</div>
      <div v-else-if="store.items.length === 0" class="text-sm text-muted-foreground">{{ t('inbox.empty') }}</div>

      <div v-else class="space-y-1.5">
        <div
          v-for="item in store.items"
          :key="item.id"
          @click="store.selectItem(item.id)"
          class="flex items-center gap-3 px-3 py-2 rounded-lg border border-border cursor-pointer hover:bg-gray-50 border-l-4"
          :style="{ borderLeftColor: levelColor(item.level) }"
        >
          <span
            class="shrink-0 inline-flex items-center gap-1 px-1.5 py-0.5 rounded-full text-white text-[11px]"
            :style="{ background: levelColor(item.level) }"
          >
            <span aria-hidden="true">{{ levelIcon(item.level) }}</span>
            {{ t('level.' + item.level) }}
          </span>
          <span class="text-xs text-muted-foreground shrink-0">{{ fmt(item.received_at) }}</span>
          <span class="text-xs font-medium shrink-0">{{ item.agent_id }}</span>
          <span class="text-sm truncate flex-1">{{ item.summary }}</span>
        </div>
      </div>
    </main>

    <DecisionDetailDrawer v-if="store.selected" :item="store.selected" @close="store.clearSelection" />
  </div>
</template>
