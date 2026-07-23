<template>
  <BaseDialog
    :show="show"
    :title="t('admin.accounts.testAccountConnection')"
    width="normal"
    @close="handleClose"
  >
    <div class="account-test-stack">
      <!-- Account Info Card -->
      <div
        v-if="account"
        class="test-account-summary"
      >
        <div class="flex items-center gap-3">
          <div
            class="test-account-icon"
          >
            <Icon name="play" size="md" class="text-primary-600 dark:text-primary-400" :stroke-width="2" />
          </div>
          <div>
            <div class="font-semibold text-gray-900 dark:text-gray-100">{{ account.name }}</div>
            <div class="flex items-center gap-1.5 text-xs text-gray-500 dark:text-gray-400">
              <span
                class="platform-badge"
              >
                {{ account.type }}
              </span>
              <span>{{ t('admin.accounts.account') }}</span>
            </div>
          </div>
        </div>
        <span
          :class="[
            'test-status',
            account.status === 'active'
              ? 'test-status-active'
              : 'test-status-neutral'
          ]"
        >
          {{ account.status }}
        </span>
      </div>

      <div class="space-y-1.5">
        <label class="text-sm font-medium text-gray-700 dark:text-gray-300">
          {{ t('admin.accounts.selectTestModel') }}
        </label>
        <Select
          v-model="selectedModelId"
          :options="availableModels"
          :disabled="loadingModels || status === 'connecting'"
          value-key="id"
          label-key="display_name"
          :placeholder="loadingModels ? t('common.loading') + '...' : t('admin.accounts.selectTestModel')"
        />
      </div>

      <div v-if="isOpenAIAccount" class="space-y-1.5">
        <label class="text-sm font-medium text-gray-700 dark:text-gray-300">
          {{ t('admin.accounts.openai.testMode') }}
        </label>
        <Select
          v-model="testMode"
          :options="openAITestModeOptions"
          :disabled="status === 'connecting'"
        />
      </div>

      <div v-if="supportsImageTest" class="space-y-1.5">
        <TextArea
          v-model="testPrompt"
          :label="t('admin.accounts.imagePromptLabel')"
          :placeholder="t('admin.accounts.imagePromptPlaceholder')"
          :hint="t('admin.accounts.imageTestHint')"
          :disabled="status === 'connecting'"
          rows="3"
        />
      </div>

      <!-- Terminal Output -->
      <div class="group relative">
        <div
          ref="terminalRef"
          class="test-terminal max-h-[240px] min-h-[120px] overflow-y-auto p-4 font-mono text-sm"
        >
          <!-- Status Line -->
          <div v-if="status === 'idle'" class="flex items-center gap-2 text-gray-500">
            <Icon name="play" size="sm" :stroke-width="2" />
            <span>{{ t('admin.accounts.readyToTest') }}</span>
          </div>
          <div v-else-if="status === 'connecting'" class="flex items-center gap-2 text-yellow-400">
            <Icon name="refresh" size="sm" class="animate-spin" :stroke-width="2" />
            <span>{{ t('admin.accounts.connectingToApi') }}</span>
          </div>

          <!-- Output Lines -->
          <div v-for="(line, index) in outputLines" :key="index" :class="line.class">
            {{ line.text }}
          </div>

          <!-- Streaming Content -->
          <div v-if="streamingContent" class="text-green-400">
            {{ streamingContent }}<span class="animate-pulse">_</span>
          </div>

          <!-- Result Status -->
          <div
            v-if="status === 'success'"
            class="mt-3 flex items-center gap-2 border-t border-gray-700 pt-3 text-green-400"
          >
            <Icon name="check" size="sm" :stroke-width="2" />
            <span>{{ t('admin.accounts.testCompleted') }}</span>
          </div>
          <div
            v-else-if="status === 'error'"
            class="mt-3 flex items-center gap-2 border-t border-gray-700 pt-3 text-red-400"
          >
            <Icon name="x" size="sm" :stroke-width="2" />
            <span>{{ errorMessage }}</span>
          </div>
        </div>

        <!-- Copy Button -->
        <button
          v-if="outputLines.length > 0"
          type="button"
          @click="copyOutput"
          class="terminal-copy-button"
          :title="t('admin.accounts.copyOutput')"
        >
          <Icon name="link" size="sm" :stroke-width="2" />
        </button>
      </div>

      <div v-if="generatedImages.length > 0" class="space-y-2">
        <div class="text-xs font-medium text-gray-600 dark:text-gray-300">
          {{ t('admin.accounts.imagePreview') }}
        </div>
        <div class="flex flex-wrap justify-center gap-3">
          <div
            v-for="(image, index) in generatedImages"
            :key="`${image.url}-${index}`"
            class="test-image-card group/img"
            @click="previewImageUrl = image.url"
          >
            <img :src="image.url" :alt="`test-image-${index + 1}`" class="max-h-[360px] w-full object-contain" />
            <div class="absolute inset-0 flex items-center justify-center bg-black/0 transition-colors group-hover/img:bg-black/20">
              <Icon name="eye" size="lg" class="text-white opacity-0 drop-shadow-lg transition-opacity group-hover/img:opacity-100" :stroke-width="2" />
            </div>
            <div class="border-t border-gray-100 px-3 py-1.5 text-xs text-gray-500 dark:border-dark-500 dark:text-gray-300">
              {{ image.mimeType || 'image/*' }}
            </div>
          </div>
        </div>
      </div>

      <!-- Image Lightbox -->
      <Teleport to="body">
        <Transition name="fade">
          <div
            v-if="previewImageUrl"
            class="fixed inset-0 z-[100] flex items-center justify-center bg-black/80 p-4"
            @click.self="previewImageUrl = ''"
          >
            <button
              type="button"
              class="lightbox-close"
              :aria-label="t('common.close')"
              @click="previewImageUrl = ''"
            >
              <Icon name="x" size="lg" :stroke-width="2" />
            </button>
            <img
              :src="previewImageUrl"
              alt="preview"
              class="max-h-[90vh] max-w-[90vw] rounded-lg object-contain shadow-lg"
            />
          </div>
        </Transition>
      </Teleport>

      <!-- Test Info -->
      <div class="test-meta">
        <div class="flex items-center gap-3">
          <span class="flex items-center gap-1">
            <Icon name="grid" size="sm" :stroke-width="2" />
            {{ t('admin.accounts.testModel') }}
          </span>
        </div>
        <span class="flex items-center gap-1">
          <Icon name="chat" size="sm" :stroke-width="2" />
          {{
            supportsImageTest
              ? t('admin.accounts.imageTestMode')
              : t('admin.accounts.testPrompt')
          }}
        </span>
      </div>
    </div>

    <template #footer>
      <div class="test-footer-actions">
        <button
          type="button"
          @click="handleClose"
          class="btn btn-secondary"
        >
          {{ t('common.close') }}
        </button>
        <button
          type="button"
          @click="startTest"
          :disabled="status === 'connecting' || !selectedModelId"
          :class="[
            'test-action-button',
            status === 'connecting' || !selectedModelId
              ? 'test-action-disabled'
              : status === 'success'
                ? 'test-action-success'
                : status === 'error'
                  ? 'test-action-warning'
                  : 'test-action-default'
          ]"
        >
          <Icon
            v-if="status === 'connecting'"
            name="refresh"
            size="sm"
            class="animate-spin"
            :stroke-width="2"
          />
          <Icon v-else-if="status === 'idle'" name="play" size="sm" :stroke-width="2" />
          <Icon v-else name="refresh" size="sm" :stroke-width="2" />
          <span>
            {{
              status === 'connecting'
                ? t('admin.accounts.testing')
                : status === 'idle'
                  ? t('admin.accounts.startTest')
                  : t('admin.accounts.retry')
            }}
          </span>
        </button>
      </div>
    </template>
  </BaseDialog>
</template>

<script setup lang="ts">
import { computed, ref, watch, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import BaseDialog from '@/components/common/BaseDialog.vue'
import Select from '@/components/common/Select.vue'
import TextArea from '@/components/common/TextArea.vue'
import { Icon } from '@/components/icons'
import { useClipboard } from '@/composables/useClipboard'
import { buildApiUrl } from '@/api/client'
import { adminAPI } from '@/api/admin'
import type { Account, ClaudeModel } from '@/types'

const { t } = useI18n()
const { copyToClipboard } = useClipboard()

interface OutputLine {
  text: string
  class: string
}

interface PreviewImage {
  url: string
  mimeType?: string
}

const props = defineProps<{
  show: boolean
  account: Account | null
}>()

const emit = defineEmits<{
  (e: 'close'): void
}>()

const terminalRef = ref<HTMLElement | null>(null)
const status = ref<'idle' | 'connecting' | 'success' | 'error'>('idle')
const outputLines = ref<OutputLine[]>([])
const streamingContent = ref('')
const errorMessage = ref('')
const availableModels = ref<ClaudeModel[]>([])
const selectedModelId = ref('')
const testPrompt = ref('')
const loadingModels = ref(false)
let abortController: AbortController | null = null
const generatedImages = ref<PreviewImage[]>([])
const testMode = ref<'default' | 'compact'>('default')
const isOpenAIAccount = computed(() => props.account?.platform === 'openai')
const openAITestModeOptions = computed(() => [
  { value: 'default', label: t('admin.accounts.openai.testModeDefault') },
  { value: 'compact', label: t('admin.accounts.openai.testModeCompact') }
])
const previewImageUrl = ref('')
const prioritizedGeminiModels = ['gemini-3.1-flash-image', 'gemini-2.5-flash-image', 'gemini-3.5-flash', 'gemini-2.5-flash', 'gemini-2.5-pro', 'gemini-3-flash-preview', 'gemini-3-pro-preview', 'gemini-2.0-flash']
const supportsGeminiImageTest = computed(() => {
  const modelID = selectedModelId.value.toLowerCase()
  if (!modelID.startsWith('gemini-') || !modelID.includes('-image')) return false

  return props.account?.platform === 'gemini' || (props.account?.platform === 'antigravity' && props.account?.type === 'apikey')
})

const supportsOpenAIImageTest = computed(() => {
  const modelID = selectedModelId.value.toLowerCase()
  if (!modelID.startsWith('gpt-image-')) return false
  return props.account?.platform === 'openai'
})

const supportsImageTest = computed(() => supportsGeminiImageTest.value || supportsOpenAIImageTest.value)

const sortTestModels = (models: ClaudeModel[]) => {
  const priorityMap = new Map(prioritizedGeminiModels.map((id, index) => [id, index]))

  return [...models].sort((a, b) => {
    const aPriority = priorityMap.get(a.id) ?? Number.MAX_SAFE_INTEGER
    const bPriority = priorityMap.get(b.id) ?? Number.MAX_SAFE_INTEGER
    if (aPriority !== bPriority) return aPriority - bPriority
    return 0
  })
}

// Load available models when modal opens
watch(
  () => props.show,
  async (newVal) => {
    if (newVal && props.account) {
      testPrompt.value = ''
      testMode.value = 'default'
      resetState()
      await loadAvailableModels()
    } else {
      abortStream()
    }
  }
)

watch(selectedModelId, () => {
  if (supportsImageTest.value && !testPrompt.value.trim()) {
    testPrompt.value = t('admin.accounts.imagePromptDefault')
  }
})

const loadAvailableModels = async () => {
  if (!props.account) return

  loadingModels.value = true
  selectedModelId.value = '' // Reset selection before loading
  try {
    const models = await adminAPI.accounts.getAvailableModels(props.account.id)
    availableModels.value = props.account.platform === 'gemini' || props.account.platform === 'antigravity'
      ? sortTestModels(models)
      : models
    // Default selection by platform
    if (availableModels.value.length > 0) {
      if (props.account.platform === 'gemini') {
        selectedModelId.value = availableModels.value[0].id
      } else {
        // Try to select Sonnet as default, otherwise use first model
        const sonnetModel = availableModels.value.find((m) => m.id.includes('sonnet'))
        selectedModelId.value = sonnetModel?.id || availableModels.value[0].id
      }
    }
  } catch (error) {
    console.error('Failed to load available models:', error)
    // Fallback to empty list
    availableModels.value = []
    selectedModelId.value = ''
  } finally {
    loadingModels.value = false
  }
}

const resetState = () => {
  status.value = 'idle'
  outputLines.value = []
  streamingContent.value = ''
  errorMessage.value = ''
  generatedImages.value = []
  previewImageUrl.value = ''
}

const handleClose = () => {
  abortStream()
  emit('close')
}

const abortStream = () => {
  if (abortController) {
    abortController.abort()
    abortController = null
  }
}

const addLine = (text: string, className: string = 'text-gray-300') => {
  outputLines.value.push({ text, class: className })
  scrollToBottom()
}

const scrollToBottom = async () => {
  await nextTick()
  if (terminalRef.value) {
    terminalRef.value.scrollTop = terminalRef.value.scrollHeight
  }
}

const startTest = async () => {
  if (!props.account || !selectedModelId.value) return

  resetState()
  status.value = 'connecting'
  addLine(t('admin.accounts.startingTestForAccount', { name: props.account.name }), 'text-blue-400')
  addLine(t('admin.accounts.testAccountTypeLabel', { type: props.account.type }), 'text-gray-400')
  addLine('', 'text-gray-300')

  abortStream()

  abortController = new AbortController()

  try {
    // Use the configured API base; EventSource does not support POST.
    const url = buildApiUrl(`/admin/accounts/${props.account.id}/test`)

    // Use fetch with streaming for SSE since EventSource doesn't support POST
    const response = await fetch(url, {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${localStorage.getItem('auth_token')}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        model_id: selectedModelId.value,
        prompt: supportsImageTest.value ? testPrompt.value.trim() : '',
        mode: isOpenAIAccount.value ? testMode.value : 'default'
      }),
      signal: abortController.signal
    })

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }

    const reader = response.body?.getReader()
    if (!reader) {
      throw new Error('No response body')
    }

    const decoder = new TextDecoder()
    let buffer = ''

    while (true) {
      const { done, value } = await reader.read()
      if (done) break

      buffer += decoder.decode(value, { stream: true })
      const lines = buffer.split('\n')
      buffer = lines.pop() || ''

      for (const line of lines) {
        if (line.startsWith('data: ')) {
          const jsonStr = line.slice(6).trim()
          if (jsonStr) {
            try {
              const event = JSON.parse(jsonStr)
              handleEvent(event)
            } catch (e) {
              console.error('Failed to parse SSE event:', e)
            }
          }
        }
      }
    }
  } catch (error: unknown) {
    if (error instanceof DOMException && error.name === 'AbortError') {
      status.value = 'idle'
      return
    }
    status.value = 'error'
    const msg = error instanceof Error ? error.message : 'Unknown error'
    errorMessage.value = msg
    addLine(`Error: ${msg}`, 'text-red-400')
  }
}

const handleEvent = (event: {
  type: string
  text?: string
  model?: string
  success?: boolean
  error?: string
  image_url?: string
  mime_type?: string
}) => {
  switch (event.type) {
    case 'test_start':
      addLine(t('admin.accounts.connectedToApi'), 'text-green-400')
      if (event.model) {
        addLine(t('admin.accounts.usingModel', { model: event.model }), 'text-cyan-400')
      }
      addLine(
        supportsImageTest.value
            ? t('admin.accounts.sendingImageRequest')
            : t('admin.accounts.sendingTestMessage'),
        'text-gray-400'
      )
      addLine('', 'text-gray-300')
      addLine(t('admin.accounts.response'), 'text-yellow-400')
      break

    case 'content':
      if (event.text) {
        streamingContent.value += event.text
        scrollToBottom()
      }
      break

    case 'status':
      if (event.text) {
        addLine(event.text, 'text-cyan-300')
      }
      break

    case 'image':
      if (event.image_url) {
        generatedImages.value.push({
          url: event.image_url,
          mimeType: event.mime_type
        })
        addLine(t('admin.accounts.imageReceived', { count: generatedImages.value.length }), 'text-purple-300')
      }
      break

    case 'test_complete':
      // Move streaming content to output lines
      if (streamingContent.value) {
        addLine(streamingContent.value, 'text-green-300')
        streamingContent.value = ''
      }
      if (event.success) {
        status.value = 'success'
      } else {
        status.value = 'error'
        errorMessage.value = event.error || 'Test failed'
      }
      break

    case 'error':
      status.value = 'error'
      errorMessage.value = event.error || 'Unknown error'
      if (streamingContent.value) {
        addLine(streamingContent.value, 'text-green-300')
        streamingContent.value = ''
      }
      break
  }
}

const copyOutput = () => {
  const text = outputLines.value.map((l) => l.text).join('\n')
  copyToClipboard(text, t('admin.accounts.outputCopied'))
}
</script>

<style scoped>
.account-test-stack {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.test-account-summary {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  padding: 0.75rem;
  border: 1px solid var(--omnio-border, #e5e7eb);
  border-radius: 0.65rem;
  background: color-mix(in srgb, var(--omnio-foreground, #111827) 2.5%, var(--omnio-surface, #fff));
}

.test-account-icon {
  display: inline-flex;
  width: 2.25rem;
  height: 2.25rem;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  border: 1px solid color-mix(in srgb, var(--omnio-primary, #3b82f6) 16%, transparent);
  border-radius: 0.5rem;
  background: color-mix(in srgb, var(--omnio-primary, #3b82f6) 8%, transparent);
}

.platform-badge,
.test-status {
  display: inline-flex;
  height: 1.25rem;
  align-items: center;
  border: 1px solid var(--omnio-border, #e5e7eb);
  border-radius: 9999px;
  padding: 0 0.4rem;
  font-size: 0.65rem;
  line-height: 1;
  font-weight: 560;
  letter-spacing: 0.035em;
  text-transform: uppercase;
}

.platform-badge,
.test-status-neutral {
  color: var(--omnio-muted, #6b7280);
  background: color-mix(in srgb, var(--omnio-foreground, #111827) 4%, transparent);
}

.test-status-active {
  color: #047857;
  border-color: rgba(16, 185, 129, 0.18);
  background: rgba(16, 185, 129, 0.09);
}

.test-terminal {
  border: 1px solid color-mix(in srgb, var(--omnio-foreground, #111827) 16%, transparent);
  border-radius: 0.65rem;
  color: #d1d5db;
  background: #111214;
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.035);
}

.terminal-copy-button {
  position: absolute;
  top: 0.5rem;
  right: 0.5rem;
  display: inline-flex;
  width: 1.75rem;
  height: 1.75rem;
  align-items: center;
  justify-content: center;
  border-radius: 0.45rem;
  color: #9ca3af;
  background: rgba(17, 24, 39, 0.82);
  opacity: 0;
  outline: none;
  transition: color 140ms ease, background-color 140ms ease, opacity 140ms ease, box-shadow 140ms ease;
}

.group:hover .terminal-copy-button,
.terminal-copy-button:focus-visible {
  opacity: 1;
}

.terminal-copy-button:hover {
  color: #fff;
  background: #1f2937;
}

.terminal-copy-button:focus-visible {
  box-shadow: 0 0 0 3px rgba(96, 165, 250, 0.3);
}

.test-image-card {
  position: relative;
  overflow: hidden;
  cursor: pointer;
  border: 1px solid var(--omnio-border, #e5e7eb);
  border-radius: 0.65rem;
  background: var(--omnio-surface, #fff);
  box-shadow: none;
  transition: border-color 140ms ease, background-color 140ms ease;
}

.test-image-card:hover {
  border-color: var(--omnio-border-strong, #d1d5db);
  background: color-mix(in srgb, var(--omnio-foreground, #111827) 2.5%, var(--omnio-surface, #fff));
}

.lightbox-close {
  position: absolute;
  top: 1rem;
  right: 1rem;
  display: inline-flex;
  width: 2rem;
  height: 2rem;
  align-items: center;
  justify-content: center;
  border-radius: 9999px;
  color: #fff;
  background: rgba(0, 0, 0, 0.48);
  outline: none;
  transition: background-color 140ms ease, box-shadow 140ms ease;
}

.lightbox-close:hover {
  background: rgba(0, 0, 0, 0.68);
}

.lightbox-close:focus-visible {
  box-shadow: 0 0 0 3px rgba(255, 255, 255, 0.3);
}

.test-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.75rem;
  padding: 0 0.25rem;
  color: var(--omnio-muted, #6b7280);
  font-size: 0.75rem;
}

.test-footer-actions {
  display: flex;
  width: 100%;
  justify-content: flex-end;
  gap: 0.5rem;
}

.test-action-button {
  display: inline-flex;
  min-height: 2.25rem;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  border: 1px solid transparent;
  border-radius: 0.55rem;
  padding: 0.45rem 0.85rem;
  color: #fff;
  font-size: 0.8rem;
  font-weight: 550;
  outline: none;
  transition: background-color 140ms ease, border-color 140ms ease, box-shadow 140ms ease;
}

.test-action-default {
  border-color: var(--omnio-primary-strong, #2563eb);
  background: var(--omnio-primary-strong, #2563eb);
}

.test-action-success {
  border-color: #059669;
  background: #059669;
}

.test-action-warning {
  border-color: #d97706;
  background: #d97706;
}

.test-action-disabled {
  cursor: not-allowed;
  border-color: color-mix(in srgb, var(--omnio-primary, #3b82f6) 55%, transparent);
  background: color-mix(in srgb, var(--omnio-primary, #3b82f6) 55%, var(--omnio-surface, #fff));
}

.test-action-button:focus-visible {
  box-shadow: 0 0 0 3px color-mix(in srgb, var(--omnio-primary, #3b82f6) 25%, transparent);
}

:global(.dark) .test-status-active {
  color: #6ee7b7;
}

@media (max-width: 640px) {
  .test-account-summary {
    align-items: flex-start;
    flex-wrap: wrap;
  }

  .test-meta {
    align-items: flex-start;
    flex-direction: column;
  }

  .test-footer-actions {
    flex-direction: column-reverse;
  }

  .test-footer-actions > button {
    width: 100%;
  }
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
