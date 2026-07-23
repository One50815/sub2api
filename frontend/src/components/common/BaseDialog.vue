<template>
  <Teleport to="body">
    <Transition name="modal">
      <div
        v-if="show"
        class="modal-overlay"
        :style="zIndexStyle"
        @click.self="handleClose"
      >
        <!-- Modal panel -->
        <div
          ref="dialogRef"
          :class="['modal-content', widthClasses]"
          role="dialog"
          aria-modal="true"
          :aria-labelledby="dialogId"
          tabindex="-1"
          @click.stop
        >
          <!-- Header -->
          <div class="modal-header">
            <h3 :id="dialogId" class="modal-title">
              {{ title }}
            </h3>
            <button
              v-if="showCloseButton"
              type="button"
              @click="emit('close')"
              class="modal-close"
              aria-label="Close modal"
            >
              <Icon name="x" size="md" />
            </button>
          </div>

          <!-- Body -->
          <div class="modal-body">
            <slot></slot>
          </div>

          <!-- Footer -->
          <div v-if="$slots.footer" class="modal-footer">
            <slot name="footer"></slot>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { computed, watch, onMounted, onUnmounted, ref, nextTick } from 'vue'
import Icon from '@/components/icons/Icon.vue'

// 生成唯一ID以避免多个对话框时ID冲突
let dialogIdCounter = 0
const dialogId = `modal-title-${++dialogIdCounter}`

// 焦点管理
const dialogRef = ref<HTMLElement | null>(null)
let previousActiveElement: HTMLElement | null = null

type DialogWidth = 'narrow' | 'normal' | 'wide' | 'extra-wide' | 'full'

interface Props {
  show: boolean
  title: string
  width?: DialogWidth
  closeOnEscape?: boolean
  closeOnClickOutside?: boolean
  showCloseButton?: boolean
  zIndex?: number
}

interface Emits {
  (e: 'close'): void
}

const props = withDefaults(defineProps<Props>(), {
  width: 'normal',
  closeOnEscape: true,
  closeOnClickOutside: false,
  showCloseButton: true,
  zIndex: 50
})

const emit = defineEmits<Emits>()

// Custom z-index style (overrides the default z-50 from CSS)
const zIndexStyle = computed(() => {
  return props.zIndex !== 50 ? { zIndex: props.zIndex } : undefined
})

const widthClasses = computed(() => {
  // Width guidance: narrow=confirm/short prompts, normal=standard forms,
  // wide=multi-section forms or rich content, extra-wide=analytics/tables,
  // full=full-screen or very dense layouts.
  const widths: Record<DialogWidth, string> = {
    narrow: 'max-w-sm',
    normal: 'max-w-lg',
    wide: 'w-full sm:max-w-2xl md:max-w-3xl lg:max-w-4xl',
    'extra-wide': 'w-full sm:max-w-3xl md:max-w-4xl lg:max-w-5xl xl:max-w-6xl',
    full: 'w-full sm:max-w-4xl md:max-w-5xl lg:max-w-6xl xl:max-w-7xl'
  }
  return widths[props.width]
})

const handleClose = () => {
  if (props.closeOnClickOutside) {
    emit('close')
  }
}

const getFocusableElements = () => {
  if (!dialogRef.value) return []
  return Array.from(
    dialogRef.value.querySelectorAll<HTMLElement>(
      'button:not([disabled]), [href], input:not([disabled]), select:not([disabled]), textarea:not([disabled]), [tabindex]:not([tabindex="-1"])'
    )
  ).filter((element) => !element.hasAttribute('hidden') && element.getAttribute('aria-hidden') !== 'true')
}

const handleKeyDown = (event: KeyboardEvent) => {
  if (!props.show) return

  if (props.closeOnEscape && event.key === 'Escape') {
    event.preventDefault()
    emit('close')
    return
  }

  if (event.key !== 'Tab') return
  const focusable = getFocusableElements()
  if (focusable.length === 0) {
    event.preventDefault()
    dialogRef.value?.focus()
    return
  }

  const first = focusable[0]
  const last = focusable[focusable.length - 1]
  if (event.shiftKey && document.activeElement === first) {
    event.preventDefault()
    last.focus()
  } else if (!event.shiftKey && document.activeElement === last) {
    event.preventDefault()
    first.focus()
  }
}

// Prevent body scroll when modal is open and manage focus
watch(
  () => props.show,
  async (isOpen) => {
    if (isOpen) {
      // 保存当前焦点元素
      previousActiveElement = document.activeElement as HTMLElement
      // 使用CSS类而不是直接操作style,更易于管理多个对话框
      document.body.classList.add('modal-open')

      // 等待DOM更新后设置焦点到对话框
      await nextTick()
      if (dialogRef.value) {
        const firstFocusable = getFocusableElements()[0]
        ;(firstFocusable ?? dialogRef.value)?.focus()
      }
    } else {
      document.body.classList.remove('modal-open')
      // 恢复之前的焦点
      if (previousActiveElement && typeof previousActiveElement.focus === 'function') {
        previousActiveElement.focus()
      }
      previousActiveElement = null
    }
  },
  { immediate: true }
)

onMounted(() => {
  document.addEventListener('keydown', handleKeyDown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeyDown)
  // 确保组件卸载时移除滚动锁定
  document.body.classList.remove('modal-open')
})
</script>

<style scoped>
.modal-overlay {
  padding: 1rem !important;
  background: rgba(0, 0, 0, 0.1) !important;
  backdrop-filter: blur(2px) !important;
}

.modal-content {
  max-height: calc(100svh - 2rem) !important;
  overflow: hidden;
  border: 1px solid var(--omnio-border, rgba(15, 23, 42, 0.12)) !important;
  border-radius: 0.75rem !important;
  background: var(--omnio-surface, #fff) !important;
  box-shadow: 0 18px 60px rgba(15, 23, 42, 0.16), 0 2px 8px rgba(15, 23, 42, 0.08) !important;
  outline: none;
}

.modal-content:focus-visible {
  box-shadow: 0 18px 60px rgba(15, 23, 42, 0.16), 0 0 0 3px color-mix(in srgb, var(--omnio-primary, #3b82f6) 24%, transparent) !important;
}

.modal-header {
  min-height: 3.25rem;
  padding: 1rem 1rem 0.5rem !important;
  border-bottom: 0 !important;
}

.modal-title {
  color: var(--omnio-foreground, #111827) !important;
  font-size: 1rem !important;
  line-height: 1.25 !important;
  font-weight: 560 !important;
  letter-spacing: 0 !important;
}

.modal-close {
  display: inline-flex;
  width: 2rem;
  height: 2rem;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  margin: -0.5rem -0.5rem -0.5rem 0;
  border-radius: 0.5rem;
  color: var(--omnio-muted, #6b7280);
  transition: color 140ms ease, background-color 140ms ease, box-shadow 140ms ease;
}

.modal-close:hover {
  color: var(--omnio-foreground, #111827);
  background: color-mix(in srgb, var(--omnio-foreground, #111827) 5%, transparent);
}

.modal-close:focus-visible {
  outline: none;
  box-shadow: 0 0 0 3px color-mix(in srgb, var(--omnio-primary, #3b82f6) 24%, transparent);
}

.modal-body {
  padding: 0.5rem 1rem 1rem !important;
  color: var(--omnio-foreground, #111827);
}

.modal-footer {
  min-height: 4rem;
  padding: 1rem !important;
  border-top: 1px solid var(--omnio-border, rgba(15, 23, 42, 0.12)) !important;
  border-radius: 0 0 0.75rem 0.75rem;
  background: color-mix(in srgb, var(--omnio-foreground, #111827) 3%, var(--omnio-surface, #fff)) !important;
}

:global(.dark) .modal-content {
  box-shadow: 0 22px 70px rgba(0, 0, 0, 0.42), 0 2px 10px rgba(0, 0, 0, 0.24) !important;
}

@media (max-width: 640px) {
  .modal-overlay {
    padding: 0.75rem !important;
  }

  .modal-content {
    max-width: 100%;
    max-height: calc(100svh - 1.5rem) !important;
  }
}

@media (prefers-reduced-motion: reduce) {
  .modal-close {
    transition-duration: 1ms;
  }
}
</style>
