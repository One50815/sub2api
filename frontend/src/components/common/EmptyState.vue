<template>
  <div class="empty-state">
    <!-- Icon -->
    <div class="empty-state-media">
      <slot name="icon">
        <component v-if="icon" :is="icon" class="empty-state-icon" aria-hidden="true" />
        <svg
          v-else
          class="empty-state-icon"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          stroke-width="1.75"
          aria-hidden="true"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4"
          />
        </svg>
      </slot>
    </div>

    <!-- Title -->
    <h3 class="empty-state-title">
      {{ displayTitle }}
    </h3>

    <!-- Description -->
    <p v-if="description" class="empty-state-description">
      {{ description }}
    </p>

    <!-- Action -->
    <div v-if="actionText || $slots.action" class="empty-state-action">
      <slot name="action">
        <component
          :is="actionTo ? 'RouterLink' : 'button'"
          v-if="actionText"
          :to="actionTo"
          @click="!actionTo && $emit('action')"
          class="btn btn-primary"
        >
          <Icon v-if="actionIcon" name="plus" size="md" class="mr-2" />
          {{ actionText }}
        </component>
      </slot>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import type { Component } from 'vue'
import Icon from '@/components/icons/Icon.vue'

const { t } = useI18n()

interface Props {
  icon?: Component | string
  title?: string
  description?: string
  actionText?: string
  actionTo?: string | object
  actionIcon?: boolean
  message?: string
}

const props = withDefaults(defineProps<Props>(), {
  description: '',
  actionIcon: true
})

const displayTitle = computed(() => props.title || t('common.noData'))

defineEmits(['action'])
</script>

<style scoped>
.empty-state {
  min-height: 18.75rem;
  padding: 2rem 1.5rem;
  gap: 0;
}

.empty-state-media {
  display: inline-flex;
  width: 2.5rem;
  height: 2.5rem;
  align-items: center;
  justify-content: center;
  margin-bottom: 0.75rem;
  border-radius: 0.55rem;
  color: var(--omnio-muted, #6b7280);
  background: color-mix(in srgb, var(--omnio-foreground, #111827) 5%, transparent);
}

.empty-state-media :deep(svg),
.empty-state-icon {
  width: 1.5rem !important;
  height: 1.5rem !important;
  margin: 0 !important;
  color: currentColor !important;
}

.empty-state-title {
  margin: 0 !important;
  color: var(--omnio-foreground, #111827) !important;
  font-size: 0.9rem !important;
  line-height: 1.35 !important;
  font-weight: 560 !important;
}

.empty-state-description {
  max-width: 24rem;
  margin-top: 0.35rem;
  font-size: 0.8rem;
  line-height: 1.55;
  text-wrap: balance;
}

.empty-state-action {
  margin-top: 1rem;
}

@media (max-width: 640px) {
  .empty-state {
    min-height: 15rem;
    padding: 1.5rem 1rem;
  }
}
</style>
