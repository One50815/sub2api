<template>
  <div class="card">
    <div class="border-b border-gray-100 px-6 py-4 dark:border-dark-700">
      <h2 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('tickets.settings.title') }}</h2>
      <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">{{ t('tickets.settings.description') }}</p>
    </div>
    <div class="p-6">
      <div class="mb-4">
        <h3 class="text-sm font-semibold text-gray-900 dark:text-white">{{ t('tickets.settings.personalArea') }}</h3>
      </div>
      <div class="divide-y divide-gray-100 overflow-hidden rounded-xl border border-gray-200 dark:divide-dark-700 dark:border-dark-600">
        <div class="flex items-center justify-between gap-5 p-4">
          <div>
            <p class="text-sm font-medium text-gray-900 dark:text-white">{{ t('tickets.settings.showCenter') }}</p>
            <p class="mt-1 text-xs leading-5 text-gray-500 dark:text-gray-400">{{ t('tickets.settings.showCenterHint') }}</p>
          </div>
          <Toggle v-model="draft.user_center_enabled" />
        </div>
        <div class="flex items-center justify-between gap-5 p-4" :class="{ 'opacity-60': !draft.user_center_enabled }">
          <div>
            <p class="text-sm font-medium text-gray-900 dark:text-white">{{ t('tickets.settings.acceptNew') }}</p>
            <p class="mt-1 text-xs leading-5 text-gray-500 dark:text-gray-400">{{ t('tickets.settings.acceptNewHint') }}</p>
          </div>
          <Toggle v-model="draft.accept_new_tickets" :disabled="!draft.user_center_enabled" />
        </div>
      </div>
      <div class="mt-5 flex justify-end">
        <button type="button" class="btn btn-primary" :disabled="saving || !changed" @click="save">
          {{ saving ? t('common.saving') : t('tickets.settings.save') }}
        </button>
      </div>
    </div>
  </div>

  <BaseDialog :show="confirmOpen" :title="t('tickets.settings.confirmTitle')" width="narrow" @close="confirmOpen = false">
    <p class="text-sm leading-6 text-gray-600 dark:text-gray-300">{{ t('tickets.settings.confirmDescription') }}</p>
    <p class="mt-3 rounded-lg border border-amber-200 bg-amber-50 px-3 py-2 text-sm font-medium text-amber-800 dark:border-amber-800 dark:bg-amber-950/35 dark:text-amber-300">
      {{ t('tickets.settings.pendingWarning', { count: pendingCount }) }}
    </p>
    <template #footer>
      <button type="button" class="btn btn-secondary" @click="confirmOpen = false">{{ t('common.cancel') }}</button>
      <button type="button" class="btn btn-danger" :disabled="saving" @click="persist">{{ t('tickets.settings.confirm') }}</button>
    </template>
  </BaseDialog>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import BaseDialog from '@/components/common/BaseDialog.vue'
import Toggle from '@/components/common/Toggle.vue'
import { ticketsAPI } from '@/api/tickets'
import { useAppStore } from '@/stores/app'
import { useTicketStore } from '@/stores/tickets'
import type { TicketConfig } from '@/types/ticket'

const { t } = useI18n()
const appStore = useAppStore()
const ticketStore = useTicketStore()
const draft = reactive<TicketConfig>({ user_center_enabled: true, accept_new_tickets: true })
const saved = ref<TicketConfig>({ user_center_enabled: true, accept_new_tickets: true })
const saving = ref(false)
const confirmOpen = ref(false)
const pendingCount = ref(0)
const changed = computed(() => draft.user_center_enabled !== saved.value.user_center_enabled || draft.accept_new_tickets !== saved.value.accept_new_tickets)

onMounted(async () => {
  try {
    const config = await ticketStore.fetchConfig(true, true)
    Object.assign(draft, config)
    saved.value = { ...config }
  } catch (error) {
    appStore.showError((error as { message?: string })?.message || t('common.error'))
  }
})

const save = async () => {
  if (saved.value.user_center_enabled && !draft.user_center_enabled) {
    try {
      const summary = await ticketsAPI.admin.summary()
      pendingCount.value = summary.open_count
      if (summary.open_count > 0) {
        confirmOpen.value = true
        return
      }
    } catch {
      pendingCount.value = ticketStore.adminSummary.open_count
      if (pendingCount.value > 0) {
        confirmOpen.value = true
        return
      }
    }
  }
  await persist()
}

const persist = async () => {
  saving.value = true
  try {
    const config = await ticketStore.updateConfig({ ...draft })
    saved.value = { ...config }
    Object.assign(draft, config)
    confirmOpen.value = false
    appStore.showSuccess(t('tickets.settings.saved'))
    await ticketStore.refresh(true)
  } catch (error) {
    appStore.showError((error as { message?: string })?.message || t('common.error'))
  } finally {
    saving.value = false
  }
}
</script>

