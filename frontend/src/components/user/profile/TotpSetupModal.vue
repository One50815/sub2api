<template>
  <BaseDialog
    :show="true"
    :title="t('profile.totp.setupTitle')"
    width="normal"
    :close-on-click-outside="true"
    @close="emit('close')"
  >
    <div class="setup-content">
      <div class="step-heading">
        <div class="step-track" aria-hidden="true">
          <span v-for="index in 3" :key="index" :class="['step-dot', { active: step >= index - 1 }]">
            {{ index }}
          </span>
        </div>
        <p>{{ stepDescription }}</p>
      </div>

      <!-- Step 0: Identity Verification -->
      <div v-if="step === 0" class="setup-panel">
        <div v-if="methodLoading" class="loading-state" role="status">
          <Icon name="refresh" size="md" class="animate-spin" />
          <span>{{ t('common.loading') }}</span>
        </div>

        <template v-else>
          <div v-if="verificationMethod === 'email'" class="field-group">
            <label for="totp-setup-email-code" class="input-label">
              {{ t('profile.totp.emailCode') }}
            </label>
            <div class="verification-row">
              <input
                id="totp-setup-email-code"
                v-model="verifyForm.emailCode"
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
            <label for="totp-setup-password" class="input-label">
              {{ t('profile.currentPassword') }}
            </label>
            <input
              id="totp-setup-password"
              v-model="verifyForm.password"
              type="password"
              autocomplete="current-password"
              class="input"
              :placeholder="t('profile.totp.enterPassword')"
            />
          </div>
        </template>
      </div>

      <!-- Step 1: Show QR Code -->
      <div v-else-if="step === 1" class="setup-panel qr-panel">
        <template v-if="setupData">
          <div class="qr-card">
            <img :src="qrCodeDataUrl" alt="QR Code" class="qr-image" />
          </div>

          <div class="secret-card">
            <div class="min-w-0 flex-1">
              <p>{{ t('profile.totp.manualEntry') }}</p>
              <code>{{ setupData.secret }}</code>
            </div>
            <button
              type="button"
              class="copy-button"
              :aria-label="t('profile.totp.manualEntry')"
              @click="copySecret"
            >
              <Icon name="copy" size="sm" />
            </button>
          </div>
        </template>
      </div>

      <!-- Step 2: Verify Code -->
      <div v-else class="setup-panel verify-panel">
        <form id="totp-setup-verify-form" @submit.prevent="handleVerify">
          <label class="input-label mb-3 block text-center">
            {{ t('profile.totp.enterCode') }}
          </label>
          <div class="totp-code" role="group" :aria-label="t('profile.totp.enterCode')">
            <input
              v-for="(_, index) in 6"
              :key="index"
              :ref="(el) => setInputRef(el, index)"
              type="text"
              maxlength="1"
              inputmode="numeric"
              pattern="[0-9]"
              autocomplete="off"
              class="totp-cell"
              :class="{ 'totp-cell-group-end': index === 1 || index === 3 }"
              :aria-label="`${t('profile.totp.enterCode')} ${index + 1}`"
              :disabled="verifying"
              @input="handleCodeInput($event, index)"
              @keydown="handleKeydown($event, index)"
              @paste="handlePaste"
            />
          </div>
        </form>
      </div>
    </div>

    <template #footer>
      <div class="setup-actions">
        <button
          v-if="step < 2"
          type="button"
          class="btn btn-secondary"
          @click="emit('close')"
        >
          {{ t('common.cancel') }}
        </button>
        <button
          v-if="step === 0 && !methodLoading"
          type="button"
          class="btn btn-primary"
          :disabled="!canProceedFromVerify || setupLoading"
          @click="handleVerifyAndSetup"
        >
          <Icon v-if="setupLoading" name="refresh" size="sm" class="animate-spin" />
          {{ setupLoading ? t('common.loading') : t('common.next') }}
        </button>
        <button
          v-else-if="step === 1"
          type="button"
          class="btn btn-primary"
          :disabled="!setupData"
          @click="step = 2"
        >
          {{ t('common.next') }}
        </button>
        <template v-else-if="step === 2">
          <button type="button" class="btn btn-secondary" @click="step = 1">
            {{ t('common.back') }}
          </button>
          <button
            type="submit"
            form="totp-setup-verify-form"
            class="btn btn-primary"
            :disabled="verifying || code.join('').length !== 6"
          >
            <Icon v-if="verifying" name="refresh" size="sm" class="animate-spin" />
            {{ verifying ? t('common.verifying') : t('profile.totp.verify') }}
          </button>
        </template>
      </div>
    </template>
  </BaseDialog>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick, watch, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'
import BaseDialog from '@/components/common/BaseDialog.vue'
import { useAppStore } from '@/stores/app'
import { totpAPI } from '@/api'
import type { TotpSetupResponse } from '@/types'
import QRCode from 'qrcode'

const emit = defineEmits<{
  close: []
  success: []
}>()

const { t } = useI18n()
const appStore = useAppStore()

// Step: 0 = verify identity, 1 = QR code, 2 = verify TOTP code
const step = ref(0)
const methodLoading = ref(true)
const verificationMethod = ref<'email' | 'password'>('password')
const verifyForm = ref({ emailCode: '', password: '' })
const sendingCode = ref(false)
const codeCooldown = ref(0)
const cooldownTimer = ref<ReturnType<typeof setInterval> | null>(null)

const setupLoading = ref(false)
const setupData = ref<TotpSetupResponse | null>(null)
const verifying = ref(false)
const code = ref<string[]>(['', '', '', '', '', ''])
const inputRefs = ref<(HTMLInputElement | null)[]>([])
const qrCodeDataUrl = ref('')

const stepDescription = computed(() => {
  switch (step.value) {
    case 0:
      return verificationMethod.value === 'email'
        ? t('profile.totp.verifyEmailFirst')
        : t('profile.totp.verifyPasswordFirst')
    case 1:
      return t('profile.totp.setupStep1')
    case 2:
      return t('profile.totp.setupStep2')
    default:
      return ''
  }
})

const canProceedFromVerify = computed(() => {
  if (verificationMethod.value === 'email') {
    return verifyForm.value.emailCode.length === 6
  }
  return verifyForm.value.password.length > 0
})

// Generate QR code as base64 when setupData changes
watch(
  () => setupData.value?.qr_code_url,
  async (url) => {
    if (url) {
      try {
        qrCodeDataUrl.value = await QRCode.toDataURL(url, {
          width: 200,
          margin: 2,
          color: {
            dark: '#000000',
            light: '#ffffff'
          }
        })
      } catch (err) {
        console.error('Failed to generate QR code:', err)
      }
    }
  },
  { immediate: true }
)

const setInputRef = (el: any, index: number) => {
  inputRefs.value[index] = el as HTMLInputElement | null
}

const handleCodeInput = (event: Event, index: number) => {
  const input = event.target as HTMLInputElement
  const value = input.value.replace(/[^0-9]/g, '')
  code.value[index] = value

  if (value && index < 5) {
    nextTick(() => {
      inputRefs.value[index + 1]?.focus()
    })
  }
}

const handleKeydown = (event: KeyboardEvent, index: number) => {
  if (event.key === 'Backspace') {
    const input = event.target as HTMLInputElement
    // If current cell is empty and not the first, move to previous cell
    if (!input.value && index > 0) {
      event.preventDefault()
      inputRefs.value[index - 1]?.focus()
    }
    // Otherwise, let the browser handle the backspace naturally
    // The input event will sync code.value via handleCodeInput
  }
}

const handlePaste = (event: ClipboardEvent) => {
  event.preventDefault()
  const pastedData = event.clipboardData?.getData('text') || ''
  const digits = pastedData.replace(/[^0-9]/g, '').slice(0, 6).split('')

  // Update both the ref and the input elements
  digits.forEach((digit, index) => {
    code.value[index] = digit
    if (inputRefs.value[index]) {
      inputRefs.value[index]!.value = digit
    }
  })

  // Clear remaining inputs if pasted less than 6 digits
  for (let i = digits.length; i < 6; i++) {
    code.value[i] = ''
    if (inputRefs.value[i]) {
      inputRefs.value[i]!.value = ''
    }
  }

  const focusIndex = Math.min(digits.length, 5)
  nextTick(() => {
    inputRefs.value[focusIndex]?.focus()
  })
}

const copySecret = async () => {
  if (setupData.value) {
    try {
      await navigator.clipboard.writeText(setupData.value.secret)
      appStore.showSuccess(t('common.copied'))
    } catch {
      appStore.showError(t('common.copyFailed'))
    }
  }
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

const handleVerifyAndSetup = async () => {
  setupLoading.value = true

  try {
    const request = verificationMethod.value === 'email'
      ? { email_code: verifyForm.value.emailCode }
      : { password: verifyForm.value.password }

    setupData.value = await totpAPI.initiateSetup(request)
    step.value = 1
  } catch (err: any) {
    appStore.showError(err.response?.data?.message || t('profile.totp.setupFailed'))
  } finally {
    setupLoading.value = false
  }
}

const handleVerify = async () => {
  const totpCode = code.value.join('')
  if (totpCode.length !== 6 || !setupData.value) return

  verifying.value = true

  try {
    await totpAPI.enable({
      totp_code: totpCode,
      setup_token: setupData.value.setup_token
    })
    appStore.showSuccess(t('profile.totp.enableSuccess'))
    emit('success')
  } catch (err: any) {
    appStore.showError(err.response?.data?.message || t('profile.totp.verifyFailed'))
    code.value = ['', '', '', '', '', '']
    nextTick(() => {
      inputRefs.value[0]?.focus()
    })
  } finally {
    verifying.value = false
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
.setup-content,
.setup-panel,
.field-group {
  display: flex;
  flex-direction: column;
}

.setup-content {
  gap: 1.25rem;
}

.setup-panel,
.field-group {
  gap: 0.5rem;
}

.step-heading {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.625rem;
  color: var(--omnio-muted, #6b7280);
  font-size: 0.8125rem;
  line-height: 1.5;
  text-align: center;
}

.step-track {
  display: flex;
  align-items: center;
  gap: 0.375rem;
}

.step-dot {
  display: inline-flex;
  width: 1.5rem;
  height: 1.5rem;
  align-items: center;
  justify-content: center;
  border: 1px solid var(--omnio-border, #e5e7eb);
  border-radius: 0.5rem;
  color: var(--omnio-muted, #6b7280);
  background: var(--omnio-surface, #fff);
  font-size: 0.6875rem;
  font-weight: 600;
}

.step-dot.active {
  border-color: color-mix(in srgb, var(--omnio-primary, #3b82f6) 48%, var(--omnio-border, #e5e7eb));
  color: var(--omnio-primary-strong, #2563eb);
  background: color-mix(in srgb, var(--omnio-primary, #3b82f6) 6%, var(--omnio-surface, #fff));
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

.qr-panel {
  align-items: center;
  gap: 0.875rem;
}

.qr-card {
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid var(--omnio-border, #e5e7eb);
  border-radius: 0.625rem;
  padding: 0.75rem;
  background: #fff;
}

.qr-image {
  width: 12rem;
  height: 12rem;
}

.secret-card {
  display: flex;
  width: 100%;
  align-items: center;
  gap: 0.75rem;
  border: 1px solid var(--omnio-border, #e5e7eb);
  border-radius: 0.625rem;
  padding: 0.75rem;
  background: color-mix(in srgb, var(--omnio-foreground, #111827) 2.5%, var(--omnio-surface, #fff));
}

.secret-card p {
  margin-bottom: 0.25rem;
  color: var(--omnio-muted, #6b7280);
  font-size: 0.75rem;
  line-height: 1.4;
}

.secret-card code {
  display: block;
  overflow-wrap: anywhere;
  color: var(--omnio-foreground, #111827);
  font-size: 0.75rem;
  line-height: 1.5;
}

.copy-button {
  display: inline-flex;
  width: 2rem;
  height: 2rem;
  flex: 0 0 auto;
  align-items: center;
  justify-content: center;
  border-radius: 0.5rem;
  color: var(--omnio-muted, #6b7280);
  transition: color 140ms ease, background-color 140ms ease, box-shadow 140ms ease;
}

.copy-button:hover {
  color: var(--omnio-foreground, #111827);
  background: color-mix(in srgb, var(--omnio-foreground, #111827) 5%, transparent);
}

.copy-button:focus-visible {
  outline: none;
  box-shadow: 0 0 0 3px color-mix(in srgb, var(--omnio-primary, #3b82f6) 20%, transparent);
}

.verify-panel {
  padding: 0.5rem 0;
}

.totp-code {
  display: flex;
  justify-content: center;
  gap: 0.5rem;
}

.totp-cell {
  width: 2.5rem;
  height: 3rem;
  border: 1px solid var(--omnio-border, #d1d5db);
  border-radius: 0.625rem;
  color: var(--omnio-foreground, #111827);
  background: var(--omnio-surface, #fff);
  font-size: 1rem;
  line-height: 1;
  font-weight: 600;
  text-align: center;
  outline: none;
  transition: border-color 140ms ease, box-shadow 140ms ease, background-color 140ms ease;
}

.totp-cell:focus {
  border-color: var(--omnio-primary, #3b82f6);
  box-shadow: 0 0 0 3px color-mix(in srgb, var(--omnio-primary, #3b82f6) 20%, transparent);
}

.totp-cell:disabled {
  cursor: wait;
  opacity: 0.58;
}

.totp-cell-group-end {
  margin-right: 0.25rem;
}

.setup-actions {
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
  .setup-actions .btn {
    width: 100%;
  }

  .setup-actions {
    flex-direction: column-reverse;
  }

  .totp-code {
    gap: 0.35rem;
  }

  .totp-cell {
    width: 2.2rem;
  }

  .totp-cell-group-end {
    margin-right: 0.1rem;
  }
}

@media (prefers-reduced-motion: reduce) {
  .copy-button,
  .totp-cell {
    transition-duration: 1ms;
  }
}
</style>
