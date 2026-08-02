<template>
  <AppLayout>
    <main class="image-studio-page">
      <section class="image-studio-hero">
        <div>
          <span class="image-studio-eyebrow">{{ t('imageGeneration.eyebrow') }}</span>
          <h1>{{ t('imageGeneration.title') }}</h1>
          <p>{{ t('imageGeneration.description') }}</p>
        </div>
        <div class="image-studio-balance">
          <span>{{ t('imageGeneration.balance') }}</span>
          <strong>{{ formattedBalance }}</strong>
          <small>{{ t('imageGeneration.billingHint') }}</small>
        </div>
      </section>

      <div class="image-studio-layout">
        <section class="image-studio-card image-studio-composer">
          <div class="image-studio-card-heading">
            <span class="image-studio-heading-icon">
              <Icon name="sparkles" size="md" />
            </span>
            <div>
              <h2>{{ t('imageGeneration.composer.title') }}</h2>
              <p>{{ t('imageGeneration.composer.asyncHint') }}</p>
            </div>
          </div>

          <div v-if="loadingKeys" class="image-studio-key-loading">
            <span></span><span></span><span></span>
          </div>

          <div v-else-if="imageApiKeys.length === 0" class="image-studio-no-key">
            <span class="image-studio-no-key-icon"><Icon name="key" size="lg" /></span>
            <h3>{{ t('imageGeneration.composer.noKeyTitle') }}</h3>
            <p>{{ t('imageGeneration.composer.noKeyDescription') }}</p>
            <router-link to="/keys" class="btn btn-primary">
              <Icon name="plus" size="sm" class="mr-1.5" />
              {{ t('imageGeneration.composer.createKey') }}
            </router-link>
          </div>

          <form v-else class="image-studio-form" @submit.prevent="generateImages">
            <div class="image-studio-field">
              <label for="image-studio-key">{{ t('imageGeneration.composer.apiKey') }}</label>
              <select id="image-studio-key" v-model.number="form.apiKeyId" class="input">
                <option :value="0" disabled>{{ t('imageGeneration.composer.apiKeyPlaceholder') }}</option>
                <option v-for="key in imageApiKeys" :key="key.id" :value="key.id">
                  {{ apiKeyLabel(key) }}
                </option>
              </select>
            </div>

            <div class="image-studio-field">
              <div class="image-studio-label-row">
                <label for="image-studio-model">{{ t('imageGeneration.composer.model') }}</label>
                <span v-if="loadingModels" class="image-studio-inline-status">
                  <i></i>{{ t('imageGeneration.composer.modelLoading') }}
                </span>
              </div>
              <select
                id="image-studio-model"
                v-model="form.model"
                class="input font-mono"
                :disabled="loadingModels"
              >
                <option value="" disabled>{{ t('imageGeneration.composer.modelPlaceholder') }}</option>
                <option v-for="model in availableModels" :key="model" :value="model">
                  {{ model }}
                </option>
              </select>
              <p class="image-studio-field-hint">{{ modelHint }}</p>
            </div>

            <div class="image-studio-field">
              <div class="image-studio-label-row">
                <label for="image-studio-prompt">{{ t('imageGeneration.composer.prompt') }}</label>
                <span>{{ form.prompt.length.toLocaleString() }} / 8,000</span>
              </div>
              <textarea
                id="image-studio-prompt"
                v-model="form.prompt"
                maxlength="8000"
                rows="7"
                :placeholder="t('imageGeneration.composer.promptPlaceholder')"
              ></textarea>
            </div>

            <fieldset class="image-studio-field">
              <legend>{{ t('imageGeneration.composer.size') }}</legend>
              <div class="image-studio-size-grid">
                <button
                  v-for="option in sizeOptions"
                  :key="option.value"
                  type="button"
                  class="image-studio-size-option"
                  :class="{ 'is-active': form.size === option.value }"
                  @click="form.size = option.value"
                >
                  <i :style="{ aspectRatio: option.ratio }"></i>
                  <span>{{ option.label }}</span>
                  <small>{{ option.detail }}</small>
                </button>
              </div>
            </fieldset>

            <fieldset class="image-studio-field">
              <legend>{{ t('imageGeneration.composer.count') }}</legend>
              <div class="image-studio-count-options">
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

            <button type="submit" class="btn btn-primary image-studio-generate" :disabled="!canGenerate">
              <Icon name="sparkles" size="md" :class="generating ? 'animate-spin' : ''" />
              {{ generating ? t('imageGeneration.composer.generating') : t('imageGeneration.composer.generate') }}
            </button>
          </form>
        </section>

        <section class="image-studio-card image-studio-results">
          <div class="image-studio-results-heading">
            <div>
              <span>{{ t('imageGeneration.result.title') }}</span>
              <strong v-if="activeRun">{{ activeRun.model }}</strong>
            </div>
            <span v-if="activeRun" class="image-studio-mode-badge">
              {{ activeRun.mode === 'async' ? t('imageGeneration.result.modeAsync') : t('imageGeneration.result.modeSync') }}
            </span>
          </div>

          <div v-if="errorMessage" class="image-studio-error" role="alert">
            <Icon name="exclamationCircle" size="md" />
            <span>{{ errorMessage }}</span>
          </div>

          <div v-if="generating" class="image-studio-generating" aria-live="polite">
            <div class="image-studio-generation-visual" aria-hidden="true">
              <span class="image-studio-generation-core"><Icon name="sparkles" size="xl" /></span>
              <i v-for="index in 4" :key="index" :style="{ '--delay': `${index * 0.18}s` }"></i>
            </div>
            <h3>{{ generationStatusLabel }}</h3>
            <p>{{ form.prompt }}</p>
          </div>

          <div v-else-if="!activeRun" class="image-studio-empty">
            <div class="image-studio-empty-canvas" aria-hidden="true">
              <span><Icon name="sparkles" size="xl" /></span>
            </div>
            <h3>{{ t('imageGeneration.result.emptyTitle') }}</h3>
            <p>{{ t('imageGeneration.result.emptyDescription') }}</p>
          </div>

          <div v-else class="image-studio-result-content">
            <div class="image-studio-result-meta">
              <div>
                <span class="badge badge-success">{{ t('imageGeneration.result.completed') }}</span>
                <small>{{ t('imageGeneration.result.generatedAt', { time: formatRunTime(activeRun.createdAt) }) }}</small>
              </div>
              <p>{{ activeRun.prompt }}</p>
            </div>

            <div class="image-studio-gallery" :class="{ 'is-single': activeRun.outputs.length === 1 }">
              <article v-for="(output, index) in activeRun.outputs" :key="output.id" class="image-studio-output">
                <a :href="output.src" target="_blank" rel="noopener noreferrer" class="image-studio-output-image">
                  <img :src="output.src" :alt="activeRun.prompt" loading="lazy" />
                </a>
                <div class="image-studio-output-actions">
                  <a :href="output.src" target="_blank" rel="noopener noreferrer" class="btn btn-secondary btn-sm">
                    <Icon name="externalLink" size="sm" />
                    {{ t('imageGeneration.result.open') }}
                  </a>
                  <button type="button" class="btn btn-primary btn-sm" @click="downloadOutput(output, index)">
                    <Icon name="download" size="sm" />
                    {{ t('imageGeneration.result.download') }}
                  </button>
                </div>
                <details v-if="output.revisedPrompt" class="image-studio-revised-prompt">
                  <summary>{{ t('imageGeneration.result.revisedPrompt') }}</summary>
                  <p>{{ output.revisedPrompt }}</p>
                </details>
              </article>
            </div>
          </div>
        </section>
      </div>

      <section v-if="runs.length > 1" class="image-studio-history image-studio-card">
        <div class="image-studio-history-heading">
          <h2>{{ t('imageGeneration.result.sessionHistory') }}</h2>
          <button type="button" class="btn btn-secondary btn-sm" @click="clearHistory">
            <Icon name="trash" size="sm" />
            {{ t('imageGeneration.result.clearHistory') }}
          </button>
        </div>
        <div class="image-studio-history-list">
          <button
            v-for="run in runs"
            :key="run.id"
            type="button"
            :class="{ 'is-active': run.id === activeRunId }"
            @click="activeRunId = run.id"
          >
            <img :src="run.outputs[0]?.src" :alt="run.prompt" loading="lazy" />
            <span>
              <strong>{{ run.model }}</strong>
              <small>{{ run.prompt }}</small>
            </span>
          </button>
        </div>
      </section>

      <router-link to="/batch-image" class="image-studio-batch-link">
        <Icon name="grid" size="sm" />
        {{ t('imageGeneration.result.batchLink') }}
        <Icon name="arrowRight" size="sm" />
      </router-link>
    </main>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'
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
  { value: '1024x1024', label: '1:1', detail: '1024 × 1024', ratio: '1 / 1' },
  { value: '1536x1024', label: '3:2', detail: '1536 × 1024', ratio: '3 / 2' },
  { value: '1024x1536', label: '2:3', detail: '1024 × 1536', ratio: '2 / 3' },
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
const activeRunId = ref('')
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

const activeRun = computed(() =>
  runs.value.find((run) => run.id === activeRunId.value) || runs.value[0] || null,
)

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

function apiKeyLabel(key: ApiKey): string {
  const groupName = key.group?.name || `#${key.group_id || '-'}`
  const platform = String(key.group?.platform || '').toUpperCase()
  return `${key.name || `API Key #${key.id}`} · ${groupName}${platform ? ` · ${platform}` : ''}`
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
    activeRunId.value = run.id
    await authStore.refreshUser().catch(() => undefined)
  } catch (error) {
    if ((error as Error)?.name === 'AbortError') return
    errorMessage.value = friendlyGenerationError(error)
  } finally {
    generating.value = false
  }
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
  const active = activeRun.value
  runs.value = active ? [active] : []
  activeRunId.value = active?.id || ''
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
.image-studio-page {
  width: 100%;
  max-width: 1600px;
  margin: 0 auto;
  padding: 24px;
}

.image-studio-hero {
  position: relative;
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 28px;
  overflow: hidden;
  margin-bottom: 20px;
  padding: 28px 30px;
  border: 1px solid rgb(224 231 255);
  border-radius: 18px;
  background:
    radial-gradient(circle at 88% 18%, rgb(129 140 248 / 18%), transparent 34%),
    linear-gradient(135deg, rgb(255 255 255), rgb(248 250 255));
}

.dark .image-studio-hero {
  border-color: rgb(51 65 85);
  background:
    radial-gradient(circle at 88% 18%, rgb(99 102 241 / 24%), transparent 34%),
    linear-gradient(135deg, rgb(15 23 42), rgb(17 24 39));
}

.image-studio-eyebrow {
  display: inline-flex;
  margin-bottom: 8px;
  color: rgb(79 70 229);
  font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.16em;
}

.dark .image-studio-eyebrow { color: rgb(165 180 252); }
.image-studio-hero h1 { margin: 0; color: rgb(17 24 39); font-size: clamp(28px, 3vw, 40px); font-weight: 750; letter-spacing: -0.035em; }
.dark .image-studio-hero h1 { color: white; }
.image-studio-hero p { max-width: 680px; margin: 8px 0 0; color: rgb(107 114 128); font-size: 15px; line-height: 1.7; }
.dark .image-studio-hero p { color: rgb(148 163 184); }

.image-studio-balance {
  display: grid;
  min-width: 220px;
  padding: 16px 18px;
  border: 1px solid rgb(224 231 255);
  border-radius: 14px;
  background: rgb(255 255 255 / 78%);
  backdrop-filter: blur(12px);
}
.dark .image-studio-balance { border-color: rgb(71 85 105); background: rgb(15 23 42 / 72%); }
.image-studio-balance span { color: rgb(107 114 128); font-size: 12px; }
.image-studio-balance strong { margin-top: 2px; color: rgb(17 24 39); font-size: 24px; }
.dark .image-studio-balance strong { color: white; }
.image-studio-balance small { margin-top: 5px; color: rgb(107 114 128); font-size: 11px; line-height: 1.45; }

.image-studio-layout {
  display: grid;
  grid-template-columns: minmax(330px, 430px) minmax(0, 1fr);
  gap: 20px;
  align-items: start;
}

.image-studio-card {
  border: 1px solid rgb(229 231 235);
  border-radius: 18px;
  background: white;
  box-shadow: 0 14px 40px rgb(15 23 42 / 5%);
}
.dark .image-studio-card { border-color: rgb(51 65 85); background: rgb(17 24 39); box-shadow: 0 16px 45px rgb(0 0 0 / 20%); }

.image-studio-composer { position: sticky; top: 20px; padding: 22px; }
.image-studio-card-heading { display: flex; gap: 12px; align-items: flex-start; padding-bottom: 18px; border-bottom: 1px solid rgb(243 244 246); }
.dark .image-studio-card-heading { border-color: rgb(51 65 85); }
.image-studio-heading-icon { display: grid; place-items: center; width: 38px; height: 38px; flex: 0 0 auto; border-radius: 11px; color: rgb(79 70 229); background: rgb(238 242 255); }
.dark .image-studio-heading-icon { color: rgb(165 180 252); background: rgb(49 46 129 / 38%); }
.image-studio-card-heading h2, .image-studio-history h2 { margin: 0; color: rgb(17 24 39); font-size: 16px; font-weight: 700; }
.dark .image-studio-card-heading h2, .dark .image-studio-history h2 { color: white; }
.image-studio-card-heading p { margin: 4px 0 0; color: rgb(107 114 128); font-size: 12px; line-height: 1.5; }

.image-studio-form { display: grid; gap: 18px; margin-top: 20px; }
.image-studio-field { display: grid; gap: 7px; min-width: 0; padding: 0; border: 0; }
.image-studio-field label, .image-studio-field legend { color: rgb(55 65 81); font-size: 12px; font-weight: 650; }
.dark .image-studio-field label, .dark .image-studio-field legend { color: rgb(203 213 225); }
.image-studio-label-row { display: flex; justify-content: space-between; gap: 12px; color: rgb(156 163 175); font-size: 11px; }
.image-studio-field-hint { margin: 0; color: rgb(107 114 128); font-size: 11px; line-height: 1.45; }
.image-studio-inline-status { display: inline-flex; align-items: center; gap: 5px; }
.image-studio-inline-status i { width: 6px; height: 6px; border-radius: 999px; background: rgb(99 102 241); animation: studio-pulse 1.2s ease-in-out infinite; }
.image-studio-field textarea { width: 100%; min-height: 150px; resize: vertical; padding: 13px 14px; border: 1px solid rgb(209 213 219); border-radius: 12px; color: rgb(17 24 39); background: white; font-size: 14px; line-height: 1.6; transition: border-color 160ms ease, box-shadow 160ms ease; }
.image-studio-field textarea:focus { outline: 0; border-color: rgb(99 102 241); box-shadow: 0 0 0 3px rgb(99 102 241 / 12%); }
.dark .image-studio-field textarea { border-color: rgb(71 85 105); color: white; background: rgb(15 23 42); }

.image-studio-size-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 8px; }
.image-studio-size-option { display: grid; justify-items: center; gap: 5px; padding: 11px 6px; border: 1px solid rgb(229 231 235); border-radius: 11px; color: rgb(75 85 99); background: transparent; transition: 150ms ease; }
.dark .image-studio-size-option { border-color: rgb(51 65 85); color: rgb(148 163 184); }
.image-studio-size-option:hover { border-color: rgb(165 180 252); background: rgb(238 242 255 / 50%); }
.dark .image-studio-size-option:hover { background: rgb(49 46 129 / 18%); }
.image-studio-size-option.is-active { border-color: rgb(99 102 241); color: rgb(67 56 202); background: rgb(238 242 255); box-shadow: inset 0 0 0 1px rgb(99 102 241 / 25%); }
.dark .image-studio-size-option.is-active { color: rgb(199 210 254); background: rgb(49 46 129 / 34%); }
.image-studio-size-option i { display: block; width: 22px; max-height: 28px; border: 1.5px solid currentColor; border-radius: 3px; }
.image-studio-size-option span { font-size: 12px; font-weight: 700; }
.image-studio-size-option small { font-size: 9px; opacity: 0.68; }

.image-studio-count-options { display: grid; grid-template-columns: repeat(4, 1fr); gap: 8px; }
.image-studio-count-options button { height: 38px; border: 1px solid rgb(229 231 235); border-radius: 10px; color: rgb(75 85 99); background: transparent; font-size: 13px; font-weight: 650; }
.dark .image-studio-count-options button { border-color: rgb(51 65 85); color: rgb(148 163 184); }
.image-studio-count-options button.is-active { border-color: rgb(99 102 241); color: rgb(67 56 202); background: rgb(238 242 255); }
.dark .image-studio-count-options button.is-active { color: rgb(199 210 254); background: rgb(49 46 129 / 34%); }
.image-studio-generate { width: 100%; min-height: 46px; gap: 8px; margin-top: 2px; }

.image-studio-key-loading { display: grid; gap: 12px; margin-top: 22px; }
.image-studio-key-loading span { height: 42px; border-radius: 10px; background: linear-gradient(90deg, rgb(243 244 246), rgb(249 250 251), rgb(243 244 246)); background-size: 200% 100%; animation: studio-shimmer 1.5s linear infinite; }
.dark .image-studio-key-loading span { background: linear-gradient(90deg, rgb(30 41 59), rgb(51 65 85), rgb(30 41 59)); background-size: 200% 100%; }
.image-studio-key-loading span:nth-child(2) { height: 160px; }

.image-studio-no-key { display: grid; justify-items: center; padding: 52px 12px 28px; text-align: center; }
.image-studio-no-key-icon { display: grid; place-items: center; width: 50px; height: 50px; margin-bottom: 14px; border-radius: 14px; color: rgb(107 114 128); background: rgb(243 244 246); }
.dark .image-studio-no-key-icon { color: rgb(148 163 184); background: rgb(30 41 59); }
.image-studio-no-key h3 { margin: 0; color: rgb(17 24 39); font-size: 15px; }
.dark .image-studio-no-key h3 { color: white; }
.image-studio-no-key p { max-width: 310px; margin: 8px 0 18px; color: rgb(107 114 128); font-size: 12px; line-height: 1.6; }

.image-studio-results { min-height: 670px; overflow: hidden; }
.image-studio-results-heading { display: flex; align-items: center; justify-content: space-between; min-height: 70px; padding: 16px 20px; border-bottom: 1px solid rgb(243 244 246); }
.dark .image-studio-results-heading { border-color: rgb(51 65 85); }
.image-studio-results-heading > div { display: grid; gap: 2px; }
.image-studio-results-heading span { color: rgb(107 114 128); font-size: 12px; }
.image-studio-results-heading strong { color: rgb(17 24 39); font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-size: 13px; }
.dark .image-studio-results-heading strong { color: white; }
.image-studio-mode-badge { padding: 5px 9px; border: 1px solid rgb(224 231 255); border-radius: 999px; color: rgb(79 70 229) !important; background: rgb(238 242 255); font-weight: 650; }
.dark .image-studio-mode-badge { border-color: rgb(67 56 202 / 50%); color: rgb(199 210 254) !important; background: rgb(49 46 129 / 28%); }
.image-studio-error { display: flex; align-items: flex-start; gap: 9px; margin: 16px 18px 0; padding: 11px 12px; border: 1px solid rgb(254 202 202); border-radius: 10px; color: rgb(185 28 28); background: rgb(254 242 242); font-size: 12px; line-height: 1.5; }
.dark .image-studio-error { border-color: rgb(127 29 29); color: rgb(252 165 165); background: rgb(69 10 10 / 32%); }

.image-studio-empty, .image-studio-generating { display: grid; min-height: 590px; place-content: center; justify-items: center; padding: 36px; text-align: center; }
.image-studio-empty-canvas { position: relative; display: grid; place-items: center; width: 138px; height: 138px; margin-bottom: 22px; border: 1px solid rgb(224 231 255); border-radius: 28px; background: radial-gradient(circle at 30% 30%, rgb(199 210 254), rgb(238 242 255) 42%, rgb(249 250 251) 72%); transform: rotate(-3deg); }
.dark .image-studio-empty-canvas { border-color: rgb(67 56 202 / 45%); background: radial-gradient(circle at 30% 30%, rgb(67 56 202), rgb(30 41 59) 46%, rgb(15 23 42) 76%); }
.image-studio-empty-canvas span { display: grid; place-items: center; width: 60px; height: 60px; border-radius: 18px; color: rgb(79 70 229); background: rgb(255 255 255 / 72%); box-shadow: 0 12px 35px rgb(79 70 229 / 16%); }
.dark .image-studio-empty-canvas span { color: rgb(199 210 254); background: rgb(15 23 42 / 72%); }
.image-studio-empty h3, .image-studio-generating h3 { margin: 0; color: rgb(17 24 39); font-size: 17px; }
.dark .image-studio-empty h3, .dark .image-studio-generating h3 { color: white; }
.image-studio-empty p, .image-studio-generating p { max-width: 520px; margin: 7px 0 0; overflow: hidden; color: rgb(107 114 128); font-size: 13px; line-height: 1.6; text-overflow: ellipsis; }

.image-studio-generation-visual { position: relative; display: grid; place-items: center; width: 150px; height: 150px; margin-bottom: 24px; }
.image-studio-generation-core { z-index: 2; display: grid; place-items: center; width: 64px; height: 64px; border-radius: 20px; color: white; background: linear-gradient(135deg, rgb(79 70 229), rgb(124 58 237)); box-shadow: 0 0 42px rgb(99 102 241 / 42%); animation: studio-core 1.8s ease-in-out infinite; }
.image-studio-generation-visual i { --delay: 0s; position: absolute; inset: 18px; border: 1px solid rgb(165 180 252 / 56%); border-radius: 34px; animation: studio-orbit 2.8s linear infinite; animation-delay: var(--delay); }
.image-studio-generation-visual i:nth-of-type(2) { inset: 6px 30px; }
.image-studio-generation-visual i:nth-of-type(3) { inset: 30px 6px; }
.image-studio-generation-visual i:nth-of-type(4) { inset: 38px; border-style: dashed; animation-direction: reverse; }

.image-studio-result-content { padding: 18px; }
.image-studio-result-meta { display: flex; align-items: flex-start; justify-content: space-between; gap: 20px; margin-bottom: 14px; }
.image-studio-result-meta > div { display: flex; align-items: center; gap: 9px; flex: 0 0 auto; }
.image-studio-result-meta small { color: rgb(107 114 128); font-size: 11px; }
.image-studio-result-meta p { max-width: 560px; margin: 0; overflow: hidden; color: rgb(75 85 99); font-size: 12px; line-height: 1.55; text-align: right; text-overflow: ellipsis; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; }
.dark .image-studio-result-meta p { color: rgb(148 163 184); }
.image-studio-gallery { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 14px; }
.image-studio-gallery.is-single { grid-template-columns: minmax(0, 1fr); }
.image-studio-output { min-width: 0; overflow: hidden; border: 1px solid rgb(229 231 235); border-radius: 14px; background: rgb(249 250 251); }
.dark .image-studio-output { border-color: rgb(51 65 85); background: rgb(15 23 42); }
.image-studio-output-image { display: grid; min-height: 260px; place-items: center; overflow: hidden; background: repeating-conic-gradient(rgb(243 244 246) 0 25%, white 0 50%) 0 / 20px 20px; }
.dark .image-studio-output-image { background: repeating-conic-gradient(rgb(30 41 59) 0 25%, rgb(15 23 42) 0 50%) 0 / 20px 20px; }
.image-studio-gallery.is-single .image-studio-output-image { min-height: 520px; }
.image-studio-output-image img { display: block; width: 100%; max-height: 680px; object-fit: contain; }
.image-studio-output-actions { display: flex; justify-content: flex-end; gap: 8px; padding: 10px; }
.image-studio-output-actions .btn { gap: 6px; }
.image-studio-revised-prompt { margin: 0 10px 10px; border-top: 1px solid rgb(229 231 235); padding-top: 8px; }
.dark .image-studio-revised-prompt { border-color: rgb(51 65 85); }
.image-studio-revised-prompt summary { cursor: pointer; color: rgb(79 70 229); font-size: 11px; font-weight: 650; }
.dark .image-studio-revised-prompt summary { color: rgb(165 180 252); }
.image-studio-revised-prompt p { margin: 7px 0 0; color: rgb(75 85 99); font-size: 11px; line-height: 1.55; }
.dark .image-studio-revised-prompt p { color: rgb(148 163 184); }

.image-studio-history { margin-top: 20px; padding: 18px; }
.image-studio-history-heading { display: flex; align-items: center; justify-content: space-between; gap: 16px; margin-bottom: 13px; }
.image-studio-history-heading .btn { gap: 6px; }
.image-studio-history-list { display: grid; grid-template-columns: repeat(auto-fill, minmax(220px, 1fr)); gap: 10px; }
.image-studio-history-list > button { display: flex; min-width: 0; align-items: center; gap: 10px; padding: 8px; border: 1px solid rgb(229 231 235); border-radius: 12px; color: inherit; background: transparent; text-align: left; }
.dark .image-studio-history-list > button { border-color: rgb(51 65 85); }
.image-studio-history-list > button:hover, .image-studio-history-list > button.is-active { border-color: rgb(99 102 241); background: rgb(238 242 255 / 55%); }
.dark .image-studio-history-list > button:hover, .dark .image-studio-history-list > button.is-active { background: rgb(49 46 129 / 22%); }
.image-studio-history-list img { width: 54px; height: 54px; flex: 0 0 auto; border-radius: 9px; object-fit: cover; background: rgb(243 244 246); }
.image-studio-history-list span { display: grid; min-width: 0; gap: 4px; }
.image-studio-history-list strong { overflow: hidden; color: rgb(17 24 39); font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-size: 11px; text-overflow: ellipsis; white-space: nowrap; }
.dark .image-studio-history-list strong { color: white; }
.image-studio-history-list small { overflow: hidden; color: rgb(107 114 128); font-size: 10px; text-overflow: ellipsis; white-space: nowrap; }

.image-studio-batch-link { display: flex; width: fit-content; align-items: center; gap: 7px; margin: 18px auto 0; color: rgb(79 70 229); font-size: 12px; font-weight: 650; }
.dark .image-studio-batch-link { color: rgb(165 180 252); }

@keyframes studio-pulse { 50% { opacity: 0.35; transform: scale(0.82); } }
@keyframes studio-shimmer { to { background-position: -200% 0; } }
@keyframes studio-core { 50% { transform: scale(1.06); box-shadow: 0 0 58px rgb(99 102 241 / 55%); } }
@keyframes studio-orbit { to { transform: rotate(360deg); } }

@media (max-width: 1100px) {
  .image-studio-layout { grid-template-columns: 360px minmax(0, 1fr); }
  .image-studio-gallery { grid-template-columns: minmax(0, 1fr); }
}

@media (max-width: 900px) {
  .image-studio-page { padding: 16px; }
  .image-studio-hero { align-items: stretch; flex-direction: column; padding: 22px; }
  .image-studio-balance { min-width: 0; }
  .image-studio-layout { grid-template-columns: minmax(0, 1fr); }
  .image-studio-composer { position: static; }
  .image-studio-results { min-height: 560px; }
  .image-studio-empty, .image-studio-generating { min-height: 480px; }
}

@media (max-width: 640px) {
  .image-studio-page { padding: 12px; }
  .image-studio-hero { padding: 19px; border-radius: 15px; }
  .image-studio-hero h1 { font-size: 28px; }
  .image-studio-composer { padding: 16px; }
  .image-studio-card { border-radius: 15px; }
  .image-studio-size-grid { grid-template-columns: 1fr; }
  .image-studio-size-option { grid-template-columns: 28px 48px 1fr; justify-items: start; align-items: center; padding: 9px 12px; }
  .image-studio-size-option i { justify-self: center; }
  .image-studio-result-meta { flex-direction: column; }
  .image-studio-result-meta p { text-align: left; }
  .image-studio-result-content { padding: 12px; }
  .image-studio-gallery.is-single .image-studio-output-image { min-height: 300px; }
  .image-studio-output-image { min-height: 240px; }
  .image-studio-output-actions { display: grid; grid-template-columns: 1fr 1fr; }
  .image-studio-output-actions .btn { justify-content: center; }
}

@media (prefers-reduced-motion: reduce) {
  .image-studio-generation-core,
  .image-studio-generation-visual i,
  .image-studio-inline-status i,
  .image-studio-key-loading span { animation: none; }
}
</style>
