<script setup lang="ts">
import { onMounted } from 'vue'
import { useDecisionsStore } from '@/stores/decisions'
import { useI18n } from '@/i18n'
import { levelColor, levelIcon } from '@/levelColors'

const { t } = useI18n()
const store = useDecisionsStore()

function load() {
  if (store.filters.status === 'pending' || store.filters.status === '') {
    store.filters.status = 'approved'
  }
  store.search()
}

function onStatusChange() {
  store.offset = 0
  store.search()
}

function fmt(v: any): string {
  const d = new Date(v)
  return isNaN(d.getTime()) ? String(v) : d.toLocaleString()
}

onMounted(load)
</script>

<template>
  <div class="p-4 space-y-4">
    <div class="flex items-center justify-between">
      <h2 class="text-sm font-semibold">{{ t('resolved.title') }}</h2>
      <select v-model="store.filters.status" @change="onStatusChange" class="text-sm border border-border rounded px-2 py-1">
        <option value="approved">{{ t('resolved.status.approved') }}</option>
        <option value="rejected">{{ t('resolved.status.rejected') }}</option>
      </select>
    </div>

    <div v-if="store.error" class="text-sm border rounded px-3 py-2 border-red-300 text-red-600">
      {{ t('error.prefix') }}{{ store.error }}
      <button @click="store.search" class="ml-2 underline">{{ t('error.retry') }}</button>
    </div>

    <div v-if="store.loading" class="text-sm text-muted-foreground">{{ t('loading') }}</div>
    <div v-else-if="store.items.length === 0" class="text-sm text-muted-foreground">{{ t('resolved.empty') }}</div>

    <div v-else class="space-y-1.5">
      <div
        v-for="item in store.items"
        :key="item.id"
        class="border border-border rounded-lg px-3 py-2 border-l-4"
        :style="{ borderLeftColor: levelColor(item.level) }"
      >
        <div class="flex items-center gap-2 text-xs">
          <span
            class="inline-flex items-center gap-1 px-1.5 py-0.5 rounded-full text-white"
            :style="{ background: levelColor(item.level) }"
          >
            <span aria-hidden="true">{{ levelIcon(item.level) }}</span>
            {{ t('level.' + item.level) }}
          </span>
          <span class="font-medium">{{ item.agent_id }}</span>
          <span class="text-muted-foreground">{{ item.subject }}</span>
          <span class="text-muted-foreground ml-auto">{{ fmt(item.resolution_resolved_at) }}</span>
        </div>
        <p class="text-sm mt-1">{{ item.summary }}</p>
        <p v-if="item.resolution_feedback" class="text-xs text-muted-foreground mt-1 italic">
          "{{ item.resolution_feedback }}"
        </p>
      </div>
    </div>

    <div v-if="store.total > 0" class="flex items-center justify-between pt-2 text-xs text-muted-foreground">
      <button :disabled="store.offset === 0" @click="store.prevPage" class="disabled:opacity-30 underline">Prev</button>
      <span>{{ store.offset + 1 }}–{{ Math.min(store.offset + store.limit, store.total) }} / {{ store.total }}</span>
      <button :disabled="store.offset + store.limit >= store.total" @click="store.nextPage" class="disabled:opacity-30 underline">Next</button>
    </div>
  </div>
</template>
