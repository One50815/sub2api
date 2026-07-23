<template>
  <div class="pagination-shell">
    <div class="mobile-pagination">
      <!-- Mobile pagination -->
      <button
        @click="goToPage(page - 1)"
        :disabled="page <= 1"
        type="button"
        class="pagination-button pagination-icon-button"
        :aria-label="t('pagination.previous')"
      >
        <Icon name="chevronLeft" size="md" />
      </button>
      <span class="mobile-page-label" aria-live="polite">
        {{ t('pagination.pageOf', { page, total: displayTotalPages }) }}
      </span>
      <button
        @click="goToPage(page + 1)"
        :disabled="page >= displayTotalPages"
        type="button"
        class="pagination-button pagination-icon-button"
        :aria-label="t('pagination.next')"
      >
        <Icon name="chevronRight" size="md" />
      </button>
    </div>

    <div class="desktop-pagination">
      <!-- Desktop pagination info -->
      <div class="pagination-meta">
        <p class="pagination-results">
          {{ t('pagination.showing') }}
          <span>{{ fromItem }}</span>
          {{ t('pagination.to') }}
          <span>{{ toItem }}</span>
          {{ t('pagination.of') }}
          <span>{{ total }}</span>
          {{ t('pagination.results') }}
        </p>

        <!-- Page size selector -->
        <div v-if="showPageSizeSelector" class="page-size-control">
          <span>{{ t('pagination.perPage') }}</span>
          <div class="page-size-select">
            <Select
              :model-value="pageSize"
              :options="pageSizeSelectOptions"
              :searchable="false"
              @update:model-value="handlePageSizeChange"
            />
          </div>
        </div>

        <div v-if="showJump" class="page-jump-control">
          <span>{{ t('pagination.jumpTo') }}</span>
          <input
            v-model="jumpPage"
            type="number"
            min="1"
            :max="displayTotalPages"
            class="page-jump-input"
            :placeholder="t('pagination.jumpPlaceholder')"
            @keyup.enter="submitJump"
          />
          <button type="button" class="pagination-button page-jump-button" @click="submitJump">
            {{ t('pagination.jumpAction') }}
          </button>
        </div>
      </div>

      <!-- Desktop pagination buttons -->
      <nav
        class="pagination-nav"
        aria-label="Pagination"
      >
        <!-- Previous button -->
        <button
          @click="goToPage(page - 1)"
          :disabled="page <= 1"
          type="button"
          class="pagination-button pagination-icon-button"
          :aria-label="t('pagination.previous')"
        >
          <Icon name="chevronLeft" size="md" />
        </button>

        <!-- Page numbers -->
        <template v-for="(pageNum, index) in visiblePages" :key="`${pageNum}-${index}`">
          <button
            v-if="typeof pageNum === 'number'"
            @click="goToPage(pageNum)"
            type="button"
            :class="[
              'pagination-button pagination-page-button',
              pageNum === page && 'pagination-page-button-active'
            ]"
            :aria-label="t('pagination.goToPage', { page: pageNum })"
            :aria-current="pageNum === page ? 'page' : undefined"
          >
            {{ pageNum }}
          </button>
          <span v-else class="pagination-ellipsis" aria-hidden="true">…</span>
        </template>

        <!-- Next button -->
        <button
          @click="goToPage(page + 1)"
          :disabled="page >= displayTotalPages"
          type="button"
          class="pagination-button pagination-icon-button"
          :aria-label="t('pagination.next')"
        >
          <Icon name="chevronRight" size="md" />
        </button>
      </nav>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'
import Select from './Select.vue'
import { getConfiguredTablePageSizeOptions, normalizeTablePageSize } from '@/utils/tablePreferences'
import { setPersistedPageSize } from '@/composables/usePersistedPageSize'

const { t } = useI18n()

interface Props {
  total: number
  page: number
  pageSize: number
  pageSizeOptions?: number[]
  showPageSizeSelector?: boolean
  showJump?: boolean
}

interface Emits {
  (e: 'update:page', page: number): void
  (e: 'update:pageSize', pageSize: number): void
}

const props = withDefaults(defineProps<Props>(), {
  pageSizeOptions: () => getConfiguredTablePageSizeOptions(),
  showPageSizeSelector: true,
  showJump: false
})

const emit = defineEmits<Emits>()

const totalPages = computed(() => Math.ceil(props.total / props.pageSize))
const displayTotalPages = computed(() => Math.max(1, totalPages.value))

const fromItem = computed(() => {
  if (props.total === 0) return 0
  return (props.page - 1) * props.pageSize + 1
})

const toItem = computed(() => {
  const to = props.page * props.pageSize
  return to > props.total ? props.total : to
})

const pageSizeSelectOptions = computed(() => {
  const options = Array.from(
    new Set([
      ...getConfiguredTablePageSizeOptions(),
      normalizeTablePageSize(props.pageSize)
    ])
  ).sort((a, b) => a - b)

  return options.map((size) => ({
    value: size,
    label: String(size)
  }))
})

const jumpPage = ref('')

const visiblePages = computed(() => {
  const pages: (number | string)[] = []
  const maxVisible = 7
  const total = displayTotalPages.value

  if (total <= maxVisible) {
    // Show all pages if total is small
    for (let i = 1; i <= total; i++) {
      pages.push(i)
    }
  } else {
    // Always show first page
    pages.push(1)

    const start = Math.max(2, props.page - 2)
    const end = Math.min(total - 1, props.page + 2)

    // Add ellipsis before if needed
    if (start > 2) {
      pages.push('...')
    }

    // Add middle pages
    for (let i = start; i <= end; i++) {
      pages.push(i)
    }

    // Add ellipsis after if needed
    if (end < total - 1) {
      pages.push('...')
    }

    // Always show last page
    pages.push(total)
  }

  return pages
})

const goToPage = (newPage: number) => {
  if (newPage >= 1 && newPage <= displayTotalPages.value && newPage !== props.page) {
    emit('update:page', newPage)
  }
}

const handlePageSizeChange = (value: string | number | boolean | null) => {
  if (value === null || typeof value === 'boolean') return
  const newPageSize = normalizeTablePageSize(typeof value === 'string' ? parseInt(value, 10) : value)
  setPersistedPageSize(newPageSize)
  emit('update:pageSize', newPageSize)
}

const submitJump = () => {
  const value = jumpPage.value.trim()
  if (!value) return
  const pageNum = Number.parseInt(value, 10)
  if (Number.isNaN(pageNum)) return
  const nextPage = Math.min(Math.max(pageNum, 1), displayTotalPages.value)
  jumpPage.value = ''
  goToPage(nextPage)
}
</script>

<style scoped>
.pagination-shell {
  width: 100%;
  padding: 0.5rem 0.25rem;
}

.mobile-pagination {
  display: flex;
  width: 100%;
  align-items: center;
  justify-content: space-between;
  gap: 0.75rem;
}

.mobile-page-label,
.pagination-results,
.page-size-control,
.page-jump-control {
  color: var(--omnio-muted, #6b7280);
  font-size: 0.8rem;
  font-weight: 500;
}

.mobile-page-label,
.pagination-results span,
.pagination-page-button,
.page-size-select,
.page-jump-input {
  font-variant-numeric: tabular-nums;
}

.desktop-pagination {
  display: none;
  min-width: 0;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
}

.pagination-meta {
  display: flex;
  min-width: 0;
  align-items: center;
  gap: 1rem;
}

.pagination-results {
  white-space: nowrap;
}

.pagination-results span {
  color: var(--omnio-foreground, #111827);
  font-weight: 600;
}

.page-size-control,
.page-jump-control {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  white-space: nowrap;
}

.page-size-select {
  width: 4.4rem;
}

.page-size-select :deep(.select-trigger) {
  min-height: 2rem;
  height: 2rem;
  padding: 0.25rem 0.5rem 0.25rem 0.625rem;
  font-size: 0.8rem;
  font-weight: 550;
}

.page-jump-input {
  width: 4rem;
  height: 2rem;
  border: 1px solid var(--omnio-border, #e5e7eb);
  border-radius: 0.5rem;
  padding: 0 0.5rem;
  color: var(--omnio-foreground, #111827);
  background: var(--omnio-surface, #fff);
  outline: none;
}

.page-jump-input:focus-visible {
  border-color: var(--omnio-primary, #3b82f6);
  box-shadow: 0 0 0 3px color-mix(in srgb, var(--omnio-primary, #3b82f6) 22%, transparent);
}

.pagination-nav {
  display: inline-flex;
  flex-shrink: 0;
  align-items: center;
  gap: 0.25rem;
}

.pagination-button {
  display: inline-flex;
  min-height: 2rem;
  align-items: center;
  justify-content: center;
  border: 1px solid var(--omnio-border, #e5e7eb);
  border-radius: 0.5rem;
  color: var(--omnio-muted, #6b7280);
  background: var(--omnio-surface, #fff);
  font-size: 0.8rem;
  line-height: 1;
  font-weight: 550;
  outline: none;
  transition: color 140ms ease, background-color 140ms ease, border-color 140ms ease, box-shadow 140ms ease;
}

.pagination-button:hover:not(:disabled) {
  color: var(--omnio-foreground, #111827);
  border-color: var(--omnio-border-strong, #d1d5db);
  background: color-mix(in srgb, var(--omnio-foreground, #111827) 5%, var(--omnio-surface, #fff));
}

.pagination-button:focus-visible {
  border-color: var(--omnio-primary, #3b82f6);
  box-shadow: 0 0 0 3px color-mix(in srgb, var(--omnio-primary, #3b82f6) 22%, transparent);
}

.pagination-button:disabled {
  cursor: not-allowed;
  opacity: 0.42;
}

.pagination-icon-button,
.pagination-page-button {
  width: 2rem;
  height: 2rem;
  padding: 0;
}

.pagination-page-button {
  min-width: 2rem;
  width: auto;
  padding: 0 0.45rem;
}

.pagination-page-button-active,
.pagination-page-button-active:hover:not(:disabled) {
  color: #fff;
  border-color: var(--omnio-primary-strong, #2563eb);
  background: var(--omnio-primary-strong, #2563eb);
}

.pagination-ellipsis {
  display: inline-flex;
  width: 1.5rem;
  height: 2rem;
  align-items: center;
  justify-content: center;
  color: var(--omnio-muted, #6b7280);
  font-size: 0.8rem;
}

.page-jump-button {
  padding: 0 0.65rem;
}

@media (max-width: 900px) {
  .pagination-results {
    display: none;
  }

  .pagination-meta {
    gap: 0.75rem;
  }
}

@media (min-width: 641px) {
  .mobile-pagination {
    display: none;
  }

  .desktop-pagination {
    display: flex;
  }
}

@media (prefers-reduced-motion: reduce) {
  .pagination-button {
    transition-duration: 1ms;
  }
}
</style>
