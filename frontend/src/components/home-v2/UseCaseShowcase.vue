<template>
  <section
    id="use-cases"
    ref="sectionRoot"
    class="ov2-section ov2-use-section"
    data-home-reveal
  >
    <div class="ov2-shell">
      <div class="ov2-section-intro">
        <span class="ov2-kicker"><i></i>02 / {{ t('homeV2.everyoneKicker') }}</span>
        <h2>{{ t('homeV2.everyoneTitle') }}</h2>
        <p>{{ t('homeV2.everyoneDescription') }}</p>
      </div>

      <div class="ov2-use-workspace">
        <div class="ov2-use-tabs" role="tablist" :aria-label="t('homeV2.useCases')">
          <button
            v-for="(item, index) in useCases"
            :key="item.key"
            type="button"
            role="tab"
            :aria-selected="activeUseCase === index"
            :class="{ 'is-active': activeUseCase === index }"
            @click="activeUseCase = index"
          >
            <span class="ov2-use-tab-index">0{{ index + 1 }}</span>
            <span><strong>{{ item.title }}</strong><small>{{ item.short }}</small></span>
            <i>→</i>
          </button>
        </div>

        <div class="ov2-use-preview">
          <div class="ov2-use-preview-bar">
            <span><i></i>{{ t('homeV2.smartWorkspace') }}</span>
            <small>{{ t('homeV2.autoMode') }}</small>
          </div>

          <div class="ov2-use-preview-body">
            <form class="ov2-use-composer" @submit.prevent="runPrompt">
              <label for="omnio-use-prompt">{{ t('homeV2.promptLabel') }}</label>
              <div>
                <input
                  id="omnio-use-prompt"
                  v-model="draftPrompt"
                  type="text"
                  :placeholder="active.prompt"
                  autocomplete="off"
                  @focus="stopAutoTyping"
                  @input="stopAutoTyping"
                />
                <button type="submit">
                  {{ isThinking ? t('homeV2.choosingBest') : t('homeV2.sendPrompt') }}
                  <span>↗</span>
                </button>
              </div>
            </form>

            <div :class="['ov2-use-decision', { 'is-thinking': isThinking, 'is-typing': isAutoTyping }]">
              <span>{{ isThinking ? t('homeV2.choosingBest') : t('homeV2.omnioSelected') }}</span>
              <div class="ov2-use-candidates">
                <button
                  v-for="(candidate, index) in active.candidates"
                  :key="candidate.name"
                  type="button"
                  :class="{ 'is-selected': !isThinking && selectedCandidate === index }"
                  @click="selectCandidate(index)"
                >
                  <ModelIcon :model="candidate.icon" size="17px" />
                  <strong>{{ candidate.name }}</strong>
                  <i></i>
                </button>
              </div>
            </div>

            <div :class="['ov2-companion-message', { 'is-thinking': isThinking || isAutoTyping }]">
              <span class="ov2-companion-avatar"><OmnioMark /></span>
              <div>
                <small>{{ t('homeV2.companionLabel') }}</small>
                <p>
                  {{
                    isThinking || isAutoTyping
                      ? t('homeV2.companionChoosing')
                      : t('homeV2.companionSelected', {
                        model: selectedModel.name,
                        reason: active.reason
                      })
                  }}
                </p>
              </div>
            </div>

            <div :class="['ov2-use-response', { 'is-thinking': isThinking }]">
              <span class="ov2-use-result-badge">{{ selectedModel.name }}</span>
              <div>
                <small>{{ isThinking ? t('homeV2.modelOrbiting') : `Omnio · ${selectedModel.name}` }}</small>
                <h3>{{ isThinking ? t('homeV2.findingModel') : active.responseTitle }}</h3>
                <p>{{ isThinking ? t('homeV2.findingModelDescription') : active.response }}</p>
                <div class="ov2-use-result-lines"><i></i><i></i><i></i></div>
              </div>
            </div>
          </div>

          <div class="ov2-use-preview-footer">
            <span><i></i>{{ isThinking ? t('homeV2.choosingBest') : t('homeV2.modelReady') }}</span>
            <strong>{{ t('homeV2.noModelResearch') }}</strong>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import ModelIcon from '@/components/common/ModelIcon.vue'
import OmnioMark from './OmnioMark.vue'

const { t } = useI18n()
const sectionRoot = ref<HTMLElement | null>(null)
const activeUseCase = ref(0)
const draftPrompt = ref('')
const isThinking = ref(false)
const isAutoTyping = ref(false)
const selectedCandidate = ref(0)
let thinkingTimer: number | undefined
let typingTimer: number | undefined
let demoObserver: IntersectionObserver | undefined

const useCases = computed(() => [
  {
    key: 'chat',
    title: t('homeV2.chat'),
    short: t('homeV2.chatShort'),
    prompt: t('homeV2.chatPrompt'),
    candidates: [
      { name: 'Claude', icon: 'claude' },
      { name: 'GPT', icon: 'gpt-5' },
      { name: 'Gemini', icon: 'gemini' }
    ],
    reason: t('homeV2.longContext'),
    responseTitle: t('homeV2.chatResultTitle'),
    response: t('homeV2.chatResult')
  },
  {
    key: 'writing',
    title: t('homeV2.writing'),
    short: t('homeV2.writingShort'),
    prompt: t('homeV2.writingPrompt'),
    candidates: [
      { name: 'GPT', icon: 'gpt-5' },
      { name: 'Claude', icon: 'claude' },
      { name: 'Gemini', icon: 'gemini' }
    ],
    reason: t('homeV2.creativeReasoning'),
    responseTitle: t('homeV2.writingResultTitle'),
    response: t('homeV2.writingResult')
  },
  {
    key: 'code',
    title: t('homeV2.code'),
    short: t('homeV2.codeShort'),
    prompt: t('homeV2.codePrompt'),
    candidates: [
      { name: 'Grok', icon: 'grok' },
      { name: 'GPT', icon: 'gpt-5' },
      { name: 'Claude', icon: 'claude' }
    ],
    reason: t('homeV2.codeReasoning'),
    responseTitle: t('homeV2.codeResultTitle'),
    response: t('homeV2.codeResult')
  },
  {
    key: 'analysis',
    title: t('homeV2.analysis'),
    short: t('homeV2.analysisShort'),
    prompt: t('homeV2.analysisPrompt'),
    candidates: [
      { name: 'GPT', icon: 'gpt-5' },
      { name: 'Gemini', icon: 'gemini' },
      { name: 'Qwen', icon: 'qwen' }
    ],
    reason: t('homeV2.dataReasoning'),
    responseTitle: t('homeV2.analysisResultTitle'),
    response: t('homeV2.analysisResult')
  }
])

const active = computed(() => useCases.value[activeUseCase.value])
const selectedModel = computed(() => (
  active.value.candidates[selectedCandidate.value] || active.value.candidates[0]
))

function selectCandidate(index: number) {
  stopAutoTyping()
  if (thinkingTimer) window.clearTimeout(thinkingTimer)
  isThinking.value = false
  selectedCandidate.value = index
}

function runPrompt() {
  stopAutoTyping()
  if (!draftPrompt.value.trim()) draftPrompt.value = active.value.prompt
  if (thinkingTimer) window.clearTimeout(thinkingTimer)
  selectedCandidate.value = -1
  isThinking.value = true
  thinkingTimer = window.setTimeout(() => {
    selectedCandidate.value = 0
    isThinking.value = false
  }, 980)
}

function stopAutoTyping() {
  isAutoTyping.value = false
  if (typingTimer) {
    window.clearTimeout(typingTimer)
    typingTimer = undefined
  }
}

function startAutoDemo() {
  stopAutoTyping()
  draftPrompt.value = ''
  selectedCandidate.value = -1
  isThinking.value = false
  isAutoTyping.value = true

  const prompt = active.value.prompt
  let characterIndex = 0
  const typeNextCharacter = () => {
    if (!isAutoTyping.value) return
    characterIndex += 1
    draftPrompt.value = prompt.slice(0, characterIndex)
    if (characterIndex < prompt.length) {
      typingTimer = window.setTimeout(typeNextCharacter, 42)
      return
    }

    isAutoTyping.value = false
    typingTimer = window.setTimeout(() => runPrompt(), 260)
  }

  typingTimer = window.setTimeout(typeNextCharacter, 180)
}

watch(active, (next) => {
  stopAutoTyping()
  if (thinkingTimer) window.clearTimeout(thinkingTimer)
  draftPrompt.value = next.prompt
  selectedCandidate.value = 0
  isThinking.value = false
}, { immediate: true })

onMounted(() => {
  if (
    !sectionRoot.value
    || window.matchMedia('(prefers-reduced-motion: reduce)').matches
    || !('IntersectionObserver' in window)
  ) return

  demoObserver = new IntersectionObserver((entries) => {
    if (!entries.some((entry) => entry.isIntersecting)) return
    demoObserver?.disconnect()
    startAutoDemo()
  }, {
    threshold: 0.45
  })
  demoObserver.observe(sectionRoot.value)
})

onBeforeUnmount(() => {
  demoObserver?.disconnect()
  stopAutoTyping()
  if (thinkingTimer) window.clearTimeout(thinkingTimer)
})
</script>
