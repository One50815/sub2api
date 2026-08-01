<template>
  <section id="network" class="ov2-section ov2-network-section" data-home-reveal>
    <div class="ov2-shell">
      <div class="ov2-section-intro ov2-section-intro-split">
        <div>
          <span class="ov2-kicker"><i></i>03 / {{ t('homeV2.networkKicker') }}</span>
          <h2>{{ t('homeV2.networkTitle') }}</h2>
        </div>
        <p>{{ t('homeV2.networkDescription') }}</p>
      </div>

      <div class="ov2-network-console">
        <div class="ov2-network-console-bar">
          <div>
            <span class="ov2-console-light"></span>
            <strong>{{ t('homeV2.liveNetwork') }}</strong>
          </div>
          <div class="ov2-network-health">
            <span><i></i>{{ t('homeV2.allSystems') }}</span>
            <span>{{ currentTime }}</span>
          </div>
        </div>

        <div class="ov2-network-console-body">
          <div :class="['ov2-network-map', { 'is-routing': isRouting }]">
            <svg viewBox="0 0 760 390" aria-hidden="true">
              <path class="ov2-map-line" d="M136 62 C 254 28, 258 218, 398 186" />
              <path class="ov2-map-line" d="M650 118 C 524 64, 544 252, 398 186" />
              <path class="ov2-map-line" d="M174 330 C 276 362, 266 224, 398 186" />
              <path class="ov2-map-line" d="M612 338 C 536 370, 518 210, 398 186" />
              <path class="ov2-map-line" d="M72 202 C 182 132, 274 240, 398 186" />
              <path class="ov2-map-flow ov2-map-flow-a" d="M136 62 C 254 28, 258 218, 398 186" />
              <path class="ov2-map-flow ov2-map-flow-b" d="M650 118 C 524 64, 544 252, 398 186" />
              <path class="ov2-map-flow ov2-map-flow-c" d="M174 330 C 276 362, 266 224, 398 186" />
              <path class="ov2-map-flow ov2-map-flow-d" d="M612 338 C 536 370, 518 210, 398 186" />
              <path class="ov2-map-flow ov2-map-flow-e" d="M72 202 C 182 132, 274 240, 398 186" />
            </svg>

            <span class="ov2-orbit-ring ov2-orbit-ring-a" aria-hidden="true"></span>
            <span class="ov2-orbit-ring ov2-orbit-ring-b" aria-hidden="true"></span>
            <span class="ov2-orbit-ring ov2-orbit-ring-c" aria-hidden="true"></span>
            <span class="ov2-galaxy-signal ov2-galaxy-signal-a" aria-hidden="true"></span>
            <span class="ov2-galaxy-signal ov2-galaxy-signal-b" aria-hidden="true"></span>
            <span class="ov2-galaxy-signal ov2-galaxy-signal-c" aria-hidden="true"></span>

            <button
              v-for="(node, index) in nodes"
              :key="node.name"
              type="button"
              :class="['ov2-map-node', `ov2-map-node-${index}`, { 'is-selected': selectedNode === index }]"
              :style="{ '--model-color': node.color }"
              @click="selectNode(index)"
            >
              <span class="ov2-map-node-top">
                <span><ModelIcon :model="node.icon" size="20px" />{{ node.name }}</span>
                <i></i>
              </span>
              <strong>{{ node.model }}</strong>
              <span class="ov2-map-node-meta"><small>{{ node.latency }}ms</small><small>{{ node.price }}</small></span>
              <span class="ov2-pixel-meter" aria-hidden="true">
                <i
                  v-for="segment in 8"
                  :key="segment"
                  :class="{ 'is-off': segment > signalLevel(node.latency) }"
                ></i>
              </span>
            </button>

            <div class="ov2-map-core">
              <OmnioMark :label="`${siteName} Core`" />
              <strong>{{ siteName }}</strong>
              <span>{{ t('homeV2.oneGateway') }}</span>
            </div>
          </div>

          <aside class="ov2-network-detail">
            <span class="ov2-detail-label">{{ t('homeV2.selectedRoute') }}</span>
            <div class="ov2-detail-provider" :style="{ '--model-color': selected.color }">
              <span><ModelIcon :model="selected.icon" size="28px" /></span>
              <div><strong>{{ selected.model }}</strong><small>{{ selected.name }}</small></div>
              <i></i>
            </div>
            <dl>
              <div><dt>{{ t('homeV2.status') }}</dt><dd>{{ t('homeV2.online') }}</dd></div>
              <div><dt>{{ t('homeV2.latency') }}</dt><dd>{{ selected.latency }}ms</dd></div>
              <div><dt>{{ t('homeV2.priceMode') }}</dt><dd>{{ t('homeV2.usageBased') }}</dd></div>
              <div><dt>{{ t('homeV2.routeQuality') }}</dt><dd>99.98%</dd></div>
            </dl>
            <div class="ov2-request-feed">
              <span v-for="(request, index) in requests" :key="request.id">
                <i :style="{ animationDelay: `${index * -0.35}s` }"></i>
                <code>{{ request.id }}</code>
                <small>{{ request.time }}</small>
                <b>200</b>
              </span>
            </div>
            <router-link to="/available-channels" class="ov2-text-link">
              {{ t('homeV2.viewAllModels') }} <span>↗</span>
            </router-link>
          </aside>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import ModelIcon from '@/components/common/ModelIcon.vue'
import OmnioMark from './OmnioMark.vue'

defineProps<{
  siteName: string
}>()

const { t } = useI18n()
const selectedNode = ref(0)
const currentTime = ref('')
const isRouting = ref(false)
let clockTimer: number | undefined
let networkTimer: number | undefined
let routeTimer: number | undefined
let routingFeedbackTimer: number | undefined

const nodes = ref([
  { name: 'OpenAI', icon: 'gpt-5', model: 'GPT-5', color: '#35c889', latency: 120, price: 'LIVE RATE' },
  { name: 'Anthropic', icon: 'claude', model: 'Claude', color: '#ff9f43', latency: 138, price: 'LIVE RATE' },
  { name: 'Google', icon: 'gemini', model: 'Gemini', color: '#7b6cff', latency: 96, price: 'LIVE RATE' },
  { name: 'xAI', icon: 'grok', model: 'Grok', color: '#3f8cff', latency: 111, price: 'LIVE RATE' },
  { name: 'Alibaba', icon: 'qwen', model: 'Qwen', color: '#9a62ff', latency: 126, price: 'LIVE RATE' }
])

const requests = ref([
  { id: 'req_7f31', time: '122ms' },
  { id: 'req_7f30', time: '98ms' },
  { id: 'req_7f2f', time: '141ms' }
])

const selected = computed(() => nodes.value[selectedNode.value])

function signalLevel(latency: number) {
  return Math.max(3, Math.min(8, 10 - Math.floor(latency / 25)))
}

function updateClock() {
  currentTime.value = new Intl.DateTimeFormat(undefined, {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false
  }).format(new Date())
}

function selectNode(index: number) {
  selectedNode.value = index
  isRouting.value = true
  if (routingFeedbackTimer) window.clearTimeout(routingFeedbackTimer)
  routingFeedbackTimer = window.setTimeout(() => {
    isRouting.value = false
  }, 560)
}

onMounted(() => {
  updateClock()
  clockTimer = window.setInterval(updateClock, 1000)
  if (window.matchMedia('(prefers-reduced-motion: reduce)').matches) return

  networkTimer = window.setInterval(() => {
    nodes.value = nodes.value.map((node) => ({
      ...node,
      latency: Math.max(78, Math.min(164, node.latency + Math.floor(Math.random() * 11) - 5))
    }))
    const nextLatency = nodes.value[selectedNode.value].latency
    requests.value = [
      {
        id: `req_${Math.floor(Math.random() * 0xffff).toString(16).padStart(4, '0')}`,
        time: `${nextLatency + Math.floor(Math.random() * 13) - 6}ms`
      },
      ...requests.value
    ].slice(0, 3)
  }, 1500)

  routeTimer = window.setInterval(() => {
    selectNode((selectedNode.value + 1) % nodes.value.length)
  }, 3600)
})

onBeforeUnmount(() => {
  if (clockTimer) window.clearInterval(clockTimer)
  if (networkTimer) window.clearInterval(networkTimer)
  if (routeTimer) window.clearInterval(routeTimer)
  if (routingFeedbackTimer) window.clearTimeout(routingFeedbackTimer)
})
</script>
