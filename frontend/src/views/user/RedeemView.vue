<template>
  <AppLayout>
    <div class="mx-auto flex w-full max-w-7xl flex-col gap-4 sm:gap-5">
      <header class="page-header !mb-0">
        <h1 class="page-title">{{ t('redeem.title') }}</h1>
        <p class="page-description mt-1 text-sm">
          {{ t('redeem.description') }}
        </p>
      </header>

      <!-- Account snapshot -->
      <section
        class="grid grid-cols-2 divide-x divide-gray-200 overflow-hidden rounded-lg border border-gray-200 bg-white dark:divide-dark-700 dark:border-dark-700 dark:bg-dark-800"
        :aria-label="t('redeem.currentBalance')"
      >
        <div class="min-w-0 px-2.5 py-2.5 sm:px-5 sm:py-4">
          <div class="flex items-center gap-1.5 sm:gap-2.5">
            <span
              class="flex h-7 w-7 shrink-0 items-center justify-center rounded-md bg-emerald-50 text-emerald-600 dark:bg-emerald-500/10 dark:text-emerald-400"
            >
              <Icon name="dollar" size="sm" />
            </span>
            <span
              class="truncate text-[11px] font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400 sm:text-xs"
            >
              {{ t('redeem.currentBalance') }}
            </span>
          </div>
          <p
            class="mt-1.5 break-all font-mono text-sm font-bold tabular-nums text-gray-950 dark:text-white sm:mt-2.5 sm:text-[1.75rem]"
          >
            ${{ user?.balance?.toFixed(2) || '0.00' }}
          </p>
        </div>

        <div class="min-w-0 px-2.5 py-2.5 sm:px-5 sm:py-4">
          <div class="flex items-center gap-1.5 sm:gap-2.5">
            <span
              class="flex h-7 w-7 shrink-0 items-center justify-center rounded-md bg-blue-50 text-blue-600 dark:bg-blue-500/10 dark:text-blue-400"
            >
              <Icon name="bolt" size="sm" />
            </span>
            <span
              class="truncate text-[11px] font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400 sm:text-xs"
            >
              {{ t('redeem.concurrency') }}
            </span>
          </div>
          <p
            class="mt-1.5 font-mono text-sm font-bold tabular-nums text-gray-950 dark:text-white sm:mt-2.5 sm:text-[1.75rem]"
          >
            {{ user?.concurrency || 0 }}
          </p>
          <p class="mt-1 hidden text-xs text-gray-400 dark:text-gray-500 md:block">
            {{ t('redeem.concurrentRequests') }}
          </p>
        </div>
      </section>

      <div
        class="grid gap-4 xl:grid-cols-[minmax(0,1.05fr)_minmax(320px,0.75fr)] xl:items-start"
      >
        <!-- Redeem form and result -->
        <section
          class="overflow-hidden rounded-lg border border-gray-200 bg-white dark:border-dark-700 dark:bg-dark-800"
        >
          <div class="flex items-center gap-2.5 border-b border-gray-200 px-4 py-3.5 dark:border-dark-700 sm:px-5">
            <span
              class="flex h-8 w-8 shrink-0 items-center justify-center rounded-md bg-amber-50 text-amber-600 dark:bg-amber-500/10 dark:text-amber-400"
            >
              <Icon name="gift" size="sm" />
            </span>
            <div class="min-w-0">
              <h2 class="text-sm font-semibold text-gray-950 dark:text-white">
                {{ t('redeem.redeemCodeLabel') }}
              </h2>
              <p class="truncate text-xs text-gray-500 dark:text-gray-400">
                {{ t('redeem.redeemCodeHint') }}
              </p>
            </div>
          </div>

          <div class="p-4 sm:p-5">
            <form class="space-y-2.5" @submit.prevent="handleRedeem">
              <label
                for="code"
                class="block text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400"
              >
                {{ t('redeem.redeemCodeLabel') }}
              </label>
              <div class="grid gap-2 sm:grid-cols-[minmax(0,1fr)_auto]">
                <input
                  id="code"
                  v-model="redeemCode"
                  type="text"
                  required
                  :placeholder="t('redeem.redeemCodePlaceholder')"
                  :disabled="submitting"
                  class="input h-10 min-w-0"
                />
                <button
                  type="submit"
                  :disabled="!redeemCode || submitting"
                  class="btn btn-secondary h-10 whitespace-nowrap px-5"
                >
                  <span
                    v-if="submitting"
                    class="mr-2 h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent"
                    aria-hidden="true"
                  ></span>
                  <Icon v-else name="checkCircle" size="sm" class="mr-2" />
                  {{ submitting ? t('redeem.redeeming') : t('redeem.redeemButton') }}
                </button>
              </div>
              <p class="text-xs text-gray-500 dark:text-gray-400">
                {{ t('redeem.redeemCodeHint') }}
              </p>
            </form>
          </div>

          <!-- Success Message -->
          <transition name="fade">
            <div
              v-if="redeemResult"
              class="border-t border-gray-200 px-4 py-4 dark:border-dark-700 sm:px-5"
            >
              <div
                class="rounded-md border border-emerald-200 bg-emerald-50/40 p-3 dark:border-emerald-800/60 dark:bg-emerald-900/10 sm:p-4"
              >
                <div class="flex items-start gap-3">
                  <span
                    class="flex h-8 w-8 shrink-0 items-center justify-center rounded-md bg-emerald-100 text-emerald-600 dark:bg-emerald-500/10 dark:text-emerald-400"
                  >
                    <Icon name="checkCircle" size="sm" />
                  </span>
                  <div class="min-w-0 flex-1">
                    <h3 class="text-sm font-semibold text-emerald-800 dark:text-emerald-300">
                      {{ t('redeem.redeemSuccess') }}
                    </h3>
                    <p class="mt-1 text-xs leading-5 text-gray-600 dark:text-gray-300">
                      {{ redeemResult.message }}
                    </p>
                    <div
                      class="mt-3 grid gap-2 border-t border-emerald-200/70 pt-3 text-xs dark:border-emerald-800/50 sm:grid-cols-2"
                    >
                      <p v-if="redeemResult.type === 'balance'" class="font-medium text-gray-800 dark:text-gray-200">
                        {{ t('redeem.added') }}: ${{ redeemResult.value.toFixed(2) }}
                      </p>
                      <p
                        v-else-if="redeemResult.type === 'concurrency'"
                        class="font-medium text-gray-800 dark:text-gray-200"
                      >
                        {{ t('redeem.added') }}: {{ redeemResult.value }}
                        {{ t('redeem.concurrentRequests') }}
                      </p>
                      <p
                        v-else-if="redeemResult.type === 'subscription'"
                        class="font-medium text-gray-800 dark:text-gray-200"
                      >
                        {{ t('redeem.subscriptionAssigned') }}
                        <span v-if="redeemResult.group_name"> · {{ redeemResult.group_name }}</span>
                        <span v-if="redeemResult.validity_days">
                          ({{ t('redeem.subscriptionDays', { days: redeemResult.validity_days }) }})
                        </span>
                      </p>
                      <p v-if="redeemResult.new_balance !== undefined" class="text-gray-600 dark:text-gray-300">
                        {{ t('redeem.newBalance') }}:
                        <span class="font-mono font-semibold tabular-nums text-gray-900 dark:text-white"
                          >${{ redeemResult.new_balance.toFixed(2) }}</span
                        >
                      </p>
                      <p v-if="redeemResult.new_concurrency !== undefined" class="text-gray-600 dark:text-gray-300">
                        {{ t('redeem.newConcurrency') }}:
                        <span class="font-mono font-semibold tabular-nums text-gray-900 dark:text-white">
                          {{ redeemResult.new_concurrency }} {{ t('redeem.requests') }}
                        </span>
                      </p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </transition>

          <!-- Error Message -->
          <transition name="fade">
            <div
              v-if="errorMessage"
              class="border-t border-gray-200 px-4 py-4 dark:border-dark-700 sm:px-5"
            >
              <div
                class="rounded-md border border-red-200 bg-red-50/40 p-3 dark:border-red-800/60 dark:bg-red-900/10 sm:p-4"
              >
                <div class="flex items-start gap-3">
                  <span
                    class="flex h-8 w-8 shrink-0 items-center justify-center rounded-md bg-red-100 text-red-600 dark:bg-red-500/10 dark:text-red-400"
                  >
                    <Icon name="exclamationCircle" size="sm" />
                  </span>
                  <div class="min-w-0 flex-1">
                    <h3 class="text-sm font-semibold text-red-800 dark:text-red-300">
                      {{ t('redeem.redeemFailed') }}
                    </h3>
                    <p class="mt-1 text-xs leading-5 text-gray-600 dark:text-gray-300">
                      {{ errorMessage }}
                    </p>
                  </div>
                </div>
              </div>
            </div>
          </transition>
        </section>

        <!-- Redeem code guidance -->
        <aside
          class="overflow-hidden rounded-lg border border-gray-200 bg-gray-50/60 dark:border-dark-700 dark:bg-dark-900/30 xl:sticky xl:top-4"
        >
          <div class="flex items-center gap-2.5 border-b border-gray-200 px-4 py-3.5 dark:border-dark-700 sm:px-5">
            <span
              class="flex h-8 w-8 shrink-0 items-center justify-center rounded-md bg-blue-50 text-blue-600 dark:bg-blue-500/10 dark:text-blue-400"
            >
              <Icon name="infoCircle" size="sm" />
            </span>
            <h2 class="text-sm font-semibold text-gray-950 dark:text-white">
              {{ t('redeem.aboutCodes') }}
            </h2>
          </div>
          <ul class="space-y-3 p-4 sm:p-5">
            <li class="flex items-start gap-2.5 text-xs leading-5 text-gray-600 dark:text-gray-300">
              <Icon name="checkCircle" size="xs" class="mt-0.5 shrink-0 text-emerald-500" />
              <span>{{ t('redeem.codeRule1') }}</span>
            </li>
            <li class="flex items-start gap-2.5 text-xs leading-5 text-gray-600 dark:text-gray-300">
              <Icon name="checkCircle" size="xs" class="mt-0.5 shrink-0 text-emerald-500" />
              <span>{{ t('redeem.codeRule2') }}</span>
            </li>
            <li class="flex items-start gap-2.5 text-xs leading-5 text-gray-600 dark:text-gray-300">
              <Icon name="checkCircle" size="xs" class="mt-0.5 shrink-0 text-emerald-500" />
              <span>
                {{ t('redeem.codeRule3') }}
                <span
                  v-if="contactInfo"
                  class="ml-1 inline-flex items-center rounded-full border border-gray-200 bg-white px-2 py-0.5 text-[11px] font-medium text-gray-700 dark:border-dark-600 dark:bg-dark-800 dark:text-gray-200"
                >
                  {{ contactInfo }}
                </span>
              </span>
            </li>
            <li class="flex items-start gap-2.5 text-xs leading-5 text-gray-600 dark:text-gray-300">
              <Icon name="checkCircle" size="xs" class="mt-0.5 shrink-0 text-emerald-500" />
              <span>{{ t('redeem.codeRule4') }}</span>
            </li>
          </ul>
        </aside>
      </div>

      <!-- Recent Activity -->
      <section
        class="overflow-hidden rounded-lg border border-gray-200 bg-white dark:border-dark-700 dark:bg-dark-800"
      >
        <div class="flex items-center gap-2.5 border-b border-gray-200 px-4 py-3.5 dark:border-dark-700 sm:px-5">
          <span
            class="flex h-8 w-8 shrink-0 items-center justify-center rounded-md bg-gray-100 text-gray-500 dark:bg-dark-700 dark:text-gray-300"
          >
            <Icon name="clock" size="sm" />
          </span>
          <h2 class="text-sm font-semibold text-gray-950 dark:text-white">
            {{ t('redeem.recentActivity') }}
          </h2>
        </div>

        <!-- Loading State -->
        <div v-if="loadingHistory" class="space-y-2 p-3 sm:p-4">
          <div
            v-for="index in 3"
            :key="index"
            class="animate-pulse rounded-md border border-gray-200 p-3 dark:border-dark-700 sm:p-4"
          >
            <div class="flex items-center justify-between gap-4">
              <div class="flex min-w-0 flex-1 items-center gap-3">
                <div class="h-8 w-8 shrink-0 rounded-md bg-gray-100 dark:bg-dark-700"></div>
                <div class="min-w-0 flex-1">
                  <div class="h-3 w-32 rounded bg-gray-100 dark:bg-dark-700"></div>
                  <div class="mt-2 h-2.5 w-24 rounded bg-gray-100 dark:bg-dark-700"></div>
                </div>
              </div>
              <div class="h-3 w-16 rounded bg-gray-100 dark:bg-dark-700"></div>
            </div>
          </div>
        </div>

        <!-- History List -->
        <div v-else-if="history.length > 0" class="space-y-2 p-3 sm:p-4">
          <article
            v-for="item in history"
            :key="item.id"
            class="rounded-md border border-gray-200 p-3 dark:border-dark-700 sm:p-4"
          >
            <div class="flex items-start justify-between gap-3 sm:items-center">
              <div class="flex min-w-0 items-start gap-3 sm:items-center">
                <span
                  :class="[
                    'flex h-8 w-8 shrink-0 items-center justify-center rounded-md',
                    isBalanceType(item.type)
                      ? item.value >= 0
                        ? 'bg-emerald-50 text-emerald-600 dark:bg-emerald-500/10 dark:text-emerald-400'
                        : 'bg-red-50 text-red-600 dark:bg-red-500/10 dark:text-red-400'
                      : isSubscriptionType(item.type)
                        ? 'bg-purple-50 text-purple-600 dark:bg-purple-500/10 dark:text-purple-400'
                        : item.value >= 0
                          ? 'bg-blue-50 text-blue-600 dark:bg-blue-500/10 dark:text-blue-400'
                          : 'bg-orange-50 text-orange-600 dark:bg-orange-500/10 dark:text-orange-400'
                  ]"
                >
                  <Icon v-if="isBalanceType(item.type)" name="dollar" size="sm" />
                  <Icon v-else-if="isSubscriptionType(item.type)" name="badge" size="sm" />
                  <Icon v-else name="bolt" size="sm" />
                </span>
                <div class="min-w-0">
                  <p class="truncate text-sm font-medium text-gray-900 dark:text-white">
                    {{ getHistoryItemTitle(item) }}
                  </p>
                  <p class="mt-0.5 text-xs text-gray-500 dark:text-dark-400">
                    {{ formatDateTime(item.used_at) }}
                  </p>
                </div>
              </div>

              <div class="min-w-0 shrink-0 text-right">
                <p
                  :class="[
                    'text-sm font-semibold tabular-nums',
                    isBalanceType(item.type)
                      ? item.value >= 0
                        ? 'text-emerald-600 dark:text-emerald-400'
                        : 'text-red-600 dark:text-red-400'
                      : isSubscriptionType(item.type)
                        ? 'text-purple-600 dark:text-purple-400'
                        : item.value >= 0
                          ? 'text-blue-600 dark:text-blue-400'
                          : 'text-orange-600 dark:text-orange-400'
                  ]"
                >
                  {{ formatHistoryValue(item) }}
                </p>
                <p
                  v-if="!isAdminAdjustment(item.type)"
                  class="mt-0.5 font-mono text-[11px] text-gray-400 dark:text-dark-500"
                >
                  {{ item.code.slice(0, 8) }}...
                </p>
                <p v-else class="mt-0.5 text-[11px] text-gray-400 dark:text-dark-500">
                  {{ t('redeem.adminAdjustment') }}
                </p>
              </div>
            </div>
            <p
              v-if="item.notes"
              class="mt-2 truncate border-t border-gray-100 pt-2 text-xs italic text-gray-500 dark:border-dark-700 dark:text-dark-400"
              :title="item.notes"
            >
              {{ item.notes }}
            </p>
          </article>
        </div>

        <!-- Empty State -->
        <div v-else class="flex min-h-40 flex-col items-center justify-center px-5 py-10 text-center">
          <span
            class="mb-3 flex h-8 w-8 items-center justify-center rounded-md bg-gray-100 text-gray-400 dark:bg-dark-700 dark:text-dark-400"
          >
            <Icon name="clock" size="sm" />
          </span>
          <p class="text-sm font-medium text-gray-800 dark:text-gray-200">
            {{ t('redeem.recentActivity') }}
          </p>
          <p class="mt-1 text-xs text-gray-500 dark:text-dark-400">
            {{ t('redeem.historyWillAppear') }}
          </p>
        </div>
      </section>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { useAppStore } from '@/stores/app'
import { useSubscriptionStore } from '@/stores/subscriptions'
import { redeemAPI, authAPI, type RedeemHistoryItem } from '@/api'
import AppLayout from '@/components/layout/AppLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import { formatDateTime } from '@/utils/format'

const { t } = useI18n()
const authStore = useAuthStore()
const appStore = useAppStore()
const subscriptionStore = useSubscriptionStore()

const user = computed(() => authStore.user)

const redeemCode = ref('')
const submitting = ref(false)
const redeemResult = ref<{
  message: string
  type: string
  value: number
  new_balance?: number
  new_concurrency?: number
  group_name?: string
  validity_days?: number
} | null>(null)
const errorMessage = ref('')

// History data
const history = ref<RedeemHistoryItem[]>([])
const loadingHistory = ref(false)
const contactInfo = ref('')

// Helper functions for history display
const isBalanceType = (type: string) => {
  return type === 'balance' || type === 'admin_balance'
}

const isSubscriptionType = (type: string) => {
  return type === 'subscription'
}

const isAdminAdjustment = (type: string) => {
  return type === 'admin_balance' || type === 'admin_concurrency'
}

const getHistoryItemTitle = (item: RedeemHistoryItem) => {
  if (item.type === 'balance') {
    return t('redeem.balanceAddedRedeem')
  } else if (item.type === 'admin_balance') {
    return item.value >= 0 ? t('redeem.balanceAddedAdmin') : t('redeem.balanceDeductedAdmin')
  } else if (item.type === 'concurrency') {
    return t('redeem.concurrencyAddedRedeem')
  } else if (item.type === 'admin_concurrency') {
    return item.value >= 0 ? t('redeem.concurrencyAddedAdmin') : t('redeem.concurrencyReducedAdmin')
  } else if (item.type === 'subscription') {
    return t('redeem.subscriptionAssigned')
  }
  return t('common.unknown')
}

const formatHistoryValue = (item: RedeemHistoryItem) => {
  if (isBalanceType(item.type)) {
    const sign = item.value >= 0 ? '+' : ''
    return `${sign}$${item.value.toFixed(2)}`
  } else if (isSubscriptionType(item.type)) {
    // 订阅类型显示有效天数和分组名称
    const days = item.validity_days || Math.round(item.value)
    const groupName = item.group?.name || ''
    return groupName ? `${days}${t('redeem.days')} - ${groupName}` : `${days}${t('redeem.days')}`
  } else {
    const sign = item.value >= 0 ? '+' : ''
    return `${sign}${item.value} ${t('redeem.requests')}`
  }
}

const fetchHistory = async () => {
  loadingHistory.value = true
  try {
    history.value = await redeemAPI.getHistory()
  } catch (error) {
    console.error('Failed to fetch history:', error)
  } finally {
    loadingHistory.value = false
  }
}

const handleRedeem = async () => {
  if (!redeemCode.value.trim()) {
    appStore.showError(t('redeem.pleaseEnterCode'))
    return
  }

  submitting.value = true
  errorMessage.value = ''
  redeemResult.value = null

  try {
    const result = await redeemAPI.redeem(redeemCode.value.trim())

    redeemResult.value = result

    // Refresh user data to get updated balance/concurrency
    await authStore.refreshUser()

    // If subscription type, immediately refresh subscription status
    if (result.type === 'subscription') {
      try {
        await subscriptionStore.fetchActiveSubscriptions(true) // force refresh
      } catch (error) {
        console.error('Failed to refresh subscriptions after redeem:', error)
        appStore.showWarning(t('redeem.subscriptionRefreshFailed'))
      }
    }

    // Clear the input
    redeemCode.value = ''

    // Refresh history
    await fetchHistory()

    // Show success toast
    appStore.showSuccess(t('redeem.codeRedeemSuccess'))
  } catch (error: any) {
    errorMessage.value = error.response?.data?.detail || t('redeem.failedToRedeem')

    appStore.showError(t('redeem.redeemFailed'))
  } finally {
    submitting.value = false
  }
}

onMounted(async () => {
  fetchHistory()
  try {
    const settings = await authAPI.getPublicSettings()
    contactInfo.value = settings.contact_info || ''
  } catch (error) {
    console.error('Failed to load contact info:', error)
  }
})
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: all 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}
</style>
