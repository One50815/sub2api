<template>
  <AppLayout>
    <div class="image-playground">
      <header class="image-playground-toolbar">
        <div class="image-playground-title">
          <span class="image-playground-title-icon" aria-hidden="true">
            <Icon name="sparkles" size="sm" />
          </span>
          <div>
            <h1>{{ t('imageGeneration.title') }}</h1>
            <p>{{ t('imageGeneration.description') }}</p>
          </div>
        </div>

        <div class="image-playground-toolbar-actions">
          <div class="image-playground-balance" :title="t('imageGeneration.billingHint')">
            <span>{{ t('imageGeneration.balance') }}</span>
            <strong>{{ formattedBalance }}</strong>
          </div>
          <router-link to="/batch-image" class="image-playground-toolbar-button">
            <Icon name="grid" size="sm" />
            <span>{{ t('imageGeneration.result.batchLink') }}</span>
          </router-link>
          <button
            v-if="runs.length > 0"
            type="button"
            class="image-playground-icon-button image-playground-clear-button"
            :title="t('imageGeneration.result.clearHistory')"
            :aria-label="t('imageGeneration.result.clearHistory')"
            @click="clearHistory"
          >
            <Icon name="trash" size="sm" />
          </button>
        </div>
      </header>

      <section ref="canvasRef" class="image-playground-canvas" aria-live="polite">
        <div v-if="loadingKeys" class="image-playground-loading">
          <span class="image-playground-loading-icon"><Icon name="sparkles" size="lg" /></span>
          <div class="image-playground-loading-copy">
            <i></i>
            <i></i>
          </div>
        </div>

        <div v-else-if="imageApiKeys.length === 0" class="image-playground-empty-state">
          <span class="image-playground-empty-icon"><Icon name="key" size="lg" /></span>
          <h2>{{ t('imageGeneration.composer.noKeyTitle') }}</h2>
          <p>{{ t('imageGeneration.composer.noKeyDescription') }}</p>
          <router-link to="/keys" class="btn btn-primary">
            <Icon name="plus" size="sm" />
            {{ t('imageGeneration.composer.createKey') }}
          </router-link>
        </div>

        <div v-else class="image-playground-stream">
          <div v-if="runs.length === 0 && !generating" class="image-playground-empty-state">
            <span class="image-playground-empty-icon"><Icon name="sparkles" size="lg" /></span>
            <h2>{{ t('imageGeneration.result.emptyTitle') }}</h2>
            <p>{{ t('imageGeneration.result.emptyDescription') }}</p>
            <div class="image-playground-starters">
              <button
                v-for="promptKey in starterPromptKeys"
                :key="promptKey"
                type="button"
                @click="useStarterPrompt(promptKey)"
              >
                <Icon name="arrowRight" size="xs" />
                <span>{{ t(promptKey) }}</span>
              </button>
            </div>
          </div>

          <article v-for="run in chronologicalRuns" :key="run.id" class="image-playground-run">
            <header class="image-playground-run-header">
              <div class="image-playground-run-prompt">
                <span>{{ t('imageGeneration.result.promptLabel') }}</span>
                <p>{{ run.prompt }}</p>
              </div>
              <div class="image-playground-run-actions">
                <button
                  type="button"
                  class="image-playground-action-button"
                  @click="reuseRun(run)"
                >
                  <Icon name="refresh" size="sm" />
                  {{ t('imageGeneration.result.reusePrompt') }}
                </button>
                <button
                  v-if="run.outputs.length > 1"
                  type="button"
                  class="image-playground-action-button"
                  @click="downloadAllOutputs(run)"
                >
                  <Icon name="download" size="sm" />
                  {{ t('imageGeneration.result.downloadAll') }}
                </button>
              </div>
            </header>

            <div class="image-playground-run-meta">
              <span>{{ run.model }}</span>
              <span>{{ run.size }}</span>
              <span>{{ run.mode === 'async' ? t('imageGeneration.result.modeAsync') : t('imageGeneration.result.modeSync') }}</span>
              <span>{{ formatRunTime(run.createdAt) }}</span>
            </div>

            <div class="image-playground-gallery" :class="{ 'is-single': run.outputs.length === 1 }">
              <article v-for="(output, index) in run.outputs" :key="output.id" class="image-playground-output">
                <a :href="output.src" target="_blank" rel="noopener noreferrer" class="image-playground-output-image">
                  <img :src="output.src" :alt="run.prompt" loading="lazy" />
                </a>
                <div class="image-playground-output-footer">
                  <span>{{ t('imageGeneration.result.imageNumber', { number: index + 1 }) }}</span>
                  <div>
                    <a
                      :href="output.src"
                      target="_blank"
                      rel="noopener noreferrer"
                      class="image-playground-icon-button"
                      :title="t('imageGeneration.result.open')"
                      :aria-label="t('imageGeneration.result.open')"
                    >
                      <Icon name="externalLink" size="sm" />
                    </a>
                    <button
                      type="button"
                      class="image-playground-icon-button"
                      :title="t('imageGeneration.result.download')"
                      :aria-label="t('imageGeneration.result.download')"
                      @click="downloadOutput(output, index)"
                    >
                      <Icon name="download" size="sm" />
                    </button>
                  </div>
                </div>
                <details v-if="output.revisedPrompt" class="image-playground-revised-prompt">
                  <summary>{{ t('imageGeneration.result.revisedPrompt') }}</summary>
                  <p>{{ output.revisedPrompt }}</p>
                </details>
              </article>
            </div>
          </article>

          <article v-if="generating" class="image-playground-run image-playground-run-generating">
            <header class="image-playground-run-header">
              <div class="image-playground-run-prompt">
                <span>{{ t('imageGeneration.result.promptLabel') }}</span>
                <p>{{ form.prompt }}</p>
              </div>
              <span class="image-playground-generation-status">
                <i></i>
                {{ generationStatusLabel }}
              </span>
            </header>
            <div class="image-playground-gallery" :class="{ 'is-single': form.count === 1 }">
              <div v-for="index in form.count" :key="index" class="image-playground-output-skeleton">
                <span><Icon name="sparkles" size="lg" /></span>
              </div>
            </div>
          </article>

          <div v-if="errorMessage" class="image-playground-error" role="alert">
            <Icon name="exclamationCircle" size="md" />
            <div>
              <strong>{{ t('imageGeneration.errors.title') }}</strong>
              <p>{{ errorMessage }}</p>
            </div>
          </div>
        </div>
      </section>

      <footer v-if="!loadingKeys && imageApiKeys.length > 0" class="image-playground-composer-shell">
        <form class="image-playground-composer" @submit.prevent="generateImages">
          <Transition name="playground-parameters">
            <div v-if="showParameters" class="image-playground-parameters">
              <fieldset>
                <legend>{{ t('imageGeneration.composer.size') }}</legend>
                <div class="image-playground-size-options">
                  <button
                    v-for="option in sizeOptions"
                    :key="option.value"
                    type="button"
                    :class="{ 'is-active': form.size === option.value }"
                    @click="form.size = option.value"
                  >
                    <i :style="{ aspectRatio: option.ratio }"></i>
                    <span>{{ option.label }}</span>
                    <small>{{ option.detail }}</small>
                  </button>
                </div>
              </fieldset>

              <fieldset>
                <legend>{{ t('imageGeneration.composer.count') }}</legend>
                <div class="image-playground-count-options">
                  <button
                    v-for="count in 4"
                    :key="count"
                    type="button"
                    :class="{ 'is-active': form.count === count }"
                    @click="form.count = count"
                  >
                    {{ count }}
                  </button>
                </div>
              </fieldset>

              <div class="image-playground-parameter-note">
                <Icon name="infoCircle" size="sm" />
                <span>{{ t('imageGeneration.composer.asyncHint') }}</span>
              </div>
            </div>
          </Transition>

          <div class="image-playground-prompt-field">
            <label for="image-playground-prompt" class="sr-only">
              {{ t('imageGeneration.composer.prompt') }}
            </label>
            <textarea
              id="image-playground-prompt"
              ref="promptInputRef"
              v-model="form.prompt"
              maxlength="8000"
              rows="3"
              :placeholder="t('imageGeneration.composer.promptPlaceholder')"
              @keydown.meta.enter.prevent="generateImages"
              @keydown.ctrl.enter.prevent="generateImages"
            ></textarea>
            <span>{{ form.prompt.length.toLocaleString() }} / 8,000</span>
          </div>

          <div class="image-playground-composer-footer">
            <div class="image-playground-tool-controls">
              <button
                type="button"
                class="image-playground-settings-button"
                :class="{ 'is-active': showParameters }"
                :aria-expanded="showParameters"
                @click="showParameters = !showParameters"
              >
                <Icon name="cog" size="sm" />
                <span>{{ t('imageGeneration.composer.parameters') }}</span>
                <small>{{ selectedSizeLabel }} / {{ form.count }}</small>
              </button>
            </div>

            <div class="image-playground-selectors">
              <label class="image-playground-select-control">
                <span class="sr-only">{{ t('imageGeneration.composer.apiKey') }}</span>
                <Icon name="key" size="sm" />
                <select v-model.number="form.apiKeyId" :aria-label="t('imageGeneration.composer.apiKey')">
                  <option :value="0" disabled>{{ t('imageGeneration.composer.apiKeyPlaceholder') }}</option>
                  <option v-for="key in imageApiKeys" :key="key.id" :value="key.id">
                    {{ apiKeyLabel(key) }}
                  </option>
                </select>
                <Icon name="chevronDown" size="xs" />
              </label>

              <label class="image-playground-select-control" :class="{ 'is-loading': loadingModels }">
                <span class="sr-only">{{ t('imageGeneration.composer.model') }}</span>
                <Icon name="cpu" size="sm" />
                <select
                  v-model="form.model"
                  :disabled="loadingModels"
                  :aria-label="t('imageGeneration.composer.model')"
                  :title="modelHint"
                >
                  <option value="" disabled>{{ t('imageGeneration.composer.modelPlaceholder') }}</option>
                  <option v-for="model in availableModels" :key="model" :value="model">
                    {{ model }}
                  </option>
                </select>
                <Icon name="chevronDown" size="xs" />
              </label>

              <button
                v-if="generating"
                type="button"
                class="image-playground-submit is-stopping"
                @click="cancelGeneration"
              >
                <Icon name="x" size="sm" />
                <span>{{ t('imageGeneration.composer.stop') }}</span>
              </button>
              <button v-else type="submit" class="image-playground-submit" :disabled="!canGenerate">
                <Icon name="sparkles" size="sm" />
                <span>{{ t('imageGeneration.composer.generate') }}</span>
              </button>
            </div>
          </div>
        </form>
        <p class="image-playground-billing-hint">{{ t('imageGeneration.billingHint') }}</p>
      </footer>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import { keysAPI } from '@/api'
import {
  ImageGenerationAPIError,
  getImageGenerationTask,
  imageOutputSource,
  imageTaskResult,
  isImageGenerationModel,
  listImageGenerationModels,
  mergeImageGenerationModels,
  submitImageGeneration,
  type GeneratedImageOutput,
  type ImageGenerationResponse,
} from '@/api/imageGeneration'
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'
import { sanitizeUrl } from '@/utils/url'
import type { ApiKey } from '@/types'

type GenerationStage = 'submitting' | 'processing' | 'finalizing'

interface StudioOutput {
  id: string
  src: string
  revisedPrompt: string
}

interface StudioRun {
  id: string
  prompt: string
  model: string
  size: string
  mode: 'async' | 'sync'
  createdAt: number
  outputs: StudioOutput[]
}

const { t, locale } = useI18n()
const appStore = useAppStore()
const authStore = useAuthStore()

const form = reactive({
  apiKeyId: 0,
  model: '',
  prompt: '',
  size: '1024x1024',
  count: 1,
})

const sizeOptions = [
  { value: '1024x1024', label: '1:1', detail: '1024 x 1024', ratio: '1 / 1' },
  { value: '1536x1024', label: '3:2', detail: '1536 x 1024', ratio: '3 / 2' },
  { value: '1024x1536', label: '2:3', detail: '1024 x 1536', ratio: '2 / 3' },
]

const starterPromptKeys = [
  'imageGeneration.starters.cinematic',
  'imageGeneration.starters.product',
  'imageGeneration.starters.illustration',
  'imageGeneration.starters.architecture',
]

const apiKeys = ref<ApiKey[]>([])
const availableModels = ref<string[]>([])
const loadingKeys = ref(true)
const loadingModels = ref(false)
const modelLoadFailed = ref(false)
const generating = ref(false)
const generationStage = ref<GenerationStage>('submitting')
const errorMessage = ref('')
const runs = ref<StudioRun[]>([])
const showParameters = ref(false)
const canvasRef = ref<HTMLElement | null>(null)
const promptInputRef = ref<HTMLTextAreaElement | null>(null)
let generationController: AbortController | null = null
let modelController: AbortController | null = null
let modelRequestSequence = 0

const imageApiKeys = computed(() => apiKeys.value.filter((key) => {
  const platform = String(key.group?.platform || '')
  return key.status === 'active' &&
    key.group?.allow_image_generation === true &&
    ['openai', 'grok', 'composite'].includes(platform)
}))

const selectedApiKey = computed(() =>
  imageApiKeys.value.find((key) => key.id === Number(form.apiKeyId)) || null,
)

const chronologicalRuns = computed(() => [...runs.value].reverse())

const canGenerate = computed(() =>
  !generating.value && !!selectedApiKey.value && !!form.model.trim() && !!form.prompt.trim(),
)

const formattedBalance = computed(() => {
  const balance = Number(authStore.user?.balance || 0)
  return `$${balance.toFixed(2)}`
})

const generationStatusLabel = computed(() => t(`imageGeneration.result.${generationStage.value}`))

const modelHint = computed(() => {
  if (modelLoadFailed.value) return t('imageGeneration.errors.loadModels')
  return t('imageGeneration.composer.modelHint')
})

const selectedSizeLabel = computed(() =>
  sizeOptions.find((option) => option.value === form.size)?.label || form.size,
)

function apiKeyLabel(key: ApiKey): string {
  const groupName = key.group?.name || `#${key.group_id || '-'}`
  const platform = String(key.group?.platform || '').toUpperCase()
  return `${key.name || `API Key #${key.id}`} / ${groupName}${platform ? ` / ${platform}` : ''}`
}

function fallbackModelsForKey(key: ApiKey | null): string[] {
  const platform = String(key?.group?.platform || '')
  if (platform === 'grok') return ['grok-imagine', 'grok-imagine-image']
  if (platform === 'composite') return ['gpt-image-2', 'grok-imagine']
  return ['gpt-image-2', 'gpt-image-1.5', 'gpt-image-1']
}

async function loadApiKeys() {
  loadingKeys.value = true
  try {
    const response = await keysAPI.list(1, 100, {
      status: 'active',
      sort_by: 'created_at',
      sort_order: 'desc',
    })
    apiKeys.value = response.items || []
    if (!selectedApiKey.value && imageApiKeys.value.length > 0) {
      form.apiKeyId = imageApiKeys.value[0].id
    }
  } catch {
    appStore.showError(t('imageGeneration.errors.loadKeys'))
  } finally {
    loadingKeys.value = false
  }
}

async function loadModels() {
  modelController?.abort()
  modelController = new AbortController()
  const key = selectedApiKey.value
  const requestSequence = ++modelRequestSequence
  modelLoadFailed.value = false
  availableModels.value = []
  if (!key) {
    form.model = ''
    return
  }

  loadingModels.value = true
  try {
    const response = await listImageGenerationModels(key.key, modelController.signal)
    if (requestSequence !== modelRequestSequence) return
    const discovered = (response.data || [])
      .map((model) => String(model?.id || '').trim())
      .filter(isImageGenerationModel)
    availableModels.value = mergeImageGenerationModels(fallbackModelsForKey(key), discovered)
  } catch (error) {
    if ((error as Error)?.name === 'AbortError') return
    if (requestSequence !== modelRequestSequence) return
    modelLoadFailed.value = true
    availableModels.value = fallbackModelsForKey(key)
  } finally {
    if (requestSequence === modelRequestSequence) loadingModels.value = false
  }

  if (!availableModels.value.includes(form.model)) {
    form.model = availableModels.value[0] || ''
  }
}

function normalizeOutputs(response: ImageGenerationResponse): StudioOutput[] {
  return (response.data || [])
    .map((output: GeneratedImageOutput, index) => ({
      id: `${Date.now()}-${index}`,
      src: sanitizeUrl(imageOutputSource(output, 'png'), { allowDataUrl: true }),
      revisedPrompt: String(output.revised_prompt || '').trim(),
    }))
    .filter((output) => !!output.src)
}

function waitForNextPoll(signal: AbortSignal): Promise<void> {
  return new Promise((resolve, reject) => {
    const handleAbort = () => {
      window.clearTimeout(timer)
      reject(new DOMException('Aborted', 'AbortError'))
    }
    const timer = window.setTimeout(() => {
      signal.removeEventListener('abort', handleAbort)
      resolve()
    }, 3000)
    signal.addEventListener('abort', handleAbort, { once: true })
  })
}

async function pollImageTask(apiKey: string, taskId: string, signal: AbortSignal): Promise<ImageGenerationResponse> {
  const deadline = Date.now() + 30 * 60 * 1000
  while (Date.now() < deadline) {
    const task = await getImageGenerationTask(apiKey, taskId, signal)
    if (task.status === 'completed') {
      const result = imageTaskResult(task)
      if (result) return result
      throw new ImageGenerationAPIError(t('imageGeneration.errors.noResult'), task.http_status || 0)
    }
    if (task.status === 'failed') {
      throw new ImageGenerationAPIError(
        task.error?.message || t('imageGeneration.errors.taskFailed'),
        task.http_status || 0,
        task.error?.code || task.error?.type || '',
      )
    }
    await waitForNextPoll(signal)
  }
  throw new ImageGenerationAPIError(t('imageGeneration.errors.taskFailed'), 408, 'IMAGE_TASK_TIMEOUT')
}

async function generateImages() {
  if (generating.value) return
  const key = selectedApiKey.value
  if (!key) {
    appStore.showError(t('imageGeneration.errors.selectKey'))
    return
  }
  if (!form.model.trim()) {
    appStore.showError(t('imageGeneration.errors.modelRequired'))
    return
  }
  if (!form.prompt.trim()) {
    appStore.showError(t('imageGeneration.errors.promptRequired'))
    return
  }

  generationController?.abort()
  generationController = new AbortController()
  generating.value = true
  generationStage.value = 'submitting'
  errorMessage.value = ''
  await scrollToLatest()

  try {
    const submission = await submitImageGeneration(key.key, {
      model: form.model.trim(),
      prompt: form.prompt.trim(),
      n: form.count,
      size: form.size,
      response_format: 'b64_json',
      output_format: 'png',
    }, {
      signal: generationController.signal,
      // Composite routing is resolved by the synchronous gateway handler.
      // The existing async task handler currently accepts concrete providers only.
      preferAsync: String(key.group?.platform || '') !== 'composite',
    })

    let result: ImageGenerationResponse
    if (submission.mode === 'async') {
      generationStage.value = 'processing'
      result = await pollImageTask(key.key, submission.task.task_id || submission.task.id, generationController.signal)
    } else {
      generationStage.value = 'finalizing'
      result = submission.result
    }

    generationStage.value = 'finalizing'
    const outputs = normalizeOutputs(result)
    if (outputs.length === 0) throw new ImageGenerationAPIError(t('imageGeneration.errors.noResult'))

    const run: StudioRun = {
      id: `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`,
      prompt: form.prompt.trim(),
      model: form.model.trim(),
      size: form.size,
      mode: submission.mode,
      createdAt: Date.now(),
      outputs,
    }
    runs.value = [run, ...runs.value].slice(0, 8)
    await authStore.refreshUser().catch(() => undefined)
    await scrollToLatest()
  } catch (error) {
    if ((error as Error)?.name === 'AbortError') return
    errorMessage.value = friendlyGenerationError(error)
    await scrollToLatest()
  } finally {
    generating.value = false
    generationController = null
  }
}

function cancelGeneration() {
  generationController?.abort()
}

function friendlyGenerationError(error: unknown): string {
  const apiError = error as ImageGenerationAPIError
  const code = String(apiError?.code || '').toUpperCase()
  const message = String(apiError?.message || '')
  if (apiError?.status === 401) return t('imageGeneration.errors.unauthorized')
  if (apiError?.status === 403) return t('imageGeneration.errors.forbidden')
  if (code.includes('INSUFFICIENT') || /insufficient balance|quota/i.test(message)) {
    return t('imageGeneration.errors.insufficientBalance')
  }
  if (message && !/^HTTP \d+$/i.test(message)) return message
  return t('imageGeneration.errors.generic')
}

async function downloadOutput(output: StudioOutput, index: number) {
  const fileName = `omnio-image-${Date.now()}-${index + 1}.png`
  try {
    const response = await fetch(output.src)
    if (!response.ok) throw new Error('download failed')
    const blob = await response.blob()
    const objectURL = URL.createObjectURL(blob)
    triggerDownload(objectURL, fileName)
    window.setTimeout(() => URL.revokeObjectURL(objectURL), 1000)
  } catch {
    window.open(output.src, '_blank', 'noopener,noreferrer')
  }
}

async function downloadAllOutputs(run: StudioRun) {
  for (const [index, output] of run.outputs.entries()) {
    await downloadOutput(output, index)
  }
}

function triggerDownload(url: string, fileName: string) {
  const link = document.createElement('a')
  link.href = url
  link.download = fileName
  link.rel = 'noopener noreferrer'
  document.body.appendChild(link)
  link.click()
  link.remove()
}

function clearHistory() {
  runs.value = []
  errorMessage.value = ''
}

function useStarterPrompt(promptKey: string) {
  form.prompt = t(promptKey)
  focusPromptInput()
}

function reuseRun(run: StudioRun) {
  form.prompt = run.prompt
  form.size = run.size
  form.count = Math.min(Math.max(run.outputs.length, 1), 4)
  if (availableModels.value.includes(run.model)) form.model = run.model
  focusPromptInput()
}

function focusPromptInput() {
  void nextTick(() => {
    promptInputRef.value?.focus()
  })
}

async function scrollToLatest() {
  await nextTick()
  if (!canvasRef.value) return
  const behavior = window.matchMedia('(prefers-reduced-motion: reduce)').matches ? 'auto' : 'smooth'
  canvasRef.value.scrollTo({ top: canvasRef.value.scrollHeight, behavior })
}

function formatRunTime(timestamp: number): string {
  return new Intl.DateTimeFormat(locale.value, {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  }).format(new Date(timestamp))
}

watch(() => form.apiKeyId, () => {
  void loadModels()
})

onMounted(() => {
  void loadApiKeys()
})

onBeforeUnmount(() => {
  generationController?.abort()
  modelController?.abort()
})
</script>

<style scoped>
.image-playground {
  --playground-accent: var(--omnio-primary-strong);
  --playground-control: color-mix(in srgb, var(--omnio-surface) 96%, var(--omnio-bg));
  display: grid;
  width: 100%;
  height: 100%;
  min-height: 620px;
  grid-template-rows: auto minmax(0, 1fr) auto;
  overflow: hidden;
  border: 1px solid var(--omnio-border);
  border-radius: 16px;
  color: var(--omnio-foreground);
  background: var(--omnio-bg);
}

.image-playground-toolbar {
  display: flex;
  min-width: 0;
  min-height: 64px;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  padding: 10px 14px 10px 16px;
  border-bottom: 1px solid var(--omnio-border);
  background: color-mix(in srgb, var(--omnio-surface) 94%, transparent);
}

.image-playground-title,
.image-playground-toolbar-actions,
.image-playground-run-actions,
.image-playground-output-footer,
.image-playground-output-footer > div,
.image-playground-composer-footer,
.image-playground-tool-controls,
.image-playground-selectors {
  display: flex;
  align-items: center;
}

.image-playground-title {
  min-width: 0;
  gap: 10px;
}

.image-playground-title-icon,
.image-playground-empty-icon,
.image-playground-loading-icon {
  display: grid;
  flex: 0 0 auto;
  place-items: center;
  border: 1px solid var(--omnio-border);
  color: var(--playground-accent);
  background: var(--omnio-surface-subtle);
}

.image-playground-title-icon {
  width: 34px;
  height: 34px;
  border-radius: 9px;
}

.image-playground-title h1 {
  margin: 0;
  color: var(--omnio-foreground);
  font-size: 15px;
  font-weight: 650;
  line-height: 1.3;
}

.image-playground-title p {
  margin: 2px 0 0;
  overflow: hidden;
  color: var(--omnio-muted);
  font-size: 11px;
  line-height: 1.35;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.image-playground-toolbar-actions {
  flex: 0 0 auto;
  gap: 8px;
}

.image-playground-balance {
  display: flex;
  align-items: baseline;
  gap: 6px;
  padding-right: 10px;
  border-right: 1px solid var(--omnio-border);
  white-space: nowrap;
}

.image-playground-balance span {
  color: var(--omnio-muted);
  font-size: 11px;
}

.image-playground-balance strong {
  color: var(--omnio-foreground);
  font-size: 13px;
  font-weight: 650;
  font-variant-numeric: tabular-nums;
}

.image-playground-toolbar-button,
.image-playground-icon-button,
.image-playground-action-button,
.image-playground-settings-button {
  border: 1px solid var(--omnio-border);
  color: color-mix(in srgb, var(--omnio-foreground) 76%, transparent);
  background: transparent;
  transition: border-color 150ms ease, color 150ms ease, background-color 150ms ease, transform 120ms ease;
}

.image-playground-toolbar-button:hover,
.image-playground-icon-button:hover,
.image-playground-action-button:hover,
.image-playground-settings-button:hover,
.image-playground-settings-button.is-active {
  border-color: var(--omnio-border-strong);
  color: var(--omnio-foreground);
  background: color-mix(in srgb, var(--omnio-foreground) 4.5%, transparent);
}

.image-playground-toolbar-button:active,
.image-playground-icon-button:active,
.image-playground-action-button:active,
.image-playground-settings-button:active,
.image-playground-submit:active:not(:disabled) {
  transform: scale(0.98);
}

.image-playground-toolbar-button {
  display: inline-flex;
  min-height: 34px;
  align-items: center;
  gap: 7px;
  padding: 0 10px;
  border-radius: 9px;
  font-size: 11px;
  font-weight: 560;
}

.image-playground-icon-button {
  display: inline-grid;
  width: 32px;
  height: 32px;
  place-items: center;
  border-radius: 8px;
}

.image-playground-clear-button:hover {
  border-color: color-mix(in srgb, #dc2626 34%, var(--omnio-border));
  color: #dc2626;
  background: color-mix(in srgb, #dc2626 8%, transparent);
}

.image-playground-canvas {
  min-height: 0;
  overflow-x: hidden;
  overflow-y: auto;
  overscroll-behavior: contain;
  background: color-mix(in srgb, var(--omnio-bg) 94%, var(--omnio-surface-subtle));
}

.image-playground-stream {
  width: min(100%, 1040px);
  min-height: 100%;
  margin: 0 auto;
  padding: 28px 28px 48px;
}

.image-playground-empty-state,
.image-playground-loading {
  display: grid;
  width: min(100% - 32px, 680px);
  min-height: 100%;
  margin: 0 auto;
  place-content: center;
  justify-items: center;
  padding: 44px 0;
  text-align: center;
}

.image-playground-empty-icon,
.image-playground-loading-icon {
  width: 48px;
  height: 48px;
  margin-bottom: 16px;
  border-radius: 12px;
}

.image-playground-empty-state h2 {
  margin: 0;
  color: var(--omnio-foreground);
  font-size: 21px;
  font-weight: 650;
  line-height: 1.25;
}

.image-playground-empty-state > p {
  max-width: 520px;
  margin: 8px 0 20px;
  color: var(--omnio-muted);
  font-size: 13px;
  line-height: 1.6;
}

.image-playground-starters {
  display: grid;
  width: min(100%, 620px);
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 8px;
}

.image-playground-starters button {
  display: flex;
  min-width: 0;
  min-height: 46px;
  align-items: center;
  gap: 8px;
  padding: 9px 12px;
  border: 1px solid var(--omnio-border);
  border-radius: 10px;
  color: color-mix(in srgb, var(--omnio-foreground) 78%, transparent);
  background: var(--omnio-surface);
  font-size: 12px;
  line-height: 1.45;
  text-align: left;
  transition: border-color 150ms ease, background-color 150ms ease, transform 120ms ease;
}

.image-playground-starters button:hover {
  border-color: var(--omnio-border-strong);
  background: color-mix(in srgb, var(--playground-accent) 5%, var(--omnio-surface));
}

.image-playground-starters button:active {
  transform: scale(0.99);
}

.image-playground-starters svg {
  flex: 0 0 auto;
  color: var(--playground-accent);
}

.image-playground-run {
  padding: 26px 0 30px;
  border-top: 1px solid var(--omnio-border);
}

.image-playground-run:first-child {
  padding-top: 8px;
  border-top: 0;
}

.image-playground-run-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 20px;
  margin-bottom: 10px;
}

.image-playground-run-prompt {
  min-width: 0;
}

.image-playground-run-prompt > span {
  display: block;
  margin-bottom: 4px;
  color: var(--omnio-muted);
  font-size: 10px;
  font-weight: 620;
}

.image-playground-run-prompt p {
  max-width: 760px;
  margin: 0;
  color: var(--omnio-foreground);
  font-size: 14px;
  line-height: 1.6;
  white-space: pre-wrap;
}

.image-playground-run-actions {
  flex: 0 0 auto;
  gap: 6px;
}

.image-playground-action-button {
  display: inline-flex;
  min-height: 32px;
  align-items: center;
  gap: 6px;
  padding: 0 9px;
  border-radius: 8px;
  font-size: 11px;
  white-space: nowrap;
}

.image-playground-run-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 6px 12px;
  margin-bottom: 14px;
  color: var(--omnio-muted);
  font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
  font-size: 10px;
}

.image-playground-run-meta span + span::before {
  content: '/';
  margin-right: 12px;
  color: color-mix(in srgb, var(--omnio-muted) 38%, transparent);
}

.image-playground-gallery {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.image-playground-gallery.is-single {
  grid-template-columns: minmax(0, 1fr);
}

.image-playground-output {
  min-width: 0;
  overflow: hidden;
  border: 1px solid var(--omnio-border);
  border-radius: 12px;
  background: var(--omnio-surface);
}

.image-playground-output-image {
  display: grid;
  min-height: 250px;
  place-items: center;
  overflow: hidden;
  background-color: var(--omnio-surface-subtle);
  background-image:
    linear-gradient(45deg, color-mix(in srgb, var(--omnio-foreground) 4%, transparent) 25%, transparent 25%),
    linear-gradient(-45deg, color-mix(in srgb, var(--omnio-foreground) 4%, transparent) 25%, transparent 25%),
    linear-gradient(45deg, transparent 75%, color-mix(in srgb, var(--omnio-foreground) 4%, transparent) 75%),
    linear-gradient(-45deg, transparent 75%, color-mix(in srgb, var(--omnio-foreground) 4%, transparent) 75%);
  background-position: 0 0, 0 10px, 10px -10px, -10px 0;
  background-size: 20px 20px;
}

.image-playground-gallery.is-single .image-playground-output-image {
  min-height: 460px;
}

.image-playground-output-image img {
  display: block;
  width: 100%;
  max-height: 660px;
  object-fit: contain;
}

.image-playground-output-footer {
  min-height: 44px;
  justify-content: space-between;
  gap: 10px;
  padding: 6px 7px 6px 11px;
  border-top: 1px solid var(--omnio-border);
}

.image-playground-output-footer > span {
  color: var(--omnio-muted);
  font-size: 10px;
}

.image-playground-output-footer > div {
  gap: 5px;
}

.image-playground-revised-prompt {
  border-top: 1px solid var(--omnio-border);
  padding: 9px 11px 11px;
}

.image-playground-revised-prompt summary {
  cursor: pointer;
  color: var(--playground-accent);
  font-size: 11px;
  font-weight: 600;
}

.image-playground-revised-prompt p {
  margin: 7px 0 0;
  color: var(--omnio-muted);
  font-size: 11px;
  line-height: 1.6;
}

.image-playground-generation-status {
  display: inline-flex;
  flex: 0 0 auto;
  align-items: center;
  gap: 7px;
  color: var(--playground-accent);
  font-size: 11px;
  font-weight: 600;
}

.image-playground-generation-status i {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: currentColor;
  animation: playground-pulse 1.2s ease-in-out infinite;
}

.image-playground-output-skeleton {
  display: grid;
  min-height: 320px;
  place-items: center;
  overflow: hidden;
  border: 1px solid var(--omnio-border);
  border-radius: 12px;
  color: color-mix(in srgb, var(--playground-accent) 58%, transparent);
  background: linear-gradient(
    100deg,
    var(--omnio-surface-subtle) 20%,
    color-mix(in srgb, var(--playground-accent) 8%, var(--omnio-surface)) 42%,
    var(--omnio-surface-subtle) 64%
  );
  background-size: 240% 100%;
  animation: playground-shimmer 1.7s linear infinite;
}

.image-playground-output-skeleton span {
  display: grid;
  width: 48px;
  height: 48px;
  place-items: center;
  border: 1px solid color-mix(in srgb, var(--playground-accent) 18%, var(--omnio-border));
  border-radius: 12px;
  background: color-mix(in srgb, var(--omnio-surface) 88%, transparent);
}

.image-playground-error {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  margin-top: 20px;
  padding: 12px 13px;
  border: 1px solid color-mix(in srgb, #dc2626 28%, var(--omnio-border));
  border-radius: 10px;
  color: #b91c1c;
  background: color-mix(in srgb, #dc2626 6%, var(--omnio-surface));
}

.dark .image-playground-error {
  color: #fca5a5;
}

.image-playground-error strong {
  display: block;
  font-size: 12px;
}

.image-playground-error p {
  margin: 2px 0 0;
  font-size: 11px;
  line-height: 1.5;
}

.image-playground-loading {
  grid-template-columns: auto minmax(180px, 300px);
  gap: 14px;
  text-align: left;
}

.image-playground-loading-icon {
  margin: 0;
}

.image-playground-loading-copy {
  display: grid;
  width: 100%;
  gap: 8px;
}

.image-playground-loading-copy i {
  height: 10px;
  border-radius: 5px;
  background: var(--omnio-border);
  animation: playground-fade 1.4s ease-in-out infinite;
}

.image-playground-loading-copy i:last-child {
  width: 68%;
  animation-delay: 160ms;
}

.image-playground-composer-shell {
  position: relative;
  z-index: 2;
  padding: 12px 18px 10px;
  border-top: 1px solid var(--omnio-border);
  background: color-mix(in srgb, var(--omnio-bg) 92%, transparent);
  backdrop-filter: blur(18px) saturate(1.1);
}

.image-playground-composer {
  width: min(100%, 960px);
  margin: 0 auto;
  overflow: hidden;
  border: 1px solid var(--omnio-border-strong);
  border-radius: 14px;
  background: var(--playground-control);
  box-shadow: 0 18px 52px -36px color-mix(in srgb, var(--omnio-foreground) 50%, transparent);
  transition: border-color 160ms ease, box-shadow 160ms ease;
}

.image-playground-composer:focus-within {
  border-color: color-mix(in srgb, var(--playground-accent) 48%, var(--omnio-border));
  box-shadow: 0 0 0 3px color-mix(in srgb, var(--playground-accent) 10%, transparent), 0 20px 58px -38px color-mix(in srgb, var(--omnio-foreground) 56%, transparent);
}

.image-playground-parameters {
  display: grid;
  grid-template-columns: minmax(0, 1.25fr) minmax(180px, 0.75fr) minmax(200px, 1fr);
  gap: 16px;
  padding: 14px 16px;
  border-bottom: 1px solid var(--omnio-border);
  background: var(--omnio-surface-subtle);
}

.image-playground-parameters fieldset {
  display: grid;
  min-width: 0;
  gap: 7px;
  padding: 0;
  border: 0;
}

.image-playground-parameters legend {
  color: color-mix(in srgb, var(--omnio-foreground) 72%, transparent);
  font-size: 10px;
  font-weight: 620;
}

.image-playground-size-options,
.image-playground-count-options {
  display: grid;
  gap: 5px;
  padding: 3px;
  border: 1px solid var(--omnio-border);
  border-radius: 10px;
  background: var(--omnio-surface);
}

.image-playground-size-options {
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

.image-playground-count-options {
  grid-template-columns: repeat(4, minmax(0, 1fr));
}

.image-playground-size-options button,
.image-playground-count-options button {
  min-width: 0;
  min-height: 38px;
  border: 0;
  border-radius: 7px;
  color: var(--omnio-muted);
  background: transparent;
  font-size: 11px;
  transition: color 140ms ease, background-color 140ms ease, box-shadow 140ms ease;
}

.image-playground-size-options button {
  display: grid;
  grid-template-columns: auto auto;
  grid-template-rows: auto auto;
  align-items: center;
  justify-content: center;
  gap: 0 6px;
}

.image-playground-size-options button i {
  display: block;
  width: 13px;
  max-height: 17px;
  grid-row: 1 / 3;
  border: 1px solid currentColor;
  border-radius: 2px;
}

.image-playground-size-options button span {
  align-self: end;
  color: inherit;
  font-weight: 650;
  line-height: 1;
}

.image-playground-size-options button small {
  align-self: start;
  overflow: hidden;
  font-size: 8px;
  line-height: 1.2;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.image-playground-size-options button.is-active,
.image-playground-count-options button.is-active {
  color: var(--omnio-foreground);
  background: color-mix(in srgb, var(--playground-accent) 9%, var(--omnio-surface));
  box-shadow: inset 0 0 0 1px color-mix(in srgb, var(--playground-accent) 24%, transparent);
}

.image-playground-parameter-note {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  align-self: end;
  padding: 9px 10px;
  color: var(--omnio-muted);
  font-size: 10px;
  line-height: 1.45;
}

.image-playground-parameter-note svg {
  flex: 0 0 auto;
  color: var(--playground-accent);
}

.image-playground-prompt-field {
  position: relative;
}

.image-playground-prompt-field textarea {
  display: block;
  width: 100%;
  min-height: 84px;
  max-height: 180px;
  resize: vertical;
  padding: 14px 16px 24px;
  border: 0;
  outline: 0;
  color: var(--omnio-foreground);
  background: transparent;
  font-size: 14px;
  line-height: 1.65;
}

.image-playground-prompt-field textarea::placeholder {
  color: color-mix(in srgb, var(--omnio-muted) 72%, transparent);
}

.image-playground-prompt-field > span {
  position: absolute;
  right: 14px;
  bottom: 5px;
  color: color-mix(in srgb, var(--omnio-muted) 72%, transparent);
  font-size: 9px;
  font-variant-numeric: tabular-nums;
}

.image-playground-composer-footer {
  min-height: 50px;
  justify-content: space-between;
  gap: 12px;
  padding: 7px;
  border-top: 1px solid var(--omnio-border);
  background: color-mix(in srgb, var(--omnio-surface-subtle) 70%, transparent);
}

.image-playground-tool-controls,
.image-playground-selectors {
  min-width: 0;
  gap: 6px;
}

.image-playground-settings-button {
  display: inline-flex;
  min-height: 34px;
  align-items: center;
  gap: 7px;
  padding: 0 9px;
  border-radius: 8px;
  font-size: 11px;
  white-space: nowrap;
}

.image-playground-settings-button small {
  color: var(--omnio-muted);
  font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
  font-size: 9px;
}

.image-playground-select-control {
  display: flex;
  min-width: 0;
  height: 34px;
  align-items: center;
  gap: 6px;
  padding: 0 7px 0 9px;
  border: 1px solid var(--omnio-border);
  border-radius: 8px;
  color: var(--omnio-muted);
  background: var(--omnio-surface);
}

.image-playground-select-control:focus-within {
  border-color: color-mix(in srgb, var(--playground-accent) 48%, var(--omnio-border));
  box-shadow: 0 0 0 2px color-mix(in srgb, var(--playground-accent) 10%, transparent);
}

.image-playground-select-control:first-of-type {
  width: min(28vw, 250px);
}

.image-playground-select-control:nth-of-type(2) {
  width: min(24vw, 220px);
}

.image-playground-select-control select {
  min-width: 0;
  flex: 1;
  overflow: hidden;
  border: 0;
  outline: 0;
  color: var(--omnio-foreground);
  background: transparent;
  font-size: 10px;
  text-overflow: ellipsis;
  white-space: nowrap;
  appearance: none;
}

.image-playground-select-control.is-loading {
  opacity: 0.62;
}

.image-playground-submit {
  display: inline-flex;
  min-width: 104px;
  min-height: 34px;
  align-items: center;
  justify-content: center;
  gap: 7px;
  padding: 0 12px;
  border: 1px solid color-mix(in srgb, var(--playground-accent) 78%, #111827);
  border-radius: 8px;
  color: #ffffff;
  background: var(--playground-accent);
  font-size: 11px;
  font-weight: 620;
  transition: filter 150ms ease, opacity 150ms ease, transform 120ms ease;
}

.image-playground-submit:hover:not(:disabled) {
  filter: brightness(0.94);
}

.image-playground-submit:disabled {
  cursor: not-allowed;
  opacity: 0.42;
}

.image-playground-submit.is-stopping {
  border-color: color-mix(in srgb, #dc2626 68%, #111827);
  background: #b91c1c;
}

.image-playground-billing-hint {
  width: min(100%, 960px);
  margin: 5px auto 0;
  color: color-mix(in srgb, var(--omnio-muted) 74%, transparent);
  font-size: 9px;
  line-height: 1.4;
  text-align: center;
}

.playground-parameters-enter-active,
.playground-parameters-leave-active {
  transition: opacity 140ms ease, transform 140ms ease;
}

.playground-parameters-enter-from,
.playground-parameters-leave-to {
  opacity: 0;
  transform: translateY(4px);
}

@keyframes playground-pulse {
  50% { opacity: 0.35; transform: scale(0.82); }
}

@keyframes playground-shimmer {
  to { background-position: -180% 0; }
}

@keyframes playground-fade {
  50% { opacity: 0.42; }
}

@media (max-width: 1040px) {
  .image-playground-stream {
    padding-inline: 20px;
  }

  .image-playground-parameters {
    grid-template-columns: minmax(0, 1fr) minmax(170px, 0.65fr);
  }

  .image-playground-parameter-note {
    display: none;
  }

  .image-playground-select-control:first-of-type {
    width: min(31vw, 220px);
  }

  .image-playground-select-control:nth-of-type(2) {
    width: min(26vw, 190px);
  }
}

@media (max-width: 760px) {
  .image-playground {
    height: auto;
    min-height: calc(100svh - var(--omnio-header-height) - 2rem);
    overflow: visible;
    border: 0;
    border-radius: 0;
  }

  .image-playground-toolbar {
    position: sticky;
    top: var(--omnio-header-height);
    z-index: 3;
    min-height: 58px;
    padding-inline: 12px;
    backdrop-filter: blur(18px);
  }

  .image-playground-title p,
  .image-playground-balance span,
  .image-playground-toolbar-button span {
    display: none;
  }

  .image-playground-balance {
    padding-right: 8px;
  }

  .image-playground-toolbar-button {
    width: 34px;
    justify-content: center;
    padding: 0;
  }

  .image-playground-canvas {
    min-height: 480px;
    overflow: visible;
  }

  .image-playground-stream {
    padding: 22px 14px 34px;
  }

  .image-playground-run-header {
    flex-direction: column;
    gap: 10px;
  }

  .image-playground-run-actions {
    width: 100%;
  }

  .image-playground-gallery {
    grid-template-columns: minmax(0, 1fr);
  }

  .image-playground-gallery.is-single .image-playground-output-image,
  .image-playground-output-image {
    min-height: 280px;
  }

  .image-playground-composer-shell {
    position: sticky;
    bottom: 0;
    padding: 9px 8px calc(8px + env(safe-area-inset-bottom));
  }

  .image-playground-parameters {
    grid-template-columns: minmax(0, 1fr);
    gap: 12px;
  }

  .image-playground-composer-footer {
    align-items: stretch;
    flex-direction: column;
  }

  .image-playground-tool-controls,
  .image-playground-selectors {
    width: 100%;
  }

  .image-playground-settings-button {
    width: 100%;
    justify-content: flex-start;
  }

  .image-playground-settings-button small {
    margin-left: auto;
  }

  .image-playground-selectors {
    display: grid;
    grid-template-columns: minmax(0, 1fr) minmax(0, 1fr) auto;
  }

  .image-playground-select-control:first-of-type,
  .image-playground-select-control:nth-of-type(2) {
    width: auto;
  }

  .image-playground-select-control svg:first-of-type {
    display: none;
  }

  .image-playground-submit {
    min-width: 88px;
  }

  .image-playground-billing-hint {
    display: none;
  }
}

@media (max-width: 520px) {
  .image-playground-toolbar-actions {
    gap: 5px;
  }

  .image-playground-title-icon {
    display: none;
  }

  .image-playground-starters {
    grid-template-columns: minmax(0, 1fr);
  }

  .image-playground-run-actions {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .image-playground-action-button {
    justify-content: center;
  }

  .image-playground-selectors {
    grid-template-columns: minmax(0, 1fr) auto;
  }

  .image-playground-select-control:first-of-type {
    grid-column: 1 / -1;
  }

  .image-playground-size-options button small {
    display: none;
  }
}

@media (prefers-reduced-motion: reduce) {
  .image-playground-generation-status i,
  .image-playground-output-skeleton,
  .image-playground-loading-copy i {
    animation: none;
  }

  .playground-parameters-enter-active,
  .playground-parameters-leave-active {
    transition: none;
  }
}
</style>
