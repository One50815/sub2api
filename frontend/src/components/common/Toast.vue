<template>
  <Teleport to="body">
    <div
      class="toast-viewport"
      aria-live="polite"
      aria-atomic="true"
    >
      <TransitionGroup
        name="toast-item"
      >
        <div
          v-for="toast in toasts"
          :key="toast.id"
          :class="['toast-card', `toast-card-${toast.type}`]"
          :role="toast.type === 'error' ? 'alert' : 'status'"
          :aria-live="toast.type === 'error' ? 'assertive' : 'polite'"
        >
          <div class="toast-card-content">
            <div class="toast-card-row">
              <!-- Icon -->
              <div class="toast-icon">
                <Icon
                  :name="getToastIconName(toast.type)"
                  size="sm"
                  aria-hidden="true"
                />
              </div>

              <!-- Content -->
              <div class="toast-copy">
                <p v-if="toast.title" class="toast-title">
                  {{ toast.title }}
                </p>
                <p :class="['toast-message', toast.title && 'toast-message-with-title']">
                  {{ toast.message }}
                </p>
              </div>

              <!-- Close button -->
              <button
                @click="removeToast(toast.id)"
                type="button"
                class="toast-close"
                aria-label="Close notification"
              >
                <Icon name="x" size="sm" />
              </button>
            </div>
          </div>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import Icon from '@/components/icons/Icon.vue'
import { useAppStore } from '@/stores/app'

const appStore = useAppStore()

const toasts = computed(() => appStore.toasts)

const getToastIconName = (type: string): 'checkCircle' | 'xCircle' | 'exclamationTriangle' | 'infoCircle' => {
  switch (type) {
    case 'success':
      return 'checkCircle'
    case 'error':
      return 'xCircle'
    case 'warning':
      return 'exclamationTriangle'
    case 'info':
    default:
      return 'infoCircle'
  }
}

const removeToast = (id: string) => {
  appStore.hideToast(id)
}
</script>

<style scoped>
.toast-viewport {
  pointer-events: none;
  position: fixed;
  top: 1rem;
  left: 50%;
  z-index: 9999;
  display: flex;
  width: min(26rem, calc(100vw - 2rem));
  flex-direction: column;
  gap: 0.5rem;
  transform: translateX(-50%);
}

.toast-card {
  pointer-events: auto;
  width: 100%;
  overflow: hidden;
  border: 1px solid var(--omnio-border, #e5e7eb);
  border-radius: 0.65rem;
  color: var(--omnio-foreground, #111827);
  background: var(--omnio-surface, #fff);
  box-shadow: 0 10px 30px rgba(15, 23, 42, 0.12), 0 2px 7px rgba(15, 23, 42, 0.07);
}

.toast-card-success {
  border-color: color-mix(in srgb, #10b981 35%, var(--omnio-border, #e5e7eb));
  background: color-mix(in srgb, #10b981 10%, var(--omnio-surface, #fff));
}

.toast-card-error {
  border-color: color-mix(in srgb, #ef4444 35%, var(--omnio-border, #e5e7eb));
  background: color-mix(in srgb, #ef4444 9%, var(--omnio-surface, #fff));
}

.toast-card-warning {
  border-color: color-mix(in srgb, #f59e0b 38%, var(--omnio-border, #e5e7eb));
  background: color-mix(in srgb, #f59e0b 11%, var(--omnio-surface, #fff));
}

.toast-card-info {
  border-color: color-mix(in srgb, #3b82f6 35%, var(--omnio-border, #e5e7eb));
  background: color-mix(in srgb, #3b82f6 9%, var(--omnio-surface, #fff));
}

.toast-card-content {
  padding: 0.75rem;
}

.toast-card-row {
  display: flex;
  align-items: flex-start;
  gap: 0.625rem;
}

.toast-icon {
  display: inline-flex;
  width: 1.25rem;
  height: 1.25rem;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  margin-top: 0.05rem;
}

.toast-card-success .toast-icon {
  color: #059669;
}

.toast-card-error .toast-icon {
  color: #dc2626;
}

.toast-card-warning .toast-icon {
  color: #d97706;
}

.toast-card-info .toast-icon {
  color: #2563eb;
}

.toast-copy {
  min-width: 0;
  flex: 1;
}

.toast-title {
  font-size: 0.8rem;
  line-height: 1.35;
  font-weight: 600;
}

.toast-message {
  overflow-wrap: anywhere;
  color: var(--omnio-foreground, #111827);
  font-size: 0.8rem;
  line-height: 1.45;
}

.toast-message-with-title {
  margin-top: 0.2rem;
  color: var(--omnio-muted, #6b7280);
}

.toast-close {
  display: inline-flex;
  width: 1.75rem;
  height: 1.75rem;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  margin: -0.25rem -0.25rem -0.25rem 0;
  border-radius: 0.45rem;
  color: var(--omnio-muted, #6b7280);
  outline: none;
  transition: color 140ms ease, background-color 140ms ease, box-shadow 140ms ease;
}

.toast-close:hover {
  color: var(--omnio-foreground, #111827);
  background: color-mix(in srgb, var(--omnio-foreground, #111827) 6%, transparent);
}

.toast-close:focus-visible {
  box-shadow: 0 0 0 3px color-mix(in srgb, var(--omnio-primary, #3b82f6) 24%, transparent);
}

:global(.dark) .toast-card {
  box-shadow: 0 12px 34px rgba(0, 0, 0, 0.38), 0 2px 8px rgba(0, 0, 0, 0.24);
}

:global(.dark) .toast-card-success .toast-icon {
  color: #6ee7b7;
}

:global(.dark) .toast-card-error .toast-icon {
  color: #fca5a5;
}

:global(.dark) .toast-card-warning .toast-icon {
  color: #fcd34d;
}

:global(.dark) .toast-card-info .toast-icon {
  color: #93c5fd;
}

.toast-item-enter-active,
.toast-item-leave-active {
  transition: opacity 160ms ease, transform 160ms ease;
}

.toast-item-enter-from,
.toast-item-leave-to {
  opacity: 0;
  transform: translateY(-0.5rem) scale(0.98);
}

@media (max-width: 640px) {
  .toast-viewport {
    top: 0.75rem;
    width: calc(100vw - 1.5rem);
  }
}

@media (prefers-reduced-motion: reduce) {
  .toast-item-enter-active,
  .toast-item-leave-active,
  .toast-close {
    transition-duration: 1ms;
  }
}
</style>
