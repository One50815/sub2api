<template>
  <AppLayout>
    <div class="mx-auto max-w-6xl space-y-6 p-4 sm:p-6">
      <div class="flex flex-wrap items-end justify-between gap-4">
        <div>
          <p class="text-xs font-semibold uppercase tracking-[0.16em] text-primary-600 dark:text-primary-400">Omnio Pro</p>
          <h1 class="mt-1 text-2xl font-bold text-gray-950 dark:text-white">Omnio Pro 权益</h1>
          <p class="mt-2 max-w-2xl text-sm text-gray-500 dark:text-gray-400">Pro 可独立购买，也可以由兼容订阅套餐赠送；有效期、专属倍率、分组访问权限和并发加成都在这里集中展示。</p>
        </div>
        <button class="btn btn-secondary" :disabled="loading" @click="load">
          <Icon name="refresh" size="md" :class="loading ? 'animate-spin' : ''" />
          刷新权益
        </button>
      </div>

      <section class="grid gap-4 lg:grid-cols-[1.2fr_0.8fr]">
        <div class="overflow-hidden rounded-xl border border-gray-200 bg-white dark:border-dark-700 dark:bg-dark-900">
          <div class="border-b border-gray-100 px-5 py-4 dark:border-dark-700">
            <div class="flex items-center justify-between gap-3">
              <h2 class="font-semibold text-gray-900 dark:text-white">当前等级</h2>
              <span v-if="summary.effective_level" class="rounded-full px-3 py-1 text-xs font-bold text-white" :style="{ backgroundColor: summary.effective_level.badge_color }">{{ summary.effective_level.name }}</span>
              <span v-else class="rounded-full bg-gray-100 px-3 py-1 text-xs font-semibold text-gray-500 dark:bg-dark-800 dark:text-gray-400">暂未开通</span>
            </div>
          </div>
          <div class="p-5">
            <template v-if="summary.effective_level">
              <p class="text-lg font-semibold text-gray-950 dark:text-white">{{ summary.effective_level.description || '已启用 Omnio Pro 权益' }}</p>
              <div class="mt-5 grid gap-3 sm:grid-cols-3">
                <div class="rounded-lg bg-gray-50 p-3 dark:bg-dark-800"><p class="text-xs text-gray-500">并发加成</p><p class="mt-1 text-xl font-bold text-gray-900 dark:text-white">+{{ summary.effective_level.concurrency_bonus }}</p></div>
                <div class="rounded-lg bg-gray-50 p-3 dark:bg-dark-800"><p class="text-xs text-gray-500">优先支持</p><p class="mt-1 text-xl font-bold text-gray-900 dark:text-white">{{ summary.effective_level.priority_support ? '已开启' : '标准' }}</p></div>
                <div class="rounded-lg bg-gray-50 p-3 dark:bg-dark-800"><p class="text-xs text-gray-500">权益组数</p><p class="mt-1 text-xl font-bold text-gray-900 dark:text-white">{{ summary.benefits.length }}</p></div>
              </div>
              <div class="mt-5 flex flex-wrap gap-2">
                <span v-for="benefit in summary.benefits" :key="benefit.id" class="inline-flex items-center gap-1.5 rounded-full border border-gray-200 px-3 py-1.5 text-xs text-gray-600 dark:border-dark-700 dark:text-gray-300">
                  {{ benefit.group_name || `分组 #${benefit.group_id}` }}<span v-if="benefit.pro_only" class="font-semibold text-amber-600 dark:text-amber-300">Pro 专属</span><span v-if="benefit.rate_multiplier != null" class="font-semibold text-primary-600 dark:text-primary-400">{{ benefit.rate_multiplier }}x</span>
                </span>
              </div>
            </template>
            <p v-else class="text-sm text-gray-500 dark:text-gray-400">选择下方方案即可开通 Omnio Pro。</p>
          </div>
        </div>

        <div class="rounded-xl border border-gray-200 bg-gray-950 p-5 text-white dark:border-dark-700">
          <p class="text-xs font-semibold uppercase tracking-[0.16em] text-primary-300">权益来源</p>
          <h2 class="mt-2 text-xl font-bold">有效权益不会互相覆盖</h2>
          <p class="mt-2 text-sm leading-6 text-gray-300">多个来源同时有效时，系统取最高等级。单个订阅的倍率和额度按购买时快照执行，后续后台调整不会改变已购买权益。</p>
          <div class="mt-6 space-y-3">
            <div v-for="grant in activeGrants" :key="grant.id" class="flex items-center justify-between gap-3 border-b border-white/10 pb-3 last:border-0 last:pb-0"><span class="text-sm text-gray-300">{{ grant.source_type === 'manual' ? '管理员赠送' : grant.source_type === 'subscription_order' ? '订阅赠送' : 'Pro 订单' }}</span><span class="text-sm font-semibold">{{ formatDate(grant.expires_at) }} 到期</span></div>
            <p v-if="activeGrants.length === 0" class="text-sm text-gray-400">暂时没有有效 Pro 来源</p>
          </div>
        </div>
      </section>

      <section v-if="summary.quota_progress.length">
        <div class="mb-3">
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white">Pro 免费额度</h2>
          <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">按北京时间重置；免费额度用完后自动扣除钱包余额。</p>
        </div>
        <div class="grid gap-4 md:grid-cols-2">
          <article
            v-for="quota in summary.quota_progress"
            :key="`${quota.level_id}:${quota.group_id}`"
            class="rounded-xl border border-gray-200 bg-white p-5 dark:border-dark-700 dark:bg-dark-900"
          >
            <div class="flex items-center justify-between gap-3">
              <div><h3 class="font-semibold text-gray-950 dark:text-white">{{ quota.group_name }}</h3><p class="mt-1 text-xs text-gray-500 dark:text-gray-400">{{ quota.level_name }}</p></div>
              <span class="rounded-full bg-primary-50 px-2.5 py-1 text-xs font-semibold text-primary-700 dark:bg-primary-950/40 dark:text-primary-300">超额扣钱包</span>
            </div>
            <div class="mt-5 space-y-4">
              <div v-if="quota.daily_limit_usd > 0">
                <div class="mb-1.5 flex items-center justify-between text-sm">
                  <span class="text-gray-500 dark:text-gray-400">今日免费额度</span>
                  <span class="font-medium text-gray-900 dark:text-white">{{ formatUSD(quota.daily_remaining_usd) }} / {{ formatUSD(quota.daily_limit_usd) }} 剩余</span>
                </div>
                <div class="h-2 overflow-hidden rounded-full bg-gray-100 dark:bg-dark-700">
                  <div class="h-full rounded-full bg-primary-500 transition-all" :style="{ width: `${quotaPercent(quota.daily_used_usd, quota.daily_limit_usd)}%` }" />
                </div>
                <p class="mt-1.5 text-xs text-gray-400">已使用 {{ formatUSD(quota.daily_used_usd) }}</p>
              </div>
              <div v-if="quota.monthly_limit_usd > 0">
                <div class="mb-1.5 flex items-center justify-between text-sm">
                  <span class="text-gray-500 dark:text-gray-400">本月免费额度</span>
                  <span class="font-medium text-gray-900 dark:text-white">{{ formatUSD(quota.monthly_remaining_usd) }} / {{ formatUSD(quota.monthly_limit_usd) }} 剩余</span>
                </div>
                <div class="h-2 overflow-hidden rounded-full bg-gray-100 dark:bg-dark-700">
                  <div class="h-full rounded-full bg-emerald-500 transition-all" :style="{ width: `${quotaPercent(quota.monthly_used_usd, quota.monthly_limit_usd)}%` }" />
                </div>
                <p class="mt-1.5 text-xs text-gray-400">已使用 {{ formatUSD(quota.monthly_used_usd) }}</p>
              </div>
            </div>
          </article>
        </div>
      </section>

      <section>
        <div class="mb-3 flex items-center justify-between"><h2 class="text-lg font-semibold text-gray-900 dark:text-white">Omnio Pro 方案</h2><button class="text-sm font-semibold text-primary-600 hover:text-primary-700 dark:text-primary-400" @click="router.push('/purchase?tab=membership')">前往购买页</button></div>
        <div class="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
          <article v-for="offer in summary.offers" :key="offer.id" class="flex flex-col rounded-xl border border-gray-200 bg-white p-5 dark:border-dark-700 dark:bg-dark-900">
            <div class="flex items-start justify-between gap-3"><div><h3 class="font-semibold text-gray-950 dark:text-white">{{ offer.name }}</h3><p class="mt-1 text-xs text-gray-500 dark:text-gray-400">{{ offer.level_name }} · {{ offer.duration_days }} 天</p></div><span class="rounded-full px-2.5 py-1 text-[11px] font-bold text-white" :style="{ backgroundColor: offer.badge_color }">Omnio Pro</span></div>
            <p class="mt-4 min-h-10 text-sm text-gray-500 dark:text-gray-400">{{ offer.description || '解锁对应等级的 Omnio 专属权益。' }}</p>
            <div class="mt-5 flex items-end justify-between gap-3"><span class="font-mono text-2xl font-bold text-gray-950 dark:text-white">{{ offer.currency }} {{ offer.price.toFixed(2) }}</span><button class="btn btn-primary" @click="router.push({ path: '/purchase', query: { tab: 'membership', offer: String(offer.id) } })">立即开通</button></div>
          </article>
          <div v-if="!loading && summary.offers.length === 0" class="rounded-xl border border-dashed border-gray-300 p-8 text-center text-sm text-gray-500 dark:border-dark-600 dark:text-gray-400">管理员暂未发布 Omnio Pro 方案</div>
        </div>
      </section>

      <section v-if="summary.grants.length" class="rounded-xl border border-gray-200 bg-white dark:border-dark-700 dark:bg-dark-900">
        <div class="border-b border-gray-100 px-5 py-4 dark:border-dark-700"><h2 class="font-semibold text-gray-900 dark:text-white">Pro 来源记录</h2></div>
        <div class="divide-y divide-gray-100 dark:divide-dark-700"><div v-for="grant in summary.grants" :key="grant.id" class="flex flex-wrap items-center justify-between gap-3 px-5 py-4"><div><p class="font-medium text-gray-900 dark:text-white">{{ grant.level_name }} <span class="ml-2 rounded-full bg-gray-100 px-2 py-0.5 text-xs text-gray-500 dark:bg-dark-800 dark:text-gray-400">{{ grant.status }}</span></p><p class="mt-1 text-xs text-gray-500 dark:text-gray-400">{{ grant.source_type }} · {{ formatDate(grant.starts_at) }} - {{ formatDate(grant.expires_at) }}</p></div><span class="text-xs text-gray-500 dark:text-gray-400">{{ grant.notes }}</span></div></div>
      </section>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import AppLayout from '@/components/layout/AppLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import membershipAPI from '@/api/membership'
import type { MembershipSummary } from '@/types/membership'

const router = useRouter()
const loading = ref(false)
const summary = reactive<MembershipSummary>({ effective_level: null, grants: [], benefits: [], offers: [], quota_progress: [] })
const activeGrants = computed(() => summary.grants.filter((grant) => grant.status === 'active' && new Date(grant.expires_at).getTime() > Date.now()))

function formatDate(value: string) {
  return new Date(value).toLocaleDateString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit' })
}

function formatUSD(value: number) {
  return `$${Number(value || 0).toFixed(2)}`
}

function quotaPercent(used: number, limit: number) {
  if (!limit || limit <= 0) return 0
  return Math.min(100, Math.max(0, (used / limit) * 100))
}

async function load() {
  loading.value = true
  try {
    const { data } = await membershipAPI.getSummary()
    Object.assign(summary, data)
  } finally {
    loading.value = false
  }
}

onMounted(() => { void load() })
</script>
