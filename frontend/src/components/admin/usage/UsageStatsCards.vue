<template>
  <div class="grid grid-cols-2 gap-px overflow-hidden rounded-lg border border-gray-200 bg-gray-200 dark:border-dark-700 dark:bg-dark-700 lg:grid-cols-4">
    <div class="min-w-0 bg-white px-3 py-3 dark:bg-dark-800 sm:px-5 sm:py-4">
      <div class="flex items-center gap-2">
        <span class="flex h-7 w-7 shrink-0 items-center justify-center rounded-md bg-blue-50 text-blue-600 dark:bg-blue-500/10 dark:text-blue-400">
          <Icon name="document" size="sm" />
        </span>
        <p class="truncate text-[11px] font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400 sm:text-xs">{{ t('usage.totalRequests') }}</p>
      </div>
      <p class="mt-2 font-mono text-lg font-bold tabular-nums text-gray-950 dark:text-white sm:text-2xl">{{ stats?.total_requests?.toLocaleString() || '0' }}</p>
      <p class="mt-1 hidden text-xs text-gray-400 sm:block">{{ t('usage.inSelectedRange') }}</p>
    </div>
    <div class="min-w-0 bg-white px-3 py-3 dark:bg-dark-800 sm:px-5 sm:py-4">
      <div class="flex items-center gap-2">
        <span class="flex h-7 w-7 shrink-0 items-center justify-center rounded-md bg-amber-50 text-amber-600 dark:bg-amber-500/10 dark:text-amber-400">
          <Icon name="cube" size="sm" />
        </span>
        <p class="truncate text-[11px] font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400 sm:text-xs">{{ t('usage.totalTokens') }}</p>
      </div>
      <p class="mt-2 font-mono text-lg font-bold tabular-nums text-gray-950 dark:text-white sm:text-2xl">{{ formatTokens(stats?.total_tokens || 0) }}</p>
      <p class="mt-1 hidden flex-wrap items-center gap-x-1 text-xs text-gray-500 sm:flex dark:text-gray-400">
          <span>{{ t('usage.in') }}: {{ formatTokens(stats?.total_input_tokens || 0) }}</span>
          <span>/</span>
          <span>{{ t('usage.out') }}: {{ formatTokens(stats?.total_output_tokens || 0) }}</span>
          <span>/</span>
          <span class="group relative inline-flex cursor-help items-center gap-0.5" tabindex="0">
            <span>{{ cacheLabel() }}: {{ formatTokens(stats?.total_cache_tokens || 0) }}</span>
            <Icon name="infoCircle" size="xs" class="text-gray-400" />
            <span
              class="pointer-events-none absolute left-1/2 top-full z-30 mt-2 w-56 -translate-x-1/2 rounded-lg border border-gray-200 bg-white p-3 text-left text-xs text-gray-700 opacity-0 shadow-lg transition-opacity duration-150 group-hover:opacity-100 group-focus:opacity-100 dark:border-dark-600 dark:bg-dark-800 dark:text-dark-200"
            >
              <span class="mb-2 block font-medium text-gray-900 dark:text-white">
                {{ cacheDetailLabel() }}
              </span>
              <span class="flex items-center justify-between gap-3">
                <span>{{ t('usage.cacheCreationTokensLabel') }}</span>
                <span class="tabular-nums">
                  {{ formatTokens(stats?.total_cache_creation_tokens || 0) }}
                </span>
              </span>
              <span class="mt-1 flex items-center justify-between gap-3">
                <span>{{ t('usage.cacheReadTokensLabel') }}</span>
                <span class="tabular-nums">
                  {{ formatTokens(stats?.total_cache_read_tokens || 0) }}
                </span>
              </span>
            </span>
          </span>
        </p>
    </div>
    <div class="min-w-0 bg-white px-3 py-3 dark:bg-dark-800 sm:px-5 sm:py-4">
      <div class="flex items-center gap-2">
        <span class="flex h-7 w-7 shrink-0 items-center justify-center rounded-md bg-emerald-50 text-emerald-600 dark:bg-emerald-500/10 dark:text-emerald-400">
          <Icon name="dollar" size="sm" />
        </span>
        <p class="truncate text-[11px] font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400 sm:text-xs">{{ t('usage.totalCost') }}</p>
      </div>
        <p class="mt-2 font-mono text-lg font-bold tabular-nums text-emerald-600 dark:text-emerald-400 sm:text-2xl">
          ${{ (stats?.total_actual_cost || 0).toFixed(4) }}
        </p>
        <p class="mt-1 hidden text-xs text-gray-400 sm:block">
          <template v-if="showAccountCost && totalAccountCost != null">
            <span class="text-orange-500">{{ t('usage.accountCost') }} ${{ totalAccountCost.toFixed(4) }}</span>
            <span> · </span>
          </template>
          <span>
            {{ t('usage.standardCost') }}
            <span :class="{ 'line-through': strikeStandardCost }">${{ (stats?.total_cost || 0).toFixed(4) }}</span>
          </span>
        </p>
    </div>
    <div class="min-w-0 bg-white px-3 py-3 dark:bg-dark-800 sm:px-5 sm:py-4">
      <div class="flex items-center gap-2">
        <span class="flex h-7 w-7 shrink-0 items-center justify-center rounded-md bg-violet-50 text-violet-600 dark:bg-violet-500/10 dark:text-violet-400">
          <Icon name="clock" size="sm" />
        </span>
        <p class="truncate text-[11px] font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400 sm:text-xs">{{ t('usage.avgDuration') }}</p>
      </div>
      <p class="mt-2 font-mono text-lg font-bold tabular-nums text-gray-950 dark:text-white sm:text-2xl">{{ formatDuration(stats?.average_duration_ms || 0) }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import type { AdminUsageStatsResponse } from '@/api/admin/usage'
import type { UsageStatsResponse } from '@/types'
import Icon from '@/components/icons/Icon.vue'

const props = withDefaults(defineProps<{
  stats: (AdminUsageStatsResponse | UsageStatsResponse) | null
  showAccountCost?: boolean
  strikeStandardCost?: boolean
}>(), {
  showAccountCost: true,
  strikeStandardCost: false,
})

const { t } = useI18n()

const totalAccountCost = computed(() => {
  const stats = props.stats as (AdminUsageStatsResponse & { total_account_cost?: number }) | null
  return stats?.total_account_cost ?? null
})
const showAccountCost = computed(() => props.showAccountCost)
const strikeStandardCost = computed(() => props.strikeStandardCost)

const formatDuration = (ms: number) =>
  ms < 1000 ? `${ms.toFixed(0)}ms` : `${(ms / 1000).toFixed(2)}s`

const formatTokens = (value: number) => {
  if (value >= 1e9) return (value / 1e9).toFixed(2) + 'B'
  if (value >= 1e6) return (value / 1e6).toFixed(2) + 'M'
  if (value >= 1e3) return (value / 1e3).toFixed(2) + 'K'
  return value.toLocaleString()
}

const cacheLabel = () => t('usage.cacheTotal')
const cacheDetailLabel = () => t('usage.cacheBreakdown')
</script>
