<template>
  <div class="grid grid-cols-2 gap-3 lg:grid-cols-4">
    <div v-for="item in cards" :key="item.key" class="card px-4 py-4">
      <p class="text-xs font-medium text-gray-500 dark:text-gray-400">{{ item.label }}</p>
      <div class="mt-2 flex items-end justify-between gap-3">
        <strong class="text-2xl font-semibold tracking-tight text-gray-950 dark:text-white">{{ item.value }}</strong>
        <span :class="['h-2.5 w-2.5 rounded-full', item.dot]"></span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import type { TicketSummary } from '@/types/ticket'

const props = defineProps<{ summary: TicketSummary }>()
const { t } = useI18n()
const cards = computed(() => [
  { key: 'total', label: t('tickets.summary.total'), value: props.summary.total, dot: 'bg-gray-400' },
  { key: 'pending_admin', label: t('tickets.summary.pendingAdmin'), value: props.summary.pending_admin, dot: 'bg-amber-400' },
  { key: 'pending_user', label: t('tickets.summary.pendingUser'), value: props.summary.pending_user, dot: 'bg-blue-400' },
  { key: 'unread', label: t('tickets.summary.unread'), value: props.summary.unread_count, dot: 'bg-primary-500' },
])
</script>

