<template>
  <div class="space-y-4">
    <article v-for="message in messages" :key="message.id" :class="['flex', message.sender_type === ownSender ? 'justify-end' : 'justify-start']">
      <div :class="['max-w-[88%] rounded-2xl border px-4 py-3 sm:max-w-[72%]', bubbleClass(message.sender_type)]">
        <div class="mb-2 flex items-center justify-between gap-4 text-xs">
          <span class="font-semibold">{{ senderLabel(message.sender_type) }}</span>
          <time class="whitespace-nowrap opacity-70">{{ formatTime(message.created_at) }}</time>
        </div>
        <p class="whitespace-pre-wrap break-words text-sm leading-6">{{ message.content }}</p>
        <p v-if="message.request_id" class="mt-3 rounded-lg bg-black/5 px-2.5 py-1.5 font-mono text-[11px] dark:bg-white/5">
          {{ t('tickets.detail.messageRequestId', { id: message.request_id }) }}
        </p>
      </div>
    </article>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import type { TicketMessage, TicketSenderType } from '@/types/ticket'

const props = defineProps<{ messages: TicketMessage[]; admin?: boolean }>()
const { t, locale } = useI18n()
const ownSender = computed<TicketSenderType>(() => (props.admin ? 'admin' : 'user'))
const senderLabel = (sender: TicketSenderType) => props.admin
  ? t(sender === 'admin' ? 'tickets.detail.adminLabel' : 'tickets.detail.userLabel')
  : t(sender === 'user' ? 'tickets.detail.userMessage' : 'tickets.detail.adminMessage')
const bubbleClass = (sender: TicketSenderType) => sender === ownSender.value
  ? 'border-primary-200 bg-primary-50 text-gray-900 dark:border-primary-800 dark:bg-primary-950/35 dark:text-gray-100'
  : 'border-gray-200 bg-white text-gray-900 dark:border-dark-600 dark:bg-dark-800 dark:text-gray-100'
const formatTime = (value: string) => new Intl.DateTimeFormat(locale.value, { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(value))
</script>

