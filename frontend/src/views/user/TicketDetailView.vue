<template>
  <AppLayout>
    <div class="mx-auto w-full max-w-5xl space-y-5">
      <button class="inline-flex items-center gap-2 text-sm font-medium text-gray-500 hover:text-primary-600 dark:text-gray-400" @click="router.push('/tickets')">← {{ t('tickets.back') }}</button>
      <div v-if="loading" class="card h-64 animate-pulse bg-gray-100 dark:bg-dark-700"></div>
      <div v-else-if="!detail" class="card px-6 py-20 text-center text-sm text-gray-500">{{ t('tickets.notFound') }}</div>
      <template v-else>
        <header class="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
          <div class="min-w-0">
            <div class="flex flex-wrap items-center gap-2"><TicketBadge kind="status" :value="detail.ticket.status" /><TicketBadge kind="category" :value="detail.ticket.category" /></div>
            <h1 class="mt-3 break-words text-2xl font-semibold tracking-tight text-gray-950 dark:text-white">#{{ detail.ticket.id }} {{ detail.ticket.subject }}</h1>
            <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">{{ t('tickets.detail.createdAt') }} {{ formatDate(detail.ticket.created_at) }}</p>
          </div>
          <div v-if="detail.ticket.status !== 'closed'" class="flex shrink-0 flex-wrap gap-2">
            <button v-if="detail.ticket.status !== 'resolved'" class="btn btn-secondary" :disabled="acting" @click="setStatus('resolved')">{{ t('tickets.actions.resolve') }}</button>
            <button class="btn btn-danger" :disabled="acting" @click="setStatus('closed')">{{ t('tickets.actions.close') }}</button>
          </div>
        </header>

        <section class="grid gap-4 sm:grid-cols-3">
          <div class="card p-4"><p class="text-xs font-medium text-gray-500">{{ t('tickets.table.status') }}</p><div class="mt-2"><TicketBadge kind="status" :value="detail.ticket.status" /></div></div>
          <div class="card p-4"><p class="text-xs font-medium text-gray-500">{{ t('tickets.detail.updatedAt') }}</p><p class="mt-2 text-sm font-semibold text-gray-900 dark:text-white">{{ formatDate(detail.ticket.last_message_at) }}</p></div>
          <div class="card p-4"><p class="text-xs font-medium text-gray-500">{{ t('tickets.detail.requestId') }}</p><p class="mt-2 truncate font-mono text-sm text-gray-900 dark:text-white">{{ detail.ticket.related_request_id || '—' }}</p></div>
        </section>

        <section class="card overflow-hidden">
          <div class="border-b border-gray-100 px-5 py-4 dark:border-dark-700"><h2 class="font-semibold text-gray-900 dark:text-white">{{ t('tickets.detail.conversation') }}</h2></div>
          <div class="p-5"><TicketThread :messages="detail.messages" /></div>
          <div class="border-t border-gray-100 p-5 dark:border-dark-700">
            <p v-if="detail.ticket.status === 'closed'" class="rounded-lg bg-gray-50 px-4 py-3 text-sm text-gray-500 dark:bg-dark-800 dark:text-gray-400">{{ t('tickets.reply.closed') }}</p>
            <template v-else><p v-if="detail.ticket.status === 'resolved'" class="mb-3 text-xs text-gray-500">{{ t('tickets.actions.reopenHint') }}</p><TicketReplyComposer :submitting="replying" @submit="reply" /></template>
          </div>
        </section>
      </template>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import AppLayout from '@/components/layout/AppLayout.vue'
import TicketBadge from '@/components/ticket/TicketBadge.vue'
import TicketReplyComposer from '@/components/ticket/TicketReplyComposer.vue'
import TicketThread from '@/components/ticket/TicketThread.vue'
import { ticketsAPI } from '@/api/tickets'
import { useAppStore } from '@/stores/app'
import { useTicketStore } from '@/stores/tickets'
import type { ReplyTicketPayload, TicketDetail } from '@/types/ticket'

const { t, locale } = useI18n(); const route = useRoute(); const router = useRouter(); const appStore = useAppStore(); const ticketStore = useTicketStore()
const id = Number(route.params.id); const detail = ref<TicketDetail | null>(null); const loading = ref(false); const replying = ref(false); const acting = ref(false)
const load = async () => { loading.value = true; try { detail.value = await ticketsAPI.get(id); await ticketsAPI.markRead(id); await ticketStore.refresh(false) } catch (e) { appStore.showError((e as { message?: string })?.message || t('tickets.loadFailed')) } finally { loading.value = false } }
const reply = async (payload: ReplyTicketPayload) => { replying.value = true; try { await ticketsAPI.reply(id, payload); appStore.showSuccess(t('tickets.reply.success')); await load() } catch (e) { appStore.showError((e as { message?: string })?.message || t('common.error')) } finally { replying.value = false } }
const setStatus = async (status: 'resolved' | 'closed') => { acting.value = true; try { await ticketsAPI.updateStatus(id, status); appStore.showSuccess(t(status === 'resolved' ? 'tickets.actions.resolved' : 'tickets.actions.closed')); await load() } catch (e) { appStore.showError((e as { message?: string })?.message || t('common.error')) } finally { acting.value = false } }
const formatDate = (value: string) => new Intl.DateTimeFormat(locale.value, { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(value))
onMounted(load)
</script>

