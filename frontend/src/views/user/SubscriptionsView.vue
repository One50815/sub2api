<template>
  <AppLayout>
    <div class="mx-auto flex w-full max-w-7xl flex-col gap-4 sm:gap-5">
      <header class="page-header !mb-0">
        <h1 class="page-title">{{ t('userSubscriptions.title') }}</h1>
        <p class="page-description mt-1 text-sm">
          {{ t('userSubscriptions.description') }}
        </p>
      </header>

      <!-- Loading State -->
      <div
        v-if="loading"
        class="flex min-h-40 items-center justify-center rounded-lg border border-gray-200 bg-white dark:border-dark-700 dark:bg-dark-800"
      >
        <div
          class="h-8 w-8 animate-spin rounded-full border-2 border-primary-500 border-t-transparent"
        ></div>
      </div>

      <!-- Empty State -->
      <div
        v-else-if="subscriptions.length === 0"
        class="flex min-h-64 flex-col items-center justify-center rounded-lg border border-gray-200 bg-white px-5 py-10 text-center dark:border-dark-700 dark:bg-dark-800"
      >
        <div
          class="mb-3 flex h-9 w-9 items-center justify-center rounded-md bg-gray-100 text-gray-400 dark:bg-dark-700 dark:text-dark-400"
        >
          <Icon name="creditCard" size="sm" />
        </div>
        <h2 class="text-sm font-semibold text-gray-900 dark:text-white">
          {{ t('userSubscriptions.noActiveSubscriptions') }}
        </h2>
        <p class="mt-1 max-w-md text-xs leading-5 text-gray-500 dark:text-dark-400">
          {{ t('userSubscriptions.noActiveSubscriptionsDesc') }}
        </p>
      </div>

      <!-- Subscriptions List -->
      <section
        v-else
        class="overflow-hidden rounded-lg border border-gray-200 bg-white dark:border-dark-700 dark:bg-dark-800"
      >
        <div
          class="flex items-center justify-between gap-3 border-b border-gray-200 px-4 py-3.5 dark:border-dark-700 sm:px-5"
        >
          <div class="flex min-w-0 items-center gap-2.5">
            <span
              class="flex h-8 w-8 shrink-0 items-center justify-center rounded-md bg-amber-50 text-amber-600 dark:bg-amber-500/10 dark:text-amber-400"
            >
              <Icon name="creditCard" size="sm" />
            </span>
            <div class="min-w-0">
              <h2 class="text-sm font-semibold text-gray-950 dark:text-white">
                {{ t('userSubscriptions.usage') }}
              </h2>
              <p class="truncate text-xs text-gray-500 dark:text-gray-400">
                {{ t('userSubscriptions.description') }}
              </p>
            </div>
          </div>
          <span
            class="inline-flex min-w-7 items-center justify-center rounded-full border border-gray-200 bg-gray-50 px-2 py-1 text-xs font-semibold tabular-nums text-gray-600 dark:border-dark-600 dark:bg-dark-700 dark:text-gray-300"
          >
            {{ subscriptions.length }}
          </span>
        </div>

        <div class="space-y-3 p-3 sm:p-4">
          <article
            v-for="subscription in subscriptions"
            :key="subscription.id"
            class="rounded-md border border-gray-200 p-3 dark:border-dark-700 sm:p-4"
          >
          <!-- Header -->
          <div
            class="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between"
          >
            <div class="flex min-w-0 items-start gap-3">
              <span
                :class="[
                  'mt-2 h-1.5 w-1.5 shrink-0 rounded-full',
                  platformAccentDotClass(subscription.group?.platform || '')
                ]"
                aria-hidden="true"
              ></span>
              <div class="min-w-0">
                <div class="flex items-center gap-2">
                  <h3 class="truncate text-sm font-semibold text-gray-900 dark:text-white">
                    {{ subscription.group?.name || `Group #${subscription.group_id}` }}
                  </h3>
                  <span
                    :class="[
                      'shrink-0 rounded-md border px-2 py-0.5 text-[11px] font-medium',
                      platformBadgeClass(subscription.group?.platform || '')
                    ]"
                  >
                    {{ platformLabel(subscription.group?.platform || '') }}
                  </span>
                </div>
                <p
                  v-if="subscription.group?.description"
                  class="mt-1 text-xs leading-5 text-gray-500 dark:text-dark-400"
                >
                  {{ subscription.group.description }}
                </p>
                <div
                  class="mt-1.5 flex flex-wrap gap-x-3 gap-y-1 text-[11px] text-gray-400 dark:text-gray-500"
                >
                  <span
                    >{{ t('payment.planCard.rate') }}: ×{{
                      subscription.entitlement?.rate_multiplier ?? subscription.group?.rate_multiplier ?? 1
                    }}</span
                  >
                  <span
                    v-if="subscriptionHasPeakRate(subscription)"
                    class="text-amber-700 dark:text-amber-300"
                  >
                    {{ t('payment.planCard.peakRate') }}: {{ subscriptionPeakRateLabel(subscription) }}
                  </span>
                </div>
              </div>
            </div>
            <div class="flex shrink-0 items-center gap-2 pl-4 sm:pl-0">
              <span
                :class="[
                  'rounded-full border px-2 py-0.5 text-[11px] font-medium',
                  subscription.status === 'active'
                    ? 'border-emerald-200 bg-emerald-50 text-emerald-700 dark:border-emerald-800/70 dark:bg-emerald-900/20 dark:text-emerald-300'
                    : subscription.status === 'expired'
                      ? 'border-gray-200 bg-gray-50 text-gray-600 dark:border-dark-600 dark:bg-dark-700 dark:text-gray-400'
                      : 'border-red-200 bg-red-50 text-red-700 dark:border-red-800/70 dark:bg-red-900/20 dark:text-red-300'
                ]"
              >
                {{ t(`userSubscriptions.status.${subscription.status}`) }}
              </span>
              <button
                v-if="subscription.status === 'active'"
                class="btn btn-secondary btn-sm"
                @click="router.push({ path: '/purchase', query: { tab: 'subscription', group: String(subscription.group_id) } })"
              >
                {{ t('payment.renewNow') }}
              </button>
            </div>
          </div>

          <!-- Usage Progress -->
          <div
            class="mt-4 grid gap-3 border-t border-gray-200 pt-4 dark:border-dark-700 md:grid-cols-3"
          >
            <!-- Expiration Info -->
            <div
              v-if="subscription.expires_at"
              class="flex items-center justify-between gap-4 rounded-md bg-gray-50 px-3 py-2 text-xs dark:bg-dark-900/40 md:col-span-3"
            >
              <span class="text-gray-500 dark:text-dark-400">{{
                t('userSubscriptions.expires')
              }}</span>
              <span :class="getExpirationClass(subscription.expires_at)">
                {{ formatExpirationDate(subscription.expires_at) }}
              </span>
            </div>
            <div
              v-else
              class="flex items-center justify-between gap-4 rounded-md bg-gray-50 px-3 py-2 text-xs dark:bg-dark-900/40 md:col-span-3"
            >
              <span class="text-gray-500 dark:text-dark-400">{{
                t('userSubscriptions.expires')
              }}</span>
              <span class="text-gray-700 dark:text-gray-300">{{
                t('userSubscriptions.noExpiration')
              }}</span>
            </div>

            <div
              v-if="subscription.entitlement"
              class="flex flex-col gap-3 rounded-md border border-gray-200 px-3 py-3 dark:border-dark-700 md:col-span-3 sm:flex-row sm:items-center sm:justify-between"
            >
              <div>
                <p class="text-xs font-medium text-gray-800 dark:text-gray-200">{{ subscription.entitlement.plan_name }}</p>
                <p class="mt-0.5 text-xs text-gray-500 dark:text-gray-400">额度用尽后</p>
              </div>
              <div class="inline-flex rounded-md border border-gray-200 bg-gray-50 p-1 dark:border-dark-600 dark:bg-dark-900">
                <button
                  class="rounded px-3 py-1.5 text-xs font-medium transition-colors"
                  :class="subscription.entitlement.overage_policy === 'block' ? 'bg-white text-gray-950 shadow-sm dark:bg-dark-700 dark:text-white' : 'text-gray-500 dark:text-gray-400'"
                  :disabled="updatingEntitlementId === subscription.entitlement.id"
                  @click="setOveragePolicy(subscription, 'block')"
                >阻止调用</button>
                <button
                  class="rounded px-3 py-1.5 text-xs font-medium transition-colors"
                  :class="subscription.entitlement.overage_policy === 'wallet' ? 'bg-white text-gray-950 shadow-sm dark:bg-dark-700 dark:text-white' : 'text-gray-500 dark:text-gray-400'"
                  :disabled="updatingEntitlementId === subscription.entitlement.id"
                  @click="setOveragePolicy(subscription, 'wallet')"
                >钱包续用</button>
              </div>
            </div>

            <!-- Daily Usage -->
            <div
              v-if="subscriptionDailyLimit(subscription)"
              class="space-y-2 rounded-md border border-gray-200 p-3 dark:border-dark-700"
            >
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-700 dark:text-gray-300">
                  {{ t('userSubscriptions.daily') }}
                </span>
                <span class="text-sm text-gray-500 dark:text-dark-400">
                  ${{ (subscription.daily_usage_usd || 0).toFixed(2) }} / ${{
                    subscriptionDailyLimit(subscription)?.toFixed(2)
                  }}
                </span>
              </div>
              <div class="relative h-1.5 overflow-hidden rounded-full bg-gray-100 dark:bg-dark-600">
                <div
                  class="absolute inset-y-0 left-0 rounded-full transition-all duration-300"
                  :class="
                    getProgressBarClass(
                      subscription.daily_usage_usd,
                      subscriptionDailyLimit(subscription)
                    )
                  "
                  :style="{
                    width: getProgressWidth(
                      subscription.daily_usage_usd,
                      subscriptionDailyLimit(subscription)
                    )
                  }"
                ></div>
              </div>
              <p
                v-if="subscription.daily_window_start"
                class="text-xs text-gray-500 dark:text-dark-400"
              >
                {{ formatDailyUsageWindow(subscription) }}
              </p>
            </div>

            <!-- Weekly Usage -->
            <div
              v-if="subscriptionWeeklyLimit(subscription)"
              class="space-y-2 rounded-md border border-gray-200 p-3 dark:border-dark-700"
            >
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-700 dark:text-gray-300">
                  {{ t('userSubscriptions.weekly') }}
                </span>
                <span class="text-sm text-gray-500 dark:text-dark-400">
                  ${{ (subscription.weekly_usage_usd || 0).toFixed(2) }} / ${{
                    subscriptionWeeklyLimit(subscription)?.toFixed(2)
                  }}
                </span>
              </div>
              <div class="relative h-1.5 overflow-hidden rounded-full bg-gray-100 dark:bg-dark-600">
                <div
                  class="absolute inset-y-0 left-0 rounded-full transition-all duration-300"
                  :class="
                    getProgressBarClass(
                      subscription.weekly_usage_usd,
                      subscriptionWeeklyLimit(subscription)
                    )
                  "
                  :style="{
                    width: getProgressWidth(
                      subscription.weekly_usage_usd,
                      subscriptionWeeklyLimit(subscription)
                    )
                  }"
                ></div>
              </div>
              <p
                v-if="subscription.weekly_window_start"
                class="text-xs text-gray-500 dark:text-dark-400"
              >
                {{
                  t('userSubscriptions.resetIn', {
                    time: formatResetTime(subscription.weekly_window_start, 168)
                  })
                }}
              </p>
            </div>

            <!-- Monthly Usage -->
            <div
              v-if="subscriptionMonthlyLimit(subscription)"
              class="space-y-2 rounded-md border border-gray-200 p-3 dark:border-dark-700"
            >
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-700 dark:text-gray-300">
                  {{ t('userSubscriptions.monthly') }}
                </span>
                <span class="text-sm text-gray-500 dark:text-dark-400">
                  ${{ (subscription.monthly_usage_usd || 0).toFixed(2) }} / ${{
                    subscriptionMonthlyLimit(subscription)?.toFixed(2)
                  }}
                </span>
              </div>
              <div class="relative h-1.5 overflow-hidden rounded-full bg-gray-100 dark:bg-dark-600">
                <div
                  class="absolute inset-y-0 left-0 rounded-full transition-all duration-300"
                  :class="
                    getProgressBarClass(
                      subscription.monthly_usage_usd,
                      subscriptionMonthlyLimit(subscription)
                    )
                  "
                  :style="{
                    width: getProgressWidth(
                      subscription.monthly_usage_usd,
                      subscriptionMonthlyLimit(subscription)
                    )
                  }"
                ></div>
              </div>
              <p
                v-if="subscription.monthly_window_start"
                class="text-xs text-gray-500 dark:text-dark-400"
              >
                {{
                  t('userSubscriptions.resetIn', {
                    time: formatResetTime(subscription.monthly_window_start, 720)
                  })
                }}
              </p>
            </div>

            <!-- No limits configured - Unlimited badge -->
            <div
              v-if="
                !subscriptionDailyLimit(subscription) &&
                !subscriptionWeeklyLimit(subscription) &&
                !subscriptionMonthlyLimit(subscription)
              "
              class="flex items-center gap-3 rounded-md border border-dashed border-gray-200 bg-gray-50 px-3 py-3 dark:border-dark-700 dark:bg-dark-900/40 md:col-span-3"
            >
              <span
                class="flex h-7 w-7 shrink-0 items-center justify-center rounded-md bg-emerald-50 text-emerald-600 dark:bg-emerald-500/10 dark:text-emerald-400"
              >
                <Icon name="checkCircle" size="sm" />
              </span>
              <div>
                <p class="text-xs font-medium text-gray-800 dark:text-gray-200">
                  {{ t('userSubscriptions.unlimited') }}
                </p>
                <p class="mt-0.5 text-xs text-gray-500 dark:text-gray-400">
                  {{ t('userSubscriptions.unlimitedDesc') }}
                </p>
              </div>
            </div>
          </div>
          </article>
        </div>
      </section>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { useAppStore } from '@/stores/app'
import subscriptionsAPI from '@/api/subscriptions'
import membershipAPI from '@/api/membership'
import type { UserSubscription } from '@/types'
import AppLayout from '@/components/layout/AppLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import { formatDateTimeToMinute } from '@/utils/format'
import { hasPeakRate, formatPeakRateWindow, serverTimezoneLabel } from '@/utils/peak-rate'
import { platformBadgeClass, platformLabel } from '@/utils/platformColors'
import { getRemainingDurationParts, isOneTimeDailyQuota, type RemainingDurationParts } from '@/utils/subscriptionQuota'

function platformAccentDotClass(p: string): string {
  switch (p) {
    case 'anthropic': return 'bg-orange-500'
    case 'openai': return 'bg-emerald-500'
    case 'antigravity': return 'bg-purple-500'
    case 'gemini': return 'bg-blue-500'
    default: return 'bg-gray-400'
  }
}

const { t } = useI18n()
const router = useRouter()
const appStore = useAppStore()

const subscriptions = ref<UserSubscription[]>([])
const loading = ref(true)
const updatingEntitlementId = ref<number | null>(null)

function subscriptionDailyLimit(subscription: UserSubscription): number | null | undefined {
  return subscription.entitlement?.daily_limit_usd ?? subscription.group?.daily_limit_usd
}

function subscriptionWeeklyLimit(subscription: UserSubscription): number | null | undefined {
  return subscription.entitlement?.weekly_limit_usd ?? subscription.group?.weekly_limit_usd
}

function subscriptionMonthlyLimit(subscription: UserSubscription): number | null | undefined {
  return subscription.entitlement?.monthly_limit_usd ?? subscription.group?.monthly_limit_usd
}

async function setOveragePolicy(subscription: UserSubscription, policy: 'block' | 'wallet') {
  const entitlement = subscription.entitlement
  if (!entitlement || entitlement.overage_policy === policy || updatingEntitlementId.value !== null) return
  updatingEntitlementId.value = entitlement.id
  try {
    await membershipAPI.setOveragePolicy(entitlement.id, policy)
    entitlement.overage_policy = policy
  } catch (error) {
    console.error('Failed to update subscription overage policy:', error)
    appStore.showError(t('common.error'))
  } finally {
    updatingEntitlementId.value = null
  }
}

function subscriptionHasPeakRate(subscription: UserSubscription): boolean {
  return hasPeakRate(subscription.group)
}

function subscriptionPeakRateLabel(subscription: UserSubscription): string {
  return formatPeakRateWindow(subscription.group, serverTimezoneLabel(appStore.cachedPublicSettings?.server_utc_offset))
}

async function loadSubscriptions() {
  try {
    loading.value = true
    subscriptions.value = await subscriptionsAPI.getMySubscriptions()
  } catch (error) {
    console.error('Failed to load subscriptions:', error)
    appStore.showError(t('userSubscriptions.failedToLoad'))
  } finally {
    loading.value = false
  }
}

function getProgressWidth(used: number | undefined, limit: number | null | undefined): string {
  if (!limit || limit === 0) return '0%'
  const percentage = Math.min(((used || 0) / limit) * 100, 100)
  return `${percentage}%`
}

function getProgressBarClass(used: number | undefined, limit: number | null | undefined): string {
  if (!limit || limit === 0) return 'bg-gray-400'
  const percentage = ((used || 0) / limit) * 100
  if (percentage >= 90) return 'bg-red-500'
  if (percentage >= 70) return 'bg-amber-500'
  return 'bg-primary-500'
}

function formatExpirationDate(expiresAt: string): string {
  const now = new Date()
  const expires = new Date(expiresAt)
  const diff = expires.getTime() - now.getTime()
  const days = Math.ceil(diff / (1000 * 60 * 60 * 24))

  if (days < 0) {
    return t('userSubscriptions.status.expired')
  }

  const dateStr = formatDateTimeToMinute(expires)

  if (days === 0) {
    return `${dateStr} (${t('common.today')})`
  }
  if (days === 1) {
    return `${dateStr} (${t('common.tomorrow')})`
  }

  return t('userSubscriptions.daysRemaining', { days }) + ` (${dateStr})`
}

function getExpirationClass(expiresAt: string): string {
  const now = new Date()
  const expires = new Date(expiresAt)
  const diff = expires.getTime() - now.getTime()
  const days = Math.ceil(diff / (1000 * 60 * 60 * 24))

  if (days <= 0) return 'text-red-600 dark:text-red-400 font-medium'
  if (days <= 3) return 'text-red-600 dark:text-red-400'
  if (days <= 7) return 'text-orange-600 dark:text-orange-400'
  return 'text-gray-700 dark:text-gray-300'
}

function formatDurationParts(parts: RemainingDurationParts): string {
  if (parts.days > 0) {
    return `${parts.days}d ${parts.hours}h`
  }

  if (parts.hours > 0) {
    return `${parts.hours}h ${parts.minutes}m`
  }

  return `${parts.minutes}m`
}

function formatDailyUsageWindow(subscription: UserSubscription): string {
  if (isOneTimeDailyQuota(subscription) && subscription.expires_at) {
    const parts = getRemainingDurationParts(subscription.expires_at)
    if (!parts) return t('userSubscriptions.windowNotActive')
    return t('userSubscriptions.quotaEndsIn', { time: formatDurationParts(parts) })
  }

  return t('userSubscriptions.resetIn', {
    time: formatResetTime(subscription.daily_window_start, 24)
  })
}

function formatResetTime(windowStart: string | null, windowHours: number): string {
  if (!windowStart) return t('userSubscriptions.windowNotActive')

  const start = new Date(windowStart)
  const end = new Date(start.getTime() + windowHours * 60 * 60 * 1000)
  const parts = getRemainingDurationParts(end)

  return parts ? formatDurationParts(parts) : t('userSubscriptions.windowNotActive')
}

onMounted(() => {
  loadSubscriptions()
})
</script>
