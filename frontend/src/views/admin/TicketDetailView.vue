<template>
  <AppLayout>
    <div class="mx-auto w-full max-w-6xl space-y-5">
      <button class="inline-flex items-center gap-2 text-sm font-medium text-gray-500 hover:text-primary-600 dark:text-gray-400" @click="router.push('/admin/tickets')">← {{ t('tickets.back') }}</button>
      <div v-if="loading" class="card h-72 animate-pulse bg-gray-100 dark:bg-dark-700"></div>
      <div v-else-if="!detail" class="card px-6 py-20 text-center text-sm text-gray-500">{{ t('tickets.notFound') }}</div>
      <template v-else>
        <header>
          <div class="flex flex-wrap items-center gap-2"><TicketBadge kind="status" :value="detail.ticket.status" /><TicketBadge kind="priority" :value="detail.ticket.priority" /><TicketBadge kind="category" :value="detail.ticket.category" /></div>
          <h1 class="mt-3 break-words text-2xl font-semibold tracking-tight text-gray-950 dark:text-white">#{{ detail.ticket.id }} {{ detail.ticket.subject }}</h1>
          <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">{{ detail.ticket.username || `#${detail.ticket.user_id}` }} · {{ formatDate(detail.ticket.created_at) }}</p>
        </header>

        <div class="grid gap-5 lg:grid-cols-[minmax(0,1fr)_320px]">
          <section class="card overflow-hidden">
            <div class="border-b border-gray-100 px-5 py-4 dark:border-dark-700"><h2 class="font-semibold text-gray-900 dark:text-white">{{ t('tickets.detail.conversation') }}</h2></div>
            <div class="p-5"><TicketThread :messages="detail.messages" admin /></div>
            <div class="border-t border-gray-100 p-5 dark:border-dark-700"><TicketReplyComposer :disabled="detail.ticket.status === 'closed'" :submitting="replying" @submit="reply" /></div>
          </section>

          <aside class="space-y-4">
            <section class="card space-y-4 p-5">
              <label class="block"><span class="mb-1.5 block text-xs font-semibold text-gray-500">{{ t('tickets.filters.status') }}</span><select v-model="form.status" class="input"><option v-for="v in TICKET_STATUSES" :key="v" :value="v">{{ t(`tickets.status.${v}`) }}</option></select></label>
              <label class="block"><span class="mb-1.5 block text-xs font-semibold text-gray-500">{{ t('tickets.filters.priority') }}</span><select v-model="form.priority" class="input"><option v-for="v in TICKET_PRIORITIES" :key="v" :value="v">{{ t(`tickets.priority.${v}`) }}</option></select></label>
              <label class="block"><span class="mb-1.5 block text-xs font-semibold text-gray-500">{{ t('tickets.filters.assignee') }}</span><select v-model.number="form.assignee_id" class="input"><option :value="0">{{ t('tickets.admin.noAssignee') }}</option><option v-for="item in assignees" :key="item.id" :value="item.id">{{ item.username }}</option></select></label>
              <div class="flex gap-2"><button class="btn btn-secondary flex-1" :disabled="!currentAdminID" @click="claim">{{ t('tickets.admin.claim') }}</button><button class="btn btn-primary flex-1" :disabled="saving" @click="save">{{ t('tickets.admin.save') }}</button></div>
              <p class="text-[11px] leading-5 text-gray-500 dark:text-gray-400">{{ t('tickets.admin.permissions') }}</p>
            </section>
            <section class="card space-y-3 p-5 text-sm">
              <div><p class="text-xs font-semibold text-gray-500">{{ t('tickets.detail.requestId') }}</p><p class="mt-1 break-all font-mono text-gray-900 dark:text-white">{{ detail.ticket.related_request_id || '—' }}</p></div>
              <div><p class="text-xs font-semibold text-gray-500">{{ t('tickets.detail.updatedAt') }}</p><p class="mt-1 text-gray-900 dark:text-white">{{ formatDate(detail.ticket.last_message_at) }}</p></div>
            </section>
          </aside>
        </div>
      </template>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import AppLayout from '@/components/layout/AppLayout.vue'
import TicketBadge from '@/components/ticket/TicketBadge.vue'
import TicketReplyComposer from '@/components/ticket/TicketReplyComposer.vue'
import TicketThread from '@/components/ticket/TicketThread.vue'
import { ticketsAPI } from '@/api/tickets'
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'
import { useTicketStore } from '@/stores/tickets'
import { TICKET_PRIORITIES, TICKET_STATUSES, type ReplyTicketPayload, type TicketAssignee, type TicketDetail, type TicketPriority, type TicketStatus } from '@/types/ticket'

const { t, locale } = useI18n(); const route = useRoute(); const router = useRouter(); const appStore = useAppStore(); const authStore = useAuthStore(); const ticketStore = useTicketStore()
const id = Number(route.params.id); const detail = ref<TicketDetail | null>(null); const assignees = ref<TicketAssignee[]>([]); const loading = ref(false); const replying = ref(false); const saving = ref(false)
const form = reactive<{ status: TicketStatus; priority: TicketPriority; assignee_id: number }>({ status: 'pending_admin', priority: 'normal', assignee_id: 0 })
const currentAdminID = computed(() => Number(authStore.user?.id || 0))
const syncForm = () => { if (!detail.value) return; form.status = detail.value.ticket.status; form.priority = detail.value.ticket.priority; form.assignee_id = detail.value.ticket.assignee_id || 0 }
const load = async () => { loading.value = true; try { detail.value = await ticketsAPI.admin.get(id); syncForm(); await ticketsAPI.admin.markRead(id); await ticketStore.refresh(true) } catch (e) { appStore.showError((e as { message?: string })?.message || t('tickets.loadFailed')) } finally { loading.value = false } }
const reply = async (payload: ReplyTicketPayload) => { replying.value = true; try { await ticketsAPI.admin.reply(id, payload); appStore.showSuccess(t('tickets.reply.success')); await load() } catch (e) { appStore.showError((e as { message?: string })?.message || t('common.error')) } finally { replying.value = false } }
const save = async () => { saving.value = true; try { await ticketsAPI.admin.update(id, { ...form }); appStore.showSuccess(t('tickets.admin.updated')); await load() } catch (e) { appStore.showError((e as { message?: string })?.message || t('common.error')) } finally { saving.value = false } }
const claim = () => { if (currentAdminID.value) form.assignee_id = currentAdminID.value; void save() }
const formatDate = (value: string) => new Intl.DateTimeFormat(locale.value, { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(value))
onMounted(async () => { try { assignees.value = await ticketsAPI.admin.assignees() } catch { assignees.value = [] } await load() })
</script>
