<template>
  <BaseDialog
    :show="true"
    :title="t('profile.totp.disableTitle')"
    width="narrow"
    :close-on-click-outside="true"
    @close="handleClose"
  >
    <div class="disable-content">
      <div class="warning-note">
        <span class="warning-icon" aria-hidden="true">
          <Icon name="exclamationTriangle" size="md" />
        </span>
        <p>{{ t('profile.totp.disableWarning') }}</p>
      </div>

      <div v-if="methodLoading" class="loading-state" role="status">
        <Icon name="refresh" size="md" class="animate-spin" />
        <span>{{ t('common.loading') }}</span>
      </div>

      <form v-else id="totp-disable-form" class="verification-form" @submit.prevent="handleDisable">
        <div v-if="verificationMethod === 'email'" class="field-group">
          <label for="totp-disable-email-code" class="input-label">
            {{ t('profile.totp.emailCode') }}
          </label>
          <div class="verification-row">
            <input
              id="totp-disable-email-code"
              v-model="form.emailCode"
              type="text"
              maxlength="6"
              inputmode="numeric"
              autocomplete="one-time-code"
              class="input min-w-0 flex-1"
              :placeholder="t('profile.totp.enterEmailCode')"
            />
            <button
              type="button"
              class="btn btn-secondary whitespace-nowrap"
              :disabled="sendingCode || codeCooldown > 0"
              @click="handleSendCode"
            >
              {{ codeCooldown > 0 ? `${codeCooldown}s` : (sendingCode ? t('common.sending') : t('profile.totp.sendCode')) }}
            </button>
          </div>
        </div>

        <div v-else class="field-group">
          <label for="totp-disable-password" class="input-label">
            {{ t('profile.currentPassword') }}
          </label>
          <input
            id="totp-disable-password"
            v-model="form.password"
            type="password"
            autocomplete="current-password"
            class="input"
            :placeholder="t('profile.totp.enterPassword')"
          />
        </div>
      </form>
    </div>

    <template #footer>
      <div class="dialog-actions">
        <button type="button" class="btn btn-secondary" @click="handleClose">
          {{ t('common.cancel') }}
        </button>
        <button
          v-if="!methodLoading"
          type="submit"
          form="totp-disable-form"
          class="btn btn-danger"
          :disabled="loading || !canSubmit"
        >
          <Icon v-if="loading" name="refresh" size="sm" class="animate-spin" />
          {{ loading ? t('common.processing') : t('profile.totp.confirmDisable') }}
        </button>
      </div>
    </template>
  </BaseDialog>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'
import BaseDialog from '@/components/common/BaseDialog.vue'
import { useAppStore } from '@/stores/app'
import { totpAPI } from '@/api'

const emit = defineEmits<{
  close: []
  success: []
}>()

const { t } = useI18n()
const appStore = useAppStore()

const methodLoading = ref(true)
const verificationMethod = ref<'email' | 'password'>('password')
const loading = ref(false)
const sendingCode = ref(false)
const codeCooldown = ref(0)
const cooldownTimer = ref<ReturnType<typeof setInterval> | null>(null)
const form = ref({
  emailCode: '',
  password: ''
})

const canSubmit = computed(() => {
  if (verificationMethod.value === 'email') {
    return form.value.emailCode.length === 6
  }
  return form.value.password.length > 0
})

const handleClose = () => {
  emit('close')
}

const loadVerificationMethod = async () => {
  methodLoading.value = true
  try {
    const method = await totpAPI.getVerificationMethod()
    verificationMethod.value = method.method
  } catch (err: any) {
    appStore.showError(err.response?.data?.message || t('common.error'))
    emit('close')
  } finally {
    methodLoading.value = false
  }
}

const handleSendCode = async () => {
  sendingCode.value = true
  try {
    await totpAPI.sendVerifyCode()
    appStore.showSuccess(t('profile.totp.codeSent'))
    // Start cooldown
    codeCooldown.value = 60
    if (cooldownTimer.value) {
      clearInterval(cooldownTimer.value)
      cooldownTimer.value = null
    }
    cooldownTimer.value = setInterval(() => {
      codeCooldown.value--
      if (codeCooldown.value <= 0) {
        if (cooldownTimer.value) {
          clearInterval(cooldownTimer.value)
          cooldownTimer.value = null
        }
      }
    }, 1000)
  } catch (err: any) {
    appStore.showError(err.response?.data?.message || t('profile.totp.sendCodeFailed'))
  } finally {
    sendingCode.value = false
  }
}

const handleDisable = async () => {
  if (!canSubmit.value) return

  loading.value = true

  try {
    const request = verificationMethod.value === 'email'
      ? { email_code: form.value.emailCode }
      : { password: form.value.password }

    await totpAPI.disable(request)
    appStore.showSuccess(t('profile.totp.disableSuccess'))
    emit('success')
  } catch (err: any) {
    appStore.showError(err.response?.data?.message || t('profile.totp.disableFailed'))
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadVerificationMethod()
})

onUnmounted(() => {
  if (cooldownTimer.value) {
    clearInterval(cooldownTimer.value)
    cooldownTimer.value = null
  }
})
</script>

<style scoped>
.disable-content,
.verification-form,
.field-group {
  display: flex;
  flex-direction: column;
}

.disable-content {
  gap: 1.25rem;
}

.verification-form,
.field-group {
  gap: 0.5rem;
}

.warning-note {
  display: flex;
  align-items: flex-start;
  gap: 0.75rem;
  border: 1px solid color-mix(in srgb, #dc2626 18%, var(--omnio-border, #e5e7eb));
  border-radius: 0.625rem;
  padding: 0.75rem;
  color: var(--omnio-muted, #6b7280);
  background: color-mix(in srgb, #dc2626 4%, var(--omnio-surface, #fff));
  font-size: 0.8125rem;
  line-height: 1.55;
}

.warning-icon {
  display: inline-flex;
  flex: 0 0 auto;
  color: #dc2626;
}

.loading-state {
  display: flex;
  min-height: 5rem;
  align-items: center;
  justify-content: center;
  gap: 0.625rem;
  color: var(--omnio-muted, #6b7280);
  font-size: 0.8125rem;
}

.verification-row {
  display: flex;
  gap: 0.5rem;
}

.dialog-actions {
  display: flex;
  width: 100%;
  justify-content: flex-end;
  gap: 0.5rem;
}

@media (max-width: 460px) {
  .verification-row {
    flex-direction: column;
  }

  .verification-row .btn,
  .dialog-actions .btn {
    width: 100%;
  }

  .dialog-actions {
    flex-direction: column-reverse;
  }
}
</style>
