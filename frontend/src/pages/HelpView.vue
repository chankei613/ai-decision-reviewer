<script setup lang="ts">
import { useI18n } from '@/i18n'
import { levelColor, levelIcon } from '@/levelColors'

const { t } = useI18n()

const levels = [
  { key: 'interrupt', textKey: 'help.levels.interrupt' },
  { key: 'urgent', textKey: 'help.levels.urgent' },
  { key: 'emergency_stop', textKey: 'help.levels.emergency' },
]
</script>

<template>
  <div class="p-6 space-y-6 max-w-xl overflow-y-auto h-full">
    <div>
      <h2 class="text-sm font-semibold">{{ t('help.title') }}</h2>
      <p class="text-xs text-muted-foreground mt-1">{{ t('help.intro') }}</p>
    </div>

    <section class="space-y-2">
      <h3 class="text-xs font-semibold text-muted-foreground">{{ t('help.what.title') }}</h3>
      <p class="text-sm">{{ t('help.what.body') }}</p>
    </section>

    <section class="space-y-2">
      <h3 class="text-xs font-semibold text-muted-foreground">{{ t('help.levels.title') }}</h3>
      <ul class="space-y-2">
        <li v-for="l in levels" :key="l.key" class="border border-border rounded px-3 py-2 text-sm flex gap-2">
          <span
            class="shrink-0 inline-flex items-center gap-1 h-fit px-1.5 py-0.5 rounded-full text-white text-[11px]"
            :style="{ background: levelColor(l.key) }"
          >
            <span aria-hidden="true">{{ levelIcon(l.key) }}</span>
            {{ t('level.' + l.key) }}
          </span>
          <span>{{ t(l.textKey) }}</span>
        </li>
      </ul>
    </section>

    <section class="space-y-2">
      <h3 class="text-xs font-semibold text-muted-foreground">{{ t('help.flow.title') }}</h3>
      <ol class="space-y-1.5 text-sm list-decimal list-inside">
        <li>{{ t('help.flow.1') }}</li>
        <li>{{ t('help.flow.2') }}</li>
        <li>{{ t('help.flow.3') }}</li>
        <li>{{ t('help.flow.4') }}</li>
      </ol>
    </section>

    <section class="space-y-2">
      <h3 class="text-xs font-semibold text-muted-foreground">{{ t('help.stuck.title') }}</h3>
      <ul class="space-y-1.5 text-sm list-disc list-inside text-muted-foreground">
        <li><span class="text-foreground">{{ t('help.stuck.1') }}</span></li>
        <li><span class="text-foreground">{{ t('help.stuck.2') }}</span></li>
        <li><span class="text-foreground">{{ t('help.stuck.3') }}</span></li>
      </ul>
    </section>
  </div>
</template>
