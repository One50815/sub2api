<template>
  <span :class="['inline-flex items-center rounded-full border px-2.5 py-1 text-xs font-semibold', colorClass]">
    {{ label }}
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps<{ kind: 'status' | 'category' | 'priority'; value: string }>()
const { t } = useI18n()

const label = computed(() => t(`tickets.${props.kind}.${props.value}`))
const colorClass = computed(() => {
  const map: Record<string, string> = {
    pending_admin: 'border-amber-200 bg-amber-50 text-amber-700 dark:border-amber-800 dark:bg-amber-950/40 dark:text-amber-300',
    pending_user: 'border-blue-200 bg-blue-50 text-blue-700 dark:border-blue-800 dark:bg-blue-950/40 dark:text-blue-300',
    resolved: 'border-emerald-200 bg-emerald-50 text-emerald-700 dark:border-emerald-800 dark:bg-emerald-950/40 dark:text-emerald-300',
    closed: 'border-gray-200 bg-gray-50 text-gray-600 dark:border-dark-600 dark:bg-dark-800 dark:text-gray-300',
    urgent: 'border-red-200 bg-red-50 text-red-700 dark:border-red-800 dark:bg-red-950/40 dark:text-red-300',
    high: 'border-orange-200 bg-orange-50 text-orange-700 dark:border-orange-800 dark:bg-orange-950/40 dark:text-orange-300',
    normal: 'border-gray-200 bg-white text-gray-600 dark:border-dark-600 dark:bg-dark-800 dark:text-gray-300',
    low: 'border-slate-200 bg-slate-50 text-slate-600 dark:border-slate-700 dark:bg-slate-900 dark:text-slate-300',
  }
  return map[props.value] || 'border-primary-200 bg-primary-50 text-primary-700 dark:border-primary-800 dark:bg-primary-950/30 dark:text-primary-300'
})
</script>

