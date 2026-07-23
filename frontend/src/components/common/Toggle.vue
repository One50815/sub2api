<template>
  <button
    type="button"
    @click="toggle"
    class="toggle-control"
    :class="[
      modelValue ? 'toggle-control-checked' : 'toggle-control-unchecked',
      disabled && 'toggle-control-disabled'
    ]"
    :disabled="disabled"
    role="switch"
    :aria-checked="modelValue"
    :aria-disabled="disabled"
  >
    <span
      class="toggle-thumb"
      :class="modelValue ? 'toggle-thumb-checked' : 'toggle-thumb-unchecked'"
    />
  </button>
</template>

<script setup lang="ts">
const props = withDefaults(defineProps<{
  modelValue: boolean
  disabled?: boolean
}>(), {
  disabled: false
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
}>()

function toggle() {
  if (props.disabled) return
  emit('update:modelValue', !props.modelValue)
}
</script>

<style scoped>
.toggle-control {
  position: relative;
  display: inline-flex;
  width: 2rem;
  height: 1.15rem;
  flex-shrink: 0;
  align-items: center;
  padding: 0;
  overflow: visible;
  cursor: pointer;
  border: 1px solid transparent;
  border-radius: 9999px;
  outline: none;
  transition: background-color 150ms ease, border-color 150ms ease, box-shadow 150ms ease;
}

.toggle-control::after {
  content: '';
  position: absolute;
  inset: -0.5rem -0.75rem;
}

.toggle-control-checked {
  background: var(--omnio-primary-strong, #2563eb);
}

.toggle-control-unchecked {
  background: color-mix(in srgb, var(--omnio-foreground, #111827) 12%, var(--omnio-surface, #fff));
}

.toggle-control:focus-visible {
  border-color: var(--omnio-primary, #3b82f6);
  box-shadow: 0 0 0 3px color-mix(in srgb, var(--omnio-primary, #3b82f6) 28%, transparent);
}

.toggle-control-disabled {
  cursor: not-allowed;
  opacity: 0.5;
}

.toggle-thumb {
  pointer-events: none;
  display: block;
  width: 1rem;
  height: 1rem;
  flex-shrink: 0;
  border-radius: 9999px;
  background: var(--omnio-surface, #fff);
  box-shadow: 0 1px 2px rgba(15, 23, 42, 0.18);
  transition: transform 150ms ease, background-color 150ms ease;
}

.toggle-thumb-unchecked {
  transform: translateX(0);
}

.toggle-thumb-checked {
  transform: translateX(0.85rem);
}

:global(.dark) .toggle-control-unchecked {
  background: color-mix(in srgb, var(--omnio-foreground, #f8fafc) 20%, var(--omnio-surface, #111827));
}

:global(.dark) .toggle-thumb {
  background: var(--omnio-foreground, #f8fafc);
}

:global(.dark) .toggle-control-checked .toggle-thumb {
  background: #fff;
}

@media (prefers-reduced-motion: reduce) {
  .toggle-control,
  .toggle-thumb {
    transition-duration: 1ms;
  }
}
</style>
