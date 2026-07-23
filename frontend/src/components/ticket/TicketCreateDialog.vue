<template>
  <BaseDialog :show="show" :title="t('tickets.create.title')" width="normal" @close="emit('close')">
    <form id="create-ticket-form" class="space-y-4" @submit.prevent="submit">
      <label class="block">
        <span class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('tickets.create.subject') }}</span>
        <input v-model="form.subject" class="input" maxlength="200" required :placeholder="t('tickets.create.subjectPlaceholder')" />
      </label>
      <label class="block">
        <span class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('tickets.create.category') }}</span>
        <select v-model="form.category" class="input" required>
          <option v-for="value in TICKET_CATEGORIES" :key="value" :value="value">{{ t(`tickets.category.${value}`) }}</option>
        </select>
      </label>
      <label class="block">
        <span class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('tickets.create.content') }}</span>
        <textarea v-model="form.content" class="input min-h-36 resize-y" maxlength="10000" required :placeholder="t('tickets.create.contentPlaceholder')"></textarea>
      </label>
      <label class="block">
        <span class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('tickets.create.requestId') }}</span>
        <input v-model="form.related_request_id" class="input font-mono text-sm" maxlength="128" />
        <span class="mt-1.5 block text-xs text-gray-500 dark:text-gray-400">{{ t('tickets.create.requestIdHint') }}</span>
      </label>
    </form>
    <template #footer>
      <button type="button" class="btn btn-secondary" @click="emit('close')">{{ t('common.cancel') }}</button>
      <button type="submit" form="create-ticket-form" class="btn btn-primary" :disabled="submitting || !form.subject.trim() || !form.content.trim()">
        {{ submitting ? t('common.submitting') : t('tickets.create.submit') }}
      </button>
    </template>
  </BaseDialog>
</template>

<script setup lang="ts">
import { reactive, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import BaseDialog from '@/components/common/BaseDialog.vue'
import { TICKET_CATEGORIES, type CreateTicketPayload, type TicketCategory } from '@/types/ticket'

const props = defineProps<{ show: boolean; submitting?: boolean }>()
const emit = defineEmits<{ close: []; submit: [payload: CreateTicketPayload] }>()
const { t } = useI18n()
const initial = (): CreateTicketPayload => ({ subject: '', category: 'api_model' as TicketCategory, content: '', related_request_id: '' })
const form = reactive<CreateTicketPayload>(initial())
watch(() => props.show, (show) => { if (show) Object.assign(form, initial()) })
const submit = () => emit('submit', {
  subject: form.subject.trim(), category: form.category, content: form.content.trim(), related_request_id: form.related_request_id?.trim() || undefined,
})
</script>

