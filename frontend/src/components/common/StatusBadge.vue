<template>
  <span :class="['status-badge', variantClass]" :title="label">
    <span class="status-badge-label">{{ label }}</span>
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  status: string
  label: string
}>()

const variantClass = computed(() => {
  switch (props.status) {
    case 'active':
    case 'success':
      return 'status-badge-success'
    case 'warning':
    case 'pending':
      return 'status-badge-warning'
    case 'disabled':
    case 'inactive':
      return 'status-badge-neutral'
    case 'error':
    case 'danger':
      return 'status-badge-danger'
    case 'info':
      return 'status-badge-info'
    default:
      return 'status-badge-neutral'
  }
})
</script>

<style scoped>
.status-badge {
  display: inline-flex;
  width: fit-content;
  max-width: 100%;
  min-width: 0;
  height: 1.25rem;
  align-items: center;
  padding: 0 0.4rem;
  overflow: hidden;
  border: 1px solid transparent;
  border-radius: 9999px;
  font-size: 0.75rem;
  line-height: 1;
  font-weight: 560;
  letter-spacing: 0;
  white-space: nowrap;
}

.status-badge-label {
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
}

.status-badge-success {
  color: #047857;
  border-color: rgba(16, 185, 129, 0.18);
  background: rgba(16, 185, 129, 0.1);
}

.status-badge-warning {
  color: #b45309;
  border-color: rgba(245, 158, 11, 0.2);
  background: rgba(245, 158, 11, 0.11);
}

.status-badge-danger {
  color: #dc2626;
  border-color: rgba(239, 68, 68, 0.18);
  background: rgba(239, 68, 68, 0.09);
}

.status-badge-info {
  color: #2563eb;
  border-color: rgba(59, 130, 246, 0.18);
  background: rgba(59, 130, 246, 0.09);
}

.status-badge-neutral {
  color: var(--omnio-muted, #6b7280);
  border-color: var(--omnio-border, #e5e7eb);
  background: color-mix(in srgb, var(--omnio-foreground, #111827) 4%, transparent);
}

:global(.dark) .status-badge-success {
  color: #6ee7b7;
}

:global(.dark) .status-badge-warning {
  color: #fcd34d;
}

:global(.dark) .status-badge-danger {
  color: #fca5a5;
}

:global(.dark) .status-badge-info {
  color: #93c5fd;
}
</style>
