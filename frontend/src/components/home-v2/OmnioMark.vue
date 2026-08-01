<template>
  <svg
    class="ov2-mark"
    viewBox="0 0 48 48"
    role="img"
    :aria-label="label"
  >
    <path class="ov2-mark-line" d="M9 11h6l7 9" />
    <path class="ov2-mark-line" d="M9 24h10" />
    <path class="ov2-mark-line" d="M9 37h6l7-9" />
    <path class="ov2-mark-line ov2-mark-line-out" d="M34 24h7" />
    <rect class="ov2-mark-node" x="4" y="7" width="8" height="8" rx="1" />
    <rect class="ov2-mark-node" x="4" y="20" width="8" height="8" rx="1" />
    <rect class="ov2-mark-node" x="4" y="33" width="8" height="8" rx="1" />
    <path class="ov2-mark-core" d="M21 15h10l6 9-6 9H21l-6-9z" />
    <rect class="ov2-mark-output" x="39" y="20" width="8" height="8" rx="1" />
    <rect class="ov2-mark-cut" x="23" y="21" width="6" height="6" rx="1" />
  </svg>
</template>

<script setup lang="ts">
withDefaults(defineProps<{
  label?: string
}>(), {
  label: 'Omnio'
})
</script>
