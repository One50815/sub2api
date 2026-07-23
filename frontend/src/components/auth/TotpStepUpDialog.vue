<template>
  <BaseDialog
    :show="controller.visible.value"
    :title="t('stepUp.title')"
    width="narrow"
    :z-index="60"
    :close-on-click-outside="true"
    :close-on-escape="!verifying"
    :show-close-button="!verifying"
    @close="handleCancel"
  >
    <div class="totp-content">
      <div class="totp-intro">
        <span class="totp-icon" aria-hidden="true">
          <Icon name="lock" size="lg" />
        </span>
        <p class="totp-hint">{{ t('stepUp.hint') }}</p>
      </div>

      <input
        ref="hiddenOtpInputRef"
        type="text"
        inputmode="numeric"
        autocomplete="one-time-code"
        maxlength="6"
        class="pointer-events-none absolute left-0 top-0 h-px w-px opacity-0"
        aria-hidden="true"
        tabindex="-1"
        @input="handleHiddenOtpInput"
      />

      <div class="totp-code" role="group" :aria-label="t('stepUp.title')">
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
          :aria-label="`${t('stepUp.title')} ${index + 1}`"
          :disabled="verifying"
          @input="handleCodeInput($event, index)"
          @keydown="handleKeydown($event, index)"
          @paste="handlePaste"
        />
      </div>

      <div v-if="verifying" class="totp-status" role="status">
        <Icon name="refresh" size="sm" class="animate-spin" />
        <span>{{ t('common.verifying') }}</span>
      </div>
    </div>

    <template #footer>
      <div class="totp-actions">
        <button
          type="button"
          class="btn btn-secondary"
          :disabled="verifying"
          @click="handleCancel"
        >
          {{ t('common.cancel') }}
        </button>
      </div>
    </template>
  </BaseDialog>
</template>

<script setup lang="ts">
import { ref, watch, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores'
import { totpAPI } from '@/api'
import type { StepUpController } from '@/composables/useStepUp'
import BaseDialog from '@/components/common/BaseDialog.vue'
import Icon from '@/components/icons/Icon.vue'

const props = defineProps<{
  controller: StepUpController
}>()

const { t } = useI18n()
const appStore = useAppStore()

const verifying = ref(false)
const code = ref<string[]>(['', '', '', '', '', ''])
const inputRefs = ref<(HTMLInputElement | null)[]>([])
const hiddenOtpInputRef = ref<HTMLInputElement | null>(null)

// Focus the first cell whenever the dialog opens.
watch(
  () => props.controller.visible.value,
  (open) => {
    if (open) {
      resetInputs()
      nextTick(() => inputRefs.value[0]?.focus())
    }
  }
)

// Auto-submit once 6 digits are entered.
watch(
  () => code.value.join(''),
  (newCode) => {
    if (newCode.length === 6 && !verifying.value) {
      submit(newCode)
    }
  }
)

async function submit(otp: string) {
  verifying.value = true
  try {
    await totpAPI.stepUp(otp)
    verifying.value = false
    resetInputs()
    props.controller.onVerified()
  } catch (err: any) {
    verifying.value = false
    appStore.showError(err?.message || t('stepUp.verifyFailed'))
    resetInputs()
    nextTick(() => inputRefs.value[0]?.focus())
  }
}

function resetInputs() {
  code.value = ['', '', '', '', '', '']
  inputRefs.value.forEach((input) => {
    if (input) input.value = ''
  })
  if (hiddenOtpInputRef.value) hiddenOtpInputRef.value.value = ''
}

function handleCancel() {
  if (verifying.value) return
  props.controller.onCancel()
}

const setInputRef = (el: any, index: number) => {
  inputRefs.value[index] = el as HTMLInputElement | null
}

const handleCodeInput = (event: Event, index: number) => {
  const input = event.target as HTMLInputElement
  const value = input.value.replace(/[^0-9]/g, '')
  code.value[index] = value
  if (value && index < 5) {
    nextTick(() => inputRefs.value[index + 1]?.focus())
  }
}

const handleHiddenOtpInput = (event: Event) => {
  const input = event.target as HTMLInputElement
  const digits = input.value.replace(/[^0-9]/g, '').slice(0, 6).split('')
  for (let i = 0; i < 6; i++) {
    code.value[i] = digits[i] || ''
    if (inputRefs.value[i]) inputRefs.value[i]!.value = digits[i] || ''
  }
}

const handleKeydown = (event: KeyboardEvent, index: number) => {
  if (event.key === 'Backspace') {
    const input = event.target as HTMLInputElement
    if (!input.value && index > 0) {
      event.preventDefault()
      inputRefs.value[index - 1]?.focus()
    }
  }
}

const handlePaste = (event: ClipboardEvent) => {
  event.preventDefault()
  const pastedData = event.clipboardData?.getData('text') || ''
  const digits = pastedData.replace(/[^0-9]/g, '').slice(0, 6).split('')
  for (let i = 0; i < 6; i++) {
    code.value[i] = digits[i] || ''
    if (inputRefs.value[i]) inputRefs.value[i]!.value = digits[i] || ''
  }
  const focusIndex = Math.min(digits.length, 5)
  nextTick(() => inputRefs.value[focusIndex]?.focus())
}
</script>

<style scoped>
.totp-content {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.totp-intro {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.625rem;
  text-align: center;
}

.totp-icon {
  display: inline-flex;
  width: 2.5rem;
  height: 2.5rem;
  align-items: center;
  justify-content: center;
  border: 1px solid color-mix(in srgb, var(--omnio-primary, #3b82f6) 18%, var(--omnio-border, #e5e7eb));
  border-radius: 0.625rem;
  color: var(--omnio-primary, #3b82f6);
  background: color-mix(in srgb, var(--omnio-primary, #3b82f6) 7%, var(--omnio-surface, #fff));
}

.totp-hint,
.totp-status {
  color: var(--omnio-muted, #6b7280);
  font-size: 0.8125rem;
  line-height: 1.5;
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
  background: color-mix(in srgb, var(--omnio-foreground, #111827) 4%, var(--omnio-surface, #fff));
}

.totp-cell-group-end {
  margin-right: 0.25rem;
}

.totp-status {
  display: flex;
  min-height: 1.25rem;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
}

.totp-actions {
  display: flex;
  width: 100%;
  justify-content: flex-end;
}

@media (max-width: 380px) {
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
  .totp-cell {
    transition-duration: 1ms;
  }
}
</style>
