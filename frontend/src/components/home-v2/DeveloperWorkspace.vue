<template>
  <section id="developers" class="ov2-section ov2-developer-section" data-home-reveal>
    <div class="ov2-shell">
      <div class="ov2-section-intro ov2-section-intro-split">
        <div>
          <span class="ov2-kicker"><i></i>04 / API</span>
          <h2>{{ t('homeV2.developerTitle') }}</h2>
        </div>
        <p>{{ t('homeV2.developerDescription') }}</p>
      </div>

      <ol class="ov2-developer-flow" :aria-label="t('homeV2.developerFlowLabel')">
        <li
          v-for="(step, index) in developerSteps"
          :key="step.title"
          :class="{ 'is-active': activeFlowStep === index }"
        >
          <button type="button" @click="activeFlowStep = index">
            <span>0{{ index + 1 }}</span>
            <div><strong>{{ step.title }}</strong><small>{{ step.description }}</small></div>
          </button>
          <div v-if="index === 0" class="ov2-flow-models">
            <button
              v-for="model in developerModels"
              :key="model.id"
              type="button"
              :aria-pressed="activeModel === model.id"
              :class="{ 'is-selected': activeModel === model.id }"
              @click="selectDeveloperModel(model.id)"
            >
              <ModelIcon :model="model.icon" size="16px" />
              {{ model.label }}
            </button>
          </div>
          <div v-else-if="index === 1" class="ov2-flow-endpoint">
            <small>{{ t('homeV2.generatedEndpoint') }}</small>
            <code>POST /v1/chat/completions</code>
            <span>model: {{ activeModel }}</span>
          </div>
          <div v-else class="ov2-flow-copy">
            <small>{{ t('homeV2.readyToCopy') }}</small>
            <button type="button" @click="copyCode">
              <Icon name="copy" size="xs" />
              {{ copied ? t('common.copied') : t('common.copy') }}
            </button>
          </div>
        </li>
      </ol>

      <div class="ov2-developer-workspace">
        <div class="ov2-code-panel">
          <div class="ov2-panel-bar">
            <div class="ov2-code-tabs" role="tablist" aria-label="Code examples">
              <button
                v-for="snippet in snippets"
                :key="snippet.id"
                type="button"
                :class="{ 'is-active': activeSnippet === snippet.id }"
                @click="activeSnippet = snippet.id"
              >
                {{ snippet.label }}
              </button>
            </div>
            <button type="button" class="ov2-copy-button" @click="copyCode">
              <Icon name="copy" size="xs" /> {{ copied ? t('common.copied') : t('common.copy') }}
            </button>
          </div>

          <div class="ov2-code-file">
            <span>quickstart.{{ selected.extension }}</span>
            <small>Omnio API</small>
          </div>
          <pre class="ov2-code-block"><code>{{ selected.code }}</code></pre>
          <div class="ov2-code-response">
            <span><i></i>200 OK</span>
            <strong>omnio-route: {{ activeModel }}</strong>
            <small>{{ latency }}ms</small>
          </div>
        </div>

        <div class="ov2-dashboard-panel">
          <div class="ov2-panel-bar">
            <div>
              <span class="ov2-dashboard-live"><i></i>{{ t('homeV2.liveDashboard') }}</span>
            </div>
            <span class="ov2-dashboard-range">{{ t('homeV2.lastHour') }}</span>
          </div>

          <div class="ov2-dashboard-metrics">
            <div>
              <span>{{ t('homeV2.requests') }}</span>
              <strong>{{ formattedRequests }}</strong>
              <small>↑ 12.4%</small>
            </div>
            <div>
              <span>{{ t('homeV2.avgLatency') }}</span>
              <strong>{{ latency }}<b>ms</b></strong>
              <small class="is-good">↓ 8.1%</small>
            </div>
            <div>
              <span>{{ t('homeV2.tokens') }}</span>
              <strong>{{ tokenTotal }}<b>k</b></strong>
              <small>↑ 6.8%</small>
            </div>
          </div>

          <div class="ov2-dashboard-chart">
            <div class="ov2-chart-labels"><span>1.2k</span><span>800</span><span>400</span><span>0</span></div>
            <div class="ov2-chart-area">
              <i v-for="(bar, index) in bars" :key="index" :style="{ height: `${bar}%` }"><span></span></i>
            </div>
          </div>

          <div class="ov2-model-status-list">
            <span v-for="model in modelStatus" :key="model.name">
              <i><ModelIcon :model="model.icon" size="16px" /></i>
              <strong>{{ model.name }}</strong>
              <small>{{ model.share }}</small>
              <b>{{ t('homeV2.healthy') }}</b>
            </span>
          </div>
        </div>
      </div>

      <div class="ov2-developer-actions">
        <router-link :to="isAuthenticated ? dashboardPath : '/register'" class="ov2-button ov2-button-primary">
          {{ isAuthenticated ? t('home.goToDashboard') : t('homeV2.getApiKey') }}
          <Icon name="arrowRight" size="xs" />
        </router-link>
        <a :href="docUrl || '/docs/usage-guide/'" class="ov2-text-link">
          {{ t('homeV2.readApiDocs') }} <span>↗</span>
        </a>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'
import ModelIcon from '@/components/common/ModelIcon.vue'

const props = defineProps<{
  gatewayOrigin: string
  docUrl: string
  isAuthenticated: boolean
  dashboardPath: string
}>()

type SnippetId = 'curl' | 'python' | 'javascript'

const { t } = useI18n()
const activeSnippet = ref<SnippetId>('curl')
const activeFlowStep = ref(0)
const copied = ref(false)
const requestTotal = ref(12842)
const latency = ref(118)
const tokenTotal = ref(482)
const activeModel = ref('gpt-5')
let metricTimer: number | undefined
let flowTimer: number | undefined
let copyTimer: number | undefined

const developerSteps = computed(() => [
  { title: t('homeV2.chooseModel'), description: t('homeV2.chooseModelDescription') },
  { title: t('homeV2.createApiKey'), description: t('homeV2.createApiKeyDescription') },
  { title: t('homeV2.makeFirstCall'), description: t('homeV2.makeFirstCallDescription') }
])

const developerModels = [
  { id: 'gpt-5', label: 'GPT-5', icon: 'gpt-5' },
  { id: 'claude-sonnet', label: 'Claude', icon: 'claude' },
  { id: 'gemini', label: 'Gemini', icon: 'gemini' },
  { id: 'grok', label: 'Grok', icon: 'grok' }
]

const snippets = computed(() => [
  {
    id: 'curl' as const,
    label: 'cURL',
    extension: 'sh',
    code: `curl ${props.gatewayOrigin}/v1/chat/completions \\
  -H "Authorization: Bearer $OMNIO_KEY" \\
  -H "Content-Type: application/json" \\
  -d '{
    "model": "${activeModel.value}",
    "messages": [{"role": "user", "content": "Hello"}]
  }'`
  },
  {
    id: 'python' as const,
    label: 'Python',
    extension: 'py',
    code: `from openai import OpenAI

client = OpenAI(
    base_url="${props.gatewayOrigin}/v1",
    api_key=OMNIO_KEY
)

response = client.chat.completions.create(
    model="${activeModel.value}",
    messages=[{"role": "user", "content": "Hello"}]
)`
  },
  {
    id: 'javascript' as const,
    label: 'JavaScript',
    extension: 'js',
    code: `const client = new OpenAI({
  baseURL: '${props.gatewayOrigin}/v1',
  apiKey: process.env.OMNIO_KEY
})

const response = await client.chat.completions.create({
  model: '${activeModel.value}',
  messages: [{ role: 'user', content: 'Hello' }]
})`
  }
])

const selected = computed(() => snippets.value.find((snippet) => snippet.id === activeSnippet.value) || snippets.value[0])
const formattedRequests = computed(() => new Intl.NumberFormat().format(requestTotal.value))
const bars = [26, 34, 31, 42, 38, 47, 44, 56, 53, 62, 58, 71, 67, 76, 72, 84, 79, 88, 82, 91]
const modelStatus = [
  { name: 'GPT', icon: 'gpt-5', share: '42%' },
  { name: 'Claude', icon: 'claude', share: '31%' },
  { name: 'Gemini', icon: 'gemini', share: '18%' }
]

function selectDeveloperModel(model: string) {
  activeModel.value = model
  activeFlowStep.value = 1
}

async function copyCode() {
  activeFlowStep.value = 2
  try {
    await navigator.clipboard.writeText(selected.value.code)
    copied.value = true
    if (copyTimer) window.clearTimeout(copyTimer)
    copyTimer = window.setTimeout(() => { copied.value = false }, 1800)
  } catch {
    copied.value = false
  }
}

onMounted(() => {
  if (window.matchMedia('(prefers-reduced-motion: reduce)').matches) return
  metricTimer = window.setInterval(() => {
    requestTotal.value += Math.floor(Math.random() * 8) + 3
    latency.value = Math.max(96, Math.min(138, latency.value + Math.floor(Math.random() * 7) - 3))
    tokenTotal.value += Math.random() > 0.55 ? 1 : 0
  }, 1800)
  flowTimer = window.setInterval(() => {
    activeFlowStep.value = (activeFlowStep.value + 1) % developerSteps.value.length
  }, 2400)
})

onBeforeUnmount(() => {
  if (metricTimer) window.clearInterval(metricTimer)
  if (flowTimer) window.clearInterval(flowTimer)
  if (copyTimer) window.clearTimeout(copyTimer)
})
</script>
