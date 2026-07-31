<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from '@/i18n'
import { levelColor, levelIcon } from '@/levelColors'
import { useDecisionsStore } from '@/stores/decisions'
import { db } from '../../wailsjs/go/models'

const props = defineProps<{ item: db.DecisionItem }>()
const emit = defineEmits<{ close: [] }>()

const { t } = useI18n()
const store = useDecisionsStore()
const feedback = ref('')

const isResolved = props.item.status !== 'pending'

async function approve() {
  await store.approve(props.item.id, feedback.value)
  emit('close')
}

async function reject() {
  await store.reject(props.item.id, feedback.value)
  emit('close')
}

function fmt(v: any): string {
  const d = new Date(v)
  return isNaN(d.getTime()) ? String(v) : d.toLocaleString()
}
</script>

<template>
  <div class="fixed inset-0 z-20 flex justify-end" @click.self="emit('close')" style="background: rgba(0,0,0,0.15)">
    <div class="w-full max-w-md h-full bg-white border-l border-border overflow-y-auto p-5 space-y-4">
      <div class="flex items-center justify-between">
        <h2 class="text-sm font-semibold">{{ t('detail.title') }}</h2>
        <button @click="emit('close')" class="text-muted-foreground hover:text-foreground">{{ t('detail.close') }}</button>
      </div>

      <div class="flex items-center gap-2 text-xs">
        <span
          class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-white"
          :style="{ background: levelColor(item.level) }"
        >
          <span aria-hidden="true">{{ levelIcon(item.level) }}</span>
          {{ t('level.' + item.level) }}
        </span>
        <span class="text-muted-foreground">{{ item.agent_id }} · {{ item.subject }}</span>
      </div>
      <div class="text-xs text-muted-foreground">{{ fmt(item.received_at) }}</div>

      <div class="space-y-1">
        <h3 class="text-xs font-semibold text-muted-foreground">{{ t('detail.summary') }}</h3>
        <p class="text-sm">{{ item.summary }}</p>
      </div>

      <div v-if="item.reason" class="space-y-1">
        <h3 class="text-xs font-semibold text-muted-foreground">{{ t('detail.reason') }}</h3>
        <code class="text-xs bg-gray-50 border border-border rounded px-2 py-1 inline-block">{{ item.reason }}</code>
      </div>

      <div class="space-y-1">
        <h3 class="text-xs font-semibold text-muted-foreground">{{ t('detail.context') }}</h3>
        <pre
          v-if="item.context && Object.keys(item.context).length"
          class="text-xs font-mono whitespace-pre-wrap border border-border rounded p-2 bg-gray-50"
        >{{ JSON.stringify(item.context, null, 2) }}</pre>
        <p v-else class="text-xs text-muted-foreground">{{ t('detail.context.empty') }}</p>
      </div>

      <template v-if="isResolved">
        <div class="text-xs border rounded px-3 py-2 border-border bg-gray-50 text-muted-foreground">
          {{ t('detail.alreadyResolved') }}
        </div>
      </template>
      <template v-else>
        <div class="space-y-2 pt-2 border-t border-border">
          <textarea
            v-model="feedback"
            :placeholder="t('detail.feedback.placeholder')"
            rows="3"
            class="w-full text-sm border border-border rounded px-2 py-1.5"
          />
          <div class="flex gap-2">
            <button
              @click="approve"
              :disabled="store.resolving"
              class="flex-1 text-sm px-3 py-1.5 rounded bg-gray-900 text-white disabled:opacity-40"
            >{{ t('detail.approve') }}</button>
            <button
              @click="reject"
              :disabled="store.resolving"
              class="flex-1 text-sm px-3 py-1.5 rounded border border-border text-red-600 hover:bg-red-50 disabled:opacity-40"
            >{{ t('detail.reject') }}</button>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>
