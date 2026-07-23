<template>
  <div class="overflow-hidden rounded-lg border border-gray-200 bg-white dark:border-dark-700 dark:bg-dark-800">
    <div class="border-b border-gray-200 px-4 py-3.5 dark:border-dark-700 sm:px-5">
      <h2 class="text-base font-semibold text-gray-950 dark:text-white">
        {{ t('profile.totp.title') }}
      </h2>
      <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
        {{ t('profile.totp.description') }}
      </p>
    </div>
    <div class="p-4 sm:p-5">
      <!-- Loading state -->
      <div v-if="loading" class="flex items-center justify-center py-8">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-500"></div>
      </div>

      <!-- Feature disabled globally -->
      <div v-else-if="status && !status.feature_enabled" class="flex items-center gap-4 py-4">
        <div class="flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-lg bg-gray-100 text-gray-400 dark:bg-dark-700">
          <Icon name="exclamationTriangle" size="md" />
        </div>
        <div>
          <p class="font-medium text-gray-700 dark:text-gray-300">
            {{ t('profile.totp.featureDisabled') }}
          </p>
          <p class="text-sm text-gray-500 dark:text-gray-400">
            {{ t('profile.totp.featureDisabledHint') }}
          </p>
        </div>
      </div>

      <!-- 2FA Enabled -->
      <div v-else-if="status?.enabled" class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div class="flex items-center gap-4">
          <div class="flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-lg bg-emerald-50 text-emerald-600 dark:bg-emerald-500/10 dark:text-emerald-400">
            <Icon name="shield" size="md" />
          </div>
          <div>
            <p class="font-medium text-gray-900 dark:text-white">
              {{ t('profile.totp.enabled') }}
            </p>
            <p v-if="status.enabled_at" class="text-sm text-gray-500 dark:text-gray-400">
              {{ t('profile.totp.enabledAt') }}: {{ formatDate(status.enabled_at) }}
            </p>
          </div>
        </div>
        <button
          type="button"
          class="btn btn-outline-danger w-full sm:w-auto"
          @click="showDisableDialog = true"
        >
          {{ t('profile.totp.disable') }}
        </button>
      </div>

      <!-- 2FA Not Enabled -->
      <div v-else class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div class="flex items-center gap-4">
          <div class="flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-lg bg-gray-100 text-gray-400 dark:bg-dark-700">
            <Icon name="shield" size="md" />
          </div>
          <div>
            <p class="font-medium text-gray-700 dark:text-gray-300">
              {{ t('profile.totp.notEnabled') }}
            </p>
            <p class="text-sm text-gray-500 dark:text-gray-400">
              {{ t('profile.totp.notEnabledHint') }}
            </p>
          </div>
        </div>
        <button
          type="button"
          class="btn btn-primary w-full sm:w-auto"
          @click="showSetupModal = true"
        >
          {{ t('profile.totp.enable') }}
        </button>
      </div>
    </div>

    <!-- Setup Modal -->
    <TotpSetupModal
      v-if="showSetupModal"
      @close="showSetupModal = false"
      @success="handleSetupSuccess"
    />

    <!-- Disable Dialog -->
    <TotpDisableDialog
      v-if="showDisableDialog"
      @close="showDisableDialog = false"
      @success="handleDisableSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { totpAPI } from '@/api'
import Icon from '@/components/icons/Icon.vue'
import type { TotpStatus } from '@/types'
import TotpSetupModal from './TotpSetupModal.vue'
import TotpDisableDialog from './TotpDisableDialog.vue'

const { t } = useI18n()

const loading = ref(true)
const status = ref<TotpStatus | null>(null)
const showSetupModal = ref(false)
const showDisableDialog = ref(false)

const loadStatus = async () => {
  loading.value = true
  try {
    status.value = await totpAPI.getStatus()
  } catch (error) {
    console.error('Failed to load TOTP status:', error)
  } finally {
    loading.value = false
  }
}

const handleSetupSuccess = () => {
  showSetupModal.value = false
  loadStatus()
}

const handleDisableSuccess = () => {
  showDisableDialog.value = false
  loadStatus()
}

const formatDate = (timestamp: number) => {
  // Backend returns Unix timestamp in seconds, convert to milliseconds
  const date = new Date(timestamp * 1000)
  return date.toLocaleDateString(undefined, {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

onMounted(() => {
  loadStatus()
})
</script>
