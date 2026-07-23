<template>
  <form class="space-y-3" @submit.prevent="submit">
    <textarea v-model="content" class="input min-h-28 resize-y" :placeholder="t('tickets.reply.placeholder')" :disabled="disabled || submitting" maxlength="10000"></textarea>
    <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
      <input v-model="requestId" class="input font-mono text-sm sm:max-w-sm" :placeholder="t('tickets.reply.requestId')" :disabled="disabled || submitting" maxlength="128" />
      <button type="submit" class="btn btn-primary shrink-0" :disabled="disabled || submitting || !content.trim()">
        {{ submitting ? t('common.submitting') : t('tickets.reply.send') }}
      </button>
    </div>
  </form>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import type { ReplyTicketPayload } from '@/types/ticket'

defineProps<{ disabled?: boolean; submitting?: boolean }>()
const emit = defineEmits<{ submit: [payload: ReplyTicketPayload] }>()
const { t } = useI18n()
const content = ref('')
const requestId = ref('')
const submit = () => {
  if (!content.value.trim()) return
  emit('submit', { content: content.value.trim(), request_id: requestId.value.trim() || undefined })
  content.value = ''
  requestId.value = ''
}
</script>

