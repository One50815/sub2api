<template>
  <div
    :class="['ov2-hero-network', `is-route-${activeIndex}`, { 'is-interacting': isInteracting }]"
    aria-label="Live AI model routing preview"
  >
    <svg class="ov2-network-lines" viewBox="0 0 640 520" aria-hidden="true">
      <path class="ov2-network-line" d="M112 58 C 220 54, 230 198, 346 234" />
      <path class="ov2-network-line" d="M540 100 C 450 80, 440 188, 346 234" />
      <path class="ov2-network-line" d="M104 390 C 220 410, 224 268, 346 234" />
      <path class="ov2-network-line" d="M548 432 C 440 414, 438 280, 346 234" />
      <path class="ov2-network-line" d="M74 236 C 180 190, 236 250, 346 234" />
      <path class="ov2-network-flow ov2-network-flow-0" d="M112 58 C 220 54, 230 198, 346 234" />
      <path class="ov2-network-flow ov2-network-flow-1" d="M540 100 C 450 80, 440 188, 346 234" />
      <path class="ov2-network-flow ov2-network-flow-2" d="M104 390 C 220 410, 224 268, 346 234" />
      <path class="ov2-network-flow ov2-network-flow-3" d="M548 432 C 440 414, 438 280, 346 234" />
      <path class="ov2-network-flow ov2-network-flow-4" d="M74 236 C 180 190, 236 250, 346 234" />
    </svg>

    <span class="ov2-pixel-module ov2-pixel-module-a" aria-hidden="true"><i></i><i></i><i></i></span>
    <span class="ov2-pixel-module ov2-pixel-module-b" aria-hidden="true"><i></i><i></i></span>
    <span class="ov2-pixel-module ov2-pixel-module-c" aria-hidden="true"><i></i><i></i><i></i></span>

    <div class="ov2-hero-live-status">
      <span><i></i>{{ t('homeV2.liveModelCount') }}</span>
      <small>{{ t('homeV2.liveRouteHealth') }}</small>
    </div>

    <button
      v-for="(model, index) in models"
      :key="model.name"
      type="button"
      :class="['ov2-hero-node', `ov2-hero-node-${index}`, { 'is-active': activeIndex === index }]"
      :style="{ '--model-color': model.color }"
      @click="selectModel(index)"
    >
      <span class="ov2-hero-node-icon">
        <ModelIcon :model="model.icon" size="22px" />
      </span>
      <span class="ov2-hero-node-copy">
        <strong>{{ model.name }}</strong>
        <small><i></i>{{ activeIndex === index ? t('homeV2.routingLive') : t('homeV2.online') }}</small>
      </span>
      <span class="ov2-hero-node-latency">{{ model.latency }}</span>
      <span class="ov2-pixel-meter" aria-hidden="true">
        <i
          v-for="segment in 8"
          :key="segment"
          :class="{ 'is-off': segment > model.signal }"
        ></i>
      </span>
    </button>

    <div class="ov2-network-core">
      <span class="ov2-network-core-halo"></span>
      <span class="ov2-core-gate" aria-hidden="true"><i></i><i></i></span>
      <OmnioMark :label="siteName" />
      <strong>{{ siteName }} Core</strong>
      <small>{{ t('homeV2.autoRouting') }}</small>
    </div>

    <button
      type="button"
      class="ov2-hero-operator"
      :aria-label="t('homeV2.mascotAction')"
      @click="surpriseRoute"
    >
      <img
        src="/assets/brand/omnio-pixel-mascot-connector-v2.png"
        :alt="t('homeV2.mascotAlt')"
        width="820"
        height="620"
      />
      <span><i></i>{{ t('homeV2.mascotWorking') }}</span>
    </button>

    <div class="ov2-route-readout" aria-live="polite">
      <span>{{ t('homeV2.routingRequest') }}</span>
      <strong>{{ models[activeIndex].name }}</strong>
      <small>{{ models[activeIndex].reason }}</small>
    </div>
  </div>
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
const activeIndex = ref(0)
const latencies = ref([120, 138, 96, 111, 126])
const isInteracting = ref(false)
let routeTimer: number | undefined
let interactionTimer: number | undefined

const models = computed(() => [
  {
    name: 'GPT',
    icon: 'gpt-5',
    color: '#35c889',
    latency: `${latencies.value[0]}ms`,
    signal: Math.max(3, Math.min(8, 10 - Math.floor(latencies.value[0] / 25))),
    reason: t('homeV2.routeReasoning')
  },
  {
    name: 'Claude',
    icon: 'claude',
    color: '#ff9f43',
    latency: `${latencies.value[1]}ms`,
    signal: Math.max(3, Math.min(8, 10 - Math.floor(latencies.value[1] / 25))),
    reason: t('homeV2.routeContext')
  },
  {
    name: 'Gemini',
    icon: 'gemini',
    color: '#7b6cff',
    latency: `${latencies.value[2]}ms`,
    signal: Math.max(3, Math.min(8, 10 - Math.floor(latencies.value[2] / 25))),
    reason: t('homeV2.routeMultimodal')
  },
  {
    name: 'Grok',
    icon: 'grok',
    color: '#3f8cff',
    latency: `${latencies.value[3]}ms`,
    signal: Math.max(3, Math.min(8, 10 - Math.floor(latencies.value[3] / 25))),
    reason: t('homeV2.routeRealtime')
  },
  {
    name: 'Qwen',
    icon: 'qwen',
    color: '#9a62ff',
    latency: `${latencies.value[4]}ms`,
    signal: Math.max(3, Math.min(8, 10 - Math.floor(latencies.value[4] / 25))),
    reason: t('homeV2.routeEfficient')
  }
])

function selectModel(index: number) {
  activeIndex.value = index
  isInteracting.value = true
  if (interactionTimer) window.clearTimeout(interactionTimer)
  interactionTimer = window.setTimeout(() => {
    isInteracting.value = false
  }, 620)
}

function surpriseRoute() {
  selectModel((activeIndex.value + 1) % models.value.length)
}

onMounted(() => {
  if (window.matchMedia('(prefers-reduced-motion: reduce)').matches) return
  routeTimer = window.setInterval(() => {
    activeIndex.value = (activeIndex.value + 1) % models.value.length
    latencies.value = latencies.value.map((latency) => (
      Math.max(82, Math.min(156, latency + Math.floor(Math.random() * 9) - 4))
    ))
  }, 2600)
})

onBeforeUnmount(() => {
  if (routeTimer) window.clearInterval(routeTimer)
  if (interactionTimer) window.clearTimeout(interactionTimer)
})
</script>
