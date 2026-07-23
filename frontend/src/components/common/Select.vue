<template>
  <div class="relative" ref="containerRef">
    <button
      ref="triggerRef"
      type="button"
      @click="toggle"
      :disabled="disabled"
      :aria-expanded="isOpen"
      aria-haspopup="listbox"
      :aria-controls="isOpen ? `${instanceId}-listbox` : undefined"
      :aria-label="placeholderText"
      :class="[
        'select-trigger',
        isOpen && 'select-trigger-open',
        error && 'select-trigger-error',
        disabled && 'select-trigger-disabled'
      ]"
      @keydown.down.prevent="onTriggerKeyDown"
      @keydown.up.prevent="onTriggerKeyDown"
    >
      <span :class="['select-value', !hasValue && 'select-value-placeholder']">
        <slot name="selected" :option="selectedOption">
          {{ selectedLabel }}
        </slot>
      </span>
      <span
        v-if="clearable && hasValue && !disabled"
        class="select-clear"
        role="button"
        tabindex="-1"
        aria-label="Clear selection"
        @click.stop="clearSelection"
        @mousedown.stop
        @keydown.enter.stop.prevent="clearSelection"
      >
        <Icon name="x" size="sm" />
      </span>
      <span class="select-icon">
        <Icon
          name="chevronDown"
          size="md"
          :class="['transition-transform duration-200', isOpen && 'rotate-180']"
        />
      </span>
    </button>

    <!-- Teleport dropdown to body to escape stacking context -->
    <Teleport to="body">
      <Transition name="select-dropdown">
        <div
          v-if="isOpen"
          ref="dropdownRef"
          :id="`${instanceId}-listbox`"
          class="select-dropdown-portal"
          :class="[instanceId]"
          :style="dropdownStyle"
          role="listbox"
          tabindex="-1"
          :aria-activedescendant="focusedIndex >= 0 ? getOptionId(focusedIndex) : undefined"
          @click.stop
          @mousedown.stop
          @keydown="onDropdownKeyDown"
        >
          <!-- Search input -->
          <div v-if="isSearchable" class="select-search">
            <Icon name="search" size="sm" class="text-gray-400" />
            <input
              ref="searchInputRef"
              v-model="searchQuery"
              type="text"
              :placeholder="searchPlaceholderText"
              class="select-search-input"
              @click.stop
            />
          </div>

          <!-- Options list -->
          <div class="select-options" ref="optionsListRef">
            <div
              v-for="(option, index) in filteredOptions"
              :key="`${typeof getOptionValue(option)}:${String(getOptionValue(option) ?? '')}`"
              :id="getOptionId(index)"
              :role="isGroupHeaderOption(option) ? 'presentation' : 'option'"
              :aria-selected="isGroupHeaderOption(option) ? undefined : isSelected(option)"
              :aria-disabled="isGroupHeaderOption(option) || isOptionDisabled(option)"
              @click.stop="isSelectableOption(option) && selectOption(option)"
              @mouseenter="handleOptionMouseEnter(option, index)"
              :class="[
                'select-option',
                isGroupHeaderOption(option) && 'select-option-group',
                isSelected(option) && 'select-option-selected',
                isOptionDisabled(option) && !isGroupHeaderOption(option) && 'select-option-disabled',
                focusedIndex === index && !isGroupHeaderOption(option) && 'select-option-focused'
              ]"
            >
              <slot name="option" :option="option" :selected="isSelected(option)">
                <Icon
                  v-if="option._creatable"
                  name="search"
                  size="sm"
                  class="flex-shrink-0 text-gray-400"
                />
                <span class="select-option-label" :class="option._creatable && 'italic text-gray-500 dark:text-dark-300'">{{ getOptionLabel(option) }}</span>
                <Icon
                  v-if="isSelected(option)"
                  name="check"
                  size="sm"
                  class="text-primary-500"
                  :stroke-width="2"
                />
              </slot>
            </div>

            <!-- Empty state -->
            <div v-if="filteredOptions.length === 0" class="select-empty">
              {{ emptyTextDisplay }}
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'

const { t } = useI18n()

// Instance ID for unique click-outside detection
const instanceId = `select-${Math.random().toString(36).substring(2, 9)}`

export interface SelectOption {
  value: string | number | boolean | null
  label: string
  disabled?: boolean
  [key: string]: unknown
}

interface Props {
  modelValue: string | number | boolean | null | undefined
  options: SelectOption[] | Array<Record<string, unknown>>
  placeholder?: string
  disabled?: boolean
  error?: boolean
  searchable?: boolean | 'auto'
  searchPlaceholder?: string
  emptyText?: string
  valueKey?: string
  labelKey?: string
  creatable?: boolean
  creatablePrefix?: string
  clearable?: boolean
}

interface Emits {
  (e: 'update:modelValue', value: string | number | boolean | null): void
  (e: 'change', value: string | number | boolean | null, option: SelectOption | null): void
}

const props = withDefaults(defineProps<Props>(), {
  disabled: false,
  error: false,
  searchable: 'auto',
  creatable: false,
  creatablePrefix: '',
  clearable: false,
  valueKey: 'value',
  labelKey: 'label'
})

const emit = defineEmits<Emits>()

const isOpen = ref(false)
const searchQuery = ref('')
const focusedIndex = ref(-1)
const containerRef = ref<HTMLElement | null>(null)
const triggerRef = ref<HTMLButtonElement | null>(null)
const searchInputRef = ref<HTMLInputElement | null>(null)
const dropdownRef = ref<HTMLElement | null>(null)
const optionsListRef = ref<HTMLElement | null>(null)
const dropdownPosition = ref<'bottom' | 'top'>('bottom')
const triggerRect = ref<DOMRect | null>(null)

// i18n placeholders
const placeholderText = computed(() => props.placeholder ?? t('common.selectOption'))
const searchPlaceholderText = computed(() => props.searchPlaceholder ?? t('common.searchPlaceholder'))
const emptyTextDisplay = computed(() => props.emptyText ?? t('common.noOptionsFound'))

const isSearchable = computed(() => {
  if (props.searchable === 'auto') return props.options.length > 5
  return props.searchable
})

// Computed style for teleported dropdown
const dropdownStyle = computed(() => {
  if (!triggerRect.value) return {}

  const rect = triggerRect.value
  const style: Record<string, string> = {
    position: 'fixed',
    left: `${rect.left}px`,
    width: `${rect.width}px`,
    minWidth: `${rect.width}px`,
    maxWidth: 'calc(100vw - 16px)',
    zIndex: '100000020'
  }

  if (dropdownPosition.value === 'top') {
    style.bottom = `${window.innerHeight - rect.top + 4}px`
  } else {
    style.top = `${rect.bottom + 4}px`
  }

  return style
})

const getOptionValue = (option: any): any => {
  if (typeof option === 'object' && option !== null) {
    return option[props.valueKey]
  }
  return option
}

const getOptionLabel = (option: any): string => {
  if (typeof option === 'object' && option !== null) {
    return String(option[props.labelKey] ?? '')
  }
  return String(option ?? '')
}

const isOptionDisabled = (option: any): boolean => {
  if (typeof option === 'object' && option !== null) {
    return !!option.disabled
  }
  return false
}

const isGroupHeaderOption = (option: any): boolean => {
  if (typeof option === 'object' && option !== null) {
    return option.kind === 'group'
  }
  return false
}

const isSelectableOption = (option: any): boolean =>
  !isGroupHeaderOption(option) && !isOptionDisabled(option)

const getOptionId = (index: number) => `${instanceId}-option-${index}`

const selectedOption = computed(() => {
  return props.options.find((opt) => getOptionValue(opt) === props.modelValue) || null
})

const selectedLabel = computed(() => {
  if (selectedOption.value) {
    return getOptionLabel(selectedOption.value)
  }
  // In creatable mode, show the raw value if no matching option
  if (props.creatable && props.modelValue) {
    return String(props.modelValue)
  }
  return placeholderText.value
})

const hasValue = computed(
  () => props.modelValue !== null && props.modelValue !== undefined && props.modelValue !== ''
)

const filteredOptions = computed(() => {
  let opts = props.options as any[]
  if (isSearchable.value && searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    opts = opts.filter((opt) => {
      // Match label
      if (getOptionLabel(opt).toLowerCase().includes(query)) return true
      // Also match description if present
      if (opt.description && String(opt.description).toLowerCase().includes(query)) return true
      return false
    })
    // In creatable mode, always prepend a fuzzy search option
    if (props.creatable && searchQuery.value.trim()) {
      const trimmed = searchQuery.value.trim()
      const prefix = props.creatablePrefix || t('common.search')
      opts = [{ [props.valueKey]: trimmed, [props.labelKey]: `${prefix} "${trimmed}"`, _creatable: true }, ...opts]
    }
  }
  return opts
})

const isSelected = (option: any): boolean => {
  return getOptionValue(option) === props.modelValue
}

const findNextEnabledIndex = (startIndex: number): number => {
  const opts = filteredOptions.value
  if (opts.length === 0) return -1
  for (let offset = 0; offset < opts.length; offset++) {
    const idx = (startIndex + offset) % opts.length
    if (isSelectableOption(opts[idx])) return idx
  }
  return -1
}

const findPrevEnabledIndex = (startIndex: number): number => {
  const opts = filteredOptions.value
  if (opts.length === 0) return -1
  for (let offset = 0; offset < opts.length; offset++) {
    const idx = (startIndex - offset + opts.length) % opts.length
    if (isSelectableOption(opts[idx])) return idx
  }
  return -1
}

const handleOptionMouseEnter = (option: any, index: number) => {
  if (isOptionDisabled(option) || isGroupHeaderOption(option)) return
  focusedIndex.value = index
}

// Update trigger rect periodically while open to follow scroll/resize
const updateTriggerRect = () => {
  if (containerRef.value) {
    triggerRect.value = containerRef.value.getBoundingClientRect()
  }
}

const calculateDropdownPosition = () => {
  if (!containerRef.value) return
  updateTriggerRect()

  nextTick(() => {
    if (!dropdownRef.value || !triggerRect.value) return
    const dropdownHeight = dropdownRef.value.offsetHeight || 240
    const spaceBelow = window.innerHeight - triggerRect.value.bottom
    const spaceAbove = triggerRect.value.top

    if (spaceBelow < dropdownHeight && spaceAbove > dropdownHeight) {
      dropdownPosition.value = 'top'
    } else {
      dropdownPosition.value = 'bottom'
    }
  })
}

const toggle = () => {
  if (props.disabled) return
  isOpen.value = !isOpen.value
}

watch(isOpen, (open) => {
  if (open) {
    calculateDropdownPosition()
    // Reset focused index to current selection or first item
    if (filteredOptions.value.length === 0) {
      focusedIndex.value = -1
    } else {
      const selectedIdx = filteredOptions.value.findIndex(isSelected)
      const initialIdx = selectedIdx >= 0 ? selectedIdx : 0
      focusedIndex.value = !isSelectableOption(filteredOptions.value[initialIdx])
        ? findNextEnabledIndex(initialIdx + 1)
        : initialIdx
    }

    if (isSearchable.value) {
      nextTick(() => searchInputRef.value?.focus())
    } else {
      nextTick(() => dropdownRef.value?.focus())
    }
    // Add scroll listener to update position
    window.addEventListener('scroll', updateTriggerRect, { capture: true, passive: true })
    window.addEventListener('resize', calculateDropdownPosition)
  } else {
    searchQuery.value = ''
    focusedIndex.value = -1
    window.removeEventListener('scroll', updateTriggerRect, { capture: true })
    window.removeEventListener('resize', calculateDropdownPosition)
  }
})

const selectOption = (option: any) => {
  const value = getOptionValue(option) ?? null
  emit('update:modelValue', value)
  emit('change', value, option)
  isOpen.value = false
  triggerRef.value?.focus()
}

const clearSelection = () => {
  if (props.disabled) return
  emit('update:modelValue', null)
  emit('change', null, null)
}

// Keyboards
const onTriggerKeyDown = () => {
  if (!isOpen.value) {
    isOpen.value = true
  }
}

const onDropdownKeyDown = (e: KeyboardEvent) => {
  switch (e.key) {
    case 'ArrowDown':
      e.preventDefault()
      focusedIndex.value = findNextEnabledIndex(focusedIndex.value + 1)
      if (focusedIndex.value >= 0) scrollToFocused()
      break
    case 'ArrowUp':
      e.preventDefault()
      focusedIndex.value = findPrevEnabledIndex(focusedIndex.value - 1)
      if (focusedIndex.value >= 0) scrollToFocused()
      break
    case 'Enter':
      e.preventDefault()
      if (focusedIndex.value >= 0 && focusedIndex.value < filteredOptions.value.length) {
        const opt = filteredOptions.value[focusedIndex.value]
        if (isSelectableOption(opt)) selectOption(opt)
      }
      break
    case 'Escape':
      e.preventDefault()
      isOpen.value = false
      triggerRef.value?.focus()
      break
    case 'Tab':
      isOpen.value = false
      break
  }
}

const scrollToFocused = () => {
  nextTick(() => {
    const list = optionsListRef.value
    if (!list) return
    const focusedEl = list.children[focusedIndex.value] as HTMLElement
    if (!focusedEl) return

    if (focusedEl.offsetTop < list.scrollTop) {
      list.scrollTop = focusedEl.offsetTop
    } else if (focusedEl.offsetTop + focusedEl.offsetHeight > list.scrollTop + list.offsetHeight) {
      list.scrollTop = focusedEl.offsetTop + focusedEl.offsetHeight - list.offsetHeight
    }
  })
}

const handleClickOutside = (event: MouseEvent) => {
  const target = event.target as HTMLElement
  // Check if click is inside THIS specific instance's dropdown or trigger
  const isInDropdown = !!target.closest(`.${instanceId}`)
  const isInTrigger = containerRef.value?.contains(target)

  if (!isInDropdown && !isInTrigger && isOpen.value) {
    isOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
  window.removeEventListener('scroll', updateTriggerRect, { capture: true })
  window.removeEventListener('resize', calculateDropdownPosition)
})
</script>

<style scoped>
.select-trigger {
  display: flex;
  width: 100%;
  min-height: var(--omnio-control-height, 2.5rem);
  align-items: center;
  justify-content: space-between;
  gap: 0.375rem;
  padding: 0.45rem 0.75rem;
  cursor: pointer;
  border: 1px solid var(--omnio-border, #e5e7eb);
  border-radius: var(--omnio-radius-control, 0.875rem);
  color: var(--omnio-foreground, #111827);
  background: transparent;
  font-size: 0.875rem;
  line-height: 1.25;
  outline: none;
  transition: color 140ms ease, background-color 140ms ease, border-color 140ms ease, box-shadow 140ms ease;
}

.select-trigger:hover:not(:disabled) {
  border-color: var(--omnio-border-strong, #d1d5db);
  background: color-mix(in srgb, var(--omnio-foreground, #111827) 2.5%, transparent);
}

.select-trigger:focus-visible {
  border-color: var(--omnio-primary, #3b82f6);
  box-shadow: 0 0 0 3px color-mix(in srgb, var(--omnio-primary, #3b82f6) 22%, transparent);
}

.select-trigger-open {
  border-color: var(--omnio-primary, #3b82f6);
  box-shadow: 0 0 0 3px color-mix(in srgb, var(--omnio-primary, #3b82f6) 20%, transparent);
}

.select-trigger-error {
  border-color: #ef4444;
  box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.16);
}

.select-trigger-disabled {
  cursor: not-allowed;
  opacity: 0.5;
  background: color-mix(in srgb, var(--omnio-foreground, #111827) 4%, transparent);
}

.select-value {
  min-width: 0;
  flex: 1;
  overflow: hidden;
  text-align: left;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.select-value-placeholder {
  color: var(--omnio-muted, #6b7280);
}

.select-icon {
  display: inline-flex;
  flex-shrink: 0;
  color: var(--omnio-muted, #6b7280);
}

.select-clear {
  display: inline-flex;
  width: 1.25rem;
  height: 1.25rem;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  border-radius: 0.35rem;
  color: var(--omnio-muted, #6b7280);
  transition: color 140ms ease, background-color 140ms ease;
}

.select-clear:hover {
  color: var(--omnio-foreground, #111827);
  background: color-mix(in srgb, var(--omnio-foreground, #111827) 6%, transparent);
}
</style>

<style>
.select-dropdown-portal {
  min-width: 9rem;
  max-height: min(20rem, calc(100svh - 1rem));
  padding: 0.25rem;
  overflow: hidden;
  pointer-events: auto !important;
  border: 0;
  border: 1px solid var(--omnio-border, #e5e7eb);
  border-radius: var(--omnio-radius-panel, 1.25rem);
  color: var(--omnio-foreground, #111827);
  background: var(--omnio-surface, #fff);
  box-shadow: var(--omnio-card-shadow-hover, 0 8px 24px rgba(15, 23, 42, 0.14));
}

.select-dropdown-portal .select-search {
  display: flex;
  min-height: 2.25rem;
  align-items: center;
  gap: 0.375rem;
  margin-bottom: 0.25rem;
  padding: 0.3rem 0.5rem;
  border: 1px solid var(--omnio-border, #e5e7eb);
  border-radius: var(--omnio-radius-control, 0.875rem);
}

.select-dropdown-portal .select-search-input {
  min-width: 0;
  flex: 1;
  color: var(--omnio-foreground, #111827);
  background: transparent;
  font-size: 0.8rem;
  outline: none;
}

.select-dropdown-portal .select-search-input::placeholder {
  color: var(--omnio-muted, #6b7280);
}

.select-dropdown-portal .select-options {
  max-height: 18rem;
  overflow-y: auto;
  outline: none;
}

.select-dropdown-portal .select-option {
  position: relative;
  display: flex;
  width: 100%;
  min-height: 2.25rem;
  align-items: center;
  justify-content: space-between;
  gap: 0.375rem;
  padding: 0.25rem 0.4rem;
  pointer-events: auto !important;
  cursor: default;
  border-radius: 0.75rem;
  color: var(--omnio-foreground, #111827);
  font-size: 0.82rem;
  line-height: 1.35;
  transition: color 100ms ease, background-color 100ms ease;
}

.select-dropdown-portal .select-option-selected {
  color: var(--omnio-primary-strong, #2563eb);
  background: color-mix(in srgb, var(--omnio-primary, #3b82f6) 9%, transparent);
}

.select-dropdown-portal .select-option:hover,
.select-dropdown-portal .select-option-focused {
  color: var(--omnio-foreground, #111827);
  background: color-mix(in srgb, var(--omnio-foreground, #111827) 6%, transparent);
}

.select-dropdown-portal .select-option-disabled {
  pointer-events: none !important;
  cursor: not-allowed;
  opacity: 0.45;
}

.select-dropdown-portal .select-option-group {
  min-height: 1.75rem;
  padding: 0.35rem 0.4rem 0.2rem;
  cursor: default;
  user-select: none;
  color: var(--omnio-muted, #6b7280);
  background: transparent;
  font-size: 0.68rem;
  line-height: 1.2;
  font-weight: 600;
  letter-spacing: 0.07em;
  text-transform: uppercase;
}

.select-dropdown-portal .select-option-group:hover {
  color: var(--omnio-muted, #6b7280);
  background: transparent;
}

.select-dropdown-portal .select-option-label {
  min-width: 0;
  flex: 1;
  overflow: hidden;
  text-align: left;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.select-dropdown-portal .select-empty {
  padding: 1.5rem 0.75rem;
  color: var(--omnio-muted, #6b7280);
  text-align: center;
  font-size: 0.8rem;
}

.dark .select-dropdown-portal {
  box-shadow: 0 10px 28px rgba(0, 0, 0, 0.38), 0 0 0 1px color-mix(in srgb, var(--omnio-foreground, #f8fafc) 12%, transparent);
}

.select-dropdown-enter-active,
.select-dropdown-leave-active {
  transform-origin: top;
  transition: opacity 100ms ease, transform 100ms ease;
}

.select-dropdown-enter-from,
.select-dropdown-leave-to {
  opacity: 0;
  transform: translateY(-0.25rem) scale(0.98);
}

@media (prefers-reduced-motion: reduce) {
  .select-dropdown-enter-active,
  .select-dropdown-leave-active,
  .select-dropdown-portal .select-option {
    transition-duration: 1ms;
  }
}
</style>
