<template>
  <BaseDialog :show="show" :title="title" width="narrow" @close="handleCancel">
    <div class="confirm-content">
      <p class="confirm-message">{{ message }}</p>
      <slot></slot>
    </div>

    <template #footer>
      <div class="confirm-actions">
        <button
          @click="handleCancel"
          type="button"
          class="confirm-button confirm-button-cancel"
        >
          {{ cancelText }}
        </button>
        <button
          @click="handleConfirm"
          type="button"
          :class="[
            'confirm-button confirm-button-primary',
            danger ? 'confirm-button-danger' : 'confirm-button-default'
          ]"
        >
          {{ confirmText }}
        </button>
      </div>
    </template>
  </BaseDialog>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import BaseDialog from './BaseDialog.vue'

const { t } = useI18n()

interface Props {
  show: boolean
  title: string
  message: string
  confirmText?: string
  cancelText?: string
  danger?: boolean
}

interface Emits {
  (e: 'confirm'): void
  (e: 'cancel'): void
}

const props = withDefaults(defineProps<Props>(), {
  danger: false
})

const confirmText = computed(() => props.confirmText || t('common.confirm'))
const cancelText = computed(() => props.cancelText || t('common.cancel'))

const emit = defineEmits<Emits>()

const handleConfirm = () => {
  emit('confirm')
}

const handleCancel = () => {
  emit('cancel')
}
</script>

<style scoped>
.confirm-content {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.confirm-message {
  color: var(--omnio-muted, #6b7280);
  font-size: 0.875rem;
  line-height: 1.55;
  text-wrap: pretty;
}

.confirm-actions {
  display: flex;
  width: 100%;
  flex-direction: column-reverse;
  gap: 0.5rem;
}

.confirm-button {
  display: inline-flex;
  min-height: 2rem;
  align-items: center;
  justify-content: center;
  border: 1px solid transparent;
  border-radius: 0.5rem;
  padding: 0.38rem 0.75rem;
  font-size: 0.8rem;
  line-height: 1.2;
  font-weight: 550;
  outline: none;
  transition: color 140ms ease, background-color 140ms ease, border-color 140ms ease, box-shadow 140ms ease;
}

.confirm-button:focus-visible {
  box-shadow: 0 0 0 3px color-mix(in srgb, var(--dialog-focus-color, var(--omnio-primary, #3b82f6)) 25%, transparent);
}

.confirm-button-cancel {
  color: var(--omnio-foreground, #111827);
  border-color: var(--omnio-border, #e5e7eb);
  background: var(--omnio-surface, #fff);
}

.confirm-button-cancel:hover {
  background: color-mix(in srgb, var(--omnio-foreground, #111827) 5%, var(--omnio-surface, #fff));
}

.confirm-button-primary {
  color: #fff;
}

.confirm-button-default {
  --dialog-focus-color: var(--omnio-primary, #3b82f6);
  border-color: color-mix(in srgb, var(--omnio-primary-strong, #2563eb) 82%, #000);
  background: var(--omnio-primary-strong, #2563eb);
}

.confirm-button-default:hover {
  background: color-mix(in srgb, var(--omnio-primary-strong, #2563eb) 90%, #000);
}

.confirm-button-danger {
  --dialog-focus-color: #ef4444;
  border-color: #dc2626;
  background: #dc2626;
}

.confirm-button-danger:hover {
  background: #b91c1c;
}

@media (min-width: 640px) {
  .confirm-actions {
    flex-direction: row;
    justify-content: flex-end;
  }
}

@media (prefers-reduced-motion: reduce) {
  .confirm-button {
    transition-duration: 1ms;
  }
}
</style>
