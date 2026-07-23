<template>
  <AppLayout>
    <div class="mx-auto w-full max-w-[1500px] space-y-5">
      <header class="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div><h1 class="text-2xl font-semibold tracking-tight text-gray-950 dark:text-white">{{ t('tickets.admin.title') }}</h1><p class="mt-1 text-sm text-gray-500 dark:text-gray-400">{{ t('tickets.admin.description') }}</p></div>
        <button class="btn btn-secondary" @click="router.push('/tickets')">{{ t('tickets.admin.myTickets') }}</button>
      </header>
      <TicketSummaryCards :summary="summary" />
      <section class="card overflow-hidden">
        <form class="grid gap-3 border-b border-gray-100 p-4 dark:border-dark-700 sm:grid-cols-2 lg:grid-cols-7" @submit.prevent="applyFilters">
          <select v-model="filters.status" class="input"><option value="">{{ t('tickets.filters.status') }} · {{ t('tickets.filters.all') }}</option><option v-for="v in TICKET_STATUSES" :key="v" :value="v">{{ t(`tickets.status.${v}`) }}</option></select>
          <select v-model="filters.category" class="input"><option value="">{{ t('tickets.filters.category') }} · {{ t('tickets.filters.all') }}</option><option v-for="v in TICKET_CATEGORIES" :key="v" :value="v">{{ t(`tickets.category.${v}`) }}</option></select>
          <select v-model="filters.priority" class="input"><option value="">{{ t('tickets.filters.priority') }} · {{ t('tickets.filters.all') }}</option><option v-for="v in TICKET_PRIORITIES" :key="v" :value="v">{{ t(`tickets.priority.${v}`) }}</option></select>
          <input v-model.number="filters.user_id" type="number" min="1" class="input" :placeholder="t('tickets.filters.userId')" />
          <select v-model.number="filters.assignee_id" class="input"><option :value="undefined">{{ t('tickets.filters.assignee') }} · {{ t('tickets.filters.all') }}</option><option :value="0">{{ t('tickets.filters.unassigned') }}</option><option v-for="item in assignees" :key="item.id" :value="item.id">{{ item.username }}</option></select>
          <input v-model="filters.related_request_id" class="input font-mono text-sm" :placeholder="t('tickets.filters.requestId')" />
          <div class="flex gap-2"><button class="btn btn-primary flex-1" type="submit">{{ t('tickets.filters.apply') }}</button><button class="btn btn-secondary" type="button" @click="resetFilters">↺</button></div>
        </form>
        <div v-if="loading" class="space-y-3 p-5"><div v-for="n in 6" :key="n" class="h-16 animate-pulse rounded-xl bg-gray-100 dark:bg-dark-700"></div></div>
        <div v-else-if="pageData.items.length === 0" class="px-6 py-16 text-center text-sm text-gray-500">{{ t('tickets.empty') }}</div>
        <div v-else class="overflow-x-auto">
          <table class="min-w-[1050px] w-full divide-y divide-gray-100 dark:divide-dark-700">
            <thead class="bg-gray-50/70 dark:bg-dark-800/70"><tr><th v-for="head in headers" :key="head" class="px-4 py-3 text-left text-xs font-semibold text-gray-500">{{ head }}</th></tr></thead>
            <tbody class="divide-y divide-gray-100 dark:divide-dark-700">
              <tr v-for="ticket in pageData.items" :key="ticket.id" class="cursor-pointer hover:bg-gray-50 dark:hover:bg-dark-800/70" @click="router.push(`/admin/tickets/${ticket.id}`)">
                <td class="px-4 py-4"><div class="max-w-md truncate font-medium text-gray-900 dark:text-white">#{{ ticket.id }} {{ ticket.subject }}</div><span v-if="ticket.unread_count" class="text-xs font-semibold text-primary-600">{{ t('tickets.table.unread', { count: ticket.unread_count }) }}</span></td>
                <td class="px-4 py-4 text-sm text-gray-600 dark:text-gray-300">{{ ticket.username || `#${ticket.user_id}` }}</td>
                <td class="px-4 py-4"><TicketBadge kind="category" :value="ticket.category" /></td>
                <td class="px-4 py-4"><TicketBadge kind="status" :value="ticket.status" /></td>
                <td class="px-4 py-4"><TicketBadge kind="priority" :value="ticket.priority" /></td>
                <td class="px-4 py-4 text-sm text-gray-600 dark:text-gray-300">{{ ticket.assignee_username || t('tickets.admin.noAssignee') }}</td>
                <td class="whitespace-nowrap px-4 py-4 text-sm text-gray-500">{{ formatDate(ticket.last_message_at) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
        <div v-if="pageData.total > pageSize" class="flex items-center justify-between border-t border-gray-100 px-4 py-3 dark:border-dark-700"><span class="text-xs text-gray-500">{{ (page - 1) * pageSize + 1 }}–{{ Math.min(page * pageSize, pageData.total) }} / {{ pageData.total }}</span><div class="flex gap-2"><button class="btn btn-secondary" :disabled="page <= 1" @click="changePage(page - 1)">‹</button><button class="btn btn-secondary" :disabled="page * pageSize >= pageData.total" @click="changePage(page + 1)">›</button></div></div>
      </section>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import AppLayout from '@/components/layout/AppLayout.vue'
import TicketBadge from '@/components/ticket/TicketBadge.vue'
import TicketSummaryCards from '@/components/ticket/TicketSummaryCards.vue'
import { ticketsAPI } from '@/api/tickets'
import { useAppStore } from '@/stores/app'
import { useTicketStore } from '@/stores/tickets'
import { TICKET_CATEGORIES, TICKET_PRIORITIES, TICKET_STATUSES, type TicketAssignee, type TicketListParams, type TicketPage, type TicketSummary } from '@/types/ticket'

const { t, locale } = useI18n(); const router = useRouter(); const appStore = useAppStore(); const ticketStore = useTicketStore()
const loading = ref(false); const page = ref(1); const pageSize = 20; const assignees = ref<TicketAssignee[]>([])
const filters = reactive<TicketListParams>({ status: '', category: '', priority: '', related_request_id: '' })
const pageData = reactive<TicketPage>({ items: [], total: 0, page: 1, page_size: pageSize })
const summary = reactive<TicketSummary>({ total: 0, pending_admin: 0, pending_user: 0, resolved: 0, closed: 0, unread_count: 0, pending_admin_count: 0, open_count: 0 })
const headers = computed(() => [t('tickets.table.subject'), t('tickets.table.user'), t('tickets.table.category'), t('tickets.table.status'), t('tickets.table.priority'), t('tickets.table.assignee'), t('tickets.table.updatedAt')])
const load = async () => { loading.value = true; try { const [list, stats] = await Promise.all([ticketsAPI.admin.list({ ...filters, page: page.value, page_size: pageSize }), ticketsAPI.admin.summary()]); Object.assign(pageData, list); Object.assign(summary, stats); ticketStore.adminSummary = stats } catch (e) { appStore.showError((e as { message?: string })?.message || t('tickets.loadFailed')) } finally { loading.value = false } }
const applyFilters = () => { page.value = 1; void load() }
const resetFilters = () => { Object.assign(filters, { status: '', category: '', priority: '', user_id: undefined, assignee_id: undefined, related_request_id: '' }); applyFilters() }
const changePage = (value: number) => { page.value = value; void load() }
const formatDate = (value: string) => new Intl.DateTimeFormat(locale.value, { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(value))
onMounted(async () => { try { assignees.value = await ticketsAPI.admin.assignees() } catch { assignees.value = [] } await load() })
</script>
