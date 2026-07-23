<template>
  <AppLayout>
    <div class="mx-auto w-full max-w-7xl space-y-5">
      <header class="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <h1 class="text-2xl font-semibold tracking-tight text-gray-950 dark:text-white">{{ t('tickets.title') }}</h1>
          <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">{{ t('tickets.description') }}</p>
        </div>
        <button class="btn btn-primary" :disabled="!ticketStore.config.accept_new_tickets" @click="createOpen = true">
          <span class="text-lg leading-none">＋</span>{{ t('tickets.newTicket') }}
        </button>
      </header>

      <div v-if="!ticketStore.config.accept_new_tickets" class="rounded-xl border border-amber-200 bg-amber-50 px-4 py-3 text-sm text-amber-800 dark:border-amber-800 dark:bg-amber-950/30 dark:text-amber-300">
        {{ t('tickets.intakePaused') }}
      </div>

      <TicketSummaryCards :summary="summary" />

      <section class="card overflow-hidden">
        <div class="flex flex-col gap-3 border-b border-gray-100 p-4 dark:border-dark-700 sm:flex-row sm:items-center">
          <select v-model="filters.status" class="input sm:w-48" @change="resetAndLoad">
            <option value="">{{ t('tickets.filters.all') }} · {{ t('tickets.filters.status') }}</option>
            <option v-for="status in TICKET_STATUSES" :key="status" :value="status">{{ t(`tickets.status.${status}`) }}</option>
          </select>
          <select v-model="filters.category" class="input sm:w-48" @change="resetAndLoad">
            <option value="">{{ t('tickets.filters.all') }} · {{ t('tickets.filters.category') }}</option>
            <option v-for="category in TICKET_CATEGORIES" :key="category" :value="category">{{ t(`tickets.category.${category}`) }}</option>
          </select>
          <button class="btn btn-secondary sm:ml-auto" :disabled="loading" @click="load">{{ t('common.refresh') }}</button>
        </div>

        <div v-if="loading" class="space-y-3 p-5">
          <div v-for="index in 5" :key="index" class="h-16 animate-pulse rounded-xl bg-gray-100 dark:bg-dark-700"></div>
        </div>
        <div v-else-if="pageData.items.length === 0" class="px-6 py-16 text-center">
          <div class="mx-auto flex h-12 w-12 items-center justify-center rounded-2xl bg-primary-50 text-2xl text-primary-600 dark:bg-primary-950/40">?</div>
          <h2 class="mt-4 font-semibold text-gray-900 dark:text-white">{{ t('tickets.empty') }}</h2>
          <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">{{ t('tickets.emptyDescription') }}</p>
        </div>
        <div v-else>
          <div class="hidden overflow-x-auto md:block">
            <table class="min-w-full divide-y divide-gray-100 dark:divide-dark-700">
              <thead class="bg-gray-50/70 dark:bg-dark-800/70">
                <tr>
                  <th class="px-5 py-3 text-left text-xs font-semibold text-gray-500">{{ t('tickets.table.subject') }}</th>
                  <th class="px-5 py-3 text-left text-xs font-semibold text-gray-500">{{ t('tickets.table.category') }}</th>
                  <th class="px-5 py-3 text-left text-xs font-semibold text-gray-500">{{ t('tickets.table.status') }}</th>
                  <th class="px-5 py-3 text-left text-xs font-semibold text-gray-500">{{ t('tickets.table.updatedAt') }}</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-100 dark:divide-dark-700">
                <tr v-for="ticket in pageData.items" :key="ticket.id" class="cursor-pointer transition-colors hover:bg-gray-50 dark:hover:bg-dark-800/70" @click="openTicket(ticket.id)">
                  <td class="px-5 py-4">
                    <div class="flex items-center gap-2 font-medium text-gray-900 dark:text-white"><span>#{{ ticket.id }}</span><span class="max-w-xl truncate">{{ ticket.subject }}</span></div>
                    <span v-if="ticket.unread_count" class="mt-1 inline-block text-xs font-semibold text-primary-600 dark:text-primary-400">{{ t('tickets.table.unread', { count: ticket.unread_count }) }}</span>
                  </td>
                  <td class="px-5 py-4"><TicketBadge kind="category" :value="ticket.category" /></td>
                  <td class="px-5 py-4"><TicketBadge kind="status" :value="ticket.status" /></td>
                  <td class="whitespace-nowrap px-5 py-4 text-sm text-gray-500 dark:text-gray-400">{{ formatDate(ticket.last_message_at) }}</td>
                </tr>
              </tbody>
            </table>
          </div>
          <div class="divide-y divide-gray-100 dark:divide-dark-700 md:hidden">
            <button v-for="ticket in pageData.items" :key="ticket.id" class="block w-full p-4 text-left" @click="openTicket(ticket.id)">
              <div class="flex items-start justify-between gap-3"><span class="font-medium text-gray-900 dark:text-white">#{{ ticket.id }} {{ ticket.subject }}</span><span v-if="ticket.unread_count" class="min-w-6 rounded-full bg-primary-600 px-1.5 py-0.5 text-center text-xs font-bold text-white">{{ ticket.unread_count }}</span></div>
              <div class="mt-3 flex flex-wrap items-center gap-2"><TicketBadge kind="status" :value="ticket.status" /><TicketBadge kind="category" :value="ticket.category" /><span class="ml-auto text-xs text-gray-500">{{ formatDate(ticket.last_message_at) }}</span></div>
            </button>
          </div>
        </div>

        <div v-if="pageData.total > pageSize" class="flex items-center justify-between border-t border-gray-100 px-4 py-3 dark:border-dark-700">
          <span class="text-xs text-gray-500">{{ (page - 1) * pageSize + 1 }}–{{ Math.min(page * pageSize, pageData.total) }} / {{ pageData.total }}</span>
          <div class="flex gap-2"><button class="btn btn-secondary" :disabled="page <= 1" @click="changePage(page - 1)">‹</button><button class="btn btn-secondary" :disabled="page * pageSize >= pageData.total" @click="changePage(page + 1)">›</button></div>
        </div>
      </section>
    </div>

    <TicketCreateDialog :show="createOpen" :submitting="creating" @close="createOpen = false" @submit="createTicket" />
  </AppLayout>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import AppLayout from '@/components/layout/AppLayout.vue'
import TicketBadge from '@/components/ticket/TicketBadge.vue'
import TicketCreateDialog from '@/components/ticket/TicketCreateDialog.vue'
import TicketSummaryCards from '@/components/ticket/TicketSummaryCards.vue'
import { ticketsAPI } from '@/api/tickets'
import { useAppStore } from '@/stores/app'
import { useTicketStore } from '@/stores/tickets'
import { TICKET_CATEGORIES, TICKET_STATUSES, type CreateTicketPayload, type TicketListParams, type TicketPage, type TicketSummary } from '@/types/ticket'

const { t, locale } = useI18n()
const router = useRouter()
const appStore = useAppStore()
const ticketStore = useTicketStore()
const loading = ref(false)
const creating = ref(false)
const createOpen = ref(false)
const page = ref(1)
const pageSize = 20
const filters = reactive<TicketListParams>({ status: '', category: '' })
const pageData = reactive<TicketPage>({ items: [], total: 0, page: 1, page_size: pageSize })
const summary = reactive<TicketSummary>({ total: 0, pending_admin: 0, pending_user: 0, resolved: 0, closed: 0, unread_count: 0, pending_admin_count: 0, open_count: 0 })

const load = async () => {
  loading.value = true
  try {
    const [list, stats] = await Promise.all([ticketsAPI.list({ ...filters, page: page.value, page_size: pageSize }), ticketsAPI.summary()])
    Object.assign(pageData, list)
    Object.assign(summary, stats)
    ticketStore.userSummary = stats
  } catch (error) {
    appStore.showError((error as { message?: string })?.message || t('tickets.loadFailed'))
  } finally { loading.value = false }
}
const resetAndLoad = () => { page.value = 1; void load() }
const changePage = (value: number) => { page.value = value; void load() }
const openTicket = (id: number) => router.push(`/tickets/${id}`)
const createTicket = async (payload: CreateTicketPayload) => {
  creating.value = true
  try {
    const detail = await ticketsAPI.create(payload)
    appStore.showSuccess(t('tickets.create.success'))
    createOpen.value = false
    await router.push(`/tickets/${detail.ticket.id}`)
  } catch (error) { appStore.showError((error as { message?: string })?.message || t('common.error')) }
  finally { creating.value = false }
}
const formatDate = (value: string) => new Intl.DateTimeFormat(locale.value, { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(value))
onMounted(async () => { await ticketStore.fetchConfig(false, true); await load() })
</script>

